package tests

import (
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate(map[string]interface{}{}, valis.Index(0, is.Any)),
		"(not_array) must be any array",
	)
	assert.EqualError(
		v.Validate(1, valis.Index(0, is.Any)),
		"(not_array) must be any array",
	)
	assert.EqualError(
		v.Validate([]int{}, valis.Index(0, is.Any)),
		"(out_of_range) requires more than 0 elements",
	)
	assert.EqualError(
		v.Validate([...]int{}, valis.Index(0, is.Any)),
		"(out_of_range) requires more than 0 elements",
	)
	assert.EqualError(
		v.Validate([]int{1, 5}, valis.Index(0, is.GreaterThan(2))),
		"(gt) [0] must be greater than 2",
	)
	assert.EqualError(
		v.Validate([...]int{1, 5}, valis.Index(0, is.GreaterThan(2))),
		"(gt) [0] must be greater than 2",
	)
	assert.NoError(
		v.Validate([]int{5}, valis.Index(0, is.GreaterThan(1))),
	)
	assert.NoError(
		v.Validate([...]int{5}, valis.Index(0, is.GreaterThan(1))),
	)
}

func TestIndexIfExist(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(
		v.Validate(map[string]interface{}{}, valis.IndexIfExist(0, is.Any)),
		"(not_array) must be any array",
	)
	assert.EqualError(
		v.Validate(1, valis.IndexIfExist(0, is.Any)),
		"(not_array) must be any array",
	)
	assert.NoError(
		v.Validate([]int{}, valis.IndexIfExist(0, is.Any)),
	)
	assert.NoError(
		v.Validate([...]int{}, valis.IndexIfExist(0, is.Any)),
	)
	assert.EqualError(
		v.Validate([]int{1, 5}, valis.IndexIfExist(0, is.GreaterThan(2))),
		"(gt) [0] must be greater than 2",
	)
	assert.EqualError(
		v.Validate([...]int{1, 5}, valis.IndexIfExist(0, is.GreaterThan(2))),
		"(gt) [0] must be greater than 2",
	)
	assert.NoError(
		v.Validate([]int{5}, valis.IndexIfExist(0, is.GreaterThan(1))),
	)
	assert.NoError(
		v.Validate([...]int{5}, valis.IndexIfExist(0, is.GreaterThan(1))),
	)
}
