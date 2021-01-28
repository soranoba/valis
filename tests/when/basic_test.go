package when_test

import (
	"github.com/soranoba/valis/is"
	"github.com/soranoba/valis/when"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIsType(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate("", is.NonZero),
		"(non_zero) can't be blank (or zero)",
	)

	assert.NoError(
		v.Validate("", when.IsType(reflect.TypeOf(0), is.NonZero)),
	)
	assert.NoError(
		v.Validate("", when.IsType(reflect.TypeOf((*string)(nil)), is.NonZero)),
	)
	assert.EqualError(
		v.Validate("", when.IsType(reflect.TypeOf(""), is.NonZero)),
		"(non_zero) can't be blank (or zero)",
	)
	assert.NoError(
		v.Validate((*string)(nil), when.IsType(reflect.TypeOf(""), is.NonZero)),
	)

	// NOTE: invalid type
	assert.NoError(
		v.Validate(nil, when.IsType(reflect.TypeOf(""), is.NonZero)),
	)
	assert.Panics(func() {
		v.Validate("", when.IsType(reflect.TypeOf(nil), is.NonZero))
	})
}

func TestIsTypeOrPtr(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate("", is.NonZero),
		"(non_zero) can't be blank (or zero)",
	)

	assert.NoError(
		v.Validate("", when.IsTypeOrPtr(reflect.TypeOf(0), is.NonZero)),
	)
	assert.NoError(
		v.Validate("", when.IsTypeOrPtr(reflect.TypeOf((*string)(nil)), is.NonZero)),
	)
	assert.EqualError(
		v.Validate("", when.IsTypeOrPtr(reflect.TypeOf(""), is.NonZero)),
		"(non_zero) can't be blank (or zero)",
	)
	assert.EqualError(
		v.Validate((*string)(nil), when.IsTypeOrPtr(reflect.TypeOf(""), is.NonZero)),
		"(non_zero) can't be blank (or zero)",
	)

	// NOTE: invalid type
	assert.NoError(
		v.Validate(nil, when.IsTypeOrPtr(reflect.TypeOf(""), is.NonZero)),
	)
	assert.Panics(func() {
		v.Validate("", when.IsTypeOrPtr(reflect.TypeOf(nil), is.NonZero))
	})
}

func TestIsTypeOrElem(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate((*string)(nil), is.NonZero),
		"(non_zero) can't be blank (or zero)",
	)

	assert.NoError(
		v.Validate((*string)(nil), when.IsTypeOrElem(reflect.TypeOf(0), is.NonZero)),
	)
	assert.NoError(
		v.Validate((*string)(nil), when.IsTypeOrElem(reflect.TypeOf(""), is.NonZero)),
	)
	assert.EqualError(
		v.Validate((*string)(nil), when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.NonZero)),
		"(non_zero) can't be blank (or zero)",
	)
	assert.EqualError(
		v.Validate("", when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.NonZero)),
		"(non_zero) can't be blank (or zero)",
	)

	// NOTE: invalid type
	assert.NoError(
		v.Validate(nil, when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.NonZero)),
	)
	assert.Panics(func() {
		v.Validate("", when.IsTypeOrElem(reflect.TypeOf(nil), is.NonZero))
	})
}

func TestIsNumeric(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate("aa", is.Zero),
		"(zero_only) must be blank",
	)

	assert.NoError(
		v.Validate("aa", when.IsNumeric(is.Zero)),
	)
	assert.EqualError(
		v.Validate(1.25, when.IsNumeric(is.Zero)),
		"(zero_only) must be blank",
	)
	assert.EqualError(
		v.Validate(1, when.IsNumeric(is.Zero)),
		"(zero_only) must be blank",
	)
}

func TestIsNil(t *testing.T) {
	assert := assert.New(t)

	var s []string
	assert.EqualError(
		v.Validate(s, is.Never),
		"(invalid) is invalid",
	)

	assert.NoError(
		v.Validate(s, when.IsNil().Else(is.Never)),
	)
}
