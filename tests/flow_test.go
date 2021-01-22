package tests

import (
	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOr(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(v.Validate("", valis.Or(is.Zero, is.In("aaa"))))
	assert.NoError(v.Validate("aaa", valis.Or(is.Zero, is.In("aaa"))))
	assert.EqualError(
		v.Validate("abc", valis.Or(is.Zero, is.In("aaa"))),
		"(invalid) is invalid",
	)
}

func TestWhen(t *testing.T) {
	assert := assert.New(t)

	Any := func(interface{}) bool { return true }
	Never := func(interface{}) bool { return false }
	assert.EqualError(
		v.Validate("", valis.When(Any, is.Required)),
		"(required) is required",
	)
	assert.NoError(v.Validate("", valis.When(Never, is.Required)))
}

func TestEach(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate([]string{"a", "", "b"}, valis.Each(is.Required)),
		"(required) [1] is required",
	)
	assert.NoError(
		v.Validate([]string{"a", "b", "c"}, valis.Each(is.Required)),
	)

	value := []*string{
		henge.New("a").StringPtr().Value(),
		henge.New("").StringPtr().Value(),
		henge.New("b").StringPtr().Value(),
	}
	assert.EqualError(
		v.Validate(value, valis.Each(is.NilOrNonZero)),
		"(nil_or_non_zero) [1] can't be blank (or zero) if specified",
	)
	value = []*string{
		henge.New("a").StringPtr().Value(),
		henge.New("b").StringPtr().Value(),
		henge.New("c").StringPtr().Value(),
	}
	assert.NoError(
		v.Validate(value, valis.Each(is.Required)),
	)

	// NOTE: value must be array or slice.
	assert.EqualError(v.Validate("", valis.Each(is.Required)), "(not_array) must be any array")

	// NOTE: CommonRules automatically check.
	v := valis.NewValidator()
	v.SetCommonRules()
	assert.NoError(v.Validate([]string{""}, valis.Each()))
	v.SetCommonRules(is.Required)
	assert.EqualError(v.Validate([]string{""}, valis.Each()), "(required) [0] is required")
}
