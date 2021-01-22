package is

import (
	"errors"
	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/code"
	"math"
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
	rangeRule struct {
		lower            interface{}
		isExcludingLower bool
		upper            interface{}
		isExcludingUpper bool
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
	// GreaterThanOrEqualTo is equiv to Min
	GreaterThanOrEqualTo = Min
	// LessThanOrEqualTo is equiv to Max
	LessThanOrEqualTo = Max
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
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.NotString, value))
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

func Min(min interface{}) valis.Rule {
	return Range(min, nil)
}

func Max(max interface{}) valis.Rule {
	return Range(nil, max)
}

func GreaterThan(num interface{}) valis.Rule {
	val := reflect.ValueOf(num)
	if !isNumeric(val) {
		panic("num must be a numeric value")
	}
	return &rangeRule{lower: num, isExcludingLower: true}
}

func LessThan(num interface{}) valis.Rule {
	val := reflect.ValueOf(num)
	if !isNumeric(val) {
		panic("num must be a numeric value")
	}
	return &rangeRule{upper: num, isExcludingUpper: true}
}

func Range(min interface{}, max interface{}) valis.Rule {
	for _, val := range []reflect.Value{reflect.ValueOf(min), reflect.ValueOf(max)} {
		if val.Kind() == reflect.Invalid {
			continue
		}
		if !isNumeric(val) {
			panic("arguments must be numeric values")
		}
	}
	return &rangeRule{lower: min, upper: max}
}

func (rule *rangeRule) Validate(validator *valis.Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var lowerOpt, upperOpt = henge.WithRoundingFunc(math.Ceil), henge.WithRoundingFunc(math.Floor)
	if rule.isExcludingLower {
		lowerOpt = henge.WithRoundingFunc(math.Floor)
	}
	if rule.isExcludingUpper {
		upperOpt = henge.WithRoundingFunc(math.Ceil)
	}

	addInvalidLower := func(lower interface{}) {
		if rule.isExcludingLower {
			validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.GreaterThan, value, rule.lower))
		} else {
			validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.GreaterThanOrEqual, value, rule.lower))
		}
	}
	addInvalidUpper := func(upper interface{}) {
		if rule.isExcludingUpper {
			validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.LessThan, value, rule.upper))
		} else {
			validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.LessThanOrEqual, value, rule.upper))
		}
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var lower int64 = math.MinInt64
		henge.New(rule.lower, lowerOpt).Convert(&lower)
		var upper int64 = math.MaxInt64
		henge.New(rule.upper, upperOpt).Convert(&upper)

		i := val.Int()
		if (rule.isExcludingLower && i <= lower) || (!rule.isExcludingLower && i < lower) {
			addInvalidLower(lower)
		}
		if (rule.isExcludingUpper && i >= upper) || (!rule.isExcludingUpper && i > upper) {
			addInvalidUpper(upper)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var lower uint64
		henge.New(rule.lower, lowerOpt).Convert(&lower)
		var upper uint64 = math.MaxUint64
		henge.New(rule.upper, upperOpt).Convert(&upper)

		u := val.Uint()
		if (rule.isExcludingLower && u <= lower) || (!rule.isExcludingLower && u < lower) {
			addInvalidLower(lower)
		}
		if (rule.isExcludingUpper && u >= upper) || (!rule.isExcludingUpper && u > upper) {
			addInvalidUpper(upper)
		}
	case reflect.Float32, reflect.Float64:
		var lower float64 = -math.MaxFloat64
		henge.New(rule.lower).Convert(&lower)
		var upper float64 = math.MaxFloat64
		henge.New(rule.upper).Convert(&upper)

		f := val.Float()
		if (rule.isExcludingLower && f <= lower) || (!rule.isExcludingLower && f < lower) {
			addInvalidLower(lower)
		}
		if (rule.isExcludingUpper && f >= upper) || (!rule.isExcludingUpper && f > upper) {
			addInvalidUpper(upper)
		}
	default:
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.NotNumeric, value))
	}
}

func isNumeric(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
