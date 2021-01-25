package when

import (
	"github.com/soranoba/valis"
	"reflect"
)

// IsStruct returns a valis.WhenRule that verifies the value meets the rules when the value is struct.
func IsStruct(rules ...valis.Rule) *valis.WhenRule {
	cond := func(ctx *valis.WhenContext) bool {
		val := reflect.ValueOf(ctx.Value())
		for val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		return val.Kind() == reflect.Struct
	}
	return valis.When(cond, rules...)
}

// HasFieldTag returns a valis.WhenRule that verifies the value meets the rules when the value is in a struct field and the field has the key.
func HasFieldTag(key string, rules ...valis.Rule) *valis.WhenRule {
	cond := func(ctx *valis.WhenContext) bool {
		loc := ctx.Location()
		if loc.Kind() != valis.LocationKindField {
			return false
		}
		_, ok := loc.Field().Tag.Lookup(key)
		return ok
	}
	return valis.When(cond, rules...)
}
