package interfaces

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/payloads"
)

type UserRepository interface {
	CreateUser(walletAddress string) error
	FindUser(walletAddress string) (*dtos.UserDto, error)

	CreateUserUploadMedia(userId string, dto *dtos.MediaMeta) error
	FindUserUploadMedia(userId string, dto *payloads.UserGetMediaRequest) (*dtos.MediaMeta, error)
}
