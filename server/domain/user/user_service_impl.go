package user

import (
	appsecurity "dalkak/internal/security"
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/payloads"
	"dalkak/pkg/utils/validateutils"
	"errors"
)

type UserServiceImpl struct {
	mode    string
	domain  string
	db      interfaces.UserRepository
	KMS     interfaces.KMS
	storage interfaces.Storage
}

func NewUserService(mode string, domain string, db interfaces.Database, KMS interfaces.KMS, storage interfaces.Storage) interfaces.UserService {
	userRepo := NewUserRepository(db)

	return &UserServiceImpl{
		mode:    mode,
		domain:  domain,
		db:      userRepo,
		KMS:     KMS,
		storage: storage,
	}
}

func (service *UserServiceImpl) GetMode() string {
	return service.mode
}

func (service *UserServiceImpl) GetDomain() string {
	return service.domain
}

func (service *UserServiceImpl) AuthAndSignUp(dto *payloads.UserAuthAndSignUpRequest) (*payloads.UserAuthAndSignUpResponse, error) {
	user, err := service.db.FindUser(dto.WalletAddress)
	if err != nil {
		return nil, err
	}

	if user == nil {
		err := service.db.CreateUser(dto.WalletAddress)
		if err != nil {
			return nil, err
		}
	}

	generateTokenDto := dtos.NewGenerateTokenDto(dto.WalletAddress)
	accessToken, err := appsecurity.GenerateAccessToken(service.domain, service.KMS, generateTokenDto)
	if err != nil {
		return nil, err
	}
	refreshToken, err := appsecurity.GenerateRefreshToken(service.domain, service.KMS, generateTokenDto)
	if err != nil {
		return nil, err
	}

	return payloads.NewUserAuthAndSignUpResponse(accessToken, refreshToken), nil
}

func (service *UserServiceImpl) ReissueAccessToken(dto *payloads.UserReissueAccessTokenRequest) (*payloads.UserReissueAccessTokenResponse, error) {
	walletAddress, err := appsecurity.ParseTokenWithPublicKey(dto.RefreshToken, service.KMS)
	if err != nil {
		return nil, err
	}

	generateTokenDto := dtos.NewGenerateTokenDto(walletAddress)
	accessToken, err := appsecurity.GenerateAccessToken(service.domain, service.KMS, generateTokenDto)
	if err != nil {
		return nil, err
	}

	return payloads.NewUserReissueAccessTokenResponse(accessToken), nil
}

func (service *UserServiceImpl) CreatePresignedURL(userInfo *dtos.UserInfo, dto *payloads.UserCreateMediaRequest) (*payloads.UserCreateMediaResponse, error) {
	err := validateutils.Validate(dto)
	if err != nil {
		return nil, err
	}

	media, err := service.readMedia(userInfo.WalletAddress, dto)
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

	media, err := service.readMedia(userInfo.WalletAddress, dto)
	if err != nil {
		return nil, err
	}

	return payloads.NewUserGetMediaResponse(media), nil
}

func (service *UserServiceImpl) ConfirmMediaUpload(userInfo *dtos.UserInfo, dto *payloads.UserConfirmMediaRequest) error {
	err := validateutils.Validate(dto)
	if err != nil {
		return err
	}

	media, err := service.readMedia(userInfo.WalletAddress, dto)
	if err != nil {
		return err
	}
	if media == nil {
		return dtos.NewAppError(dtos.ErrCodeBadRequest, dtos.ErrMsgMediaNotFound, errors.New("media not found"))
	}

	err = service.confirmMedia(userInfo, dto, media)
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

	media, err := service.readMedia(userInfo.WalletAddress, dto)
	if err != nil {
		return err
	}
	if media == nil {
		return dtos.NewAppError(dtos.ErrCodeNotFound, dtos.ErrMsgMediaNotFound, errors.New("media not found"))
	}

	if ok := dto.Verify(media); !ok {
		return dtos.NewAppError(dtos.ErrCodeBadRequest, dtos.ErrMsgRequestInvalid, errors.New("invalid request"))
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

func (service *UserServiceImpl) readMedia(userId string, finder interfaces.MediaFinder) (*dtos.MediaMeta, error) {
	findDto, err := finder.ToFindUserUploadMediaDto()
	if err != nil {
		return nil, err
	}

	media, err := service.db.FindUserUploadMedia(userId, findDto)
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (service *UserServiceImpl) confirmMedia(userInfo *dtos.UserInfo, dto *payloads.UserConfirmMediaRequest, media *dtos.MediaMeta) error {
	mediaHeadDto, err := service.storage.GetHeadObject(dto.Key)
	if err != nil {
		return err
	}

	if ok := mediaHeadDto.Verify(userInfo.WalletAddress, media); !ok {
		err := service.storage.DeleteObject(dto.Key)
		if err != nil {
			return err
		}
		return nil
	}

	err = service.db.UpdateUserUploadMedia(userInfo.WalletAddress, media, &dtos.UpdateUserUploadMediaDto{
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
