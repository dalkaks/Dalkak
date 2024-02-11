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
	Prefix      string `json:"prefix"`
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
