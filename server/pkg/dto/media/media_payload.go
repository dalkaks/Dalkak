package mediadto

type CreateMediaTempRequest struct {
	MediaType string `json:"mediaType" validate:"required"`
	Ext       string `json:"ext" validate:"required"`
	Prefix    string `json:"prefix" validate:"required"`
}

type CreateMediaTempResponse struct {
	Id        string `json:"id"`
	AccessUrl string `json:"accessUrl"`
	UploadUrl string `json:"uploadUrl"`
}

func NewCreateMediaTempResponse(id string, accessUrl string, uploadUrl string) *CreateMediaTempResponse {
	return &CreateMediaTempResponse{
		Id:        id,
		AccessUrl: accessUrl,
		UploadUrl: uploadUrl,
	}
}

type GetMediaTempRequest struct {
	MediaType string `query:"media-type" validate:"required"`
	Prefix    string `query:"prefix" validate:"required"`
}

type GetMediaTempResponse struct {
	Id          string `json:"id"`
	ContentType string `json:"contentType"`
	AccessUrl   string `json:"accessUrl"`
}

func NewGetMediaTempResponse(id string, contentType string, accessUrl string) *GetMediaTempResponse {
	return &GetMediaTempResponse{
		Id:          id,
		ContentType: contentType,
		AccessUrl:   accessUrl,
	}
}
