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
		"(required) cannot be blank, but got \"\"",
	)

	assert.NoError(
		v.Validate("", when.IsType(reflect.TypeOf(0), is.Required)),
	)
	assert.NoError(
		v.Validate("", when.IsType(reflect.TypeOf((*string)(nil)), is.Required)),
	)
	assert.EqualError(
		v.Validate("", when.IsType(reflect.TypeOf(""), is.Required)),
		"(required) cannot be blank, but got \"\"",
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
		"(required) cannot be blank, but got \"\"",
	)

	assert.NoError(
		v.Validate("", when.IsTypeOrPtr(reflect.TypeOf(0), is.Required)),
	)
	assert.NoError(
		v.Validate("", when.IsTypeOrPtr(reflect.TypeOf((*string)(nil)), is.Required)),
	)
	assert.EqualError(
		v.Validate("", when.IsTypeOrPtr(reflect.TypeOf(""), is.Required)),
		"(required) cannot be blank, but got \"\"",
	)
	assert.EqualError(
		v.Validate((*string)(nil), when.IsTypeOrPtr(reflect.TypeOf(""), is.Required)),
		"(required) cannot be blank, but got (*string)(nil)",
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
		"(required) cannot be blank, but got (*string)(nil)",
	)

	assert.NoError(
		v.Validate((*string)(nil), when.IsTypeOrElem(reflect.TypeOf(0), is.Required)),
	)
	assert.NoError(
		v.Validate((*string)(nil), when.IsTypeOrElem(reflect.TypeOf(""), is.Required)),
	)
	assert.EqualError(
		v.Validate((*string)(nil), when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.Required)),
		"(required) cannot be blank, but got (*string)(nil)",
	)
	assert.EqualError(
		v.Validate("", when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.Required)),
		"(required) cannot be blank, but got \"\"",
	)

	// NOTE: invalid type
	assert.NoError(
		v.Validate(nil, when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.Required)),
	)
	assert.Panics(func() {
		v.Validate("", when.IsTypeOrElem(reflect.TypeOf(nil), is.Required))
	})
}
