package valis

import (
	"errors"
	"reflect"
)

type (
	fieldRule struct {
		fieldValue reflect.Value
		rules      []Rule
	}
)

func Field(filedValue interface{}, rules ...Rule) Rule {
	return &fieldRule{fieldValue: reflect.ValueOf(filedValue), rules: rules}
}

func (r *fieldRule) Validate(validator *Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		validator.ErrorCollector().Add(validator.Location(), &ErrorDetail{
			InvalidTypeCode,
			r,
			value,
			nil,
			errors.New("must be struct"),
		})
	}

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if r.fieldValue == f {
			for _, rule := range r.rules {
				loc := validator.Location().FieldLocation(val.Type().Field(i))
				rule.Validate(validator.WithLocation(loc), r.fieldValue.Interface())
			}
		}
	}
}
