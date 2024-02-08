package dtos

import (
	"errors"
	"io"
)

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
		return 0, errors.New("invalid media type")
	}
}

type MediaDto struct {
	Meta MediaMeta
	File io.Reader
}

func (m *MediaMeta) ToBoardImageDto() *BoardImageDto {
	return &BoardImageDto{
		Id:          m.ID,
		Extension:   m.Extension,
		ContentType: m.ContentType,
		Url:         m.URL,
	}
}
