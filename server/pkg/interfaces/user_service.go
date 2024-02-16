package interfaces

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/payloads"
)

type UserService interface {
	GetMode() string
	GetDomain() string
	AuthAndSignUp(walletAddress string, signature string) (*dtos.AuthTokens, int64, error)
	ReissueRefresh(refreshToken string) (*dtos.AuthTokens, int64, error)

	CreatePresignedURL(userInfo *dtos.UserInfo, req *payloads.UserCreateMediaRequest) (*payloads.UserCreateMediaResponse, error)
	GetUserMedia(userInfo *dtos.UserInfo, req *payloads.UserGetMediaRequest) (*payloads.UserGetMediaResponse, error)
	ConfirmMediaUpload(userInfo *dtos.UserInfo, req *payloads.UserConfirmMediaRequest) (bool, error)
}

type MediaFinder interface {
	ToFindUserUploadMediaDto() (*dtos.FindUserUploadMediaDto, error)
}
