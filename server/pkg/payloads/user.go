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

type UserUploadMediaRequest struct {
	MediaType string `json:"mediaType"`
	Ext       string `json:"ext"`
	Prefix    string `json:"prefix"`
}

func (req *UserUploadMediaRequest) IsValid() bool {
	return req.isSupportedMediaType() && req.hasValidPrefix() && req.isExtensionAllowed()
}

func (req *UserUploadMediaRequest) isSupportedMediaType() bool {
	return req.MediaType == "image"
}

func (req *UserUploadMediaRequest) hasValidPrefix() bool {
	return req.Prefix == "board"
}

func (req *UserUploadMediaRequest) isExtensionAllowed() bool {
	_, ok := config.AllowedImageExtensions[req.Ext]
	return ok
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

type UserBoardImagePresignedResponse struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}
