package valis

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	KeyNotFoundCode = "not_found"
)

type (
	eachKeyRule struct {
		rules []Rule
	}
	eachValueRule struct {
		rules []Rule
	}
	keyRule struct {
		key   interface{}
		rules []Rule
	}
)

func Key(key interface{}, rules ...Rule) Rule {
	if !reflect.ValueOf(key).IsValid() {
		panic("key is an invalid value")
	}
	return &keyRule{key: key, rules: rules}
}

func (rule *keyRule) Validate(validator *Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Map {
		validator.ErrorCollector().Add(validator.Location(), &ErrorDetail{
			InvalidTypeCode,
			rule,
			value,
			nil,
			errors.New("must be a map"),
		})
		return
	}

	keyVal := reflect.ValueOf(rule.key)
	if !keyVal.Type().AssignableTo(val.Type().Key()) {
		validator.ErrorCollector().Add(validator.Location(), &ErrorDetail{
			InvalidTypeCode,
			rule,
			value,
			nil,
			fmt.Errorf("cannot assignable %s to key", keyVal.Type().String()),
		})
		return
	}

	mapValue := val.MapIndex(keyVal)
	if mapValue.IsValid() && mapValue.CanInterface() {
		validator.DiveMapKey(rule.key, func(v *Validator) {
			And(rule.rules...).Validate(v, mapValue.Interface())
		})
	} else {
		validator.ErrorCollector().Add(validator.Location(), &ErrorDetail{
			KeyNotFoundCode,
			rule,
			value,
			nil,
			fmt.Errorf("does not have the key (%v)", rule.key),
		})
	}
}

func EachKeys(rules ...Rule) Rule {
	return &eachKeyRule{rules: rules}
}

func (rule *eachKeyRule) Validate(validator *Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Map {
		validator.ErrorCollector().Add(validator.Location(), &ErrorDetail{
			InvalidTypeCode,
			rule,
			value,
			nil,
			errors.New("must be a map"),
		})
		return
	}

	for _, keyVal := range val.MapKeys() {
		k := keyVal.Interface()
		validator.DiveMapKey(k, func(v *Validator) {
			And(rule.rules...).Validate(v, k)
		})
	}
}

func EachValues(rules ...Rule) Rule {
	return &eachValueRule{rules: rules}
}

func (rule *eachValueRule) Validate(validator *Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Map {
		validator.ErrorCollector().Add(validator.Location(), &ErrorDetail{
			InvalidTypeCode,
			rule,
			value,
			nil,
			errors.New("must be a map"),
		})
		return
	}

	iter := val.MapRange()
	for iter.Next() {
		validator.DiveMapValue(iter.Key().Interface(), func(v *Validator) {
			And(rule.rules...).Validate(v, iter.Value().Interface())
		})
	}
}
