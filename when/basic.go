// Package when implements some valis.WhenRule.
package when

import (
	"reflect"

	"github.com/soranoba/valis"
	valishelpers "github.com/soranoba/valis/helpers"
)

// IsType returns a valis.WhenRule that verifies the value meets the rules when the type of the validating value has the same type.
func IsType(ty reflect.Type, rules ...valis.Rule) *valis.WhenRule {
	if ty == nil {
		panic("invalid type")
	}
	cond := func(ctx *valis.WhenContext) bool {
		val := reflect.ValueOf(ctx.Value())
		if !val.IsValid() {
			return false
		}
		return val.Type() == ty
	}
	return valis.When(cond, rules...)
}

// IsTypeOrPtr returns a valis.WhenRule that verifies the value meets the rules when the type of the validating value has the same type or pointer type.
func IsTypeOrPtr(ty reflect.Type, rules ...valis.Rule) *valis.WhenRule {
	if ty == nil {
		panic("invalid type")
	}
	cond := func(ctx *valis.WhenContext) bool {
		val := reflect.ValueOf(ctx.Value())
		if !val.IsValid() {
			return false
		}
		return val.Type() == ty || val.Type() == reflect.PtrTo(ty)
	}
	return valis.When(cond, rules...)
}

// IsTypeOrElem returns a valis.WhenRule that verifies the value meets the rules when the type of the validating value has the same type or elem type.
func IsTypeOrElem(ty reflect.Type, rules ...valis.Rule) *valis.WhenRule {
	if ty == nil {
		panic("invalid type")
	}
	cond := func(ctx *valis.WhenContext) bool {
		val := reflect.ValueOf(ctx.Value())
		if !val.IsValid() {
			return false
		}
		return val.Type() == ty || reflect.PtrTo(val.Type()) == ty
	}
	return valis.When(cond, rules...)
}

// IsNumeric returns a valis.WhenRule that verifies the value meets the rules when the value is numeric.
func IsNumeric(rules ...valis.Rule) *valis.WhenRule {
	cond := func(ctx *valis.WhenContext) bool {
		return valishelpers.IsNumeric(ctx.Value())
	}
	return valis.When(cond, rules...)
}
