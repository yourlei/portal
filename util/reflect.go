package util

import (
  "fmt"
	"reflect"
)
// struct to map
func Struct2Map(obj interface{}) (data map[string] interface{}, err error) {
    data = make(map[string] interface{})
    keys := reflect.TypeOf(obj)
    values := reflect.ValueOf(obj)
    for i := 0; i < keys.NumField(); i++ {
      data[keys.Field(i).Name] = values.Field(i).Interface()
    }
    err = nil
    return
}
// Set struct field
func setField(obj interface{}, name string, value interface{}) error {
	structData := reflect.ValueOf(obj).Elem()
    fieldValue := structData.FieldByName(name)

    if !fieldValue.IsValid() {
        return fmt.Errorf("utils.setField() No such field: %s in obj ", name)
    }

    if !fieldValue.CanSet() {
        return fmt.Errorf("Cannot set %s field value ", name)
    }

   	fieldType := fieldValue.Type()
    val       := reflect.ValueOf(value)

	valTypeStr   := val.Type().String()
	fieldTypeStr := fieldType.String()
	if valTypeStr == "float64" && fieldTypeStr == "int" {
		val = val.Convert(fieldType)
	} else if fieldType != val.Type() {
        return fmt.Errorf("Provided value type " + valTypeStr + " didn't match obj field type " + fieldTypeStr)
    }
    fieldValue.Set(val)
    return nil
}
// SetStructByJSON 由json对象生成 struct
func SetStructByJSON(obj interface{}, mapData map[string]interface{}) error {
	for key, value := range mapData {
        if err := setField(obj, key, value); err != nil {
			fmt.Println(err.Error())
            return err
        }
    }
	return nil
}