package reflectutils

import (
	"fmt"
	"reflect"
)

func MapToStruct(data map[string]interface{}, result interface{}) error {
	for key, value := range data {
		err := setField(result, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structType := structValue.Type()

	var field reflect.Value

  for i := 0; i < structValue.NumField(); i++ {
    structField := structType.Field(i)
    if structField.Name == name {
        field = structValue.Field(i)
        break
    }
  }

	if !field.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}
	if !field.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

  fieldType := field.Type()
  val := reflect.ValueOf(value)
	if fieldType != val.Type() {
		return fmt.Errorf("Provided value type didn't match obj field type")
	}

	field.Set(val)
	return nil
}
