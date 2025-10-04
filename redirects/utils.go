package redirects

import (
	"net/url"

	"golang.org/x/net/publicsuffix"
)

func topdomain(u *url.URL) string {
	domain, _ := publicsuffix.EffectiveTLDPlusOne(u.Host)
	return domain
}
