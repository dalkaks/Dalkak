package dtos

import "io"

type MediaMeta struct {
	ID          string `json:"id"`
	Extension   string `json:"extension"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

type MediaDto struct {
	Meta MediaMeta
	File io.Reader
}

type ImageDto struct {
	MediaMeta MediaMeta
	File      io.Reader
}
