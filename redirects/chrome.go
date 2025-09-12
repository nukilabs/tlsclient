package redirects

import (
	"net/url"
	"slices"
	"strings"

	"github.com/nukilabs/http"
)

func Chrome(req *http.Request, via []*http.Request) error {
	req.Header[http.HeaderOrderKey] = []string{"host", "connection", "cache-control", "upgrade-insecure-requests", "user-agent", "accept", "x-browser-channel", "x-browser-year", "x-browser-validation", "x-browser-copyright", "x-chrome-id-consistency-request", "x-client-data", "sec-fetch-site", "sec-fetch-mode", "sec-fetch-user", "sec-fetch-dest", "sec-ch-device-memory", "sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-full-version", "sec-ch-ua-arch", "sec-ch-ua-platform", "sec-ch-ua-platform-version", "sec-ch-ua-model", "sec-ch-ua-bitness", "sec-ch-ua-wow64", "sec-ch-ua-full-version-list", "sec-ch-ua-form-factors", "referer", "accept-encoding", "accept-language", "cookie", "priority"}
	for key := range req.Header {
		if key == http.HeaderOrderKey {
			continue
		}
		if !slices.ContainsFunc(req.Header[http.HeaderOrderKey], func(s string) bool {
			return strings.EqualFold(s, key)
		}) {
			delete(req.Header, key)
		}
	}
	if len(via) > 0 {
		init := via[0]
		prev := via[len(via)-1]
		src := prev.URL
		dst := req.URL

		key, value := header(prev.Header, "Sec-Fetch-Site")
		if value == "none" {
			req.Header[key] = []string{"none"}
		} else if strings.EqualFold(src.Scheme, dst.Scheme) && strings.EqualFold(src.Host, dst.Host) {
			req.Header[key] = []string{"same-origin"}
		} else if strings.EqualFold(topdomain(src), topdomain(dst)) {
			req.Header[key] = []string{"same-site"}
		} else {
			req.Header[key] = []string{"cross-site"}
		}

		key, value = header(init.Header, "Referer")
		if value == "" {
			delete(req.Header, key)
			return nil
		}

		key, _ = header(req.Header, "Referer")
		if strings.EqualFold(src.Scheme, "https") && strings.EqualFold(dst.Scheme, "http") {
			delete(req.Header, key)
		} else if strings.EqualFold(src.Scheme, dst.Scheme) && strings.EqualFold(src.Host, dst.Host) {
			ref := &url.URL{
				Scheme:   src.Scheme,
				Host:     src.Host,
				Path:     src.Path,
				RawQuery: src.RawQuery,
			}
			req.Header[key] = []string{ref.String()}
		} else {
			ref := &url.URL{
				Scheme: src.Scheme,
				Host:   src.Host,
			}
			req.Header[key] = []string{ref.String()}
		}
	}
	return nil
}
