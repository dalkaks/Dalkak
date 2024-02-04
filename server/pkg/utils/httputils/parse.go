package httputils

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/utils/reflectutils"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// https://dev-api.dalkak.com -> dalkak.com
func ParseDomain(u string) (string, error) {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return "", err
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

func GetRequestData(r *http.Request, target interface{}) error {
	reqMap, ok := r.Context().Value("request").(map[string]interface{})
	if !ok {
		return errors.New("invalid request")
	}

	err := reflectutils.MapToStruct(reqMap, target)
	if err != nil {
		return err
	}

	return nil
}

func GetUserInfoData(r *http.Request) (*dtos.UserInfo, error) {
	userInfo, ok := r.Context().Value("user").(*dtos.UserInfo)
	if !ok {
		return nil, errors.New("invalid user info")
	}

	return userInfo, nil
}
