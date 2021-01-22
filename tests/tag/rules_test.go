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
		LastName  string `required:"false"`
		Age       int    `required:"True"`
	}

	assert.EqualError(
		v.Validate(User{}, valistag.Required),
		"(required) .FirstName is required\n"+
			"(required) .Age is required",
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
		"(required) .Name is required",
	)
	assert.NoError(
		v.Validate(&User{Name: "Alice"}, valistag.Validate),
	)
}

func TestValidate_min(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name   string            `validate:"min=2"`
		Tags   []string          `validate:"min=2"`
		Params map[string]string `validate:"min=2"`
	}
	assert.EqualError(
		v.Validate(User{
			Name:   "🍺",
			Tags:   []string{""},
			Params: map[string]string{"": ""},
		}, valistag.Validate),
		`(too_short_length) .Name is too short length (minimum is 2 characters)
(too_short_len) .Tags is too few elements (minimum is 2 elements)
(too_short_len) .Params is too few elements (minimum is 2 elements)`,
	)
	assert.NoError(
		v.Validate(User{
			Name:   "Alice",
			Tags:   []string{"", ""},
			Params: map[string]string{"a": "", "b": ""},
		}, valistag.Validate),
	)
}

func TestValidate_max(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name   string            `validate:"max=2"`
		Tags   []string          `validate:"max=2"`
		Params map[string]string `validate:"max=2"`
	}
	assert.EqualError(
		v.Validate(User{
			Name:   "abc",
			Tags:   []string{"a", "b", "c"},
			Params: map[string]string{"a": "", "b": "", "c": ""},
		}, valistag.Validate),
		`(too_long_length) .Name is too long length (maximum is 2 characters)
(too_long_len) .Tags is too many elements (maximum is 2 elements)
(too_long_len) .Params is too many elements (maximum is 2 elements)`,
	)
	assert.NoError(
		v.Validate(User{
			Name:   "🍺🍺",
			Tags:   []string{"", ""},
			Params: map[string]string{"a": "", "b": ""},
		}, valistag.Validate),
	)
}
