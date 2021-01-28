package tagrule_test

import (
	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/tagrule"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	v = valis.NewValidator()
)

func TestRequired(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		FirstName *string `required:"true"`
		LastName  *string `required:"false"`
		Nickname  string  `required:"true"`
		Age       *int64  `required:"True"`
	}

	assert.EqualError(
		v.Validate(User{}, valis.EachFields(tagrule.Required)),
		"(required) .FirstName is required\n"+
			"(required) .Age is required",
	)
	assert.NoError(
		v.Validate(&User{
			FirstName: henge.ToStringPtr("Taro"),
			Age:       henge.ToIntPtr(20),
		}, valis.EachFields(tagrule.Required)),
	)
}

func TestPattern(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name *string `pattern:"^[a-z]+$"`
	}

	assert.NoError(
		v.Validate(&User{}, valis.EachFields(tagrule.Pattern)),
	)
	assert.EqualError(
		v.Validate(&User{
			Name: henge.ToStringPtr(""),
		}, valis.EachFields(tagrule.Pattern)),
		"(regexp) .Name is a mismatch with the regular expression. (^[a-z]+$)",
	)
	assert.EqualError(
		v.Validate(&User{
			Name: henge.ToStringPtr("abc0123"),
		}, valis.EachFields(tagrule.Pattern)),
		"(regexp) .Name is a mismatch with the regular expression. (^[a-z]+$)",
	)
	assert.NoError(
		v.Validate(&User{
			Name: henge.ToStringPtr("abcdef"),
		}, valis.EachFields(tagrule.Pattern)),
	)
}

func TestEnums(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name *string `enums:"alice,bob"`
		Age  *int64  `enums:"20,30"`
	}

	assert.NoError(
		v.Validate(&User{}, valis.EachFields(tagrule.Enums)),
	)
	assert.EqualError(
		v.Validate(&User{
			Name: henge.ToStringPtr(""),
			Age:  henge.ToIntPtr(0),
		}, valis.EachFields(tagrule.Enums)),
		"(inclusion) .Name is not included in [alice bob]\n(inclusion) .Age is not included in [20 30]",
	)
	assert.EqualError(
		v.Validate(&User{
			Name: henge.ToStringPtr("a"),
			Age:  henge.ToIntPtr(10),
		}, valis.EachFields(tagrule.Enums)),
		"(inclusion) .Name is not included in [alice bob]\n(inclusion) .Age is not included in [20 30]",
	)
	assert.NoError(
		v.Validate(&User{
			Name: henge.ToStringPtr("alice"),
			Age:  henge.ToIntPtr(20),
		}, valis.EachFields(tagrule.Enums)),
	)
}

func TestValidate_required(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		FirstName string  `validate:"required"`
		LastName  *string `validate:"required"`
	}
	assert.EqualError(
		v.Validate(User{}, valis.EachFields(tagrule.Validate)),
		"(required) .LastName is required",
	)
	assert.NoError(
		v.Validate(&User{
			LastName: henge.ToStringPtr("Tanaka"),
		}, valis.EachFields(tagrule.Validate)),
	)
}

func TestValidate_lte(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Age *int64 `validate:"lte=2"`
	}
	assert.NoError(
		v.Validate(User{}, valis.EachFields(tagrule.Validate)),
	)
	assert.NoError(
		v.Validate(User{Age: henge.ToIntPtr(1)}, valis.EachFields(tagrule.Validate)),
	)
	assert.NoError(
		v.Validate(User{Age: henge.ToIntPtr(2)}, valis.EachFields(tagrule.Validate)),
	)
	assert.EqualError(
		v.Validate(User{Age: henge.ToIntPtr(20)}, valis.EachFields(tagrule.Validate)),
		"(lte) .Age must be less than or equal to 2",
	)
}

func TestValidate_lt(t *testing.T) {
	assert := assert.New(t)
	type User struct {
		Age *int64 `validate:"lt=2"`
	}
	assert.NoError(
		v.Validate(User{}, valis.EachFields(tagrule.Validate)),
	)
	assert.NoError(
		v.Validate(User{Age: henge.ToIntPtr(1)}, valis.EachFields(tagrule.Validate)),
	)
	assert.EqualError(
		v.Validate(User{Age: henge.ToIntPtr(2)}, valis.EachFields(tagrule.Validate)),
		"(lt) .Age must be less than 2",
	)
	assert.EqualError(
		v.Validate(User{Age: henge.ToIntPtr(20)}, valis.EachFields(tagrule.Validate)),
		"(lt) .Age must be less than 2",
	)
}

func TestValidate_gte(t *testing.T) {
	assert := assert.New(t)
	type User struct {
		Age *int64 `validate:"gte=2"`
	}
	assert.NoError(
		v.Validate(User{}, valis.EachFields(tagrule.Validate)),
	)
	assert.EqualError(
		v.Validate(User{Age: henge.ToIntPtr(1)}, valis.EachFields(tagrule.Validate)),
		"(gte) .Age must be greater than or equal to 2",
	)
	assert.NoError(
		v.Validate(User{Age: henge.ToIntPtr(2)}, valis.EachFields(tagrule.Validate)),
	)
	assert.NoError(
		v.Validate(User{Age: henge.ToIntPtr(20)}, valis.EachFields(tagrule.Validate)),
	)
}

func TestValidate_gt(t *testing.T) {
	assert := assert.New(t)
	type User struct {
		Age *int64 `validate:"gt=2"`
	}
	assert.NoError(
		v.Validate(User{}, valis.EachFields(tagrule.Validate)),
	)
	assert.EqualError(
		v.Validate(User{Age: henge.ToIntPtr(1)}, valis.EachFields(tagrule.Validate)),
		"(gt) .Age must be greater than 2",
	)
	assert.EqualError(
		v.Validate(User{Age: henge.ToIntPtr(2)}, valis.EachFields(tagrule.Validate)),
		"(gt) .Age must be greater than 2",
	)
	assert.NoError(
		v.Validate(User{Age: henge.ToIntPtr(20)}, valis.EachFields(tagrule.Validate)),
	)
}

func TestValidate_min(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name   *string           `validate:"min=2"`
		Age    *int64            `validate:"min=2"`
		Tags   []string          `validate:"min=2"`
		Params map[string]string `validate:"min=2"`
	}
	assert.NoError(
		v.Validate(User{}, valis.EachFields(tagrule.Validate)),
	)
	assert.EqualError(
		v.Validate(User{
			Name:   henge.ToStringPtr("🍺"),
			Age:    henge.ToIntPtr(1),
			Tags:   []string{""},
			Params: map[string]string{"": ""},
		}, valis.EachFields(tagrule.Validate)),
		`(too_short_length) .Name is too short length (minimum is 2 characters)
(gte) .Age must be greater than or equal to 2
(too_short_len) .Tags is too few elements (minimum is 2 elements)
(too_short_len) .Params is too few elements (minimum is 2 elements)`,
	)
	assert.NoError(
		v.Validate(User{
			Name:   henge.ToStringPtr("Alice"),
			Age:    henge.ToIntPtr(10),
			Tags:   []string{"", ""},
			Params: map[string]string{"a": "", "b": ""},
		}, valis.EachFields(tagrule.Validate)),
	)
}

func TestValidate_max(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name   *string           `validate:"max=2"`
		Age    *int64            `validate:"max=2"`
		Tags   []string          `validate:"max=2"`
		Params map[string]string `validate:"max=2"`
	}
	assert.NoError(
		v.Validate(User{}, valis.EachFields(tagrule.Validate)),
	)
	assert.EqualError(
		v.Validate(User{
			Name:   henge.ToStringPtr("abc"),
			Age:    henge.ToIntPtr(3),
			Tags:   []string{"a", "b", "c"},
			Params: map[string]string{"a": "", "b": "", "c": ""},
		}, valis.EachFields(tagrule.Validate)),
		`(too_long_length) .Name is too long length (maximum is 2 characters)
(lte) .Age must be less than or equal to 2
(too_long_len) .Tags is too many elements (maximum is 2 elements)
(too_long_len) .Params is too many elements (maximum is 2 elements)`,
	)
	assert.NoError(
		v.Validate(User{
			Name:   henge.ToStringPtr("🍺🍺"),
			Age:    henge.ToIntPtr(2),
			Tags:   []string{"", ""},
			Params: map[string]string{"a": "", "b": ""},
		}, valis.EachFields(tagrule.Validate)),
	)
}

func TestValidate_oneof(t *testing.T) {
	assert := assert.New(t)

	type Model struct {
		SP *string `validate:"oneof=a b c"`
		S  string  `validate:"oneof=a b c"`
		I  int     `validate:"oneof=1 2 3"`
		U  uint    `validate:"oneof=1 2 3"`
		B  bool    `validate:"oneof=true"`
		F  float64 `validate:"oneof=1 2.5 3"`
	}
	assert.EqualError(
		v.Validate(&Model{}, valis.EachFields(tagrule.Validate)),
		`(inclusion) .S is not included in [a b c]
(inclusion) .I is not included in [1 2 3]
(inclusion) .U is not included in [1 2 3]
(inclusion) .B is not included in [true]
(inclusion) .F is not included in [1 2.5 3]`,
	)
	assert.EqualError(
		v.Validate(&Model{SP: henge.ToStringPtr("")}, valis.EachFields(tagrule.Validate)),
		`(inclusion) .SP is not included in [a b c]
(inclusion) .S is not included in [a b c]
(inclusion) .I is not included in [1 2 3]
(inclusion) .U is not included in [1 2 3]
(inclusion) .B is not included in [true]
(inclusion) .F is not included in [1 2.5 3]`,
	)
	assert.NoError(
		v.Validate(&Model{
			SP: henge.ToStringPtr("a"),
			S:  "a",
			I:  2,
			U:  3,
			B:  true,
			F:  3,
		}, valis.EachFields(tagrule.Validate)),
	)
}

func TestValidate_url(t *testing.T) {
	assert := assert.New(t)

	type Model struct {
		U1 *string `validate:"url=http https"`
		U2 string  `validate:"url"`
		U3 string  `validate:"url=scp"`
	}
	assert.NoError(
		v.Validate(&Model{
			U1: henge.ToStringPtr("https://example.com/path?q=1"),
			U2: "http://example.com:1234/path?q=1",
			U3: "scp://example.com:8888",
		}, valis.EachFields(tagrule.Validate)),
	)
	assert.EqualError(
		v.Validate(&Model{
			U1: henge.ToStringPtr("rtmp://example.com/path?q=1"),
			U2: "http://example.com:1234/path?q=1",
			U3: "http://example.com:8888",
		}, valis.EachFields(tagrule.Validate)),
		`(invalid_scheme) .U1 which scheme is not included in [http https]
(invalid_scheme) .U3 which scheme is not included in [scp]`,
	)
}
