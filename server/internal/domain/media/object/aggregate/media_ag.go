package mediaaggregate

import (
	mediaentity "dalkak/internal/domain/media/object/entity"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
)

type MediaTempAggregate struct {
	MediaEntity  *mediaentity.MediaEntity
	Prefix       mediavalueobject.Prefix
	ContentType  mediavalueobject.ContentType
	MediaTempUrl *mediavalueobject.MediaTempUrl
}

type MediaTempAggregateOption func(*MediaTempAggregate)

func WithMediaEntity(mediaEntity *mediaentity.MediaEntity) MediaTempAggregateOption {
	return func(aggregate *MediaTempAggregate) {
		aggregate.MediaEntity = mediaEntity
	}
}

func WithPrefix(prefix mediavalueobject.Prefix) MediaTempAggregateOption {
	return func(aggregate *MediaTempAggregate) {
		aggregate.Prefix = prefix
	}
}

func WithContentType(contentType mediavalueobject.ContentType) MediaTempAggregateOption {
	return func(aggregate *MediaTempAggregate) {
		aggregate.ContentType = contentType
	}
}

func WithMediaTempUrl(mediaTempUrl *mediavalueobject.MediaTempUrl) MediaTempAggregateOption {
	return func(aggregate *MediaTempAggregate) {
		aggregate.MediaTempUrl = mediaTempUrl
	}
}

func NewMediaTempAggregate(options ...MediaTempAggregateOption) *MediaTempAggregate {
	aggregate := &MediaTempAggregate{}
	for _, option := range options {
		option(aggregate)
	}
	return aggregate
}

func (m *MediaTempAggregate) CheckPublic() *MediaTempAggregate {
	if m.MediaEntity.IsPublic() {
		return m
	}
	return nil
}

func (m *MediaTempAggregate) SetUploadUrl(uploadUrl string) {
	if m.MediaTempUrl != nil {
		m.MediaTempUrl.UploadUrl = &uploadUrl
	} else {
		m.MediaTempUrl = &mediavalueobject.MediaTempUrl{UploadUrl: &uploadUrl}
	}
}

// package appdto

// type MediaHeadDto struct {
// 	Key         string
// 	ContentType string
// 	Length      int64
// 	URL         string
// 	MetaUserId  string
// }

// func (m *MediaHeadDto) Verify(userId string, meta *MediaMeta) bool {
// 	return int64(config.MaxUploadSize) > m.Length && m.ContentType == meta.ContentType && m.URL == meta.URL && userId == m.MetaUserId
// }

// func WithConfirm(isConfirm bool) FindUserUploadMediaOption {
// 	return func(f *FindUserUploadMediaDto) {
// 		f.IsConfirm = &isConfirm
// 	}
// }
