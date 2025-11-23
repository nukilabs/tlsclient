package redirects

import (
	"maps"
	"net/url"
	"strings"

	"github.com/nukilabs/http"
)

var chromeallowed = map[string]struct{}{
	"host":                            {},
	"connection":                      {},
	"cache-control":                   {},
	"upgrade-insecure-requests":       {},
	"user-agent":                      {},
	"accept":                          {},
	"x-browser-channel":               {},
	"x-browser-year":                  {},
	"x-browser-validation":            {},
	"x-browser-copyright":             {},
	"x-chrome-id-consistency-request": {},
	"x-client-data":                   {},
	"sec-fetch-site":                  {},
	"sec-fetch-mode":                  {},
	"sec-fetch-user":                  {},
	"sec-fetch-dest":                  {},
	"sec-ch-device-memory":            {},
	"sec-ch-ua":                       {},
	"sec-ch-ua-mobile":                {},
	"sec-ch-ua-full-version":          {},
	"sec-ch-ua-arch":                  {},
	"sec-ch-ua-platform":              {},
	"sec-ch-ua-platform-version":      {},
	"sec-ch-ua-model":                 {},
	"sec-ch-ua-bitness":               {},
	"sec-ch-ua-wow64":                 {},
	"sec-ch-ua-full-version-list":     {},
	"sec-ch-ua-form-factors":          {},
	"referer":                         {},
	"accept-encoding":                 {},
	"accept-language":                 {},
	"cookie":                          {},
	"priority":                        {},
}

func Chrome(req *http.Request, via []*http.Request) error {
	maps.DeleteFunc(req.Header, func(k string, v []string) bool {
		_, ok := chromeallowed[strings.ToLower(k)]
		return !ok
	})

	req.Header[http.HeaderOrderKey] = []string{"host", "connection", "cache-control", "upgrade-insecure-requests", "user-agent", "accept", "x-browser-channel", "x-browser-year", "x-browser-validation", "x-browser-copyright", "x-chrome-id-consistency-request", "x-client-data", "sec-fetch-site", "sec-fetch-mode", "sec-fetch-user", "sec-fetch-dest", "sec-ch-device-memory", "sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-full-version", "sec-ch-ua-arch", "sec-ch-ua-platform", "sec-ch-ua-platform-version", "sec-ch-ua-model", "sec-ch-ua-bitness", "sec-ch-ua-wow64", "sec-ch-ua-full-version-list", "sec-ch-ua-form-factors", "referer", "accept-encoding", "accept-language", "cookie", "priority"}

	if len(via) > 0 {
		init := via[0]
		prev := via[len(via)-1]
		src := prev.URL
		dst := req.URL

		if prev.Header.Get("Sec-Fetch-Site") == "none" {
			req.Header.Set("Sec-Fetch-Site", "none")
		} else if strings.EqualFold(src.Scheme, dst.Scheme) && strings.EqualFold(src.Host, dst.Host) {
			req.Header.Set("Sec-Fetch-Site", "same-origin")
		} else if strings.EqualFold(topdomain(src), topdomain(dst)) {
			req.Header.Set("Sec-Fetch-Site", "same-site")
		} else {
			req.Header.Set("Sec-Fetch-Site", "cross-site")
		}

		if !init.Header.Has("Referer") {
			req.Header.Del("Referer")
			return nil
		}

		if strings.EqualFold(src.Scheme, "https") && strings.EqualFold(dst.Scheme, "http") {
			req.Header.Del("Referer")
		} else if strings.EqualFold(src.Scheme, dst.Scheme) && strings.EqualFold(src.Host, dst.Host) {
			ref := &url.URL{
				Scheme:   src.Scheme,
				Host:     src.Host,
				Path:     src.Path,
				RawQuery: src.RawQuery,
			}
			req.Header.Set("Referer", ref.String())
		} else {
			ref := &url.URL{
				Scheme: src.Scheme,
				Host:   src.Host,
				Path:   "/",
			}
			req.Header.Set("Referer", ref.String())
		}
	}
	return nil
}
