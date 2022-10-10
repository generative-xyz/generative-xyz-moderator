package nftdata

import (
	"net/url"
	"strings"
)

// getNameFromAddress get name for a given address
// full name is 6 first characters of wallet address
func getNameFromAddress(address string) string {
	return address[:6]
}

// GetDomainFromURI get domain from URI
func GetDomainFromURI(uri string) string {
	if uri == "" {
		return CompanyDomainDefault
	}

	if !strings.HasPrefix(uri, "https://") {
		uri = "https://" + uri
	}

	u, err := url.Parse(uri)
	if err != nil {
		return CompanyDomainDefault
	}

	parts := strings.Split(u.Hostname(), ".")
	if len(parts) < 2 {
		return CompanyDomainDefault
	}
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
	return domain
}
