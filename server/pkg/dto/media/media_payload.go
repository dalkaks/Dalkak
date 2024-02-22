package mediadto

type CreateTempMediaRequest struct {
	MediaType string `json:"mediaType"`
	Ext       string `json:"ext"`
	Prefix    string `json:"prefix"`
}

type CreateTempMediaResponse struct {
	Id        string `json:"id"`
	AccessUrl string `json:"accessUrl"`
	UploadUrl string `json:"uploadUrl"`
}

func NewUserCreateMediaResponse(id string, accessUrl string, uploadUrl string) *CreateTempMediaResponse {
	return &CreateTempMediaResponse{
		Id:        id,
		AccessUrl: accessUrl,
		UploadUrl: uploadUrl,
	}
}
