package tag_test

import (
	"github.com/soranoba/valis"
	valistag "github.com/soranoba/valis/tag"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
		"(required) .FirstName cannot be blank, but got \"\". "+
			"(required) .Age cannot be blank, but got 0",
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
		"(too_short_length) .Name is too short length (min: 2), but got \"🍺\". "+
			"(too_short_length) .Tags is too short length (min: 2), but got []string{\"\"}. "+
			"(too_short_length) .Params is too short length (min: 2), but got map[string]string{\"\":\"\"}",
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
		"(too_long_length) .Name is too long length (max: 2), but got \"abc\". "+
			"(too_long_length) .Tags is too long length (max: 2), but got []string{\"a\", \"b\", \"c\"}. "+
			"(too_long_length) .Params is too long length (max: 2), but got map[string]string{\"a\":\"\", \"b\":\"\", \"c\":\"\"}",
	)
	assert.NoError(
		v.Validate(User{
			Name:   "🍺🍺",
			Tags:   []string{"", ""},
			Params: map[string]string{"a": "", "b": ""},
		}, valistag.Validate),
	)

	type Company struct {
		StartedTime      *time.Time `validate:"max=2"`
		NumberOfEmployee int        `validate:"max=2"`
	}
}
