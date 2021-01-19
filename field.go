package valis

import (
	"errors"
	valishelpers "github.com/soranoba/valis/helpers"
	"reflect"
)

type (
	fieldRule struct {
		fieldPtr interface{}
		rules      []Rule
	}
)

// Field returns a new rule that verifies the filed value meets the rules and all common rules.
func Field(fieldPtr interface{}, rules ...Rule) Rule {
	val := reflect.ValueOf(fieldPtr)
	if val.Kind() != reflect.Ptr {
		panic("fieldPtr must be a pointer of any field")
	}
	return &fieldRule{fieldPtr: fieldPtr, rules: rules}
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
			errors.New("must be a struct"),
		})
		return
	}

	field := valishelpers.GetField(value, r.fieldPtr)
	And(r.rules...).Validate(validator.WithField(field), reflect.ValueOf(r.fieldPtr).Elem().Interface())
}
