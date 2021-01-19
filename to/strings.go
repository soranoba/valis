package to

import (
	"errors"
	"reflect"
	"strings"

	"github.com/soranoba/valis"
)

func Split(sep string, rules ...valis.Rule) valis.Rule {
	return valis.To(func(value interface{}) (interface{}, error) {
		val := reflect.ValueOf(value)
		for val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() != reflect.String || !val.CanInterface() {
			return nil, errors.New("cannot be split")
		}
		return strings.Split(val.Interface().(string), sep), nil
	}, rules...)
}
