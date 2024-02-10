package httputils

import (
	"dalkak/pkg/dtos"
	"net/http"
	"net/url"
	"strings"
)

// https://dev-api.dalkak.com -> dalkak.com
func ParseDomain(u string) (string, error) {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return "", &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "failed to parse url",
		}
	}

	host := parsedUrl.Hostname()

	if host == "localhost" {
		return "localhost", nil
	}

	host = strings.Split(host, ":")[0]

	parts := strings.Split(host, ".")
	if len(parts) >= 2 {
		host = parts[len(parts)-2] + "." + parts[len(parts)-1]
	}

	return host, nil
}

func GetUserInfoData(r *http.Request) (*dtos.UserInfo, error) {
	userInfo, ok := r.Context().Value("user").(dtos.UserInfo)
	if !ok {
		return nil, &dtos.AppError{
			Code:    http.StatusUnauthorized,
			Message: "user info not found",
		}
	}

	return &userInfo, nil
}
