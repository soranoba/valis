package valis

import (
	"reflect"
)

type (
	Rule interface {
		Validate(validator *Validator, value interface{})
	}
	CombinationRule func(rules ...Rule) Rule

	Validatable interface {
		Validate() error
	}
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
	StandardRules = [...]Rule{
		ValidatableRule,
	}
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
