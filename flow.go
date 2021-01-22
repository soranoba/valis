package valis

import (
	"github.com/soranoba/valis/code"
	"reflect"
)

type (
	ConditionFunc func(interface{}) bool
)

type (
	andRule struct {
		rules []Rule
	}
	orRule struct {
		Rules []Rule
	}
	eachRule struct {
		rules []Rule
	}
	whenRule struct {
		cond          ConditionFunc
		whenRules     []Rule
		didCalledElse bool
		elseRules     []Rule
	}
)

// And returns a new rule that verifies the value meets the rules and all common rules.
func And(rules ...Rule) Rule {
	return &andRule{rules: rules}
}

func (r *andRule) Validate(validator *Validator, value interface{}) {
	if len(validator.commonRules) > 0 {
		rules := append(validator.commonRules, r.rules...)
		for _, rule := range rules {
			rule.Validate(validator, value)
		}
	} else {
		for _, rule := range r.rules {
			rule.Validate(validator, value)
		}
	}
}

// Or returns a new rule that verifies the value meets the rules at least one.
func Or(rules ...Rule) Rule {
	return &orRule{Rules: rules}
}

func (r *orRule) Validate(validator *Validator, value interface{}) {
	for _, rule := range r.Rules {
		newValidator := validator.Clone(&CloneOpts{InheritLocation: true})
		rule.Validate(newValidator, value)
		if !newValidator.ErrorCollector().HasError() {
			return
		}
	}
	validator.ErrorCollector().Add(validator.Location(), NewError(code.Invalid, value))
}

// When returns a new rule that verifies the value meets the rules and all common rules when cond returns true.
func When(cond ConditionFunc, rules ...Rule) *whenRule {
	return &whenRule{cond: cond, whenRules: rules}
}

// Else set the rules that verified the value when cond returns false.
func (r *whenRule) Else(rules ...Rule) Rule {
	r.elseRules = rules
	r.didCalledElse = true
	return r
}

func (r *whenRule) Validate(validator *Validator, value interface{}) {
	if r.cond(value) {
		And(r.whenRules...).Validate(validator, value)
	} else if r.didCalledElse {
		And(r.elseRules...).Validate(validator, value)
	}
}

// Each returns a new rule that verifies all elements of the array or slice meet the rules and all common rules.
func Each(rules ...Rule) Rule {
	return &eachRule{rules: rules}
}

func (rule *eachRule) Validate(validator *Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			indexValue := val.Index(i).Interface()
			validator.DiveIndex(i, func(v *Validator) {
				And(rule.rules...).Validate(v, indexValue)
			})
		}
	default:
		validator.ErrorCollector().Add(validator.Location(), NewError(code.NotArray, value))
	}
}
