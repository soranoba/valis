package tests

import (
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestField(t *testing.T) {
	assert := assert.New(t)
	type User struct {
		Name string
		Age  int
	}

	user := User{}
	assert.EqualError(
		v.Validate(
			&user,
			valis.Field(&user.Name, is.Required),
		),
		"(required) .Name cannot be blank, but got \"\"",
	)
	assert.EqualError(
		v.Validate(
			&user,
			valis.Field(&user.Name, is.Required),
			valis.Field(&user.Age, is.Required),
		),
		"(required) .Name cannot be blank, but got \"\". (required) .Age cannot be blank, but got 0",
	)

	alice := User{Name: "Alice", Age: 22}
	assert.NoError(
		v.Validate(
			&alice,
			valis.Field(&alice.Name, is.Required),
			valis.Field(&alice.Age, is.Required),
		),
	)

	// NOTE: it returns an error when the value is not struct
	assert.EqualError(v.Validate("a", valis.Field(&user.Name)), "(invalid_type) must be a struct, but got \"a\"")

	// NOTE: it must be specified with pointers
	assert.Panics(func() {
		v.Validate(&alice, valis.Field(&user.Name))
	})
	assert.Panics(func() {
		v.Validate(alice, valis.Field(&alice.Name))
	})
	assert.Panics(func() {
		v.Validate(&alice, valis.Field(alice.Name))
	})

	// NOTE: CommonRules automatically check.
	v := valis.NewValidator()
	v.SetCommonRules()
	user.Name = "Bob"
	assert.NoError(
		v.Validate(&user, valis.Field(&user.Name), valis.Field(&user.Age)),
	)
	v.SetCommonRules(is.Required)
	assert.EqualError(
		v.Validate(&user, valis.Field(&user.Name), valis.Field(&user.Age)),
		"(required) .Age cannot be blank, but got 0",
	)

	// NOTE: CommonRules only check to the specified Field.
	assert.NoError(
		v.Validate(&user, valis.Field(&user.Name)),
	)
}
