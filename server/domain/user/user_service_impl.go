package user

import (
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/kmsutils"
)

type UserServiceImpl struct {
	db     interfaces.UserRepository
	kmsSet *kmsutils.KmsSet
}

func NewUserService(db interfaces.Database, kmsSet *kmsutils.KmsSet) interfaces.UserService {
	userRepo := NewUserRepository(db)

	return &UserServiceImpl{
		db:     userRepo,
		kmsSet: kmsSet,
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
