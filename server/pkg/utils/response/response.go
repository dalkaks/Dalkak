package responseutil

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

const (
	DataCodeSuccess = http.StatusOK
	DataCodeCreated = http.StatusCreated

	ErrCodeBadRequest   = http.StatusBadRequest
	ErrCodeUnauthorized = http.StatusUnauthorized
	ErrCodeForbidden    = http.StatusForbidden
	ErrCodeNotFound     = http.StatusNotFound
	ErrCodeTimeout      = http.StatusRequestTimeout
	ErrCodeInternal     = http.StatusInternalServerError

	ErrMsgRequestInvalid   = "REQUEST:INVALID_REQUEST"
	ErrMsgRequestUnauth    = "REQUEST:UNAUTHORIZED"
	ErrMsgRequestForbidden = "REQUEST:FORBIDDEN"
	ErrMsgRequestNotFound  = "REQUEST:NOT_FOUND"
	ErrMsgRequestTimeout   = "REQUEST:TIMEOUT"

	ErrMsgDataNotFound = "DATA:NOT_FOUND"

	ErrMsgTokenParseFailed      = "TOKEN:FAILED_TO_PARSE_TOKEN"
	ErrMsgTokenExpired          = "TOKEN:EXPIRED_TOKEN"
	ErrMsgTokenInvalidSignature = "TOKEN:INVALID_SIGNATURE"

	ErrMsgMetaMaskInvalidSignature = "METAMASK:INVALID_SIGNATURE"
	ErrMsgMetaMaskNotMatchAddress  = "METAMASK:INVALID_ADDRESS"

	ErrMsgServerInternal  = "SERVER:INTERNAL_ERROR"
	ErrMsgDBInternal      = "DB:INTERNAL_ERROR"
	ErrMsgStorageInternal = "STORAGE:INTERNAL_ERROR"
)

type AppData struct {
	Code int
	Data interface{}
}

type AppError struct {
	Code    int
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppData(data interface{}, code ...int) *AppData {
	statusCode := DataCodeSuccess
	if len(code) > 0 {
		statusCode = code[0]
	}
	return &AppData{
		Code: statusCode,
		Data: data,
	}
}

func NewAppError(code int, message string, cause ...error) *AppError {
	var errCauese error
	if len(cause) > 0 && cause[0] != nil {
		errCauese = cause[0]
	} else {
		errCauese = errors.New(message)
	}

	return &AppError{
		Code:    code,
		Message: message,
		Cause:   errCauese,
	}
}

func WriteToResponse(c fiber.Ctx, response interface{}) {
	switch resp := response.(type) {
	case *AppError:
		c.Status(resp.Code).JSON(fiber.Map{
			"error": resp.Message,
		})
	case *AppData:
		c.Status(resp.Code).JSON(fiber.Map{
			"data": resp.Data,
		})
	default:
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": response,
		})
	}
}
