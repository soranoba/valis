package is_test

import (
	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/stretchr/testify/assert"
	"testing"
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

	assert.EqualError(v.Validate("", is.Required), "(required) cannot be blank, but got \"\"")
	assert.EqualError(v.Validate(0, is.Required), "(required) cannot be blank, but got 0")
	assert.EqualError(v.Validate(struct{}{}, is.Required), "(required) cannot be blank, but got struct {}{}")
	assert.EqualError(v.Validate(nil, is.Required), "(required) cannot be blank, but got <nil>")
	assert.EqualError(v.Validate([...]string{}, is.Required), "(required) cannot be blank, but got [0]string{}")
	assert.EqualError(v.Validate([...]string{""}, is.Required), "(required) cannot be blank, but got [1]string{\"\"}")

	assert.NoError(v.Validate("a", is.Required))
	assert.NoError(v.Validate(1, is.Required))
	assert.NoError(v.Validate(struct{ Name string }{Name: "name"}, is.Required))
	assert.NoError(v.Validate([...]string{"aaaa"}, is.Required))
}

func TestZero(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(v.Validate("a", is.Zero), "(zero) must be nil or zero, but got \"a\"")
	assert.EqualError(v.Validate(1, is.Zero), "(zero) must be nil or zero, but got 1")
	assert.EqualError(v.Validate([...]string{"a"}, is.Zero), "(zero) must be nil or zero, but got [1]string{\"a\"}")
	assert.EqualError(
		v.Validate(struct{ Name string }{Name: "aaa"}, is.Zero),
		"(zero) must be nil or zero, but got struct { Name string }{Name:\"aaa\"}",
	)

	assert.NoError(v.Validate("", is.Zero))
	assert.NoError(v.Validate(nil, is.Zero))
	assert.NoError(v.Validate([...]string{""}, is.Zero))
	assert.NoError(v.Validate(struct{ Name string }{}, is.Zero))
}

func TestNilOrNonZero(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(v.Validate("", is.NilOrNonZero), "(nil_or_non_zero) must be nil or non-zero, but got \"\"")
	assert.EqualError(v.Validate([...]string{""}, is.NilOrNonZero), "(nil_or_non_zero) must be nil or non-zero, but got [1]string{\"\"}")
	assert.EqualError(v.Validate(&struct{}{}, is.NilOrNonZero), "(nil_or_non_zero) must be nil or non-zero, but got &struct {}{}")
	assert.Error(v.Validate(henge.New(0).IntPtr().Value(), is.NilOrNonZero))

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
	assert.EqualError(v.Validate("d", is.In("a", "b", "c")), "(inclusion) is not included in [a b c], but got \"d\"")
	assert.NoError(v.Validate(2, is.In(1, 2, 3)))
	assert.EqualError(v.Validate(5, is.In(1, 2, 3)), "(inclusion) is not included in [1 2 3], but got 5")

	// NOTE: Does not match if the types are different
	assert.EqualError(v.Validate(int64(1), is.In(1, 2, 3)), "(inclusion) is not included in [1 2 3], but got 1")

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
