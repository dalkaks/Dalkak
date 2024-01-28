package user

import (
	"dalkak/pkg/interfaces"
	"log"
)

type UserServiceImpl struct {
	db interfaces.Database
}

func NewUserService(db interfaces.Database) interfaces.UserService {
	return &UserServiceImpl{db: db}
}

func (service *UserServiceImpl) AuthAndSignUp(walletAddress string, signature string) (string, error) {
	log.Printf("Wallet Address: %s, Signature: %s", walletAddress, signature)
	result, err := service.db.FindOrCreateUser(walletAddress)
	if err != nil {
		return "", err
	}
	log.Printf("Result: %s", result)

	return "Authentication and Sign-Up Successful", nil
}
