package valis

import (
	"reflect"

	"github.com/soranoba/valis/code"
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
		validator.ErrorCollector().Add(validator.Location(), NewError(code.NotStruct, value))
		return
	}

	field := valishelpers.GetField(value, r.fieldPtr)
	validator.DiveField(field, func(v *Validator) {
		And(r.rules...).Validate(v, reflect.ValueOf(r.fieldPtr).Elem().Interface())
	})
}

// EachFields returns a new rule that verifies all field values of the struct meet the rules and all common rules.
func EachFields(rules ...Rule) Rule {
	return &eachFieldsRule{rules: rules}
}

func (rule *eachFieldsRule) Validate(validator *Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			fieldVal := val.Field(i)
			if !fieldVal.IsValid() {
				return
			}
			field := val.Type().Field(i)
			validator.DiveField(&field, func(v *Validator) {
				And(rule.rules...).Validate(v, fieldVal.Interface())
			})
		}
	default:
		validator.ErrorCollector().Add(validator.Location(), NewError(code.NotStruct, value))
	}
}
