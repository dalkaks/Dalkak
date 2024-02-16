package parseutils

import (
	"dalkak/pkg/dtos"
	"errors"
	"strings"
)

func ConvertContentTypeToMediaType(contentType string) (string, error) {
	parts := strings.Split(contentType, "/")
	if len(parts) < 2 {
		return "", dtos.NewAppError(dtos.ErrCodeBadRequest, dtos.ErrMsgRequestInvalid, errors.New("invalid content type"))
	}
	return parts[0], nil
}
