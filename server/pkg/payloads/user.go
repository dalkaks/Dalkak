package payloads

import (
	"dalkak/config"
	"dalkak/pkg/dtos"
	"net/http"
	"strings"
)

type UserAuthAndSignUpRequest struct {
	WalletAddress string
	Signature     string
}

type UserAccessTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type UserGetMediaRequest struct {
	MediaType string `query:"mediaType" required:"true"`
	Prefix    string `query:"prefix" required:"true"`
}

func (req *UserGetMediaRequest) ToFindUserUploadMediaDto() (*dtos.FindUserUploadMediaDto, error) {
	mediaType, err := dtos.ToMediaType(req.MediaType)
	if err != nil {
		return nil, &dtos.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid media type",
		}
	}

	trueValue := true
	return &dtos.FindUserUploadMediaDto{
		MediaType: mediaType,
		Prefix:    req.Prefix,
		IsConfirm: &trueValue,
	}, nil
}

type UserGetMediaResponse struct {
	Id          string `json:"id"`
	ContentType string `json:"contentType"`
	Url         string `json:"url"`
}

func (req *UserGetMediaRequest) IsValid() bool {
	return isSupportedMediaType(req.MediaType) && hasValidPrefix(req.Prefix)
}

type UserUploadMediaRequest struct {
	MediaType string `json:"mediaType"`
	Ext       string `json:"ext"`
	Prefix    string `json:"prefix"`
}

type UserUploadMediaResponse struct {
	Id           string `json:"id"`
	Url          string `json:"url"`
	PresignedUrl string `json:"presignedUrl"`
}

func (req *UserUploadMediaRequest) IsValid() bool {
	return isSupportedMediaType(req.MediaType) && hasValidPrefix(req.Prefix) && isExtensionAllowed(req.Ext)
}

func (req *UserUploadMediaRequest) ToUploadMediaDto() (*dtos.UploadMediaDto, error) {
	mediaType, err := dtos.ToMediaType(req.MediaType)
	if err != nil {
		return nil, err
	}

	return &dtos.UploadMediaDto{
		MediaType: mediaType,
		Ext:       req.Ext,
		Prefix:    req.Prefix,
	}, nil
}

type UserConfirmMediaRequest struct {
	UserId    string `json:"userId"`
	Key       string `json:"key"`
	MediaType string `json:"mediaType"`
}

func (req *UserConfirmMediaRequest) IsValid() bool {
	return isSupportedMediaType(req.MediaType)
}

func (req *UserConfirmMediaRequest) ToFindUserUploadMediaDto() (*dtos.FindUserUploadMediaDto, error) {
	mediaType, err := dtos.ToMediaType(req.MediaType)
	if err != nil {
		return nil, &dtos.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid media type",
		}
	}

	path := strings.Split(req.Key, "/")
	if len(path) < 2 {
		return nil, &dtos.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid key: " + req.Key,
		}
	}
	prefix := path[1]

	return &dtos.FindUserUploadMediaDto{
		MediaType: mediaType,
		Prefix:    prefix,
	}, nil
}

func isSupportedMediaType(mediaType string) bool {
	return mediaType == "image"
}

func hasValidPrefix(prefix string) bool {
	return prefix == "board"
}

func isExtensionAllowed(ext string) bool {
	_, ok := config.AllowedImageExtensions[ext]
	return ok
}
