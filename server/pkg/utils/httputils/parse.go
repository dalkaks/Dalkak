package httputils

import (
	"dalkak/config"
	"dalkak/pkg/dtos"
	"dalkak/pkg/utils/reflectutils"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
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

func GetUploadImageRequest(r *http.Request) (*dtos.ImageData, error) {
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		return nil, err
	}
	defer file.Close()

  extension := filepath.Ext(strings.ToLower(fileHeader.Filename))
  if len(extension) > 1 {
    extension = extension[1:]
  }
  if !config.AllowedImageExtensions[extension] {
    return nil, errors.New("invalid image extension")
  }

	imageData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &dtos.ImageData{
		Meta: dtos.Image{
			URL:       nil,
			Extension: extension,
		},
		Data: imageData,
	}, nil
}
