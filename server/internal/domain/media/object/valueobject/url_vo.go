package mediavalueobject

import parseutil "dalkak/pkg/utils/parse"

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

func GenerateMediaTempKey(userId string, prefix Prefix, contentType ContentType) string {
	mediaTypeStr := contentType.ConvertToMediaType()
	extensionStr := contentType.ConvertToExtension()
	return "temp/" + userId + "/" + prefix.String() + "/" + mediaTypeStr + "/" + mediaTypeStr + "." + extensionStr
}

func (mu *MediaTempUrl) GetUrlKey(staticLink string) string {
	return parseutil.ConvertStaticLinkToKey(staticLink, mu.AccessUrl)
}
