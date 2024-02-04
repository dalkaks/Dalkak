package dtos

type MediaMeta struct {
	ID        string `json:"id"`
	Extension string `json:"extension"`
	URL       string `json:"url"`
}

type MediaDto struct {
	Meta MediaMeta
	Data []byte
}

type ImageDto struct {
	MediaMeta MediaMeta
	Data      []byte
}
