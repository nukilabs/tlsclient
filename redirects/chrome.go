package redirects

import (
	"maps"
	"slices"
	"strings"

	"github.com/nukilabs/http"
)

func Chrome(req *http.Request, via []*http.Request) error {
	req.Header[http.HeaderOrderKey] = []string{"cache-control", "upgrade-insecure-requests", "user-agent", "accept", "sec-fetch-site", "sec-fetch-mode", "sec-fetch-user", "sec-fetch-dest", "sec-ch-device-memory", "sec-ch-ua", "sec-ch-ua-mobile", "sec-ch-ua-arch", "sec-ch-ua-platform", "sec-ch-ua-model", "sec-ch-ua-full-version-list", "referer", "accept-encoding", "accept-language", "cookie", "priority"}
	for key := range maps.Keys(req.Header) {
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
		last := via[len(via)-1]
		location := last.URL
		target := req.URL

		key, value := header(last.Header, "Sec-Fetch-Site")
		if value == "none" {
			req.Header[key] = []string{"none"}
		} else if strings.EqualFold(location.Scheme, target.Scheme) && strings.EqualFold(location.Host, target.Host) {
			req.Header[key] = []string{"same-origin"}
		} else if strings.EqualFold(topdomain(location), topdomain(target)) {
			req.Header[key] = []string{"same-site"}
		} else {
			req.Header[key] = []string{"cross-site"}
		}
	}
	return nil
}
