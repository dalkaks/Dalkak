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
	dto        mediadto.CreateMediaNftDto
	mediaImage *mediaaggregate.MediaTempAggregate
	mediaVideo *mediaaggregate.MediaTempAggregate
}

func NewCreateMediaNftDtoFactory(dto *mediadto.CreateMediaNftDto, mediaImage *mediaaggregate.MediaTempAggregate, mediaVideo *mediaaggregate.MediaTempAggregate) MediaNftAggregateFactory {
	return &CreateMediaNftDtoFactory{
		dto:        *dto,
		mediaImage: mediaImage,
		mediaVideo: mediaVideo,
	}
}

func (factory *CreateMediaNftDtoFactory) CreateMediaNftAggregate() (*mediaaggregate.MediaNftAggregate, error) {
	media := mediaentity.ConvertMediaEntity(factory.dto.PrefixId, true, factory.dto.Timestamp)

	var options []mediaaggregate.MediaNftAggregateOption
	if factory.mediaImage != nil {
		options = append(options, mediaaggregate.WithMediaNftImageResource(&factory.mediaImage.MediaResource))
		if factory.mediaImage.MediaUrl != nil {
			options = append(options, mediaaggregate.WithMediaNftImageUrl(factory.mediaImage.MediaUrl))
		}
	}

	if factory.mediaVideo != nil {
		options = append(options, mediaaggregate.WithMediaNftVideoResource(&factory.mediaVideo.MediaResource))
		if factory.mediaVideo.MediaUrl != nil {
			options = append(options, mediaaggregate.WithMediaNftVideoUrl(factory.mediaVideo.MediaUrl))
		}
	}

	mediaNftAggregate, err := mediaaggregate.NewMediaNftAggregate(media, options...)
	if err != nil {
		return nil, err
	}
	return mediaNftAggregate, nil
}
