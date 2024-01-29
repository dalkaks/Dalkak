package interfaces

import "dalkak/pkg/dtos"

type UserService interface {
	AuthAndSignUp(walletAddress string, signature string) (*dtos.AuthTokens, error)
}
