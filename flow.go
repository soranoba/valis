package valis

import (
	"github.com/soranoba/valis/code"
	"reflect"
)

type (
	// WhenRule is a rule that verifies the value meets the rules when conditions return true.
	WhenRule struct {
		condAndRules []*condAndRule
	}
	// WhenContext is an argument of conditions used by WhenRule.
	WhenContext struct {
		value interface{}
		loc   *Location
	}
)

type (
	andRule struct {
		rules []Rule
	}
	orRule struct {
		rules []Rule
	}
	condAndRule struct {
		cond  func(ctx *WhenContext) bool
		rules []Rule
	}
	eachRule struct {
		rules []Rule
	}
	eachFieldsRule struct {
		rules []Rule
	}
)

// And returns a new rule that verifies the value meets the rules and all common rules.
// Should only use it in your own rules, because to avoid validating common rules multiple times.
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
func If(cond func(ctx *WhenContext) bool, rules ...Rule) *WhenRule {
	return When(cond, rules...)
}

// When returns a new Rule verify the value meets the rules when cond returns true.
func When(cond func(ctx *WhenContext) bool, rules ...Rule) *WhenRule {
	return &WhenRule{condAndRules: []*condAndRule{{cond: cond, rules: rules}}}
}

// ElseWhen set the WhenRule that verified when all before conditions, and returns self.
func (r *WhenRule) ElseWhen(rule *WhenRule) *WhenRule {
	r.condAndRules = append(r.condAndRules, rule.condAndRules...)
	return r
}

// ElseIf set some Rule verified when all before conditions return false and cond returns true. And it returns self.
func (r *WhenRule) ElseIf(cond func(ctx *WhenContext) bool, rules ...Rule) *WhenRule {
	r.condAndRules = append(r.condAndRules, &condAndRule{cond: cond, rules: rules})
	return r
}

// Else set the rules verified when all conditions return false, and returns self.
func (r WhenRule) Else(rules ...Rule) Rule {
	return r.ElseIf(func(ctx *WhenContext) bool { return true }, rules...)
}

// See Rule.Validate
func (r WhenRule) Validate(validator *Validator, value interface{}) {
	ctx := &WhenContext{loc: validator.loc, value: value}
	for _, condAndRule := range r.condAndRules {
		if condAndRule.cond(ctx) {
			for _, rule := range condAndRule.rules {
				rule.Validate(validator, value)
			}
			return
		}
	}
}

// Location returns a current location.
func (ctx *WhenContext) Location() *Location {
	return ctx.loc
}

// Value returns the validating value.
func (ctx *WhenContext) Value() interface{} {
	return ctx.value
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
