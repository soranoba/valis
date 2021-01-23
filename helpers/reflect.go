package valishelpers

import (
	"reflect"
)

// GetField returns the *reflect.StructField of the fieldPointer. When it is not found, it panics.
func GetField(structPointer interface{}, fieldPointer interface{}) *reflect.StructField {
	structVal := reflect.ValueOf(structPointer)
	fieldVal := reflect.ValueOf(fieldPointer)

	if structVal.Kind() != reflect.Ptr {
		panic("structPointer must be a pointer of struct")
	}
	if fieldVal.Kind() != reflect.Ptr {
		panic("fieldPointer must be a pointer")
	}

	for structVal.Kind() == reflect.Ptr {
		structVal = structVal.Elem()
	}
	if structVal.Kind() != reflect.Struct {
		panic("structPointer must be a pointer of struct")
	}

	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i).Addr()
		if field == fieldVal {
			strField := structVal.Type().Field(i)
			return &strField
		}
	}
	panic("invalid fieldPointer")
}

func IsNumeric(v interface{}) bool {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
