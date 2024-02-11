package user

import (
	appsecurity "dalkak/internal/security"
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/payloads"
	"dalkak/pkg/utils/validateutils"
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
		mode:    mode,
		domain:  domain,
		db:      userRepo,
		kmsSet:  kmsSet,
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

func (service *UserServiceImpl) CreatePresignedURL(userInfo *dtos.UserInfo, dto *payloads.UserUploadMediaRequest) (*payloads.UserUploadMediaResponse, error) {
	err := validateutils.Validate(dto)
	if err != nil {
		return nil, err
	}

	uploadMediaDto, err := dto.ToUploadMediaDto()
	if err != nil {
		return nil, err
	}

	mediaMeta, presignedUrl, err := service.storage.CreatePresignedURL(uploadMediaDto)
	if err != nil {
		return nil, err
	}

	// Todo : 기한이 지남에 따라 upload media 삭제
	err = service.db.CreateUserUploadMedia(userInfo.WalletAddress, mediaMeta)
	if err != nil {
		return nil, err
	}

	return &payloads.UserUploadMediaResponse{
		Id:           mediaMeta.ID,
		Url:          mediaMeta.URL,
		PresignedUrl: presignedUrl,
	}, nil
}

func (service *UserServiceImpl) GetUserMedia(userInfo *dtos.UserInfo, dto *payloads.UserGetMediaRequest) (*payloads.UserGetMediaResponse, error) {
	err := validateutils.Validate(dto)
	if err != nil {
		return nil, err
	}

	media, err := service.db.FindUserUploadMedia(userInfo.WalletAddress, dto)
	if err != nil {
		return nil, err
	}

	return &payloads.UserGetMediaResponse{
		Id:          media.ID,
		ContentType: media.ContentType,
		Url:         media.URL,
	}, nil
}
