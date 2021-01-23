package when

import (
	"github.com/soranoba/valis"
	"reflect"
)

// Returns a valis.WhenRule that verifies the value meets the rules when the value is struct.
func IsStruct(rules ...valis.Rule) *valis.WhenRule {
	cond := func(value interface{}) bool {
		val := reflect.ValueOf(value)
		for val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		return val.Kind() == reflect.Struct
	}
	return valis.When(cond, rules...)
}
