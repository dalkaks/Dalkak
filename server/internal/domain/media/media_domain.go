package mediadomain

import (
	"dalkak/config"
	"dalkak/internal/core"
	mediaaggregate "dalkak/internal/domain/media/object/aggregate"
	mediaentity "dalkak/internal/domain/media/object/entity"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
	mediadto "dalkak/pkg/dto/media"
	responseutil "dalkak/pkg/utils/response"
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
	prefix, err := mediavalueobject.NewPrefix(dto.Prefix)
	if err != nil {
		return nil, err
	}

	contentType, err := mediavalueobject.NewContentType(dto.MediaType, dto.Ext)
	if err != nil {
		return nil, err
	}

	media := mediaentity.NewMediaEntity()

	mediaKey := mediavalueobject.GenerateMediaTempKey(dto.UserInfo.GetUserId(), prefix, contentType)
	uploadUrl, err := service.Storage.CreatePresignedURL(mediaKey, contentType.String())
	if err != nil {
		return nil, err
	}

	mediaTempUrl := mediavalueobject.NewMediaTempUrl(service.StaticLink, mediaKey, uploadUrl)

	mediaTempAggregate := mediaaggregate.NewMediaTempAggregate(media, prefix, contentType, mediaTempUrl)
	return mediaTempAggregate, nil
}

func (service *MediaDomainServiceImpl) GetMediaTemp(dto *mediadto.GetMediaTempDto) (*mediaaggregate.MediaTempAggregate, error) {
	prefix, err := mediavalueobject.NewPrefix(dto.Prefix)
	if err != nil {
		return nil, err
	}

	if !mediavalueobject.IsAllowedMediaType(dto.MediaType) {
		return nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}

	mediaTempDao, err := service.Database.FindMediaTemp(dto.UserInfo.GetUserId(), dto.MediaType, prefix.String())
	if mediaTempDao == nil || err != nil {
		return nil, err
	}

	media := mediaentity.ConvertMediaEntity(mediaTempDao.Id, mediaTempDao.IsConfirm, mediaTempDao.Timestamp)
	contentType, err := mediavalueobject.NewContentType(mediavalueobject.SplitContentType(mediaTempDao.ContentType))
	if err != nil {
		return nil, err
	}
	mediaTempUrl := mediavalueobject.NewMediaTempUrl(service.StaticLink, mediaTempDao.Prefix, mediaTempDao.Url)

	mediaTemp := mediaaggregate.NewMediaTempAggregate(media, prefix, contentType, mediaTempUrl)
	return mediaTemp.CheckPublic(), nil
}
