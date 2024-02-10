package dtos

import (
	"net/http"
)

type UploadMediaDto struct {
	MediaType MediaType
	Ext       string
	Prefix    string
}

type MediaMeta struct {
	ID          string `json:"id"`
	Extension   string `json:"extension"`
	ContentType string `json:"contentType"`
	URL         string `json:"url"`
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
		return 0, &AppError{
			Code:    http.StatusBadRequest,
			Message: "invalid media type",
		}
	}
}

func (m *MediaMeta) ToBoardImageDto() *BoardImageDto {
	return &BoardImageDto{
		Id:          m.ID,
		Extension:   m.Extension,
		ContentType: m.ContentType,
		Url:         m.URL,
	}
}
