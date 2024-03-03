package mediavalueobject

import (
	parseutil "dalkak/pkg/utils/parse"
	responseutil "dalkak/pkg/utils/response"
	"fmt"
	"strings"
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

func ConvertMediaUrl(staticLink, prefix, id, mediaType, extensionStr string) *MediaUrl {
	accessUrl := parseutil.ConvertKeyToStaticLink(staticLink, fmt.Sprintf("%s/%s/%s/%s.%s", prefix, id, mediaType, mediaType, extensionStr))
	return NewMediaUrlWithOnlyAccessUrl(accessUrl)
}

func GenerateMediaTempKey(userId string, resource *MediaResource) (string, error) {
	if resource == nil {
		return "", responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}

	mediaTypeStr := resource.GetMediaType()
	extensionStr := resource.GetExtension()
	return "temp/" + userId + "/" + resource.Prefix.String() + "/" + mediaTypeStr + "/" + mediaTypeStr + "." + extensionStr, nil
}

// temp/{userId}/{prefix}/{mediaType}/{contentType} -> {prefix}/{id}/{mediaType}/{contentType}
func (mu *MediaUrl) ConvertMediaTempToFormalUrl(staticLink, id string) (*MediaUrl, error) {
	tempKey := parseutil.ConvertStaticLinkToKey(staticLink, mu.AccessUrl)
	parts := strings.Split(tempKey, "/")
	if len(parts) < 5 {
		return nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}

	return &MediaUrl{
		AccessUrl: parseutil.ConvertKeyToStaticLink(staticLink, fmt.Sprintf("%s/%s/%s/%s", parts[2], id, parts[3], parts[4])),
	}, nil
}

func (mu *MediaUrl) GetUrlKey(staticLink string) string {
	return parseutil.ConvertStaticLinkToKey(staticLink, mu.AccessUrl)
}
