package dtos

import "io"

type MediaMeta struct {
	ID          string `json:"id"`
	Extension   string `json:"extension"`
	ContentType string `json:"contentType"`
	URL         string `json:"url"`
}

type MediaDto struct {
	Meta MediaMeta
	File io.Reader
}

func (m *MediaMeta) ToBoardImageDto(boardId *string) *BoardImageDto {
	return &BoardImageDto{
		Id:          m.ID,
		BoardId:     boardId,
		Extension:   m.Extension,
		ContentType: m.ContentType,
		Url:         m.URL,
	}
}
