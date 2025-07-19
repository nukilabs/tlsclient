# TLS Client

[![Go Reference](https://pkg.go.dev/badge/github.com/nukilabs/tlsclient.svg)](https://pkg.go.dev/github.com/nukilabs/tlsclient)

A powerful Go HTTP client library that provides advanced TLS fingerprinting capabilities, allowing you to emulate various browsers and devices for web scraping, testing, and automation tasks.

## Features

- 🌐 **Browser Emulation**: Mimic TLS fingerprints of Chrome, Safari, OkHttp4, and HttpUrlConnection
- 🚀 **HTTP Protocol Support**: HTTP/1.1, HTTP/2, and HTTP/3
- 🗜️ **Automatic Decompression**: Built-in support for gzip, brotli, zstd, and deflate
- 🪝 **Hooks**: Injectable hooks that are performed after each request
- 🔄 **Smart Redirects**: Configurable redirect handling with custom policies
- 📊 **Bandwidth Tracking**: Monitor network usage and performance
- 🎯 **Certificate Pinning**: Enhanced security with certificate validation
- 🌍 **Proxy Support**: HTTP, HTTPS, and SOCKS proxy support
- ⚡ **High Performance**: Optimized for speed and efficiency

## Installation

```bash
go get github.com/nukilabs/tlsclient
```

## Quick Start

```go
package main

import (
    "fmt"
    "io"
    "log"
    "net/url"

    "github.com/nukilabs/http"
    "github.com/nukilabs/tlsclient"
    "github.com/nukilabs/tlsclient/profiles"
    "github.com/nukilabs/tlsclient/redirects"
)

func main() {
    // Create a new client with Chrome profile
    client := tlsclient.New(profiles.Chrome138)
    client.SetRedirectFunc(redirects.Chrome) // Use Chrome's redirect policy
    client.SetProxy(&url.URL{
        Scheme: "http",
        Host:   "localhost:8888",
    })
    
    // Make a GET request
    req, err := client.NewRequest(http.MethodGet, "https://example.com", nil)
    if err != nil {
        log.Fatal(err)
    }

    // Add chrome-like headers
    req.Header = http.Header{
		"sec-ch-ua":                 {"\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\""},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {"\"Windows\""},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"sec-fetch-site":            {"none"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-user":            {"?1"},
		"sec-fetch-dest":            {"document"},
		"accept-encoding":           {"gzip, deflate, br, zstd"},
		"accept-language":           {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"priority":                  {"u=0, i"},
		http.HeaderOrderKey:         {"sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-platform", "upgrade-insecure-requests", "user-agent", "accept", "sec-fetch-site", "sec-fetch-mode", "sec-fetch-user", "sec-fetch-dest", "accept-encoding", "accept-language", "priority"},
	}

    // Perform the request
    res, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    
    // Read and print the response
    body, err := io.ReadAll(res.Body)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(string(body))
}
```

## Available Profiles

### Chrome Browsers
```go
// Chrome versions 120-138
client := tlsclient.New(profiles.Chrome138)  // Latest Chrome
client := tlsclient.New(profiles.Chrome120)  // Older Chrome version
```

### Safari Browsers
```go
// Safari 17
client := tlsclient.New(profiles.Safari17)
```

### Android HttpUrlConnection
```go
// Android API levels 21-35
client := tlsclient.New(profiles.HttpUrlConnectionAndroid(33))  // Android 13
client := tlsclient.New(profiles.HttpUrlConnectionAndroid(21))  // Android 5.0
```

### Android OkHttp4
```go
// Android API levels 21-35
client := tlsclient.New(profiles.Okhttp4Android(33))  // Android 13
client := tlsclient.New(profiles.Okhttp4Android(29))  // Android 10
```

## Configuration Options

### Basic Options

```go
client := tlsclient.New(profiles.Chrome138,
    tlsclient.WithTimeout(30*time.Second),  // Set request timeout
    tlsclient.WithNoCookieJar(),            // Disable cookie jar
    tlsclient.WithNoAutoDecompress(),       // Disable automatic decompression
    tlsclient.WithNoFollowRedirects(),      // Disable automatic redirects
)
```

### Advanced Transport Options

```go
client := tlsclient.New(profiles.Chrome138,
    tlsclient.WithTransportOptions(tlsclient.TransportOptions{
        ServerNameOverride: "example.com",
        InsecureSkipVerify: true,
        DisableKeepAlives:  false,
        IdleConnTimeout:    90 * time.Second,
        DisableIPV4:        false,
        DisableIPV6:        false,
    }),
)
```

### Certificate Pinning

```go
client := tlsclient.New(profiles.Chrome138,
    tlsclient.WithAutoPinning(),
)
```

## Redirect Handling

### Disable Redirects

```go
client.SetFollowRedirects(false)
```

### Custom Redirect Policy

```go
client.SetRedirectFunc(func(req *http.Request, via []*http.Request) error {
    // Maximum of 5 redirects
    if len(via) >= 5 {
        return errors.New("too many redirects")
    }
    
    // Don't follow redirects to different hosts
    if req.URL.Host != via[0].URL.Host {
        return errors.New("cross-host redirect not allowed")
    }
    
    return nil
})
```

### Built-in Redirect Policies

```go
client.SetRedirectFunc(redirects.Chrome) // Use Chrome's redirect policy
```

## Proxy Support

```go
// HTTP Proxy
client.SetProxy(&url.URL{
    Scheme: "http",
    Host:   "proxy.example.com:8080",
})

// SOCKS5 Proxy
client.SetProxy(&url.URL{
    Scheme: "socks5",
    Host:   "proxy.example.com:1080",
})

// Authenticated Proxy
client.SetProxy(&url.URL{
    Scheme: "http",
    Host:   "proxy.example.com:8080",
    User:   url.UserPassword("username", "password"),
})
```

## Response Hooks

```go
func GotBlockedHook(client *tlsclient.Client, res *http.Response) (*http.Response, error) {
    if res.StatusCode == http.StatusTooManyRequests {
        // Switch to a diffrent proxy if blocked
        client.SetProxy(&url.URL{
            Scheme: "http",
            Host:   "localhost:8889",
        })
    }
}

client.AddHook(GotBlockedHook)
```

## Best Practices

### Connection Reuse

```go
// Reuse the same client for multiple requests
client := tlsclient.New(profiles.Chrome138)

// Make multiple requests with the same client
for _, url := range urls {
    res, err := client.Get(url)
    if err != nil {
        continue
    }
    // Process response...
    res.Body.Close()
}
```

### Resource Management

```go
res, err := client.Get("https://example.com")
if err != nil {
    return err
}
defer res.Body.Close() // Always close response body

// Read response body
body, err := io.ReadAll(res.Body)
if err != nil {
    return err
}
```

## Contributing

We welcome contributions from the community! This TLS Client is constantly evolving, and we'd love your help in making it even better. Here are some ways you can contribute:

### 🌐 TLS Profiles

We're always looking to expand our browser and device emulation capabilities. If you have access to a browser, device, or HTTP client that isn't currently supported, we'd love to add it!

**How to contribute a new TLS profile:**

1. **Capture the TLS fingerprint** using tools like:
   - [Powhttp HTTP Proxy App](https://powhttp.com/)
   - [Peets TLS fingerprinting API](https://tls.peet.ws/api/full)
   - Network analysis tools (Wireshark, etc.)

2. **Create the profile** by adding it to the appropriate file in `/profiles/`:
   ```go
   // profiles/newbrowser.go
   var NewBrowser = ClientProfile{
       ClientHelloSpec: func() *tls.ClientHelloSpec {
           return &tls.ClientHelloSpec{
               CipherSuites: []uint16{
                   // Your cipher suites here
               },
               Extensions: []tls.TLSExtension{
                   // Your extensions here
               },
               // ... other TLS parameters
           }
       },
       Settings: []http2.Setting{
           // HTTP/2 settings
       },
       // ... other profile parameters
   }
   ```

3. **Add tests** in `/tests/` to verify the profile works correctly

4. **Update documentation** to include the new profile in the README

**Browsers/Clients we'd love to see:**
- Firefox (all versions)
- Edge (Chromium-based)
- Opera
- Brave
- Mobile browsers (iOS Safari, Chrome Mobile, Samsung Internet)
- HTTP clients (curl, wget, Python requests, etc.)

### 🔄 Redirect Functions

Different browsers handle redirects in unique ways. We're building a comprehensive library of redirect behaviors!

**How to contribute redirect functions:**

1. **Study browser behavior** by testing how your target browser handles:
   - Different redirect status codes (301, 302, 303, 307, 308)
   - Cross-origin redirects
   - Protocol changes (HTTP to HTTPS)
   - Maximum redirect limits
   - Header preservation during redirects

2. **Implement the function** in `/redirects/`:
   ```go
   // redirects/newbrowser.go
   func NewBrowser(req *http.Request, via []*http.Request) error {
       // Implement browser-specific redirect logic
       if len(via) >= 20 { // Browser's max redirect limit
           return errors.New("too many redirects")
       }
       
       // Add any browser-specific redirect handling
       return nil
   }
   ```

3. **Add comprehensive tests** to verify the redirect behavior matches the real browser

4. **Document the behavior** with comments explaining the browser's specific quirks

We appreciate all contributions, no matter how small! 🙏

## Acknowledgments

- Uses [nukilabs/utls](https://github.com/nukilabs/utls) for TLS fingerprinting
- Uses [nukilabs/http](https://github.com/nukilabs/http) for custom HTTP1.1 and HTTP/2 support
- Uses [nukilabs/quic-qo](https://github.com/nukilabs/quic-qo) for HTTP/3 support
- Inspired by: 
    -   [refraction-networking/utls](https://github.com/refraction-networking/utls)
    -   [bogdanfinn/tls-client](https://github.com/bogdanfinn/tls-client)
    -   [useflyent/fhttp](https://github.com/useflyent/fhttp)
---

© 2025 nukilabs
