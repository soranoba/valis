package valis

import (
	"reflect"
)

type (
	// Rule is an interface where verification contents are defined.
	Rule interface {
		Validate(validator *Validator, value interface{})
	}
	// CombinationRule is a high-order function that returns a new Rule.
	CombinationRule func(rules ...Rule) Rule

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
	validatableRule struct {
	}
)

const (
	ValidatableCode = "validatable"
)

var (
	// StandardRules is common rules by default.
	StandardRules = [...]Rule{
		ValidatableRule,
	}
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
					validator.ErrorCollector().Add(validator.Location(), &ErrorDetail{
						ValidatableCode,
						rule,
						value,
						nil,
						err,
					})
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
