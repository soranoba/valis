package when_test

import (
	"github.com/soranoba/valis/is"
	"github.com/soranoba/valis/when"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsSliceOrArray(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name string
	}
	assert.EqualError(
		v.Validate("", is.NonZero),
		"(non_zero) can't be blank (or zero)",
	)
	assert.NoError(
		v.Validate("", when.IsSliceOrArray(is.NonZero)),
	)

	assert.EqualError(
		v.Validate([3]User{}, when.IsSliceOrArray(is.NonZero)),
		"(non_zero) can't be blank (or zero)",
	)
	assert.EqualError(
		v.Validate(([]User)(nil), when.IsSliceOrArray(is.NonZero)),
		"(non_zero) can't be blank (or zero)",
	)
	assert.NoError(
		v.Validate([]User{}, when.IsSliceOrArray(is.NonZero)),
	)
}
