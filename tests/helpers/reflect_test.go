package helpers_test

import (
	valishelpers "github.com/soranoba/valis/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsNumeric(t *testing.T) {
	assert := assert.New(t)

	assert.True(valishelpers.IsNumeric(1))
	assert.True(valishelpers.IsNumeric(uint(1)))
	assert.True(valishelpers.IsNumeric(1.25))

	assert.False(valishelpers.IsNumeric((*int)(nil)))
	assert.False(valishelpers.IsNumeric("1.25"))
}
