package reflectutils

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"
)

func MapToStruct[T any](data map[string]interface{}) (*T, error) {
	var result T
	t := reflect.TypeOf(result)

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("MapToStruct only supports structs")
	}

	ptr := reflect.New(t)
	for key, value := range data {
		field := ptr.Elem().FieldByName(key)
		if field.IsValid() && field.CanSet() {
			val := reflect.ValueOf(value)
			if val.Type().AssignableTo(field.Type()) {
				field.Set(val)
			} else {
				convertedValue, err := tryConvertType(val, field.Type())
				if err != nil {
					return nil, fmt.Errorf("error converting types for field %s: %v", key, err)
				}
				field.Set(convertedValue)
			}
		}
	}

	return ptr.Interface().(*T), nil
}

// suport time type is RFC3339
func tryConvertType(value reflect.Value, targetType reflect.Type) (reflect.Value, error) {
	if value.CanConvert(targetType) {
		return value.Convert(targetType), nil
	}

	if value.Kind() == reflect.Float64 && targetType.Kind() == reflect.Int {
		convertedValue := int(value.Float())
		return reflect.ValueOf(convertedValue), nil
	}

	if value.Kind() == reflect.String && targetType == reflect.TypeOf(time.Time{}) {
		str := value.String()
		parsedTime, err := time.Parse(time.RFC3339, str)
		if err == nil {
			return reflect.ValueOf(parsedTime), nil
		}
	}

	log.Printf("cannot convert %v to %v", value.Type(), targetType)
	return reflect.Value{}, errors.New("cannot convert types")
}
