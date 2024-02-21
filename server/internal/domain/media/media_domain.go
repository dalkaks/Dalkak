package mediadomain

import (
	"dalkak/config"
	"dalkak/internal/core"
	mediaobject "dalkak/internal/domain/media/object"
	mediadto "dalkak/pkg/dto/media"
)

type MediaDomainService interface {
	CreateMediaTemp(dto *mediadto.CreateTempMediaDto) (*mediaobject.MediaTempAggregate, error)
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

func (service *MediaDomainServiceImpl) CreateMediaTemp(dto *mediadto.CreateTempMediaDto) (*mediaobject.MediaTempAggregate, error) {
	prefix, err := mediaobject.NewPrefix(dto.Prefix)
	if err != nil {
		return nil, err
	}

	contentType, err := mediaobject.NewContentType(dto.MediaType, dto.Ext)
	if err != nil {
		return nil, err
	}

	media := mediaobject.NewMediaEntity()

	mediaKey := mediaobject.GenerateMediaTempKey(dto.UserInfo.GetUserId(), prefix, contentType)
	uploadUrl, err := service.Storage.CreatePresignedURL(mediaKey, contentType.String())
	if err != nil {
		return nil, err
	}

	mediaTempUrl := mediaobject.NewMediaTempUrl(service.StaticLink, mediaKey, uploadUrl)

	mediaTempAggregate := mediaobject.NewMediaTempAggregate(media, prefix, contentType, mediaTempUrl)
	return mediaTempAggregate, nil
}
