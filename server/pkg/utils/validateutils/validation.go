package validateutils

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"net/http"
)

func Validate(req interfaces.ValidatableRequest) error {
	if !req.IsValid() {
		return &dtos.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
		}
	}
	return nil
}
