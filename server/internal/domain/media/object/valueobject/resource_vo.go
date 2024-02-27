package mediavalueobject

import (
	responseutil "dalkak/pkg/utils/response"
	"strings"
)

var AllowedPrefixes = map[string]bool{
	"board": true,
}

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

type MediaResource struct {
	Prefix      MediaPrefix
	ContentType MediaContentType
}

type MediaPrefix string

type MediaContentType string

func NewMediaResource(prefixStr string, contentTypeStr string) (*MediaResource, error) {
	prefix, err := NewPrefix(prefixStr)
	if err != nil {
		return nil, err
	}
	contentType, err := NewContentType(contentTypeStr)
	if err != nil {
		return nil, err
	}
	return &MediaResource{
		Prefix:      prefix,
		ContentType: contentType,
	}, nil
}

func NewPrefix(prefixStr string) (MediaPrefix, error) {
	prefix := MediaPrefix(prefixStr)
	if !isAllowedPrefix(prefix) {
		return "", responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
	return prefix, nil
}

func NewContentType(contentTypeStr string) (MediaContentType, error) {
	mediaType, extension := SplitContentType(contentTypeStr)
	if mediaType == "" || extension == "" {
		return "", responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
	
	contentType := MediaContentType(mediaType + "/" + extension)
	if !isAllowedContentType(contentType) {
		return "", responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
	return contentType, nil
}

func (prefix MediaPrefix) String() string {
	return string(prefix)
}

func isAllowedPrefix(prefix MediaPrefix) bool {
	_, ok := AllowedPrefixes[prefix.String()]
	return ok
}

func (contentType MediaContentType) String() string {
	return string(contentType)
}

func isAllowedContentType(contentType MediaContentType) bool {
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

func MergeContentType(mediaType, extension string) string {
	return mediaType + "/" + extension
}

func (contentType MediaContentType) ConvertToMediaType() string {
	mediaType, _ := SplitContentType(contentType.String())
	return mediaType
}

func (contentType MediaContentType) ConvertToExtension() string {
	_, extension := SplitContentType(contentType.String())
	return extension
}
