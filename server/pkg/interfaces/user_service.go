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

	CreatePresignedURL(req *payloads.UserBoardImagePresignedRequest, userInfo *dtos.UserInfo) (*payloads.UserBoardImagePresignedResponse, error)
}
