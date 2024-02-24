package mediadomain

import (
	"dalkak/config"
	"dalkak/internal/core"
	mediafactory "dalkak/internal/domain/media/factory"
	mediaaggregate "dalkak/internal/domain/media/object/aggregate"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
	mediadto "dalkak/pkg/dto/media"
	responseutil "dalkak/pkg/utils/response"
)

type MediaDomainService interface {
	CreateMediaTemp(dto *mediadto.CreateMediaTempDto) (*mediaaggregate.MediaTempAggregate, error)
	GetMediaTemp(dto *mediadto.GetMediaTempDto, options ...GetMediaTempOptions) (*mediaaggregate.MediaTempAggregate, error)
	ConfirmMediaTemp(dto *mediadto.ConfirmMediaTempDto) (*mediaaggregate.MediaTempUpdate, error)
}

type MediaDomainServiceImpl struct {
	StaticLink   string
	Database     MediaRepository
	Storage      core.StorageManager
	EventManager core.EventManager
}

func NewMediaDomainService(appConfig *config.AppConfig, database MediaRepository, storage core.StorageManager, eventManager core.EventManager) MediaDomainService {
	return &MediaDomainServiceImpl{
		StaticLink:   appConfig.StaticLink,
		Database:     database,
		Storage:      storage,
		EventManager: eventManager,
	}
}

func (service *MediaDomainServiceImpl) CreateMediaTemp(dto *mediadto.CreateMediaTempDto) (*mediaaggregate.MediaTempAggregate, error) {
	factory := mediafactory.NewCreateMediaTempDtoFactory(dto, service.StaticLink)
	mediaTemp, err := factory.CreateMediaTempAggregate()
	if err != nil {
		return nil, err
	}

	uploadUrl, err := service.Storage.CreatePresignedURL(mediaTemp.MediaTempUrl.GetUrlKey(service.StaticLink), mediaTemp.MediaResource.ContentType.String())
	if err != nil {
		return nil, err
	}

	mediaTemp.SetUploadUrl(uploadUrl)
	return mediaTemp, nil
}

type GetMediaTempOptions struct {
	CheckPublic bool
}

func (service *MediaDomainServiceImpl) GetMediaTemp(dto *mediadto.GetMediaTempDto, options ...GetMediaTempOptions) (*mediaaggregate.MediaTempAggregate, error) {
	prefix, err := mediavalueobject.NewPrefix(dto.Prefix)
	if err != nil {
		return nil, err
	}

	mediaTempDao, err := service.Database.FindMediaTemp(dto.UserInfo.GetUserId(), dto.MediaType, prefix.String())
	if mediaTempDao == nil || err != nil {
		return nil, err
	}

	factory := mediafactory.NewMediaTempDaoFactory(mediaTempDao, service.StaticLink)
	mediaTemp, err := factory.CreateMediaTempAggregate()
	if err != nil {
		return nil, err
	}

	if len(options) > 0 && options[0].CheckPublic {
		return mediaTemp.CheckPublic(), nil
	}

	return mediaTemp, nil
}

func (service *MediaDomainServiceImpl) ConfirmMediaTemp(dto *mediadto.ConfirmMediaTempDto) (*mediaaggregate.MediaTempUpdate, error) {
	mediaTemp, err := service.GetMediaTemp(mediadto.NewGetMediaTempDto(dto.UserInfo, dto.MediaType, dto.Prefix))
	if err != nil {
		return nil, err
	}
	if mediaTemp == nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeNotFound, responseutil.ErrMsgDataNotFound)
	}

	if mediaTemp.CheckConfirm() {
		return nil, responseutil.NewAppError(responseutil.ErrCodeConflict, responseutil.ErrMsgDataConflict)
	}

	mediaHead, err := service.Storage.GetHeadObject(mediaTemp.MediaTempUrl.GetUrlKey(service.StaticLink))
	if err != nil {
		return nil, err
	}
	if mediaHead == nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeNotFound, responseutil.ErrMsgDataNotFound)
	}

	mediaTempUpdate, err := mediaTemp.ConfirmMediaTemp(dto.Id, mediaHead.ContentType, mediaHead.Length)
	if err != nil {
		return nil, err
	}
	return mediaTempUpdate, nil
}
