package parseutils

import (
	"dalkak/pkg/dtos"
	"errors"
	"net/http"
	"reflect"
	"strconv"
)

func GetQuery(r *http.Request, data interface{}) error {
	queryParams := r.URL.Query()
	v := reflect.ValueOf(data).Elem()
	t := v.Type()

	usedQueryParams := make(map[string]bool)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		queryTag := fieldType.Tag.Get("query")
		requiredTag := fieldType.Tag.Get("required")

		if queryTag == "" || !field.CanSet() {
			continue // 태그가 없거나 설정 불가능한 필드는 건너뜀
		}

		paramValue, ok := queryParams[queryTag]
		if requiredTag == "true" && !ok {
			return dtos.NewAppError(dtos.ErrCodeBadRequest, dtos.ErrMsgRequestInvalid, errors.New("required query parameter not found"))
		}

		if ok && len(paramValue) > 0 {
			err := reflectValueSet(field, paramValue[0])
			if err != nil {
				return err
			}
		}
		usedQueryParams[queryTag] = true
	}

	for key := range queryParams {
		if _, ok := usedQueryParams[key]; !ok {
			return dtos.NewAppError(dtos.ErrCodeBadRequest, dtos.ErrMsgRequestInvalid, errors.New("unsupported query parameter"))
		}
	}
	return nil
}

// reflect.Value와 해당 타입에 맞는 값을 설정하는 함수
func reflectValueSet(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return dtos.NewAppError(dtos.ErrCodeBadRequest, dtos.ErrMsgRequestInvalid, err)
		}
		field.SetInt(intValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return dtos.NewAppError(dtos.ErrCodeBadRequest, dtos.ErrMsgRequestInvalid, err)
		}
		field.SetBool(boolValue)
	default:
		return dtos.NewAppError(dtos.ErrCodeBadRequest, dtos.ErrMsgRequestInvalid, errors.New("unsupported type"))
	}
	return nil
}

func GetUserInfoData(r *http.Request) (*dtos.UserInfo, error) {
	userInfo, ok := r.Context().Value("user").(dtos.UserInfo)
	if !ok {
		return nil, dtos.NewAppError(dtos.ErrCodeUnauthorized, dtos.ErrMsgTokenAccessNotFound, errors.New("user info not found"))
	}

	return &userInfo, nil
}
