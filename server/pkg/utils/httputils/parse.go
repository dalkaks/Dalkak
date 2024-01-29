package httputils

import (
	"net/url"
	"strings"
)

func ParseDomain(u string) (string, error) {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	host := parsedUrl.Hostname()
	host = strings.Split(host, ":")[0]

	parts := strings.Split(host, ".")
	if len(parts) > 2 {
		host = parts[len(parts)-2]
	}

	return host, nil
}
