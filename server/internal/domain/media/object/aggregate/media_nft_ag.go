package mediaaggregate

import (
	mediaentity "dalkak/internal/domain/media/object/entity"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
)

type MediaNftAggregate struct {
	MediaEntity        *mediaentity.MediaEntity
	MediaImageResource *mediavalueobject.MediaResource
	MediaImageUrl      *mediavalueobject.MediaUrl
	MediaVideoResource *mediavalueobject.MediaResource
	MediaVideoUrl      *mediavalueobject.MediaUrl
}

type MediaNftAggregateOption func(*MediaNftAggregate)

func WithMediaNftEntity(mediaEntity *mediaentity.MediaEntity) MediaNftAggregateOption {
	return func(aggregate *MediaNftAggregate) {
		aggregate.MediaEntity = mediaEntity
	}
}

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

func NewMediaNftAggregate(options ...MediaNftAggregateOption) *MediaNftAggregate {
	aggregate := &MediaNftAggregate{}
	for _, option := range options {
		option(aggregate)
	}
	return aggregate
}
