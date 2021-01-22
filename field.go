package valis

import (
	"github.com/soranoba/valis/code"
	"reflect"

	valishelpers "github.com/soranoba/valis/helpers"
)

type (
	fieldRule struct {
		fieldPtr interface{}
		rules    []Rule
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
		validator.ErrorCollector().Add(validator.Location(), NewError(code.StructOnly, value))
		return
	}

	field := valishelpers.GetField(value, r.fieldPtr)
	validator.DiveField(field, func(v *Validator) {
		And(r.rules...).Validate(v, reflect.ValueOf(r.fieldPtr).Elem().Interface())
	})
}
