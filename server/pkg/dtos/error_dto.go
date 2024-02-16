package dtos

import "net/http"

type AppError struct {
	Code    int
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	return e.Message
}

const (
	ErrCodeBadRequest   = http.StatusBadRequest
	ErrCodeUnauthorized = http.StatusUnauthorized
	ErrCodeNotFound     = http.StatusNotFound
	ErrCodeInternal     = http.StatusInternalServerError

	ErrMsgJsonInvalid = "json: Invalid JSON"

	ErrMsgTokenParseFailed      = "token: Failed to parse token"
	ErrMsgTokenExpired          = "token: Token is expired"
	ErrMsgTokenInvalidClaim     = "token: Invalid token claim"
	ErrMsgTokenInvalidSignature = "token: Invalid token signature"
	ErrMsgTokenAccessNotFound   = "token: Access token not found"
	ErrMsgTokenRefeshNotFound   = "token: Refresh token not found"
	ErrMsgTokenSignFailed       = "token: Failed to sign token"

	ErrMsgMediaInvalidType = "media: Invalid media type"

	ErrMsgMediaInvalidKey = "media: Invalid key"
	ErrMsgMediaNotFound   = "media: Media not found"

	ErrMsgDBInternal = "db: Internal error"

	ErrMsgStorageInvalidURL = "storage: Invalid url"
	ErrMsgStorageNoSuchKey  = "storage: No such key"
	ErrMsgStorageInternal   = "storage: Internal error"

	ErrMsgRequestInvalid = "request: Invalid request"

	ErrMsgServerInternal = "server: Internal error"
)

func NewAppError(code int, message string, cause ...error) *AppError {
	var errCauese error
	if len(cause) > 0 || cause[0] != nil {
		errCauese = cause[0]
	}

	return &AppError{
		Code:    code,
		Message: message,
		Cause:   errCauese,
	}
}
