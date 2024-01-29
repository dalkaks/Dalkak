package user

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/securityutils"
	"dalkak/pkg/utils/timeutils"
)

type UserServiceImpl struct {
	mode   string
	domain string
	db     interfaces.UserRepository
	kmsSet *securityutils.KmsSet
}

func NewUserService(mode string, domain string, db interfaces.Database, kmsSet *securityutils.KmsSet) interfaces.UserService {
	userRepo := NewUserRepository(db)

	return &UserServiceImpl{
    mode:   mode,
		domain: domain,
		db:     userRepo,
		kmsSet: kmsSet,
	}
}

func (service *UserServiceImpl) AuthAndSignUp(walletAddress string, signature string) (*dtos.AuthTokens, int64, error) {
	user, err := service.db.FindUser(walletAddress)
	if err != nil {
		return nil, 0, err
	}

	if user == nil {
		_, err := service.db.CreateUser(walletAddress)
		if err != nil {
			return nil, 0, err
		}
	}

	nowTime := timeutils.GetTimestamp()
	generateTokenDto := dtos.GenerateTokenDto{
		WalletAddress: walletAddress,
		NowTime:       nowTime,
	}
	accessToken, err := securityutils.GenerateAccessToken(service.domain, service.kmsSet, &generateTokenDto)
	if err != nil {
		return nil, 0, err
	}
	refreshToken, err := securityutils.GenerateRefreshToken(service.domain, service.kmsSet, &generateTokenDto)
	if err != nil {
		return nil, 0, err
	}

	return &dtos.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nowTime, nil
}
