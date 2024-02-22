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

func NewMediaTempAggregate(mediaEntity *mediaentity.MediaEntity, prefix mediavalueobject.Prefix, contentType mediavalueobject.ContentType, mediaTempUrl *mediavalueobject.MediaTempUrl) *MediaTempAggregate {
	return &MediaTempAggregate{
		MediaEntity:  mediaEntity,
		Prefix:       prefix,
		ContentType:  contentType,
		MediaTempUrl: mediaTempUrl,
	}
}

func (m *MediaTempAggregate) CheckPublic() *MediaTempAggregate {
	if m.MediaEntity.IsPublic() {
		return m
	}
	return nil
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
