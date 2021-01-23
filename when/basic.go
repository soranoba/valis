package when

import (
	valishelpers "github.com/soranoba/valis/helpers"
	"reflect"

	"github.com/soranoba/valis"
)

// IsType returns a new rule that verifies the value meets the rules when the type of the value and the ty are the same.
func IsType(ty reflect.Type, rules ...valis.Rule) valis.Rule {
	if ty == nil {
		panic("invalid ty")
	}
	cond := func(value interface{}) bool {
		val := reflect.ValueOf(value)
		if !val.IsValid() {
			return false
		}
		return val.Type() == ty
	}
	return valis.When(cond, rules...)
}

// IsTypeOrPtr returns a new rule that verifies the value meets the rules when the type of the value is ty or pointer of ty.
func IsTypeOrPtr(ty reflect.Type, rules ...valis.Rule) valis.Rule {
	if ty == nil {
		panic("invalid ty")
	}
	cond := func(value interface{}) bool {
		val := reflect.ValueOf(value)
		if !val.IsValid() {
			return false
		}
		return val.Type() == ty || val.Type() == reflect.PtrTo(ty)
	}
	return valis.When(cond, rules...)
}

// IsTypeOrPtr returns a new rule that verifies the value meets the rules when the type of the value is ty or elem of ty.
func IsTypeOrElem(ty reflect.Type, rules ...valis.Rule) valis.WhenRule {
	if ty == nil {
		panic("invalid ty")
	}
	cond := func(value interface{}) bool {
		val := reflect.ValueOf(value)
		if !val.IsValid() {
			return false
		}
		return val.Type() == ty || reflect.PtrTo(val.Type()) == ty
	}
	return valis.When(cond, rules...)
}

// IsNumeric returns a new valis.WhenRule that verifies the value meets the rules when the value is numeric.
func IsNumeric(rules ...valis.Rule) valis.WhenRule {
	return valis.When(valishelpers.IsNumeric, rules...)
}
