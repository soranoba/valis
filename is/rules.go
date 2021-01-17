package is

import (
	"errors"
	"fmt"
	"github.com/soranoba/valis"
	"reflect"
)

type (
	requiredRule  struct{}
	anyRule       struct{}
	inclusionRule struct {
		values []interface{}
	}
)

const (
	RequiredCode  = "required"
	InclusionCode = "inclusion"
)

var (
	// Required is a rule to verify non-zero value.
	// See reflect.IsZero
	Required valis.Rule = &requiredRule{}
	// Any is a rule indicating that any value is acceptable.
	Any valis.Rule = &anyRule{}
)

func (rule *requiredRule) Validate(validator *valis.Validator, value interface{}) {
	val := reflect.ValueOf(value)
	if !val.IsValid() || val.IsZero() {
		validator.ErrorCollector().Add(validator.Location(), &valis.ErrorDetail{
			RequiredCode,
			rule,
			value,
			nil,
			errors.New("cannot be blank"),
		})
	}
}

func (rule *anyRule) Validate(validator *valis.Validator, value interface{}) {
	// NOP
}

// In returns a rule to verify inclusion in the values.
// For the pointer type, the Elem value is validated. Otherwise, It needs to same types.
func In(values ...interface{}) valis.Rule {
	return &inclusionRule{values: values}
}

func (rule *inclusionRule) Validate(validator *valis.Validator, value interface{}) {
	y := reflect.ValueOf(value)
	for _, val := range rule.values {
		x := reflect.ValueOf(val)
		if y.Kind() == reflect.Ptr {
			if x.Kind() == reflect.Ptr {
				if reflect.DeepEqual(x.Interface(), y.Interface()) {
					return
				}
			} else {
				y := y
				for y.Kind() == reflect.Ptr {
					y = y.Elem()
				}
				if y.CanInterface() && reflect.DeepEqual(x.Interface(), y.Interface()) {
					return
				}
			}
		} else {
			if reflect.DeepEqual(x.Interface(), y.Interface()) {
				return
			}
		}
	}
	validator.ErrorCollector().Add(validator.Location(), &valis.ErrorDetail{
		InclusionCode,
		rule,
		value,
		nil,
		fmt.Errorf("is not included in %v", rule.values),
	})
}
