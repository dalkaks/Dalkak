package mediadto

type CreateTempMediaRequest struct {
	MediaType string `json:"mediaType" validate:"required"`
	Ext       string `json:"ext" validate:"required"`
	Prefix    string `json:"prefix" validate:"required"`
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
