package mediavalueobject

import (
	parseutil "dalkak/pkg/utils/parse"
	responseutil "dalkak/pkg/utils/response"
)

type MediaTempUrl struct {
	AccessUrl string  `json:"accessUrl"`
	UploadUrl *string `json:"uploadUrl,omitempty"`
}

func NewMediaTempUrl(staticLink, key string, uploadUrl ...string) *MediaTempUrl {
	var uploadUrlPtr *string
	if len(uploadUrl) > 0 {
		uploadUrlPtr = &uploadUrl[0]
	}

	return &MediaTempUrl{
		AccessUrl: parseutil.ConvertKeyToStaticLink(staticLink, key),
		UploadUrl: uploadUrlPtr,
	}
}

func NewMediaTempUrlWithOnlyAccessUrl(accessUrl string) *MediaTempUrl {
	return &MediaTempUrl{
		AccessUrl: accessUrl,
	}
}

func GenerateMediaTempKey(userId string, resource *MediaResource) (string, error) {
	if resource == nil {
		return "", responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}

	mediaTypeStr := resource.ContentType.ConvertToMediaType()
	extensionStr := resource.ContentType.ConvertToExtension()
	return "temp/" + userId + "/" + resource.Prefix.String() + "/" + mediaTypeStr + "/" + mediaTypeStr + "." + extensionStr, nil
}

func (mu *MediaTempUrl) GetUrlKey(staticLink string) string {
	return parseutil.ConvertStaticLinkToKey(staticLink, mu.AccessUrl)
}
