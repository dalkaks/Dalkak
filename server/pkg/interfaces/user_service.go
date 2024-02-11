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

	GetUserMedia(userInfo *dtos.UserInfo, req *payloads.UserGetMediaRequest) (string, error)
	CreatePresignedURL(userInfo *dtos.UserInfo, req *payloads.UserUploadMediaRequest) (*payloads.UserUploadMediaResponse, error)
}
