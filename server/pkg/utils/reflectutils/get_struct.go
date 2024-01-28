package reflectutils

import (
	"errors"
	"net/http"
)

func GetRequestData(r *http.Request, target interface{}) error {
	reqMap, ok := r.Context().Value("request").(map[string]interface{})
	if !ok {
		return errors.New("invalid request")
	}

	err := MapToStruct(reqMap, target)
	if err != nil {
		return err
	}

	return nil
}
