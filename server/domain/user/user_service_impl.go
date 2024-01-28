package user

import (
	"dalkak/pkg/interfaces"

	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type UserServiceImpl struct {
	db  interfaces.UserRepository
	kms *kms.Client
}

func NewUserService(db interfaces.Database, kms *kms.Client) interfaces.UserService {
	userRepo := NewUserRepository(db)

	return &UserServiceImpl{
		db:  userRepo,
		kms: kms,
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
