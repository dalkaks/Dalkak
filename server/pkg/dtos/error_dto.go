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

	ErrMsgJsonInvalid    = "REQUEST:INVALID_JSON"
	ErrMsgRequestInvalid = "REQUEST:INVALID_REQUEST"

	ErrMsgTokenParseFailed      = "TOKEN:FAILED_TO_PARSE_TOKEN"
	ErrMsgTokenExpired          = "TOKEN:EPIRED_TOKEN"
	ErrMsgTokenInvalidClaim     = "TOKEN:INVALID_CLAIM"
	ErrMsgTokenInvalidSignature = "TOKEN:INVALID_SIGNATURE"
	ErrMsgTokenAccessNotFound   = "TOKEN:NO_ACCESS_TOKEN"
	ErrMsgTokenRefeshNotFound   = "TOKEN:NO_REFRESH_TOKEN"
	ErrMsgTokenSignFailed       = "TOKEN:FAILED_TO_SIGN_TOKEN"

	ErrMsgMediaInvalidType = "MEDIA:INVALID_TYPE"
	ErrMsgMediaInvalidKey  = "MEDIA:INVALID_KEY"
	ErrMsgMediaNotFound    = "MEDIA:NOT_FOUND"

	ErrMsgStorageInvalidURL = "STORAGE:INVALID_URL"
	ErrMsgStorageNoSuchKey  = "STORAGE:NO_SUCH_KEY"
	ErrMsgStorageInternal   = "STORAGE:INTERNAL_ERROR"

	ErrMsgDBInternal     = "DB:INTERNAL_ERROR"
	ErrMsgServerInternal = "SERVER:INTERNAL_ERROR"
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
