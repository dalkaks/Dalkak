package appdto

import (
	"dalkak/config"
	responseutil "dalkak/pkg/utils/response"
)

type MediaMeta struct {
	ID          string `json:"id"`
	Prefix      string `json:"prefix"`
	Extension   string `json:"extension"`
	ContentType string `json:"contentType"`
	URL         string `json:"url"`
}

type UploadMediaDto struct {
	MediaType MediaType
	Ext       string
	Prefix    string
}

type FindUserUploadMediaDto struct {
	MediaType MediaType
	Prefix    string
	IsConfirm *bool
}

type UpdateUserUploadMediaDto struct {
	IsConfirm bool
}

type MediaHeadDto struct {
	Key         string
	ContentType string
	Length      int64
	URL         string
	MetaUserId  string
}

func (m *MediaHeadDto) Verify(userId string, meta *MediaMeta) bool {
	return int64(config.MaxUploadSize) > m.Length && m.ContentType == meta.ContentType && m.URL == meta.URL && userId == m.MetaUserId
}


type MediaType int

const (
	Image MediaType = iota + 1
	Video
)

func (m MediaType) String() string {
	return [...]string{"image", "video"}[m-1]
}

func ToMediaType(s string) (MediaType, error) {
	switch s {
	case "image":
		return Image, nil
	case "video":
		return Video, nil
	default:
		return 0, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgMediaInvalidType)
	}
}

func NewFindUserUploadMediaDto(mediaType MediaType, prefix string, opts ...FindUserUploadMediaOption) *FindUserUploadMediaDto {
	dto := &FindUserUploadMediaDto{
		MediaType: mediaType,
		Prefix:    prefix,
	}

	for _, opt := range opts {
		opt(dto)
	}

	return dto
}

type FindUserUploadMediaOption func(*FindUserUploadMediaDto)

func WithConfirm(isConfirm bool) FindUserUploadMediaOption {
	return func(f *FindUserUploadMediaDto) {
		f.IsConfirm = &isConfirm
	}
}
