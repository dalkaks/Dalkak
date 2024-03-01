package mediaaggregate

import (
	mediaentity "dalkak/internal/domain/media/object/entity"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
	responseutil "dalkak/pkg/utils/response"
)

type MediaTempAggregate struct {
	MediaEntity   mediaentity.MediaEntity
	MediaResource mediavalueobject.MediaResource
	MediaUrl      *mediavalueobject.MediaUrl
	Length        *mediavalueobject.Length
}

type MediaTempAggregateOption func(*MediaTempAggregate)

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

func NewMediaTempAggregate(media *mediaentity.MediaEntity, resource *mediavalueobject.MediaResource, options ...MediaTempAggregateOption) (*MediaTempAggregate, error) {
	if media == nil || resource == nil {
		return nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}

	aggregate := &MediaTempAggregate{
		MediaEntity:   *media,
		MediaResource: *resource,
	}
	
	for _, option := range options {
		option(aggregate)
	}
	return aggregate, nil
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
		MediaResource: m.MediaResource,
		MediaUrl:      m.MediaUrl,
	}, nil
}

type MediaTempUpdate struct {
	MediaEntity   *mediaentity.MediaEntity
	MediaResource mediavalueobject.MediaResource
	MediaUrl      *mediavalueobject.MediaUrl
}
