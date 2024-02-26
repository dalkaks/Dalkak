package mediafactory

import (
	mediaaggregate "dalkak/internal/domain/media/object/aggregate"
	mediaentity "dalkak/internal/domain/media/object/entity"
	mediadto "dalkak/pkg/dto/media"
)

type MediaNftAggregateFactory interface {
	CreateMediaNftAggregate() (*mediaaggregate.MediaNftAggregate, error)
}

type CreateMediaNftDtoFactory struct {
	dto *mediadto.CreateMediaNftDto
	mediaImage *mediaaggregate.MediaTempAggregate
	mediaVideo *mediaaggregate.MediaTempAggregate
}

func NewCreateMediaNftDtoFactory(dto *mediadto.CreateMediaNftDto, mediaImage *mediaaggregate.MediaTempAggregate, mediaVideo *mediaaggregate.MediaTempAggregate) MediaNftAggregateFactory {
	return &CreateMediaNftDtoFactory{
		dto: dto,
		mediaImage: mediaImage,
		mediaVideo: mediaVideo,
	}
}

func (factory *CreateMediaNftDtoFactory) CreateMediaNftAggregate() (*mediaaggregate.MediaNftAggregate, error) {
	media := mediaentity.NewMediaEntity(mediaentity.WithID(factory.dto.PrefixId))
	
	mediaNftAggregate := mediaaggregate.NewMediaNftAggregate(
		mediaaggregate.WithMediaNftEntity(media),
		mediaaggregate.WithMediaNftImageResource(factory.mediaImage.MediaResource),
		mediaaggregate.WithMediaNftImageUrl(factory.mediaImage.MediaUrl),
		mediaaggregate.WithMediaNftVideoResource(factory.mediaVideo.MediaResource),
		mediaaggregate.WithMediaNftVideoUrl(factory.mediaVideo.MediaUrl),
	)
	return mediaNftAggregate, nil
}
