package validateutils

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"errors"
)

func Validate(req interfaces.ValidatableRequest) error {
	if !req.IsValid() {
		return dtos.NewAppError(dtos.ErrCodeBadRequest, dtos.ErrMsgRequestInvalid, errors.New("invalid request"))
	}
	return nil
}
