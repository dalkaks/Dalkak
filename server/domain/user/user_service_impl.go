package user

import (
	appsecurity "dalkak/internal/security"
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/payloads"
	"net/http"
)

type UserServiceImpl struct {
	mode    string
	domain  string
	db      interfaces.UserRepository
	kmsSet  *appsecurity.KmsSet
	storage interfaces.Storage
}

func NewUserService(mode string, domain string, db interfaces.Database, kmsSet *appsecurity.KmsSet, storage interfaces.Storage) interfaces.UserService {
	userRepo := NewUserRepository(db)

	return &UserServiceImpl{
		mode:   mode,
		domain: domain,
		db:     userRepo,
		kmsSet: kmsSet,
		storage: storage,
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
	authTokens, nowTime, err := appsecurity.GenerateAuthTokens(service.domain, service.kmsSet, &generateTokenDto)
	if err != nil {
		return nil, 0, err
	}

	return authTokens, nowTime, nil
}

func (service *UserServiceImpl) ReissueRefresh(refreshToken string) (*dtos.AuthTokens, int64, error) {
	walletAddress, err := appsecurity.ParseTokenWithPublicKey(refreshToken, service.kmsSet)
	if err != nil {
		return nil, 0, err
	}

	generateTokenDto := dtos.GenerateTokenDto{
		WalletAddress: walletAddress,
	}
	authTokens, nowTime, err := appsecurity.GenerateAuthTokens(service.domain, service.kmsSet, &generateTokenDto)
	if err != nil {
		return nil, 0, err
	}

	return authTokens, nowTime, nil
}

func (service *UserServiceImpl) CreatePresignedURL(dto *payloads.UserBoardImagePresignedRequest, userInfo *dtos.UserInfo) (*payloads.UserBoardImagePresignedResponse, error) {
	if dto.IsValid() == false {
		return nil, &dtos.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
		}
	}

	mediaType, err := dtos.ToMediaType(dto.MediaType)
	if err != nil {
		return nil, err
	}

	mediaMeta, err := service.storage.CreatePresignedURL(mediaType, dto.Ext)
	if err != nil {
		return nil, err
	}

	return &payloads.UserBoardImagePresignedResponse{
		Id:  mediaMeta.ID,
		Url: mediaMeta.URL,
	}, nil
}
