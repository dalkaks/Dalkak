package payloads

import (
	"dalkak/config"
	"dalkak/pkg/dtos"
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
