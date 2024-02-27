package mediaaggregate

import (
	mediaentity "dalkak/internal/domain/media/object/entity"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
	responseutil "dalkak/pkg/utils/response"
)

type MediaTempAggregate struct {
	MediaEntity   *mediaentity.MediaEntity
	MediaResource *mediavalueobject.MediaResource
	MediaUrl      *mediavalueobject.MediaUrl
	Length        *mediavalueobject.Length
}

type MediaTempAggregateOption func(*MediaTempAggregate)

func WithMediaEntity(mediaEntity *mediaentity.MediaEntity) MediaTempAggregateOption {
	return func(aggregate *MediaTempAggregate) {
		aggregate.MediaEntity = mediaEntity
	}
}

func WithMediaResource(mediaResource *mediavalueobject.MediaResource) MediaTempAggregateOption {
	return func(aggregate *MediaTempAggregate) {
		aggregate.MediaResource = mediaResource
	}
}

func WithMediaUrl(mediaUrl *mediavalueobject.MediaUrl) MediaTempAggregateOption {
	return func(aggregate *MediaTempAggregate) {
		aggregate.MediaUrl = mediaUrl
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

func (m *MediaTempAggregate) CheckConfirm() bool {
	return m.MediaEntity.CheckConfirm()
}

func (m *MediaTempAggregate) SetUploadUrl(uploadUrl string) {
	if m.MediaUrl != nil {
		m.MediaUrl.UploadUrl = &uploadUrl
	} else {
		m.MediaUrl = &mediavalueobject.MediaUrl{UploadUrl: &uploadUrl}
	}
}

func (m *MediaTempAggregate) ConfirmMediaTemp(Id string, contentTypeStr string, lengthInt int64) (*MediaTempUpdate, error) {
	_, err := mediavalueobject.NewContentType(contentTypeStr)
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

	m.MediaEntity.SetConfirm()

	return &MediaTempUpdate{
		MediaEntity: &mediaentity.MediaEntity{
			Id:        m.MediaEntity.Id,
			IsConfirm: m.MediaEntity.IsConfirm,
			Timestamp: m.MediaEntity.Timestamp,
		},
		MediaResource: *m.MediaResource,
		MediaUrl:      m.MediaUrl,
	}, nil
}

type MediaTempUpdate struct {
	MediaEntity   *mediaentity.MediaEntity
	MediaResource mediavalueobject.MediaResource
	MediaUrl      *mediavalueobject.MediaUrl
}
