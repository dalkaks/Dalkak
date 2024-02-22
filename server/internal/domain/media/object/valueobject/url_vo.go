package mediavalueobject

import parseutil "dalkak/pkg/utils/parse"

type MediaTempUrl struct {
	AccessUrl string  `json:"accessUrl"`
	UploadUrl *string `json:"uploadUrl,omitempty"`
}

func NewMediaTempUrl(staticLink, key, uploadUrl string) *MediaTempUrl {
	var uploadUrlPtr *string
	if uploadUrl != "" {
		uploadUrlPtr = &uploadUrl
	}
	return &MediaTempUrl{
		AccessUrl: parseutil.ConvertKeyToStaticLink(staticLink, key),
		UploadUrl: uploadUrlPtr,
	}
}

func GenerateMediaTempKey(userId string, prefix Prefix, contentType ContentType) string {
	mediaTypeStr := contentType.ConvertToMediaType()
	extensionStr := contentType.ConvertToExtension()
	return "temp/" + userId + "/" + prefix.String() + "/" + mediaTypeStr + "/" + mediaTypeStr + "." + extensionStr
}
