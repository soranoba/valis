package when

import (
	"github.com/soranoba/valis"
	"reflect"
)

// IsSliceOrArray returns a valis.WhenRule that verifies the value meets the rules when the value is array or slice.
func IsSliceOrArray(rules ...valis.Rule) *valis.WhenRule {
	cond := func(ctx *valis.WhenContext) bool {
		val := reflect.ValueOf(ctx.Value())
		for val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		return val.Kind() == reflect.Array || val.Kind() == reflect.Slice
	}
	return valis.When(cond, rules...)
}
