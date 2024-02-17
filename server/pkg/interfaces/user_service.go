package interfaces

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/payloads"
)

type UserService interface {
	GetMode() string
	GetDomain() string
	AuthAndSignUp(req *payloads.UserAuthAndSignUpRequest) (*payloads.UserAuthAndSignUpResponse, error)
	ReissueAccessToken(req *payloads.UserReissueAccessTokenRequest) (*payloads.UserReissueAccessTokenResponse, error)

	CreatePresignedURL(userInfo *dtos.UserInfo, req *payloads.UserCreateMediaRequest) (*payloads.UserCreateMediaResponse, error)
	GetUserMedia(userInfo *dtos.UserInfo, req *payloads.UserGetMediaRequest) (*payloads.UserGetMediaResponse, error)
	ConfirmMediaUpload(userInfo *dtos.UserInfo, req *payloads.UserConfirmMediaRequest) error
	DeleteUserMedia(userInfo *dtos.UserInfo, req *payloads.UserDeleteMediaRequest) error
}

type MediaFinder interface {
	ToFindUserUploadMediaDto() (*dtos.FindUserUploadMediaDto, error)
}
