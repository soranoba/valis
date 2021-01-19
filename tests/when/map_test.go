package when_test

import (
	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/soranoba/valis/when"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var (
	v *valis.Validator
)

func init() {
	v = valis.NewValidator()
	v.SetCommonRules()
}

func TestHasKey(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate(map[string]string{}, valis.Key("a", is.Required)),
		"(not_found) does not have the key (a), but got map[string]string{}",
	)
	assert.NoError(
		v.Validate(map[string]string{}, when.HasKey("a", valis.Key("a", is.Required))),
	)

	assert.EqualError(
		v.Validate(&map[string]string{}, valis.Key("a", is.Required)),
		"(not_found) does not have the key (a), but got &map[string]string{}",
	)
	assert.NoError(
		v.Validate(&map[string]string{}, when.HasKey("a", valis.Key("a", is.Required))),
	)

	// NOTE: If the key type is a pointer, it only matches when it specifies the same pointer
	a, b, c := "a", "b", "c"
	m := map[*string]string{&a: "A", &b: "B", &c: "C"}
	assert.NoError(
		v.Validate(&m, when.HasKey(henge.New("a").StringPtr().Value(), is.Zero)),
	)
	if err := v.Validate(&m, when.HasKey(&a, is.Zero)); assert.Error(err) {
		assert.True(strings.HasPrefix(err.Error(), "(zero) must be nil or zero, but got"), err.Error())
	}

	// NOTE: cond.Func of HasKey returns false, when the value is not a map.
	assert.NoError(
		v.Validate("", when.HasKey("a", is.Required)),
	)
}
