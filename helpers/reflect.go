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
