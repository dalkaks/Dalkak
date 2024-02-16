package user

import (
	appsecurity "dalkak/internal/security"
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/payloads"
	"dalkak/pkg/utils/validateutils"
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

func (service *UserServiceImpl) CreatePresignedURL(userInfo *dtos.UserInfo, dto *payloads.UserCreateMediaRequest) (*payloads.UserCreateMediaResponse, error) {
	err := validateutils.Validate(dto)
	if err != nil {
		return nil, err
	}

	findDto, err := dto.ToFindUserUploadMediaDto()
	if err != nil {
		return nil, err
	}

	media, err := service.readMedia(userInfo.WalletAddress, findDto)
	if err != nil {
		return nil, err
	}

	if media != nil {
		err = service.deleteExistingMedia(userInfo, media)
		if err != nil {
			return nil, err
		}
	}

	media, presignedUrl, err := service.createMedia(userInfo, dto)
	if err != nil {
		return nil, err
	}

	return payloads.NewUserCreateMediaResponse(media, presignedUrl), nil
}

func (service *UserServiceImpl) GetUserMedia(userInfo *dtos.UserInfo, dto *payloads.UserGetMediaRequest) (*payloads.UserGetMediaResponse, error) {
	err := validateutils.Validate(dto)
	if err != nil {
		return nil, err
	}

	findDto, err := dto.ToFindUserUploadMediaDto()
	if err != nil {
		return nil, err
	}

	media, err := service.readMedia(userInfo.WalletAddress, findDto)
	if err != nil {
		return nil, err
	}

	return payloads.NewUserGetMediaResponse(media), nil
}

func (service *UserServiceImpl) ConfirmMediaUpload(dto *payloads.UserConfirmMediaRequest) error {
	err := validateutils.Validate(dto)
	if err != nil {
		return err
	}

	findDto, err := dto.ToFindUserUploadMediaDto()
	if err != nil {
		return err
	}

	media, err := service.readMedia(dto.UserId, findDto)
	if err != nil {
		return err
	}
	if media == nil {
		return &dtos.AppError{
			Code:    http.StatusNotFound,
			Message: "media not found",
		}
	}

	err = service.confirmMedia(dto, media)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserServiceImpl) DeleteUserMedia(userInfo *dtos.UserInfo, dto *payloads.UserDeleteMediaRequest) error {
	err := validateutils.Validate(dto)
	if err != nil {
		return err
	}

	findDto, err := dto.ToFindUserUploadMediaDto()
	if err != nil {
		return err
	}

	media, err := service.readMedia(userInfo.WalletAddress, findDto)
	if err != nil {
		return err
	}
	if media == nil {
		return &dtos.AppError{
			Code:    http.StatusNotFound,
			Message: "media not found",
		}
	}

	if ok := dto.Verify(media); !ok {
		return &dtos.AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		}
	}

	err = service.deleteExistingMedia(userInfo, media)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserServiceImpl) createMedia(userInfo *dtos.UserInfo, dto *payloads.UserCreateMediaRequest) (*dtos.MediaMeta, string, error) {
	uploadMediaDto, err := dto.ToUploadMediaDto()
	if err != nil {
		return nil, "", err
	}

	mediaMeta, presignedUrl, err := service.storage.CreatePresignedURL(userInfo.WalletAddress, uploadMediaDto)
	if err != nil {
		return nil, "", err
	}

	err = service.db.CreateUserUploadMedia(userInfo.WalletAddress, mediaMeta)
	if err != nil {
		return nil, "", err
	}

	return mediaMeta, presignedUrl, nil
}

func (service *UserServiceImpl) readMedia(userId string, dto *dtos.FindUserUploadMediaDto) (*dtos.MediaMeta, error) {
	media, err := service.db.FindUserUploadMedia(userId, dto)
	if err != nil {
		return nil, err
	}
	
	return media, nil
}

func (service *UserServiceImpl) confirmMedia(dto *payloads.UserConfirmMediaRequest, media *dtos.MediaMeta) error {
	mediaHeadDto, err := service.storage.GetHeadObject(dto.Key)
	if err != nil {
		return err
	}

	if ok := mediaHeadDto.Verify(media); !ok {
		err := service.storage.DeleteObject(dto.Key)
		if err != nil {
			return err
		}
		return nil
	}

	err = service.db.UpdateUserUploadMedia(dto.UserId, media, &dtos.UpdateUserUploadMediaDto{
		IsConfirm: true,
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *UserServiceImpl) deleteExistingMedia(userInfo *dtos.UserInfo, media *dtos.MediaMeta) error {
	key, err := service.storage.ConvertStaticLinkToKey(media.URL)
	if err != nil {
		return err
	}
	err = service.storage.DeleteObject(key)
	if err != nil {
		return err
	}
	err = service.db.DeleteUserUploadMedia(userInfo.WalletAddress, media)
	if err != nil {
		return err
	}
	return nil
}
