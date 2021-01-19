package tests

import (
	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestKey(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate(
			map[string]string{"a": "A", "b": "B"},
			valis.Key("a", is.In("B")),
			valis.Key("b", is.In("A")),
		),
		"(inclusion) [key: a] is not included in [B], but got \"A\". (inclusion) [key: b] is not included in [A], but got \"B\"",
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
	if err := v.Validate(m, valis.Key("a", is.In("A"))); assert.Error(err) {
		assert.True(strings.HasPrefix(err.Error(), "(invalid_type) cannot assignable string to key, "), err.Error())
	}
	if err := v.Validate(m, valis.Key(henge.New("a").StringPtr().Value())); assert.Error(err) {
		assert.True(strings.HasPrefix(err.Error(), "(not_found) does not have the key"), err.Error())
	}
	assert.NoError(v.Validate(m, valis.Key(&a, is.Required), valis.Key(&b, is.Required)))

	// NOTE: it returns an error when the value is not map
	assert.EqualError(v.Validate("a", valis.Key(&a)), "(invalid_type) must be a map, but got \"a\"")

	// NOTE: CommonRules automatically check.
	v := valis.NewValidator()
	v.SetCommonRules()
	assert.NoError(v.Validate(m, valis.Key(&a)))
	v.SetCommonRules(is.Required)
	if err := v.Validate(m, valis.Key(&c)); assert.Error(err) {
		assert.True(strings.HasPrefix(err.Error(), "(required) "), err.Error())
		assert.True(strings.HasSuffix(err.Error(), "cannot be blank, but got (*string)(nil)"), err.Error())
	}

	// NOTE: CommonRules only check to the specified Field.
	assert.NoError(v.Validate(m, valis.Key(&a)))
}

func TestEachKeys(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(
		v.Validate(map[string]string{"a": "A", "b": "B"}, valis.EachKeys(is.Required)),
	)
	assert.EqualError(
		v.Validate(map[string]string{"": "A", "b": "B"}, valis.EachKeys(is.Required)),
		"(required) [key: ] cannot be blank, but got \"\"",
	)
	assert.NoError(
		v.Validate(&map[string]string{"a": "A", "b": "B"}, valis.EachKeys(is.Required)),
	)

	// NOTE: it returns an error when the value is not map
	assert.EqualError(v.Validate("a", valis.EachKeys(is.Required)), "(invalid_type) must be a map, but got \"a\"")

	// NOTE: CommonRules automatically check, but values is not validated.
	v := valis.NewValidator()
	v.SetCommonRules()
	assert.NoError(v.Validate(map[string]string{"": ""}, valis.EachKeys(is.Any)))
	v.SetCommonRules(is.Required)
	assert.EqualError(
		v.Validate(map[string]string{"": "a", "b": ""}, valis.EachKeys(is.Any)),
		"(required) [key: ] cannot be blank, but got \"\"",
	)
}

func TestEachValues(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(
		v.Validate(map[string]string{"a": "A", "b": "B"}, valis.EachValues(is.Required)),
	)
	assert.EqualError(
		v.Validate(map[string]string{"a": "A", "b": ""}, valis.EachValues(is.Required)),
		"(required) [b] cannot be blank, but got \"\"",
	)
	assert.NoError(
		v.Validate(&map[string]string{"a": "A", "b": "B"}, valis.EachValues(is.Required)),
	)

	// NOTE: it returns an error when the value is not map
	assert.EqualError(v.Validate("a", valis.EachValues(is.Required)), "(invalid_type) must be a map, but got \"a\"")

	// NOTE: CommonRules automatically check, but values is not validated.
	v := valis.NewValidator()
	v.SetCommonRules()
	assert.NoError(v.Validate(map[string]string{"": ""}, valis.EachValues(is.Any)))
	v.SetCommonRules(is.Required)
	assert.EqualError(
		v.Validate(map[string]string{"": "a", "b": ""}, valis.EachValues(is.Any)),
		"(required) [b] cannot be blank, but got \"\"",
	)
}
