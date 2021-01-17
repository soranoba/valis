package tests

import (
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator_Validate(t *testing.T) {
	assert := assert.New(t)

	assert.Error(valis.NewValidator().Validate("", is.Required))
	assert.NoError(valis.NewValidator().Validate("a"), is.Required)

	// NOTE: it returns an error if any rules are not met.
	assert.Error(valis.NewValidator().Validate("", is.Any, is.Required))
	assert.NoError(valis.NewValidator().Validate("a", is.Any, is.Required))
}
