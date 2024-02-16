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

type UserCreateMediaRequest struct {
	MediaType string `json:"mediaType"`
	Ext       string `json:"ext"`
	Prefix    string `json:"prefix"`
}

type UserCreateMediaResponse struct {
	Id           string `json:"id"`
	Url          string `json:"url"`
	PresignedUrl string `json:"presignedUrl"`
}

type UserGetMediaRequest struct {
	MediaType string `query:"mediaType" required:"true"`
	Prefix    string `query:"prefix" required:"true"`
}

type UserGetMediaResponse struct {
	Id          string `json:"id"`
	ContentType string `json:"contentType"`
	Url         string `json:"url"`
}

type UserConfirmMediaRequest struct {
	UserId    string `json:"userId"`
	Key       string `json:"key"`
	MediaType string `json:"mediaType"`
}

type UserDeleteMediaRequest struct {
	Url       string `query:"url" required:"true"`
	MediaType string `query:"mediaType" required:"true"`
	Prefix    string `query:"prefix" required:"true"`
}

func NewUserCreateMediaResponse(meta *dtos.MediaMeta, presignedUrl string) *UserCreateMediaResponse {
	return &UserCreateMediaResponse{
		Id:           meta.ID,
		Url:          meta.URL,
		PresignedUrl: presignedUrl,
	}
}

func NewUserGetMediaResponse(meta *dtos.MediaMeta) *UserGetMediaResponse {
	return &UserGetMediaResponse{
		Id:          meta.ID,
		ContentType: meta.ContentType,
		Url:         meta.URL,
	}
}

func (req *UserCreateMediaRequest) IsValid() bool {
	return isSupportedMediaType(req.MediaType) && hasValidPrefix(req.Prefix) && isExtensionAllowed(req.Ext)
}

func (req *UserCreateMediaRequest) ToFindUserUploadMediaDto() (*dtos.FindUserUploadMediaDto, error) {
	mediaType, err := dtos.ToMediaType(req.MediaType)
	if err != nil {
		return nil, err
	}

	return &dtos.FindUserUploadMediaDto{
		MediaType: mediaType,
		Prefix:    req.Prefix,
	}, nil
}

func (req *UserCreateMediaRequest) ToUploadMediaDto() (*dtos.UploadMediaDto, error) {
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

func (req *UserGetMediaRequest) IsValid() bool {
	return isSupportedMediaType(req.MediaType) && hasValidPrefix(req.Prefix)
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

func (req *UserDeleteMediaRequest) IsValid() bool {
	return isSupportedMediaType(req.MediaType)
}

func (req *UserDeleteMediaRequest) Verify(meta *dtos.MediaMeta) bool {
	return meta.URL == req.Url
}

func (req *UserDeleteMediaRequest) ToFindUserUploadMediaDto() (*dtos.FindUserUploadMediaDto, error) {
	mediaType, err := dtos.ToMediaType(req.MediaType)
	if err != nil {
		return nil, err
	}

	return &dtos.FindUserUploadMediaDto{
		MediaType: mediaType,
		Prefix:    req.Prefix,
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
