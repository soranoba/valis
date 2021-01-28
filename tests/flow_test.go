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

	Any := func(ctx *valis.WhenContext) bool { return true }
	Never := func(ctx *valis.WhenContext) bool { return false }
	assert.EqualError(
		v.Validate("", valis.When(Any, is.NonZero)),
		"(non_zero) can't be blank (or zero)",
	)
	assert.NoError(v.Validate("", valis.When(Never, is.Required)))
}

func TestWhen_ElseIf(t *testing.T) {
	assert := assert.New(t)

	Any := func(ctx *valis.WhenContext) bool { return true }
	Never := func(ctx *valis.WhenContext) bool { return false }
	assert.NoError(
		v.Validate("", valis.When(Never, is.Zero)),
	)
	assert.EqualError(
		v.Validate("", valis.When(Never, is.Zero).ElseIf(Any, is.NonZero)),
		"(non_zero) can't be blank (or zero)",
	)
	assert.NoError(
		v.Validate("", valis.When(Any, is.Zero).ElseIf(Any, is.NonZero)),
	)
}

func TestWhen_ElseWhen(t *testing.T) {
	assert := assert.New(t)

	Any := func(ctx *valis.WhenContext) bool { return true }
	Never := func(ctx *valis.WhenContext) bool { return false }
	assert.NoError(
		v.Validate("", valis.When(Never, is.Zero)),
	)
	assert.EqualError(
		v.Validate("",
			valis.When(Never, is.Zero).
				ElseWhen(valis.When(Any, is.NonZero))),
		"(non_zero) can't be blank (or zero)",
	)
	assert.EqualError(
		v.Validate("",
			valis.When(Never, is.Zero).
				ElseIf(Never, is.Zero).
				ElseWhen(valis.When(Any, is.NonZero))),
		"(non_zero) can't be blank (or zero)",
	)
	assert.EqualError(
		v.Validate("",
			valis.When(Never, is.Zero).
				ElseWhen(valis.When(Never, is.Zero)).
				ElseIf(Any, is.NonZero)),
		"(non_zero) can't be blank (or zero)",
	)
	assert.EqualError(
		v.Validate("",
			valis.When(Never, is.Zero).
				ElseWhen(valis.When(Never, is.Zero).ElseIf(Any, is.NonZero)).
				ElseIf(Any, is.In(0))),
		"(non_zero) can't be blank (or zero)",
	)
	assert.NoError(
		v.Validate("",
			valis.When(Any, is.Zero).
				ElseWhen(valis.When(Never, is.NonZero))),
	)
}

func TestWhen_Else(t *testing.T) {
	assert := assert.New(t)

	Any := func(ctx *valis.WhenContext) bool { return true }
	Never := func(ctx *valis.WhenContext) bool { return false }
	assert.NoError(
		v.Validate("", valis.When(Never, is.Zero)),
	)
	assert.EqualError(
		v.Validate("",
			valis.When(Never, is.Zero).
				Else(valis.When(Any, is.NonZero))),
		"(non_zero) can't be blank (or zero)",
	)
	assert.EqualError(
		v.Validate("",
			valis.When(Never, is.Zero).
				ElseIf(Never, is.Zero).
				Else(is.NonZero)),
		"(non_zero) can't be blank (or zero)",
	)
	assert.EqualError(
		v.Validate("",
			valis.When(Never, is.Zero).
				ElseWhen(valis.When(Any, is.NonZero)).
				Else(is.In())),
		"(non_zero) can't be blank (or zero)",
	)
	assert.EqualError(
		v.Validate("",
			valis.When(Never, is.Zero).
				ElseWhen(valis.When(Never, is.Zero).ElseIf(Any, is.NonZero)).
				Else(is.In())),
		"(non_zero) can't be blank (or zero)",
	)
	assert.EqualError(
		v.Validate("",
			valis.When(Never, is.Zero).
				ElseWhen(valis.When(Never, is.Zero).ElseIf(Never, is.In())).
				Else(is.NonZero)),
		"(non_zero) can't be blank (or zero)",
	)
}

func TestEach(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate([]string{"a", "", "b"}, valis.Each(is.NonZero)),
		"(non_zero) [1] can't be blank (or zero)",
	)
	assert.NoError(
		v.Validate([]string{"a", "b", "c"}, valis.Each(is.NonZero)),
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
		v.Validate(value, valis.Each(is.NonZero)),
	)

	// NOTE: value must be array or slice.
	assert.EqualError(v.Validate("", valis.Each(is.NonZero)), "(not_array) must be any array")

	// NOTE: CommonRules automatically check.
	v := valis.NewValidator()
	v.SetCommonRules()
	assert.NoError(v.Validate([]string{""}, valis.Each()))
	v.SetCommonRules(is.NonZero)
	assert.EqualError(v.Validate([]string{""}, valis.Each()), "(non_zero) [0] can't be blank (or zero)")
}
