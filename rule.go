package valis

import (
	"github.com/soranoba/valis/code"
	"reflect"
)

type (
	// Rule is an interface where verification contents are defined.
	Rule interface {
		Validate(validator *Validator, value interface{})
	}
	// CombinationRule is a high-order function that returns a new Rule.
	CombinationRule func(rules ...Rule) Rule
)

type (
	// Validatable will be delegated the validation by the ValidatableRule if implemented.
	Validatable interface {
		Validate() error
	}
	// ValidatableWithValidator will be delegated the validation by the ValidatableRule if implemented.
	ValidatableWithValidator interface {
		Validate(validator *Validator)
	}
)

type (
	validatableRule struct{}
)

var (
	// ValidatableRule is a rule that delegates to Validate methods when verifying.
	// See also Validatable and ValidatableWithValidator.
	ValidatableRule Rule = &validatableRule{}
)

func (rule *validatableRule) Validate(validator *Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for {
		if !(val.IsValid() && val.CanInterface()) {
			return
		}
		if v, ok := val.Interface().(ValidatableWithValidator); ok {
			v.Validate(validator)
			return
		}
		if v, ok := val.Interface().(Validatable); ok {
			if err := v.Validate(); err != nil {
				validator.ErrorCollector().Add(validator.Location(), NewError(code.Custom, value, err))
			}
			return
		}
		if val.Kind() != reflect.Ptr {
			break
		}
		val = val.Elem()
	}
}
