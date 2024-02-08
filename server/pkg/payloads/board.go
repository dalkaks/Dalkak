package payloads

type BoardUploadMediaRequest struct {
	MediaType string `json:"mediaType"`
	Ext       string `json:"ext"`
}

func (req *BoardUploadMediaRequest) IsValid() bool {
	switch req.MediaType {
	case "image", "video":
		return true
	default:
		return false
	}
}

type BoardUploadMediaResponse struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}
