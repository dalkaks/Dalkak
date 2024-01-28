package user

import (
	"dalkak/pkg/interfaces"
)

type UserServiceImpl struct {
	db interfaces.UserRepository
}

func NewUserService(db interfaces.Database) interfaces.UserService {
	userRepo := NewUserRepository(db)

	return &UserServiceImpl{
		db: userRepo,
	}
}

func (service *UserServiceImpl) AuthAndSignUp(walletAddress string, signature string) (string, error) {
	user, err := service.db.FindUser(walletAddress)
	if err != nil {
		return "", err
	}

	if user == nil {
		_, err := service.db.CreateUser(walletAddress)
		if err != nil {
			return "", err
		}
	}

	return walletAddress, nil
}
