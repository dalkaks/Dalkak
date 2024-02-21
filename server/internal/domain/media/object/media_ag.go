package mediaobject

type MediaTempAggregate struct {
	MediaEntity  *MediaEntity
	Prefix       Prefix
	ContentType  ContentType
	MediaTempUrl *MediaTempUrl
}

func NewMediaTempAggregate(mediaEntity *MediaEntity, prefix Prefix, contentType ContentType, mediaTempUrl *MediaTempUrl) *MediaTempAggregate {
	return &MediaTempAggregate{
		MediaEntity:  mediaEntity,
		Prefix:      prefix,
		ContentType: contentType,
		MediaTempUrl: mediaTempUrl,
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
