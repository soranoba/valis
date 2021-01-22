package is_test

import (
	"math"
	"testing"

	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/stretchr/testify/assert"
)

var (
	v *valis.Validator
)

func init() {
	v = valis.NewValidator()
	v.SetCommonRules()
}

func TestRequired(t *testing.T) {
	assert := assert.New(t)

	errString := "(required) is required"
	assert.EqualError(v.Validate("", is.Required), errString)
	assert.EqualError(v.Validate(0, is.Required), errString)
	assert.EqualError(v.Validate(struct{}{}, is.Required), errString)
	assert.EqualError(v.Validate(nil, is.Required), errString)
	assert.EqualError(v.Validate([...]string{}, is.Required), errString)
	assert.EqualError(v.Validate([...]string{""}, is.Required), errString)

	assert.NoError(v.Validate("a", is.Required))
	assert.NoError(v.Validate(1, is.Required))
	assert.NoError(v.Validate(struct{ Name string }{Name: "name"}, is.Required))
	assert.NoError(v.Validate([...]string{"abc"}, is.Required))
}

func TestZero(t *testing.T) {
	assert := assert.New(t)

	errString := "(zero_only) must be blank"
	assert.EqualError(v.Validate("a", is.Zero), errString)
	assert.EqualError(v.Validate(1, is.Zero), errString)
	assert.EqualError(v.Validate([...]string{"a"}, is.Zero), errString)
	assert.EqualError(v.Validate(struct{ Name string }{Name: "aaa"}, is.Zero), errString)

	assert.NoError(v.Validate("", is.Zero))
	assert.NoError(v.Validate(nil, is.Zero))
	assert.NoError(v.Validate([...]string{""}, is.Zero))
	assert.NoError(v.Validate(struct{ Name string }{}, is.Zero))
}

func TestNilOrNonZero(t *testing.T) {
	assert := assert.New(t)

	errString := "(nil_or_non_zero) can't be blank (or zero) if specified"
	assert.EqualError(v.Validate("", is.NilOrNonZero), errString)
	assert.EqualError(v.Validate([...]string{""}, is.NilOrNonZero), errString)
	assert.EqualError(v.Validate(&struct{}{}, is.NilOrNonZero), errString)
	assert.EqualError(v.Validate(henge.New(0).IntPtr().Value(), is.NilOrNonZero), errString)

	assert.NoError(v.Validate(nil, is.NilOrNonZero))
	assert.NoError(v.Validate(([]string)(nil), is.NilOrNonZero))
	assert.NoError(v.Validate((*int)(nil), is.NilOrNonZero))
	assert.NoError(v.Validate((*struct{})(nil), is.NilOrNonZero))

	assert.NoError(v.Validate("a", is.NilOrNonZero))
	assert.NoError(v.Validate([]string{"aaa"}, is.NilOrNonZero))
	assert.NoError(v.Validate(henge.New(1).IntPtr().Value(), is.NilOrNonZero))
	assert.NoError(v.Validate(&struct{ Name string }{Name: "aaa"}, is.NilOrNonZero))
}

func TestAny(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(v.Validate("", is.Any))
	assert.NoError(v.Validate(nil, is.Any))
}

func TestIn(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(v.Validate("b", is.In("a", "b", "c")))
	assert.EqualError(v.Validate("d", is.In("a", "b", "c")), "(inclusion) is not included in [a b c]")
	assert.NoError(v.Validate(2, is.In(1, 2, 3)))
	assert.EqualError(v.Validate(5, is.In(1, 2, 3)), "(inclusion) is not included in [1 2 3]")

	// NOTE: Does not match if the types are different
	assert.EqualError(v.Validate(int64(1), is.In(1, 2, 3)), "(inclusion) is not included in [1 2 3]")

	// NOTE: For pointers, it compares Elem values
	i := 2
	assert.NoError(v.Validate(&i, is.In(1, 2, 3)))
	i = 5
	assert.Error(v.Validate(&i, is.In(1, 2, 3)))

	i = 2
	x, y, z := 1, 2, 3
	assert.NoError(v.Validate(&i, is.In(&x, &y, &z)))
	i = 5
	assert.Error(v.Validate(&i, is.In(&x, &y, &z)))
}

func TestLengthRange(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate("abc", is.LengthBetween(4, math.MaxInt64)),
		"(too_short_length) is too short length (minimum is 4 characters)",
	)
	assert.EqualError(
		v.Validate("abc", is.LengthBetween(0, 2)),
		"(too_long_length) is too long length (maximum is 2 characters)",
	)
	assert.EqualError(
		v.Validate("abc", is.LengthBetween(4, 3)),
		"(too_short_length) is too short length (minimum is 4 characters)",
	)
	assert.EqualError(
		v.Validate(0, is.LengthBetween(0, 10)),
		"(string_only) must be a string",
	)
	assert.NoError(v.Validate("abc", is.LengthBetween(0, 10)))
}
