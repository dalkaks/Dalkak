package user

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/jwtutils"
	"dalkak/pkg/utils/kmsutils"
)

type UserServiceImpl struct {
	domain string
	db     interfaces.UserRepository
	kmsSet *kmsutils.KmsSet
}

func NewUserService(domain string, db interfaces.Database, kmsSet *kmsutils.KmsSet) interfaces.UserService {
	userRepo := NewUserRepository(db)

	return &UserServiceImpl{
		domain: domain,
		db:     userRepo,
		kmsSet: kmsSet,
	}
}

func (service *UserServiceImpl) AuthAndSignUp(walletAddress string, signature string) (*dtos.AuthTokens, error) {
	user, err := service.db.FindUser(walletAddress)
	if err != nil {
		return nil, err
	}

	if user == nil {
		_, err := service.db.CreateUser(walletAddress)
		if err != nil {
			return nil, err
		}
	}

	accessToken, err := jwtutils.GenerateAccessToken(service.domain, service.kmsSet, walletAddress)
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwtutils.GenerateRefreshToken(service.domain, service.kmsSet, walletAddress)
	if err != nil {
		return nil, err
	}

	return &dtos.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
