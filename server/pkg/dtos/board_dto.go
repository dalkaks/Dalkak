package dtos

type BoardImageDto struct {
	Id          string  `json:"id"`
	Extension   string  `json:"extension"`
	ContentType string  `json:"contentType"`
	Url         string  `json:"url"`
}
