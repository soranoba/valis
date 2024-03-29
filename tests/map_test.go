package tests

import (
	"strings"
	"testing"

	"github.com/soranoba/henge/v2"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/code"
	"github.com/soranoba/valis/is"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate(
			map[string]string{"a": "A", "b": "B"},
			valis.Key("a", is.In("B")),
			valis.Key("b", is.In("A")),
		),
		"(inclusion) [key: a] is not included in [B]\n(inclusion) [key: b] is not included in [A]",
	)
	assert.NoError(
		v.Validate(
			map[string]string{"a": "A", "b": "B"},
			valis.Key("a", is.In("A")),
			valis.Key("b", is.In("B")),
		),
	)
	assert.NoError(
		v.Validate(
			&map[string]string{"a": "A", "b": "B"},
			valis.Key("a", is.In("A")),
			valis.Key("b", is.In("B")),
		),
	)

	// NOTE: If the key type is a pointer, it only matches when it specifies the same pointer
	a, b, c := "a", "b", "c"
	m := map[*string]*string{
		&a: henge.New("A").StringPtr().Value(),
		&b: henge.New("B").StringPtr().Value(),
		&c: nil,
	}
	assert.EqualError(
		v.Validate(m, valis.Key("a", is.In("A"))),
		"(not_assignable) can't assign to string",
	)
	if err := v.Validate(m, valis.Key(henge.New("a").StringPtr().Value())); assert.Error(err) {
		assert.True(strings.HasPrefix(err.Error(), "(no_key) requires the value at the key"), err.Error())
	}
	assert.NoError(v.Validate(m, valis.Key(&a, is.NonZero), valis.Key(&b, is.NonZero)))

	// NOTE: it returns an error when the value is not map
	assert.EqualError(v.Validate("a", valis.Key(&a)), "(not_map) must be any map")

	// NOTE: CommonRules automatically check.
	v := valis.NewValidator()
	v.SetCommonRules()
	assert.NoError(v.Validate(m, valis.Key(&a)))
	v.SetCommonRules(is.NonZero)
	if err := v.Validate(m, valis.Key(&c)); assert.Error(err) {
		assert.Equal(err.(*valis.ValidationError).Details()[0].Code(), code.NonZero)
	}

	// NOTE: CommonRules only check to the specified Field.
	assert.NoError(v.Validate(m, valis.Key(&a)))
}

func TestEachKeys(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(
		v.Validate(map[string]string{"a": "A", "b": "B"}, valis.EachKeys(is.NonZero)),
	)
	assert.EqualError(
		v.Validate(map[string]string{"": "A", "b": "B"}, valis.EachKeys(is.NonZero)),
		"(non_zero) [key: ] can't be blank (or zero)",
	)
	assert.NoError(
		v.Validate(&map[string]string{"a": "A", "b": "B"}, valis.EachKeys(is.NonZero)),
	)

	// NOTE: it returns an error when the value is not map
	assert.EqualError(v.Validate("a", valis.EachKeys(is.NonZero)), "(not_map) must be any map")

	// NOTE: CommonRules automatically check, but values is not validated.
	v := valis.NewValidator()
	v.SetCommonRules()
	assert.NoError(v.Validate(map[string]string{"": ""}, valis.EachKeys(is.Any)))
	v.SetCommonRules(is.NonZero)
	assert.EqualError(
		v.Validate(map[string]string{"": "a", "b": ""}, valis.EachKeys(is.Any)),
		"(non_zero) [key: ] can't be blank (or zero)",
	)
}

func TestEachValues(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(
		v.Validate(map[string]string{"a": "A", "b": "B"}, valis.EachValues(is.NonZero)),
	)
	assert.EqualError(
		v.Validate(map[string]string{"a": "A", "b": ""}, valis.EachValues(is.NonZero)),
		"(non_zero) [b] can't be blank (or zero)",
	)
	assert.NoError(
		v.Validate(&map[string]string{"a": "A", "b": "B"}, valis.EachValues(is.NonZero)),
	)

	// NOTE: it returns an error when the value is not map
	assert.EqualError(v.Validate("a", valis.EachValues(is.NonZero)), "(not_map) must be any map")

	// NOTE: CommonRules automatically check, but values is not validated.
	v := valis.NewValidator()
	v.SetCommonRules()
	assert.NoError(v.Validate(map[string]string{"": ""}, valis.EachValues(is.Any)))
	v.SetCommonRules(is.NonZero)
	assert.EqualError(
		v.Validate(map[string]string{"": "a", "b": ""}, valis.EachValues(is.Any)),
		"(non_zero) [b] can't be blank (or zero)",
	)
}
