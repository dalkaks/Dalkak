package parseutil

import (
	responseutil "dalkak/pkg/utils/response"
	"strings"
)

func ConvertContentTypeToMediaType(contentType string) (string, error) {
	parts := strings.Split(contentType, "/")
	if len(parts) < 2 {
		return "", responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
	return parts[0], nil
}
