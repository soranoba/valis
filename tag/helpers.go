package valistag

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func SplitAndParseTagValues(s string, sep string, outs ...interface{}) (count int, err error) {
	for _, elem := range strings.Split(s, sep) {
		if elem == "" {
			continue
		}
		if len(outs) <= count {
			return count, errors.New("insufficient number of tag parameters")
		}
		val := reflect.ValueOf(outs[count])
		if val.Kind() != reflect.Ptr || !val.Elem().CanSet() {
			panic("can not to assign result to outs. please specify pointers")
		}
		val = val.Elem()
		switch val.Kind() {
		case reflect.Int:
			i, err := strconv.ParseInt(elem, 10, 64)
			if err != nil {
				return count, err
			}
			val.Set(reflect.ValueOf(int(i)))
		case reflect.Uint:
			u, err := strconv.ParseUint(elem, 10, 64)
			if err != nil {
				return count, err
			}
			val.Set(reflect.ValueOf(uint(u)))
		case reflect.Float64:
			f, err := strconv.ParseFloat(elem, 64)
			if err != nil {
				return count, err
			}
			val.Set(reflect.ValueOf(f))
		case reflect.Bool:
			b, err := strconv.ParseBool(elem)
			if err != nil {
				return count, err
			}
			val.Set(reflect.ValueOf(b))
		case reflect.String:
			val.Set(reflect.ValueOf(elem))
		default:
			panic(fmt.Sprintf("%s is unsupported type", val.Elem().Type()))
		}
	}
	return count, nil
}
