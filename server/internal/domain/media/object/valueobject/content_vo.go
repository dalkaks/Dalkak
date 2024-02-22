package mediavalueobject

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
	mediaType, extension := SplitContentType(contentType.String())
	if mediaType == "" || extension == "" {
		return false
	}

	return IsAllowedExtension(mediaType, extension)
}

func IsAllowedMediaType(mediaType string) bool {
	_, exists := AllowedContentType[mediaType]
	return exists
}

func IsAllowedExtension(mediaType, extension string) bool {
	if allowedExtensions, ok := AllowedContentType[mediaType]; ok {
		return allowedExtensions[extension]
	}
	return false
}

func SplitContentType(contentTypeStr string) (string, string) {
	parts := strings.Split(contentTypeStr, "/")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func (contentType ContentType) ConvertToMediaType() string {
	mediaType, _ := SplitContentType(contentType.String())
	return mediaType
}

func (contentType ContentType) ConvertToExtension() string {
	_, extension := SplitContentType(contentType.String())
	return extension
}
