package when

import (
	"github.com/soranoba/valis"
	"reflect"
)

func IsStruct(rules ...valis.Rule) valis.Rule {
	cond := func(value interface{}) bool {
		val := reflect.ValueOf(value)
		for val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		return val.Kind() == reflect.Struct
	}
	return valis.When(cond, rules...)
}
