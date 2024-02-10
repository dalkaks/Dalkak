package httputils

import (
	"net/http"
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func HandleAppError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*AppError); ok {
		ErrorJSON(w, appErr, appErr.Code)
	} else {
		ErrorJSON(w, err, http.StatusInternalServerError)
	}
}
