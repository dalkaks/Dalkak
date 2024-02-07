package dtos

type BoardImageDto struct {
	Id          string  `json:"id"`
	BoardId     *string `json:"boardId,omitempty"`
	Extension   string  `json:"extension"`
	ContentType string  `json:"contentType"`
	Url         string  `json:"url"`
}
