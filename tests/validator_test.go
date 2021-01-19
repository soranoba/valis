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

func TestValidator_AddCommonRules(t *testing.T) {
	assert := assert.New(t)

	v1 := valis.NewValidator()
	v1.SetCommonRules(is.Required)
	v1.AddCommonRules(is.In("aaa"))

	assert.Error(v1.Validate("abc"))
	assert.NoError(v1.Validate("aaa"))

	v1.AddCommonRules(is.In("bbb"))
	assert.Error(v1.Validate("aaa"))
}

func TestValidator_SetCommonRules(t *testing.T) {
	assert := assert.New(t)

	v1 := valis.NewValidator()
	v1.SetCommonRules(is.Required)

	// NOTE: CommonRules automatically check.
	assert.Error(v1.Validate(""))
	assert.NoError(v1.Validate("a"))

	// NOTE: Clone copy common rules.
	v2 := v1.Clone(&valis.CloneOpts{})
	assert.Error(v2.Validate(""))

	v1.SetCommonRules()
	assert.NoError(v1.Validate(""))
	assert.Error(v2.Validate(""))

	v1.SetCommonRules(is.Required)
	v2.SetCommonRules()
	assert.Error(v1.Validate(""))
	assert.NoError(v2.Validate(""))
}
