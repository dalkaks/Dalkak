package mediaaggregate

import (
	mediaentity "dalkak/internal/domain/media/object/entity"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
	responseutil "dalkak/pkg/utils/response"
)

type MediaNftAggregate struct {
	MediaEntity        mediaentity.MediaEntity
	MediaImageResource *mediavalueobject.MediaResource
	MediaImageUrl      *mediavalueobject.MediaUrl
	MediaVideoResource *mediavalueobject.MediaResource
	MediaVideoUrl      *mediavalueobject.MediaUrl
}

type MediaNftAggregateOption func(*MediaNftAggregate)

func WithMediaNftImageResource(mediaResource *mediavalueobject.MediaResource) MediaNftAggregateOption {
	return func(aggregate *MediaNftAggregate) {
		aggregate.MediaImageResource = mediaResource
	}
}

func WithMediaNftImageUrl(mediaUrl *mediavalueobject.MediaUrl) MediaNftAggregateOption {
	return func(aggregate *MediaNftAggregate) {
		aggregate.MediaImageUrl = mediaUrl
	}
}

func WithMediaNftVideoResource(mediaResource *mediavalueobject.MediaResource) MediaNftAggregateOption {
	return func(aggregate *MediaNftAggregate) {
		aggregate.MediaVideoResource = mediaResource
	}
}

func WithMediaNftVideoUrl(mediaUrl *mediavalueobject.MediaUrl) MediaNftAggregateOption {
	return func(aggregate *MediaNftAggregate) {
		aggregate.MediaVideoUrl = mediaUrl
	}
}

func NewMediaNftAggregate(media *mediaentity.MediaEntity, options ...MediaNftAggregateOption) (*MediaNftAggregate, error) {
	if media == nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}

	aggregate := &MediaNftAggregate{
		MediaEntity: *media,
	}

	for _, option := range options {
		option(aggregate)
	}
	return aggregate, nil
}
