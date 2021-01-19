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
		key interface{}
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
		And(rule.rules...).Validate(validator.WithMapKey(rule.key), mapValue.Interface())
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
		loc := validator.Location().MapKeyLocation(keyVal.Interface())
		newValidator := validator.Clone(&CloneOpts{Location: loc, InheritErrorCollector: true})
		And(rule.rules...).Validate(newValidator, keyVal.Interface())
	}
}

func EachValues(rules ...Rule) Rule {
	return &eachValueRule{rules: rules}
}

func (rule *eachValueRule) Validate(validator *Validator, value interface{}){
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
		And(rule.rules...).Validate(validator.WithMapKey(iter.Key().Interface()), iter.Value().Interface())
	}
}
