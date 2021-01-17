package valis

import (
	"errors"
	"reflect"
)

type (
	ConditionFunc func(value interface{}) bool

	OrRule struct {
		Rules []Rule
	}
	EachRule struct {
		Rules []Rule
	}
)

type (
	whenRule struct {
		cond  ConditionFunc
		rules []Rule
	}
)

const (
	OrCode          = "or"
	InvalidTypeCode = "invalid_type"
)

func Or(rules ...Rule) Rule {
	return &OrRule{Rules: rules}
}

func (r *OrRule) Validate(validator *Validator, value interface{}) {
	for _, rule := range r.Rules {
		newValidator := validator.Clone(&CloneOpts{KeepLocation: true})
		rule.Validate(newValidator, value)
		if !newValidator.ErrorCollector().HasError() {
			return
		}
	}
	validator.ErrorCollector().Add(validator.Location(), &ErrorDetail{
		OrCode,
		r,
		value,
		nil,
		errors.New("cannot meet either rule"),
	})
}

func When(cond ConditionFunc, rules ...Rule) Rule {
	return &whenRule{cond: cond, rules: rules}
}

func (r *whenRule) Validate(validator *Validator, value interface{}) {
	if r.cond(value) {
		for _, rule := range r.rules {
			rule.Validate(validator, value)
		}
	}
}

func Each(rules ...Rule) Rule {
	return &EachRule{Rules: rules}
}

func (rule *EachRule) Validate(validator *Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			for _, r := range rule.Rules {
				r.Validate(validator.WithLocation(validator.Location().IndexLocation(i)), v.Interface())
			}
		}
	default:
		validator.ErrorCollector().Add(validator.Location(), &ErrorDetail{
			InvalidTypeCode,
			rule,
			value,
			nil,
			errors.New("must be array or slice"),
		})
	}
}
