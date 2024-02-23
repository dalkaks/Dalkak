package mediaaggregate

import (
	mediaentity "dalkak/internal/domain/media/object/entity"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
	responseutil "dalkak/pkg/utils/response"
)

type MediaTempAggregate struct {
	MediaEntity  *mediaentity.MediaEntity
	Prefix       mediavalueobject.Prefix
	ContentType  mediavalueobject.ContentType
	MediaTempUrl *mediavalueobject.MediaTempUrl
	Length       *mediavalueobject.Length
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

func WithLength(length mediavalueobject.Length) MediaTempAggregateOption {
	return func(aggregate *MediaTempAggregate) {
		aggregate.Length = &length
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

func (m *MediaTempAggregate) ConfirmMediaTemp(Id string, contentTypeStr string, lengthInt int64) (*MediaTempUpdate, error) {
	contentType, err := mediavalueobject.NewContentType(mediavalueobject.SplitContentType(contentTypeStr))
	if err != nil {
		return nil, err
	}

	_, err = mediavalueobject.NewLength(lengthInt)
	if err != nil {
		return nil, err
	}

	if !m.MediaEntity.CheckId(Id) {
		return nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}

	return &MediaTempUpdate{
		MediaEntity: &mediaentity.MediaEntity{
			Id:        m.MediaEntity.Id,
			IsConfirm: true,
			Timestamp: m.MediaEntity.Timestamp,
		},
		Prefix:       m.Prefix,
		ContentType:  contentType,
		MediaTempUrl: m.MediaTempUrl,
	}, nil
}

type MediaTempUpdate struct {
	MediaEntity  *mediaentity.MediaEntity
	Prefix       mediavalueobject.Prefix
	ContentType  mediavalueobject.ContentType
	MediaTempUrl *mediavalueobject.MediaTempUrl
}
