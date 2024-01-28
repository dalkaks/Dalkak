package interfaces

import (
	"dalkak/pkg/dtos"
)

type UserRepository interface {
	FindUser(walletAddress string) (*dtos.UserDto, error)
	CreateUser(walletAddress string) (string, error)
}
