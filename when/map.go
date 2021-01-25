package when

import (
	"github.com/soranoba/valis"
	"reflect"
)

// HasKey returns a new rule that verifies the value meets the rules when the value has the key.
func HasKey(key interface{}, rules ...valis.Rule) *valis.WhenRule {
	cond := func(ctx *valis.WhenContext) bool {
		val := reflect.ValueOf(ctx.Value())
		for val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() != reflect.Map {
			return false
		}

		keyVal := reflect.ValueOf(key)
		if !keyVal.Type().AssignableTo(val.Type().Key()) {
			return false
		}
		return val.MapIndex(keyVal).IsValid()
	}
	return valis.When(cond, rules...)
}
