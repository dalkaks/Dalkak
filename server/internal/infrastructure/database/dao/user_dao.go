package dao

type MediaTempDao struct {
	Id          string `json:"id"`
	Prefix      string `json:"prefix"`
	Extension   string `json:"extension"`
	ContentType string `json:"content_type"`
	Url         string `json:"url"`
	IsConfirm   bool   `json:"is_confirm"`
	Timestamp   int64  `json:"timestamp"`
}
