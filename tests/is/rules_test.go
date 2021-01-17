package is_test

import (
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
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

	assert.EqualError(v.Validate("", is.Required), "(required) cannot be blank, but got \"\"")
	assert.EqualError(v.Validate(0, is.Required), "(required) cannot be blank, but got 0")
	assert.EqualError(v.Validate(struct{}{}, is.Required), "(required) cannot be blank, but got struct {}{}")
	assert.EqualError(v.Validate(nil, is.Required), "(required) cannot be blank, but got <nil>")

	assert.NoError(v.Validate("a", is.Required))
	assert.NoError(v.Validate(1, is.Required))
	assert.NoError(v.Validate(struct{ Name string }{Name: "name"}, is.Required))
}

func TestAny(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(v.Validate("", is.Any))
	assert.NoError(v.Validate(nil, is.Any))
}

func TestIn(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(v.Validate("b", is.In("a", "b", "c")))
	assert.EqualError(v.Validate("d", is.In("a", "b", "c")), "(inclusion) is not included in [a b c], but got \"d\"")
	assert.NoError(v.Validate(2, is.In(1, 2, 3)))
	assert.EqualError(v.Validate(5, is.In(1, 2, 3)), "(inclusion) is not included in [1 2 3], but got 5")

	// NOTE: Does not match if the types are different
	assert.EqualError(v.Validate(int64(1), is.In(1, 2, 3)), "(inclusion) is not included in [1 2 3], but got 1")

	// NOTE: For pointers, it compares Elem values
	i := 2
	assert.NoError(v.Validate(&i, is.In(1, 2, 3)))
	i = 5
	assert.Error(v.Validate(&i, is.In(1, 2, 3)))

	i = 2
	x, y, z := 1, 2, 3
	assert.NoError(v.Validate(&i, is.In(&x, &y, &z)))
	i = 5
	assert.Error(v.Validate(&i, is.In(&x, &y, &z)))
}
