package valis

import (
	"github.com/soranoba/valis/code"
	"reflect"
)

type (
	// WhenRule is a rule that has any condition.
	WhenRule = *whenRule
	// ConditionFunc is a condition function.
	ConditionFunc func(interface{}) bool
)

type (
	andRule struct {
		rules []Rule
	}
	orRule struct {
		rules []Rule
	}
	condAndRule struct {
		cond  ConditionFunc
		rules []Rule
	}
	whenRule struct {
		condAndRules []*condAndRule
	}
	eachRule struct {
		rules []Rule
	}
	eachFieldsRule struct {
		rules []Rule
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
	return &orRule{rules: rules}
}

func (r *orRule) Validate(validator *Validator, value interface{}) {
	for _, rule := range r.rules {
		newValidator := validator.Clone(&CloneOpts{InheritLocation: true})
		rule.Validate(newValidator, value)
		if !newValidator.ErrorCollector().HasError() {
			return
		}
	}
	validator.ErrorCollector().Add(validator.Location(), NewError(code.Invalid, value))
}

// If is equiv to When
func If(cond ConditionFunc, rules ...Rule) WhenRule {
	return When(cond, rules...)
}

// When returns a WhenRule that verified the value meets the rules and all common rules when cond returns true.
func When(cond ConditionFunc, rules ...Rule) WhenRule {
	return &whenRule{condAndRules: []*condAndRule{{cond: cond, rules: rules}}}
}

func (r *whenRule) ElseWhen(rule WhenRule) WhenRule {
	r.condAndRules = append(r.condAndRules, rule.condAndRules...)
	return r
}

// ElseIf set the rules that verified when all before conditions return false and cond returns true.
func (r *whenRule) ElseIf(cond ConditionFunc, rules ...Rule) WhenRule {
	r.condAndRules = append(r.condAndRules, &condAndRule{cond: cond, rules: rules})
	return r
}

// Else set the rules that verified when all conditions return false.
func (r *whenRule) Else(rules ...Rule) Rule {
	return r.ElseIf(func(interface{}) bool { return true }, rules...)
}

func (r *whenRule) Validate(validator *Validator, value interface{}) {
	for _, condAndRule := range r.condAndRules {
		if condAndRule.cond(value) {
			for _, rule := range condAndRule.rules {
				rule.Validate(validator, value)
			}
			return
		}
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

func EachFields(rules ...Rule) Rule {
	return &eachFieldsRule{rules: rules}
}

func (rule *eachFieldsRule) Validate(validator *Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			fieldVal := val.Field(i)
			if !fieldVal.IsValid() {
				return
			}
			field := val.Type().Field(i)
			validator.DiveField(&field, func(v *Validator) {
				And(rule.rules...).Validate(v, fieldVal.Interface())
			})
		}
	default:
		validator.ErrorCollector().Add(validator.Location(), NewError(code.NotStruct, value))
	}
}
