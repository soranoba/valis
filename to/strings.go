package to

import (
	"errors"
	"github.com/soranoba/valis"
	"reflect"
	"strings"
)

func Split(sep string) valis.CombinationRule {
	return valis.To(func(value interface{}) (interface{}, error) {
		val := reflect.ValueOf(value)
		for val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() != reflect.String || !val.CanInterface() {
			return nil, errors.New("cannot be split")
		}
		return strings.Split(val.Interface().(string), sep), nil
	})
}
