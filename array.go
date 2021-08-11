package valis

import (
	"github.com/soranoba/valis/code"
	"reflect"
)

type (
	indexRule struct {
		idx      int
		rules    []Rule
		optional bool
	}
)

// Index returns a new rule that verifies the value at the index meets the rules and all common rules.
func Index(idx int, rules ...Rule) Rule {
	return &indexRule{idx: idx, rules: rules, optional: false}
}

// IndexIfExist returns a new rule.
// When the value at the index does not exist, the rule only checks the type of value.
// Otherwise, the rule same as the rule returned by the Index method.
func IndexIfExist(idx int, rules ...Rule) Rule {
	return &indexRule{idx: idx, rules: rules, optional: true}
}

func (rule *indexRule) Validate(validator *Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		validator.ErrorCollector().Add(validator.Location(), NewError(code.NotArray, value))
		return
	}

	if rule.idx >= val.Len() {
		if !rule.optional {
			validator.ErrorCollector().Add(validator.Location(), NewError(code.OutOfRange, value, rule.idx, val.Len()))
		}
		return
	}

	idxValue := val.Index(rule.idx)
	if idxValue.IsValid() && idxValue.CanInterface() {
		validator.DiveIndex(rule.idx, func(v *Validator) {
			And(rule.rules...).Validate(v, idxValue.Interface())
		})
	} else {
		validator.ErrorCollector().Add(validator.Location(), NewError(code.Invalid, value))
	}
}
