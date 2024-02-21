package mediaobject

import (
	responseutil "dalkak/pkg/utils/response"
	"strings"
)

const MaxUploadSize = 32 << 20 // 32 MB

var AllowedContentType = map[string]map[string]bool{
	"image": {
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"gif":  true,
		"bmp":  true,
		"webp": true,
	},
	"video": {
		"mp4":  true,
		"avi":  true,
		"webm": true,
	},
}

type ContentType string

func NewContentType(mediaType, extension string) (ContentType, error) {
	contentType := ContentType(mediaType + "/" + extension)
	if !contentType.IsAllowedContentType() {
		return "", responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
	return contentType, nil
}

func (contentType ContentType) String() string {
	return string(contentType)
}

func (contentType ContentType) IsAllowedContentType() bool {
	part := strings.Split(contentType.String(), "/")
	if len(part) != 2 {
		return false
	}

	mediaType, extension := part[0], part[1]
	if allowedExtensions, ok := AllowedContentType[mediaType]; ok {
		for ext, allowed := range allowedExtensions {
			if allowed && ext == extension {
				return true
			}
		}
	}
	return false
}

func (contentType ContentType) ConvertToMediaType() string {
	parts := strings.Split(contentType.String(), "/")
	return parts[0]
}

func (contentType ContentType) ConvertToExtension() string {
	parts := strings.Split(contentType.String(), "/")
	return parts[1]
}
