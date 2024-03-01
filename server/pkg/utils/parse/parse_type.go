package parseutil

import (
	responseutil "dalkak/pkg/utils/response"
	"reflect"
	"strconv"
)

func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(value).Int()), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
}

func ToUnixTimestamp(value interface{}) (int64, error) {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		ts := reflect.ValueOf(value).Int()
		if ts < 0 {
			return 0, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
		}
		return int64(ts), nil
	case float64, float32:
		ts := int64(reflect.ValueOf(value).Float())
		if ts < 0 {
			return 0, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
		}
		return ts, nil
	default:
		return 0, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
}
