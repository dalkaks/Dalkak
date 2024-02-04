package dtos

type Image struct {
	Extension string  `json:"extension"`
	URL       *string `json:"url,omitempty"`
}

type ImageData struct {
	Meta Image
	Data []byte
}
