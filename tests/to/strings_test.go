package to_test

import (
	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/soranoba/valis/to"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplit(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(
		v.Validate("a,b,c", to.Split(",", valis.Each(is.In("a", "b", "c")))),
	)
	assert.EqualError(
		v.Validate("a,b,c", to.Split(",", valis.Each(is.In("a", "b")))),
		"(inclusion) [2] is not included in [a b], but got \"c\"",
	)

	// NOTE: to.Split can convert only from string or *string
	assert.NoError(
		v.Validate(henge.New("a,b,c").StringPtr().Value(), to.Split(",", valis.Each(is.In("a", "b", "c")))),
	)
	assert.EqualError(v.Validate(0, to.Split(",")), "(convert_to) cannot be split, but got 0")
}
