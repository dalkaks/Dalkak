package httputils

import (
	"dalkak/pkg/dtos"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// https://dev-api.dalkak.com -> dalkak.com
func ParseDomain(u string) (string, error) {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return "", &dtos.AppError{
			Code:    http.StatusInternalServerError,
			Message: "failed to parse url",
		}
	}

	host := parsedUrl.Hostname()

	if host == "localhost" {
		return "localhost", nil
	}

	host = strings.Split(host, ":")[0]

	parts := strings.Split(host, ".")
	if len(parts) >= 2 {
		host = parts[len(parts)-2] + "." + parts[len(parts)-1]
	}

	return host, nil
}

func GetUserInfoData(r *http.Request) (*dtos.UserInfo, error) {
	userInfo, ok := r.Context().Value("user").(dtos.UserInfo)
	if !ok {
		return nil, &dtos.AppError{
			Code:    http.StatusUnauthorized,
			Message: "user info not found",
		}
	}

	return &userInfo, nil
}

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
			return &dtos.AppError{
				Code:    http.StatusBadRequest,
				Message: "required query parameter not found",
			}
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
			return &dtos.AppError{
				Code:    http.StatusBadRequest,
				Message: "unexpected query parameter",
			}
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
			return &dtos.AppError{
				Code:    http.StatusBadRequest,
				Message: "failed to parse int",
			}
		}
		field.SetInt(intValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return &dtos.AppError{
				Code:    http.StatusBadRequest,
				Message: "failed to parse bool",
			}
		}
		field.SetBool(boolValue)
	default:
		return &dtos.AppError{
			Code:    http.StatusBadRequest,
			Message: "unsupported type",
		}
	}
	return nil
}
