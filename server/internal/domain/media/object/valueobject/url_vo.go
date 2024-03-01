package mediavalueobject

import (
	parseutil "dalkak/pkg/utils/parse"
	responseutil "dalkak/pkg/utils/response"
)

type MediaUrl struct {
	AccessUrl string  `json:"accessUrl"`
	UploadUrl *string `json:"uploadUrl,omitempty"`
}

func NewMediaUrl(staticLink, key string, uploadUrl ...string) *MediaUrl {
	var uploadUrlPtr *string
	if len(uploadUrl) > 0 {
		uploadUrlPtr = &uploadUrl[0]
	}

	return &MediaUrl{
		AccessUrl: parseutil.ConvertKeyToStaticLink(staticLink, key),
		UploadUrl: uploadUrlPtr,
	}
}

func NewMediaUrlWithOnlyAccessUrl(accessUrl string) *MediaUrl {
	return &MediaUrl{
		AccessUrl: accessUrl,
	}
}

func GenerateMediaTempKey(userId string, resource *MediaResource) (string, error) {
	if resource == nil {
		return "", responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}

	mediaTypeStr := resource.GetMediaType()
	extensionStr := resource.GetExtension()
	return "temp/" + userId + "/" + resource.Prefix.String() + "/" + mediaTypeStr + "/" + mediaTypeStr + "." + extensionStr, nil
}

func (mu *MediaUrl) GetUrlKey(staticLink string) string {
	return parseutil.ConvertStaticLinkToKey(staticLink, mu.AccessUrl)
}
