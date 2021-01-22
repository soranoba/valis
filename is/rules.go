package is

import (
	"errors"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/code"
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
	lenRule struct {
		min int
		max int
	}
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
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.Required, value))
	}
}

func (rule *zeroRule) Validate(validator *valis.Validator, value interface{}) {
	val := reflect.ValueOf(value)
	if val.IsValid() && !val.IsZero() {
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.ZeroOnly, value))
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
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.NilOrNonZero, value))
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
	validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.Inclusion, value, rule.values))
}

func LengthBetween(min int, max int) valis.Rule {
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
	default:
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.StringOnly, value))
		return
	}

	if length < rule.min {
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.TooShortLength, value, rule.min))
	}
	if length > rule.max {
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.TooLongLength, value, rule.max))
	}
}

func LenBetween(min int, max int) valis.Rule {
	return &lenRule{min: min, max: max}
}

func (rule *lenRule) Validate(validator *valis.Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var length int

	switch val.Kind() {
	case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
		length = val.Len()
	default:
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.NotIterable, value))
		return
	}

	if length < rule.min {
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.TooShortLen, value, rule.min))
	}
	if length > rule.max {
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.TooLongLen, value, rule.max))
	}
}
