package valis

import (
	"github.com/soranoba/valis/code"
	"reflect"
)

type (
	// Rule is an interface where verification contents are defined.
	// Needs to define it in Public member variables so that you can refer to the value of the parameter using text/template package.
	Rule interface {
		Validate(validator *Validator, value interface{})
	}
	// CombinationRule is a high-order function that returns a new Rule.
	CombinationRule func(rules ...Rule) Rule
)

type (
	// If implemented, Validate will be executed by the ValidatableRule.
	Validatable interface {
		Validate() error
	}
	// If implemented, ValidatableWithValidator will be executed by the ValidatableRule.
	ValidatableWithValidator interface {
		Validate(validator *Validator)
	}
)

type (
	validatableRule struct{}
)

var (
	// ValidatableRule is a rule that executes Validate methods if the value implements Validatable or ValidatableWithValidator at verifying.
	ValidatableRule Rule = &validatableRule{}
)

func (rule *validatableRule) Validate(validator *Validator, value interface{}) {
	val := reflect.ValueOf(value)
	if !val.IsValid() {
		return
	}

	for {
		if val.CanInterface() {
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
		}
		if val.Kind() != reflect.Ptr {
			break
		}
		val = val.Elem()
	}
}
