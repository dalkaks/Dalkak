package payloads

import (
	"dalkak/config"
	"strings"
)

type BoardUploadMediaRequest struct {
	MediaType string `json:"mediaType"`
	Ext       string `json:"ext"`
}

func (req *BoardUploadMediaRequest) IsValid() bool {
	switch req.MediaType {
	case "image":
		ext := strings.ToLower(req.Ext)
		if _, ok := config.AllowedImageExtensions[ext]; ok {
			return true
		}
		// Todo: video
	}
	return false
}

type BoardUploadMediaResponse struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}
