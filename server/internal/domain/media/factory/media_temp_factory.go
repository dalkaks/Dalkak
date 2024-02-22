package mediafactory

import (
	mediaaggregate "dalkak/internal/domain/media/object/aggregate"
	mediaentity "dalkak/internal/domain/media/object/entity"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
	"dalkak/internal/infrastructure/database/dao"
	mediadto "dalkak/pkg/dto/media"
)

type MediaTempAggregateFactory interface {
	CreateMediaTempAggregate() (*mediaaggregate.MediaTempAggregate, error)
}

type CreateMediaTempDtoFactory struct {
	dto        *mediadto.CreateMediaTempDto
	staticLink string
}

func NewCreateMediaTempDtoFactory(dto *mediadto.CreateMediaTempDto, staticLink string) *CreateMediaTempDtoFactory {
	return &CreateMediaTempDtoFactory{
		dto:        dto,
		staticLink: staticLink,
	}
}

type MediaTempDaoFactory struct {
	dao        *dao.MediaTempDao
	staticLink string
}

func NewMediaTempDaoFactory(dao *dao.MediaTempDao, staticLink string) *MediaTempDaoFactory {
	return &MediaTempDaoFactory{
		dao:        dao,
		staticLink: staticLink,
	}
}

func (factory *CreateMediaTempDtoFactory) CreateMediaTempAggregate() (*mediaaggregate.MediaTempAggregate, error) {
	prefix, err := mediavalueobject.NewPrefix(factory.dto.Prefix)
	if err != nil {
		return nil, err
	}

	contentType, err := mediavalueobject.NewContentType(factory.dto.MediaType, factory.dto.Ext)
	if err != nil {
		return nil, err
	}

	media := mediaentity.NewMediaEntity()

	mediaKey := mediavalueobject.GenerateMediaTempKey(factory.dto.UserInfo.GetUserId(), prefix, contentType)
	mediaTempUrl := mediavalueobject.NewMediaTempUrl(factory.staticLink, mediaKey)

	mediaTempAggregate := mediaaggregate.NewMediaTempAggregate(
		mediaaggregate.WithMediaEntity(media),
		mediaaggregate.WithPrefix(prefix),
		mediaaggregate.WithContentType(contentType),
		mediaaggregate.WithMediaTempUrl(mediaTempUrl),
	)
	return mediaTempAggregate, nil
}

func (factory *MediaTempDaoFactory) CreateMediaTempAggregate() (*mediaaggregate.MediaTempAggregate, error) {
	media := mediaentity.ConvertMediaEntity(factory.dao.Id, factory.dao.IsConfirm, factory.dao.Timestamp)

	prefix, err := mediavalueobject.NewPrefix(factory.dao.Prefix)
	if err != nil {
		return nil, err
	}

	contentType, err := mediavalueobject.NewContentType(mediavalueobject.SplitContentType(factory.dao.ContentType))
	if err != nil {
		return nil, err
	}

	mediaTempUrl := mediavalueobject.NewMediaTempUrlWithOnlyAccessUrl(factory.dao.Url)

	mediaTempAggregate := mediaaggregate.NewMediaTempAggregate(
		mediaaggregate.WithMediaEntity(media),
		mediaaggregate.WithPrefix(prefix),
		mediaaggregate.WithContentType(contentType),
		mediaaggregate.WithMediaTempUrl(mediaTempUrl),
	)

	return mediaTempAggregate, nil
}
