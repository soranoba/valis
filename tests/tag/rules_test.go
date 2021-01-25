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
		v.Validate(User{}, valis.EachFields(valistag.Required)),
		"(required) .FirstName is required\n"+
			"(required) .Age is required",
	)
	assert.NoError(
		v.Validate(&User{FirstName: "Taro", LastName: "Soto", Age: 20}, valis.EachFields(valistag.Required)),
	)
}

func TestFormat(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name string `format:"^[a-z]+$"`
	}

	assert.EqualError(
		v.Validate(&User{}, valis.EachFields(valistag.Format)),
		"(regexp) .Name is a mismatch with the regular expression. (^[a-z]+$)",
	)
	assert.EqualError(
		v.Validate(&User{
			Name: "abc0123",
		}, valis.EachFields(valistag.Format)),
		"(regexp) .Name is a mismatch with the regular expression. (^[a-z]+$)",
	)
	assert.NoError(
		v.Validate(&User{
			Name: "abcdef",
		}, valis.EachFields(valistag.Format)),
	)
}

func TestValidate_required(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name string `validate:"required"`
	}
	assert.EqualError(
		v.Validate(User{}, valis.EachFields(valistag.Validate)),
		"(required) .Name is required",
	)
	assert.NoError(
		v.Validate(&User{Name: "Alice"}, valis.EachFields(valistag.Validate)),
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
		}, valis.EachFields(valistag.Validate)),
		`(too_short_length) .Name is too short length (minimum is 2 characters)
(too_short_len) .Tags is too few elements (minimum is 2 elements)
(too_short_len) .Params is too few elements (minimum is 2 elements)`,
	)
	assert.NoError(
		v.Validate(User{
			Name:   "Alice",
			Tags:   []string{"", ""},
			Params: map[string]string{"a": "", "b": ""},
		}, valis.EachFields(valistag.Validate)),
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
		}, valis.EachFields(valistag.Validate)),
		`(too_long_length) .Name is too long length (maximum is 2 characters)
(too_long_len) .Tags is too many elements (maximum is 2 elements)
(too_long_len) .Params is too many elements (maximum is 2 elements)`,
	)
	assert.NoError(
		v.Validate(User{
			Name:   "🍺🍺",
			Tags:   []string{"", ""},
			Params: map[string]string{"a": "", "b": ""},
		}, valis.EachFields(valistag.Validate)),
	)
}

func TestValidate_oneof(t *testing.T) {
	assert := assert.New(t)

	type Model struct {
		S string  `validate:"oneof=a b c"`
		I int     `validate:"oneof=1 2 3"`
		U uint    `validate:"oneof=1 2 3"`
		B bool    `validate:"oneof=true"`
		F float64 `validate:"oneof=1 2.5 3"`
	}
	assert.EqualError(
		v.Validate(&Model{}, valis.EachFields(valistag.Validate)),
		`(inclusion) .S is not included in [a b c]
(inclusion) .I is not included in [1 2 3]
(inclusion) .U is not included in [1 2 3]
(inclusion) .B is not included in [true]
(inclusion) .F is not included in [1 2.5 3]`,
	)
	assert.NoError(
		v.Validate(&Model{
			S: "a",
			I: 2,
			U: 3,
			B: true,
			F: 3,
		}, valis.EachFields(valistag.Validate)),
	)
}
