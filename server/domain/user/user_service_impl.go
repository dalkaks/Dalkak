package user

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/securityutils"
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

func (service *UserServiceImpl) GetMode() string {
	return service.mode
}

func (service *UserServiceImpl) GetDomain() string {
	return service.domain
}

func (service *UserServiceImpl) AuthAndSignUp(walletAddress string, signature string) (*dtos.AuthTokens, int64, error) {
	user, err := service.db.FindUser(walletAddress)
	if err != nil {
		return nil, 0, err
	}

	if user == nil {
		err := service.db.CreateUser(walletAddress)
		if err != nil {
			return nil, 0, err
		}
	}

	generateTokenDto := dtos.GenerateTokenDto{
		WalletAddress: walletAddress,
	}
	authTokens, nowTime, err := securityutils.GenerateAuthTokens(service.domain, service.kmsSet, &generateTokenDto)
	if err != nil {
		return nil, 0, err
	}

	return authTokens, nowTime, nil
}

func (service *UserServiceImpl) ReissueRefresh(refreshToken string) (*dtos.AuthTokens, int64, error) {
	walletAddress, err := securityutils.ParseTokenWithPublicKey(refreshToken, service.kmsSet)
	if err != nil {
		return nil, 0, err
	}

	generateTokenDto := dtos.GenerateTokenDto{
		WalletAddress: walletAddress,
	}
	authTokens, nowTime, err := securityutils.GenerateAuthTokens(service.domain, service.kmsSet, &generateTokenDto)
	if err != nil {
		return nil, 0, err
	}

	return authTokens, nowTime, nil
}
