// The package implements some valis.Rule.
package is

import (
	"math"
	"net/url"
	"reflect"
	"regexp"
	"unicode/utf8"

	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/code"
	valishelpers "github.com/soranoba/valis/helpers"
)

type (
	requiredRule     struct{}
	zeroRule         struct{}
	nilOrNonZeroRule struct{}
	anyRule          struct{}
	neverRule        struct{}
	urlRule          struct {
		schemes []string
	}
	inclusionRule struct {
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
	matchRule struct {
		re *regexp.Regexp
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
	// Never is a rule indicating that any value is not acceptable.
	Never valis.Rule = &neverRule{}
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

func (rule *neverRule) Validate(validator *valis.Validator, value interface{}) {
	validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.Invalid, value))
}

// URL returns a rule to verify URL format.
//
// When you specify schemes, it verifies the scheme inclusion in those schemes.
// When you not specified, all schemes are allowed.
func URL(schemes ...string) valis.Rule {
	return &urlRule{schemes: schemes}
}

func (rule *urlRule) Validate(validator *valis.Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.String {
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.NotString, value))
		return
	}
	if val.CanInterface() {
		txt := val.Interface().(string)
		u, err := url.Parse(txt)
		if err != nil {
			validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.InvalidURLFormat, value, err))
		} else if len(rule.schemes) == 0 {
			// valid
		} else {
			for _, scheme := range rule.schemes {
				if scheme == u.Scheme {
					return // valid
				}
			}
			validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.InvalidScheme, value, rule.schemes))
		}
	}
}

// In returns a rule to verify inclusion in the values.
//
// When the validating value is a pointer and the values are not a pointer, the Elem value of the validating value is validated.
// Otherwise, It needs to same types.
func In(values ...interface{}) valis.Rule {
	return &inclusionRule{values: values}
}

func (rule *inclusionRule) Validate(validator *valis.Validator, value interface{}) {
	y := reflect.ValueOf(value)
	for _, val := range rule.values {
		x := reflect.ValueOf(val)
		if x.Kind() != reflect.Ptr && y.Kind() == reflect.Ptr {
			y := y
			for y.Kind() == reflect.Ptr {
				y = y.Elem()
			}
			if y.CanInterface() && reflect.DeepEqual(x.Interface(), y.Interface()) {
				return
			}
		} else {
			if reflect.DeepEqual(x.Interface(), y.Interface()) {
				return
			}
		}
	}
	validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.Inclusion, value, rule.values))
}

// LengthBetween returns a rule to verify the length of the value is between min and max.
// if the verifying value is not a string, the rule considers that the value is invalid.
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

// LenBetween returns a rule to verify the len(value) is between min and max.
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

// Min returns a rule to verify that the value >= min.
func Min(min interface{}) valis.Rule {
	return Range(min, nil)
}

// Max returns a rule to verify that the value <= max.
func Max(max interface{}) valis.Rule {
	return Range(nil, max)
}

// GreaterThan returns a rule to verify that the value > num.
func GreaterThan(num interface{}) valis.Rule {
	if !valishelpers.IsNumeric(num) {
		panic("num must be a numeric value")
	}
	return &rangeRule{lower: num, isExcludingLower: true}
}

// LessThan returns a rule to verify that the value < num.
func LessThan(num interface{}) valis.Rule {
	if !valishelpers.IsNumeric(num) {
		panic("num must be a numeric value")
	}
	return &rangeRule{upper: num, isExcludingUpper: true}
}

// GreaterThanOrEqualTo is equiv to Min
func GreaterThanOrEqualTo(num interface{}) valis.Rule {
	return Min(num)
}

// LessThanOrEqualTo is equiv to Max
func LessThanOrEqualTo(num interface{}) valis.Rule {
	return Max(num)
}

// Range returns a rule to verify that the value is between min and max.
func Range(min interface{}, max interface{}) valis.Rule {
	for _, val := range []interface{}{min, max} {
		if !reflect.ValueOf(val).IsValid() {
			continue
		}
		if !valishelpers.IsNumeric(val) {
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

func MatchString(r string) valis.Rule {
	re := regexp.MustCompile(r)
	return &matchRule{re: re}
}

func Match(re *regexp.Regexp) valis.Rule {
	return &matchRule{re: re}
}

func (rule *matchRule) Validate(validator *valis.Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.String {
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.NotString, value))
		return
	}
	if !rule.re.MatchString(val.Interface().(string)) {
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.RegexpMismatch, value, rule.re.String()))
	}
}
