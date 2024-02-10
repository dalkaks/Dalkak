package httputils

import (
	"dalkak/pkg/dtos"
	"net/http"
)

func HandleAppError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*dtos.AppError); ok {
		ErrorJSON(w, appErr, appErr.Code)
	} else {
		ErrorJSON(w, err, http.StatusInternalServerError)
	}
}
