package httputils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}, wrapKey ...string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var wrappedData map[string]interface{}

	key := "data" // 기본값
	if len(wrapKey) > 0 && wrapKey[0] != "" {
		key = wrapKey[0]
	}

	wrappedData = map[string]interface{}{key: data}

	return json.NewEncoder(w).Encode(wrappedData)
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{
		Message: err.Error(),
	}

	_ = WriteJSON(w, statusCode, theError, "error")
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024 // one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// attempt to decode the data
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// make sure only one JSON value in payload
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}