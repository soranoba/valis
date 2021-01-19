package tag_test

import (
	"github.com/soranoba/valis"
	valistag "github.com/soranoba/valis/tag"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	v *valis.Validator
)

func init() {
	v = valis.NewValidator()
	v.SetCommonRules()
}

func TestRequired(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		FirstName string `required:"true"`
		LastName string `required:"false"`
		Age int `required:"True"`
	}

	assert.EqualError(
		v.Validate(User{}, valistag.Required),
		"(required) .FirstName cannot be blank, but got \"\". (required) .Age cannot be blank, but got 0",
	)
	assert.NoError(
		v.Validate(&User{FirstName: "Taro", LastName: "Soto", Age: 20}, valistag.Required),
	)
}

func TestValidate_required(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name string `validate:"required"`
	}
	assert.EqualError(
		v.Validate(User{}, valistag.Validate),
		"(required) .Name cannot be blank, but got \"\"",
	)
	assert.NoError(
		v.Validate(&User{Name: "Alice"}, valistag.Validate),
	)
}
