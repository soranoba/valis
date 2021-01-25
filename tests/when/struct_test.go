package when_test

import (
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/soranoba/valis/when"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsStruct(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name string
	}
	assert.EqualError(
		v.Validate("", is.Required),
		"(required) is required",
	)
	assert.NoError(
		v.Validate("", when.IsStruct(is.Required)),
	)

	assert.EqualError(
		v.Validate(User{}, when.IsStruct(is.Required)),
		"(required) is required",
	)
	assert.NoError(
		v.Validate(User{Name: "alice"}, when.IsStruct(is.Required)),
	)
}

func TestHasFieldTag(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name string `required:"true"`
	}

	assert.EqualError(
		v.Validate("", is.Required),
		"(required) is required",
	)
	assert.NoError(
		v.Validate("", when.HasFieldTag("required", is.Required)),
	)

	assert.NoError(
		v.Validate(User{}, when.HasFieldTag("required", is.Required)),
	)

	u := User{}
	assert.EqualError(
		v.Validate(&u, valis.Field(&u.Name, when.HasFieldTag("required", is.Required))),
		"(required) .Name is required",
	)
	u = User{Name: "alice"}
	assert.NoError(
		v.Validate(&u, valis.Field(&u.Name, when.IsStruct(is.Required))),
	)
}
