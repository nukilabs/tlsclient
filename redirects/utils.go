package redirects

import (
	"net/url"
	"strings"

	"github.com/nukilabs/http"
	"golang.org/x/net/publicsuffix"
)

func topdomain(u *url.URL) string {
	domain, _ := publicsuffix.EffectiveTLDPlusOne(u.Host)
	return domain
}

func header(h http.Header, key string) (string, string) {
	if values, ok := h[key]; ok && len(values) > 0 {
		return key, values[0]
	}
	lowerKey := strings.ToLower(key)
	if values, ok := h[lowerKey]; ok && len(values) > 0 {
		return lowerKey, values[0]
	}
	return key, ""
}
