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
		v.Validate("", is.Required),
		"(required) is required",
	)

	assert.NoError(
		v.Validate("", when.IsType(reflect.TypeOf(0), is.Required)),
	)
	assert.NoError(
		v.Validate("", when.IsType(reflect.TypeOf((*string)(nil)), is.Required)),
	)
	assert.EqualError(
		v.Validate("", when.IsType(reflect.TypeOf(""), is.Required)),
		"(required) is required",
	)
	assert.NoError(
		v.Validate((*string)(nil), when.IsType(reflect.TypeOf(""), is.Required)),
	)

	// NOTE: invalid type
	assert.NoError(
		v.Validate(nil, when.IsType(reflect.TypeOf(""), is.Required)),
	)
	assert.Panics(func() {
		v.Validate("", when.IsType(reflect.TypeOf(nil), is.Required))
	})
}

func TestIsTypeOrPtr(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate("", is.Required),
		"(required) is required",
	)

	assert.NoError(
		v.Validate("", when.IsTypeOrPtr(reflect.TypeOf(0), is.Required)),
	)
	assert.NoError(
		v.Validate("", when.IsTypeOrPtr(reflect.TypeOf((*string)(nil)), is.Required)),
	)
	assert.EqualError(
		v.Validate("", when.IsTypeOrPtr(reflect.TypeOf(""), is.Required)),
		"(required) is required",
	)
	assert.EqualError(
		v.Validate((*string)(nil), when.IsTypeOrPtr(reflect.TypeOf(""), is.Required)),
		"(required) is required",
	)

	// NOTE: invalid type
	assert.NoError(
		v.Validate(nil, when.IsTypeOrPtr(reflect.TypeOf(""), is.Required)),
	)
	assert.Panics(func() {
		v.Validate("", when.IsTypeOrPtr(reflect.TypeOf(nil), is.Required))
	})
}

func TestIsTypeOrElem(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate((*string)(nil), is.Required),
		"(required) is required",
	)

	assert.NoError(
		v.Validate((*string)(nil), when.IsTypeOrElem(reflect.TypeOf(0), is.Required)),
	)
	assert.NoError(
		v.Validate((*string)(nil), when.IsTypeOrElem(reflect.TypeOf(""), is.Required)),
	)
	assert.EqualError(
		v.Validate((*string)(nil), when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.Required)),
		"(required) is required",
	)
	assert.EqualError(
		v.Validate("", when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.Required)),
		"(required) is required",
	)

	// NOTE: invalid type
	assert.NoError(
		v.Validate(nil, when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.Required)),
	)
	assert.Panics(func() {
		v.Validate("", when.IsTypeOrElem(reflect.TypeOf(nil), is.Required))
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
