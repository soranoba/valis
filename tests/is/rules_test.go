package is_test

import (
	"fmt"
	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/code"
	"github.com/soranoba/valis/is"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

type SimpleTestCase struct {
	Value   interface{}
	IsValid bool
}

func TestRequired(t *testing.T) {
	assert := assert.New(t)
	testCases := []SimpleTestCase{
		{"", false},
		{0, false},
		{nil, false},
		{"a", true},
		{1, true},
		{[...]string{}, false},
		{[...]string{""}, false},
		{[...]string{"abc"}, true},
		{map[string]string(nil), false},
		{map[string]string{}, true},
		{map[string]string{"": ""}, true},
		{struct{}{}, false},
		{struct{ Name string }{"name"}, true},
	}

	for _, testCase := range testCases {
		err := valis.Validate(testCase.Value, is.Required)
		msg := fmt.Sprintf("%#v", testCase)
		if testCase.IsValid {
			assert.NoError(err, msg)
		} else {
			if assert.Error(err, msg) {
				details := err.(*valis.ValidationError).Details()
				assert.Len(details, 1)
				assert.Equal(code.Required, details[0].Code())
			}
		}
	}
}

func TestZero(t *testing.T) {
	assert := assert.New(t)
	testCases := []SimpleTestCase{
		{"", true},
		{0, true},
		{nil, true},
		{"a", false},
		{1, false},
		{[]string(nil), true},
		{[]string{}, false},
		{[...]string{}, true},
		{[...]string{""}, true},
		{[...]string{"abc"}, false},
		{map[string]string(nil), true},
		{map[string]string{}, false},
		{map[string]string{"": ""}, false},
		{struct{}{}, true},
		{struct{ Name string }{"name"}, false},
	}

	for _, testCase := range testCases {
		err := valis.Validate(testCase.Value, is.Zero)
		msg := fmt.Sprintf("%#v", testCase)
		if testCase.IsValid {
			assert.NoError(err, msg)
		} else {
			if assert.Error(err, msg) {
				details := err.(*valis.ValidationError).Details()
				assert.Len(details, 1)
				assert.Equal(code.ZeroOnly, details[0].Code())
			}
		}
	}
}

func TestNilOrNonZero(t *testing.T) {
	assert := assert.New(t)
	testCases := []SimpleTestCase{
		{"", false},
		{0, false},
		{nil, true},
		{"a", true},
		{1, true},
		{[...]string{}, false},
		{[...]string{""}, false},
		{[...]string{"abc"}, true},
		{map[string]string(nil), true},
		{map[string]string{}, true},
		{map[string]string{"": ""}, true},
		{struct{}{}, false},
		{struct{ Name string }{"name"}, true},
	}

	for _, testCase := range testCases {
		err := valis.Validate(testCase.Value, is.NilOrNonZero)
		msg := fmt.Sprintf("%#v", testCase)
		if testCase.IsValid {
			assert.NoError(err, msg)
		} else {
			if assert.Error(err, msg) {
				details := err.(*valis.ValidationError).Details()
				assert.Len(details, 1)
				assert.Equal(code.NilOrNonZero, details[0].Code())
			}
		}
	}
}

func TestAny(t *testing.T) {
	assert := assert.New(t)
	testCases := []SimpleTestCase{
		{"", true},
		{0, true},
		{nil, true},
		{"a", true},
		{1, true},
		{[...]string{}, true},
		{[...]string{""}, true},
		{[...]string{"abc"}, true},
		{map[string]string(nil), true},
		{map[string]string{}, true},
		{map[string]string{"": ""}, true},
		{struct{}{}, true},
		{struct{ Name string }{"name"}, true},
	}

	for _, testCase := range testCases {
		err := valis.Validate(testCase.Value, is.Any)
		msg := fmt.Sprintf("%#v", testCase)
		if testCase.IsValid {
			assert.NoError(err, msg)
		} else {
			assert.Error(err, msg)
		}
	}
}

func TestNever(t *testing.T) {
	assert := assert.New(t)
	testCases := []SimpleTestCase{
		{"", false},
		{0, false},
		{nil, false},
		{"a", false},
		{1, false},
		{[...]string{}, false},
		{[...]string{""}, false},
		{[...]string{"abc"}, false},
		{map[string]string(nil), false},
		{map[string]string{}, false},
		{map[string]string{"": ""}, false},
		{struct{}{}, false},
		{struct{ Name string }{"name"}, false},
	}

	for _, testCase := range testCases {
		err := valis.Validate(testCase.Value, is.Never)
		msg := fmt.Sprintf("%#v", testCase)
		if testCase.IsValid {
			assert.NoError(err, msg)
		} else {
			if assert.Error(err, msg) {
				details := err.(*valis.ValidationError).Details()
				assert.Len(details, 1)
				assert.Equal(code.Invalid, details[0].Code())
			}
		}
	}
}

func TestIn(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(valis.Validate("b", is.In("a", "b", "c")))
	assert.EqualError(valis.Validate("d", is.In("a", "b", "c")), "(inclusion) is not included in [a b c]")
	assert.NoError(valis.Validate(2, is.In(1, 2, 3)))
	assert.EqualError(valis.Validate(5, is.In(1, 2, 3)), "(inclusion) is not included in [1 2 3]")

	// NOTE: Does not match if the types are different
	assert.EqualError(valis.Validate(int64(1), is.In(1, 2, 3)), "(inclusion) is not included in [1 2 3]")

	// NOTE: For pointers, it compares Elem values
	i := 2
	assert.NoError(valis.Validate(&i, is.In(1, 2, 3)))
	i = 5
	assert.Error(valis.Validate(&i, is.In(1, 2, 3)))

	i = 2
	x, y, z := 1, 2, 3
	assert.NoError(valis.Validate(&i, is.In(&x, &y, &z)))
	i = 5
	assert.Error(valis.Validate(&i, is.In(&x, &y, &z)))
}

func TestLengthBetween(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		valis.Validate("🍣🍺🍣", is.LengthBetween(4, math.MaxInt64)),
		"(too_short_length) is too short length (minimum is 4 characters)",
	)
	assert.EqualError(
		valis.Validate("abc", is.LengthBetween(0, 2)),
		"(too_long_length) is too long length (maximum is 2 characters)",
	)
	assert.EqualError(
		valis.Validate("abc", is.LengthBetween(4, 3)),
		"(too_short_length) is too short length (minimum is 4 characters)",
	)
	assert.EqualError(
		valis.Validate(0, is.LengthBetween(0, 10)),
		"(not_string) must be any string",
	)
	assert.NoError(valis.Validate("🍣🍺🍣", is.LengthBetween(3, 3)))
}

func TestLenBetween(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		valis.Validate("abc", is.LenBetween(4, math.MaxInt64)),
		"(too_short_len) is too few elements (minimum is 4 elements)",
	)
	assert.NoError(
		valis.Validate("🍣🍺🍣", is.LenBetween(4, math.MaxInt64)),
	)
	assert.EqualError(
		valis.Validate("abc", is.LenBetween(0, 2)),
		"(too_long_len) is too many elements (maximum is 2 elements)",
	)
	assert.EqualError(
		valis.Validate("abc", is.LenBetween(4, 3)),
		"(too_short_len) is too few elements (minimum is 4 elements)",
	)
	assert.EqualError(
		valis.Validate(0, is.LenBetween(0, 10)),
		"(not_iterable) not_iterable",
	)
	assert.NoError(valis.Validate("abc", is.LenBetween(0, 10)))
}

func TestMin(t *testing.T) {
	assert := assert.New(t)

	f := func(v interface{}, err error, isValid bool, errMsg string) {
		msg := fmt.Sprintf("%v", v)
		if isValid {
			assert.NoError(err, msg)
		} else {
			if assert.EqualError(err, errMsg, msg) {
				details := err.(*valis.ValidationError).Details()
				assert.Len(details, 1, msg)
				assert.Equal(code.GreaterThanOrEqual, details[0].Code(), msg)
			}
		}
	}

	for _, i := range []int{-2, -1, 0, 1, 2} {
		f(i, valis.Validate(i, is.Min(-1.5)), float64(i) >= -1.5, "(gte) must be greater than or equal to -1.5")
		f(i, valis.Validate(i, is.Min(1.5)), float64(i) >= 1.5, "(gte) must be greater than or equal to 1.5")
		f(i, valis.Validate(i, is.Min(-1)), i >= -1, "(gte) must be greater than or equal to -1")
		f(i, valis.Validate(i, is.Min(1)), i >= 1, "(gte) must be greater than or equal to 1")
	}
	for _, v := range []float64{-1.5, -0.5, 0.5, 1.5} {
		f(v, valis.Validate(v, is.Min(-0.5)), v >= -0.5, "(gte) must be greater than or equal to -0.5")
		f(v, valis.Validate(v, is.Min(0.5)), v >= 0.5, "(gte) must be greater than or equal to 0.5")
		f(v, valis.Validate(v, is.Min(-1)), v >= -1, "(gte) must be greater than or equal to -1")
		f(v, valis.Validate(v, is.Min(1)), v >= 1, "(gte) must be greater than or equal to 1")
	}
}

func TestMax(t *testing.T) {
	assert := assert.New(t)

	f := func(v interface{}, err error, isValid bool, errMsg string) {
		msg := fmt.Sprintf("%v", v)
		if isValid {
			assert.NoError(err, msg)
		} else {
			if assert.EqualError(err, errMsg, msg) {
				details := err.(*valis.ValidationError).Details()
				assert.Len(details, 1, msg)
				assert.Equal(code.LessThanOrEqual, details[0].Code(), msg)
			}
		}
	}

	for _, i := range []int{-2, -1, 0, 1, 2} {
		f(i, valis.Validate(i, is.Max(-1.5)), float64(i) <= -1.5, "(lte) must be less than or equal to -1.5")
		f(i, valis.Validate(i, is.Max(1.5)), float64(i) <= 1.5, "(lte) must be less than or equal to 1.5")
		f(i, valis.Validate(i, is.Max(-1)), i <= -1, "(lte) must be less than or equal to -1")
		f(i, valis.Validate(i, is.Max(1)), i <= 1, "(lte) must be less than or equal to 1")
	}
	for _, v := range []float64{-1.5, -0.5, 0.5, 1.5} {
		f(v, valis.Validate(v, is.Max(-0.5)), v <= -0.5, "(lte) must be less than or equal to -0.5")
		f(v, valis.Validate(v, is.Max(0.5)), v <= 0.5, "(lte) must be less than or equal to 0.5")
		f(v, valis.Validate(v, is.Max(-1)), v <= -1, "(lte) must be less than or equal to -1")
		f(v, valis.Validate(v, is.Max(1)), v <= 1, "(lte) must be less than or equal to 1")
	}
}

func TestGreaterThan(t *testing.T) {
	assert := assert.New(t)

	f := func(v interface{}, err error, isValid bool, errMsg string) {
		msg := fmt.Sprintf("%v", v)
		if isValid {
			assert.NoError(err, msg)
		} else {
			if assert.EqualError(err, errMsg, msg) {
				details := err.(*valis.ValidationError).Details()
				assert.Len(details, 1, msg)
				assert.Equal(code.GreaterThan, details[0].Code(), msg)
			}
		}
	}

	for _, i := range []int{-2, -1, 0, 1, 2} {
		f(i, valis.Validate(i, is.GreaterThan(-1.5)), float64(i) > -1.5, "(gt) must be greater than -1.5")
		f(i, valis.Validate(i, is.GreaterThan(1.5)), float64(i) > 1.5, "(gt) must be greater than 1.5")
		f(i, valis.Validate(i, is.GreaterThan(-1)), i > -1, "(gt) must be greater than -1")
		f(i, valis.Validate(i, is.GreaterThan(1)), i > 1, "(gt) must be greater than 1")
	}
	for _, v := range []float64{-1.5, -0.5, 0.5, 1.5} {
		f(v, valis.Validate(v, is.GreaterThan(-0.5)), v > -0.5, "(gt) must be greater than -0.5")
		f(v, valis.Validate(v, is.GreaterThan(0.5)), v > 0.5, "(gt) must be greater than 0.5")
		f(v, valis.Validate(v, is.GreaterThan(-1)), v > -1, "(gt) must be greater than -1")
		f(v, valis.Validate(v, is.GreaterThan(1)), v > 1, "(gt) must be greater than 1")
	}
}

func TestLessThan(t *testing.T) {
	assert := assert.New(t)

	f := func(v interface{}, err error, isValid bool, errMsg string) {
		msg := fmt.Sprintf("%v", v)
		if isValid {
			assert.NoError(err, msg)
		} else {
			if assert.EqualError(err, errMsg, msg) {
				details := err.(*valis.ValidationError).Details()
				assert.Len(details, 1, msg)
				assert.Equal(code.LessThan, details[0].Code(), msg)
			}
		}
	}

	for _, i := range []int{-2, -1, 0, 1, 2} {
		f(i, valis.Validate(i, is.LessThan(-1.5)), float64(i) < -1.5, "(lt) must be less than -1.5")
		f(i, valis.Validate(i, is.LessThan(1.5)), float64(i) < 1.5, "(lt) must be less than 1.5")
		f(i, valis.Validate(i, is.LessThan(-1)), i < -1, "(lt) must be less than -1")
		f(i, valis.Validate(i, is.LessThan(1)), i < 1, "(lt) must be less than 1")
	}
	for _, v := range []float64{-1.5, -0.5, 0.5, 1.5} {
		f(v, valis.Validate(v, is.LessThan(-0.5)), v < -0.5, "(lt) must be less than -0.5")
		f(v, valis.Validate(v, is.LessThan(0.5)), v < 0.5, "(lt) must be less than 0.5")
		f(v, valis.Validate(v, is.LessThan(-1)), v < -1, "(lt) must be less than -1")
		f(v, valis.Validate(v, is.LessThan(1)), v < 1, "(lt) must be less than 1")
	}
}

func TestRange(t *testing.T) {
	assert := assert.New(t)

	f := func(v interface{}, err error, isValid bool) {
		msg := fmt.Sprintf("%v", v)
		if isValid {
			assert.NoError(err, msg)
		} else {
			if assert.Error(err, msg) {
				details := err.(*valis.ValidationError).Details()
				assert.Len(details, 1, msg)
				if henge.New(v).Float().Value() < 0 {
					assert.Equal(code.GreaterThanOrEqual, details[0].Code(), msg)
				} else {
					assert.Equal(code.LessThanOrEqual, details[0].Code(), msg)
				}
			}
		}
	}

	for _, i := range []int{-2, -1, 0, 1, 2} {
		f(i, valis.Validate(i, is.Range(-1.5, 1.5)), float64(i) >= -1.5 && float64(i) <= 1.5)
		f(i, valis.Validate(i, is.Range(-1.5, 1.5)), float64(i) >= -1.5 && float64(i) <= 1.5)
		f(i, valis.Validate(i, is.Range(-1, 1)), i >= -1 && i <= 1)
		f(i, valis.Validate(i, is.Range(-1, 1)), i >= -1 && i <= 1)
	}
	for _, v := range []float64{-1.5, -0.5, 0.5, 1.5} {
		f(v, valis.Validate(v, is.Range(-0.5, 0.5)), v >= -0.5 && v <= 0.5)
		f(v, valis.Validate(v, is.Range(-0.5, 0.5)), v >= -0.5 && v <= 0.5)
		f(v, valis.Validate(v, is.Range(-1, 1)), v >= -1 && v <= 1)
		f(v, valis.Validate(v, is.Range(-1, 1)), v >= -1 && v <= 1)
	}
}

func TestMatch(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		valis.Validate(0, is.MatchString("^[0-9]+$")),
		"(not_string) must be any string",
	)
	assert.EqualError(
		valis.Validate("123456abc", is.MatchString("^[0-9]+$")),
		"(regexp) is a mismatch with the regular expression. (^[0-9]+$)",
	)
	assert.NoError(
		valis.Validate("123456789", is.MatchString("^[0-9]+$")),
	)
	assert.NoError(
		valis.Validate(henge.New("123456789").StringPtr().Value(), is.MatchString("^[0-9]+$")),
	)
}
