package is

import (
	"errors"
	"fmt"
	"github.com/soranoba/valis"
	"reflect"
	"unicode/utf8"
)

type (
	requiredRule     struct{}
	zeroRule         struct{}
	nilOrNonZeroRule struct{}
	anyRule          struct{}
	inclusionRule    struct {
		values []interface{}
	}
	lengthRule struct {
		min int
		max int
	}
)

const (
	RequiredCode     = "required"
	ZeroCode         = "zero"
	NilOrNonZeroCode = "nil_or_non_zero"
	InclusionCode    = "inclusion"
	TooShortLength   = "too_short_length"
	TooLongLength    = "too_long_length"
)

var (
	// Required is a rule to verify non-zero value.
	// See reflect.IsZero
	Required valis.Rule = &requiredRule{}
	// Zero is a rule to verify zero value.
	Zero valis.Rule = &zeroRule{}
	// NilOrNonZero is a rule to verify nil or non-zero value.
	NilOrNonZero valis.Rule = &nilOrNonZeroRule{}
	// Any is a rule indicating that any value is acceptable.
	Any valis.Rule = &anyRule{}
)

var (
	ErrCannotBeBlank = errors.New("cannot be blank")
)

func (rule *requiredRule) Validate(validator *valis.Validator, value interface{}) {
	val := reflect.ValueOf(value)
	if !val.IsValid() || val.IsZero() {
		validator.ErrorCollector().Add(validator.Location(), &valis.ErrorDetail{
			RequiredCode,
			rule,
			value,
			nil,
			ErrCannotBeBlank,
		})
	}
}

func (rule *zeroRule) Validate(validator *valis.Validator, value interface{}) {
	val := reflect.ValueOf(value)
	if val.IsValid() && !val.IsZero() {
		validator.ErrorCollector().Add(validator.Location(), &valis.ErrorDetail{
			ZeroCode,
			rule,
			value,
			nil,
			errors.New("must be nil or zero"),
		})
	}
}

func (rule *nilOrNonZeroRule) Validate(validator *valis.Validator, value interface{}) {
	val := reflect.ValueOf(value)
	isValid := false

	switch val.Kind() {
	case reflect.Ptr:
		isValid = val.IsNil() || !val.Elem().IsZero()
	case reflect.Chan, reflect.Func, reflect.Map, reflect.UnsafePointer, reflect.Slice, reflect.Interface:
		isValid = val.IsNil() || !val.IsZero()
	case reflect.Invalid:
		isValid = true // treat as nil
	default:
		isValid = !val.IsZero()
	}

	if !isValid {
		validator.ErrorCollector().Add(validator.Location(), &valis.ErrorDetail{
			NilOrNonZeroCode,
			rule,
			value,
			nil,
			errors.New("must be nil or non-zero"),
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

func LengthRange(min int, max int) valis.Rule {
	return &lengthRule{min: min, max: max}
}

func (rule *lengthRule) Validate(validator *valis.Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var length int

	switch val.Kind() {
	case reflect.String:
		str := val.Interface().(string)
		length = utf8.RuneCountInString(str)
	case reflect.Slice, reflect.Array, reflect.Map:
		length = val.Len()
	default:
		validator.ErrorCollector().Add(validator.Location(), &valis.ErrorDetail{
			valis.InvalidTypeCode,
			rule,
			value,
			nil,
			fmt.Errorf("must be string"),
		})
		return
	}

	if length < rule.min {
		validator.ErrorCollector().Add(validator.Location(), &valis.ErrorDetail{
			TooShortLength,
			rule,
			value,
			nil,
			fmt.Errorf("is too short length (min: %d)", rule.min),
		})
	}
	if length > rule.max {
		validator.ErrorCollector().Add(validator.Location(), &valis.ErrorDetail{
			TooLongLength,
			rule,
			value,
			nil,
			fmt.Errorf("is too long length (max: %d)", rule.max),
		})
	}
}
