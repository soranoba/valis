package tests

import (
	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestOr(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(v.Validate("", valis.Or(is.Zero, is.In("aaa"))))
	assert.NoError(v.Validate("aaa", valis.Or(is.Zero, is.In("aaa"))))
	assert.EqualError(
		v.Validate("abc", valis.Or(is.Zero, is.In("aaa"))),
		"(or) cannot meet either rule, but got \"abc\"",
	)
}

func TestWhen(t *testing.T) {
	assert := assert.New(t)

	Any := func(interface{}) bool { return true }
	Never := func(interface{}) bool { return false }
	assert.EqualError(
		v.Validate("", valis.When(Any, is.Required)),
		"(required) cannot be blank, but got \"\"",
	)
	assert.NoError(v.Validate("", valis.When(Never, is.Required)))
}

func TestEach(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate([]string{"a", "", "b"}, valis.Each(is.Required)),
		"(required) [1] cannot be blank, but got \"\"",
	)
	assert.NoError(
		v.Validate([]string{"a", "b", "c"}, valis.Each(is.Required)),
	)

	value := []*string{
		henge.New("a").StringPtr().Value(),
		henge.New("").StringPtr().Value(),
		henge.New("b").StringPtr().Value(),
	}
	if err := v.Validate(value, valis.Each(is.NilOrNonZero)); assert.Error(err) {
		assert.True(strings.HasPrefix(err.Error(), "(nil_or_non_zero) [1] must be nil or non-zero"), err.Error())
	}
	value = []*string{
		henge.New("a").StringPtr().Value(),
		henge.New("b").StringPtr().Value(),
		henge.New("c").StringPtr().Value(),
	}
	assert.NoError(
		v.Validate(value, valis.Each(is.Required)),
	)

	// NOTE: value must be array or slice.
	assert.EqualError(v.Validate("", valis.Each(is.Required)), "(invalid_type) must be array or slice, but got \"\"")

	// NOTE: CommonRules automatically check.
	v := valis.NewValidator()
	v.SetCommonRules()
	assert.NoError(v.Validate([]string{""}, valis.Each()))
	v.SetCommonRules(is.Required)
	assert.EqualError(v.Validate([]string{""}, valis.Each()), "(required) [0] cannot be blank, but got \"\"")
}
