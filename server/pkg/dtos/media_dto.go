package dtos

import (
	"dalkak/config"
	"errors"
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
}

func GenerateTempKey(prefix string, mediaType string, id string, ext string) string {
	return "temp/" + prefix + "/" + mediaType + "/" + id + "." + ext
}

func GenerateMediaPath(prefix string, mediaType string, id string) string {
	return prefix + "/" + mediaType + "/" + id
}

func GenerateContentType(mediaType string, ext string) string {
	return mediaType + "/" + ext
}

func (m *MediaHeadDto) Verify(meta *MediaMeta) bool {
	return int64(config.MaxUploadSize) > m.Length && m.ContentType == meta.ContentType && m.URL == meta.URL
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
		return 0, NewAppError(ErrCodeBadRequest, ErrMsgMediaInvalidType, errors.New("invalid media type"))
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
