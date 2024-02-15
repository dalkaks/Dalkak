package interfaces

import (
	"dalkak/pkg/dtos"
)

type UserRepository interface {
	CreateUser(walletAddress string) error
	FindUser(walletAddress string) (*dtos.UserDto, error)

	CreateUserUploadMedia(userId string, dto *dtos.MediaMeta) error
	FindUserUploadMedia(userId string, dto *dtos.FindUserUploadMediaDto) (*dtos.MediaMeta, error)
	UpdateUserUploadMedia(userId string, findDto *dtos.MediaMeta, updateDto *dtos.UpdateUserUploadMediaDto) error
	DeleteUserUploadMedia(userId string, dto *dtos.MediaMeta) error
}
