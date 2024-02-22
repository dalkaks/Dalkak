package mediadomain

import (
	"dalkak/config"
	"dalkak/internal/core"
	mediafactory "dalkak/internal/domain/media/factory"
	mediaaggregate "dalkak/internal/domain/media/object/aggregate"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
	mediadto "dalkak/pkg/dto/media"
)

type MediaDomainService interface {
	CreateMediaTemp(dto *mediadto.CreateMediaTempDto) (*mediaaggregate.MediaTempAggregate, error)
	GetMediaTemp(dto *mediadto.GetMediaTempDto) (*mediaaggregate.MediaTempAggregate, error)
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

	uploadUrl, err := service.Storage.CreatePresignedURL(mediaTemp.MediaTempUrl.GetUrlKey(service.StaticLink), mediaTemp.ContentType.String())
	if err != nil {
		return nil, err
	}

	mediaTemp.SetUploadUrl(uploadUrl)
	return mediaTemp, nil
}

func (service *MediaDomainServiceImpl) GetMediaTemp(dto *mediadto.GetMediaTempDto) (*mediaaggregate.MediaTempAggregate, error) {
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

	return mediaTemp.CheckPublic(), nil
}
