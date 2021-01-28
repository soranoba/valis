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
		v.Validate("", is.NonZero),
		"(non_zero) can't be blank (or zero)",
	)
	assert.NoError(
		v.Validate("", when.IsStruct(is.NonZero)),
	)

	assert.EqualError(
		v.Validate(User{}, when.IsStruct(is.NonZero)),
		"(non_zero) can't be blank (or zero)",
	)
	assert.NoError(
		v.Validate(User{Name: "alice"}, when.IsStruct(is.NonZero)),
	)
}

func TestHasFieldTag(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name string `required:"true"`
	}

	assert.EqualError(
		v.Validate("", is.NonZero),
		"(non_zero) can't be blank (or zero)",
	)
	assert.NoError(
		v.Validate("", when.HasFieldTag("required", is.NonZero)),
	)

	assert.NoError(
		v.Validate(User{}, when.HasFieldTag("required", is.NonZero)),
	)

	u := User{}
	assert.EqualError(
		v.Validate(&u, valis.Field(&u.Name, when.HasFieldTag("required", is.NonZero))),
		"(non_zero) .Name can't be blank (or zero)",
	)
	u = User{Name: "alice"}
	assert.NoError(
		v.Validate(&u, valis.Field(&u.Name, when.IsStruct(is.NonZero))),
	)
}
