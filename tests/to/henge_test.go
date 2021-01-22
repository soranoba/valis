package to_test

import (
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/soranoba/valis/to"
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

func TestString(t *testing.T) {
	assert := assert.New(t)

	assert.EqualError(v.Validate(nil, to.String(is.Required)), "(conversion) can not convert to string")
	assert.EqualError(v.Validate("", to.String(is.Required)), "(required) is required")
	assert.NoError(v.Validate(0, to.String(is.Required, is.In("0"))))
}

func TestInt(t *testing.T) {
	assert := assert.New(t)

	// NOTE: to.Int convert to int64.
	assert.NoError(v.Validate("0", to.Int(is.In(int64(0)))))
	assert.EqualError(v.Validate("0", to.Int(is.In(0))), "(inclusion) is not included in [0]")

	assert.NoError(v.Validate("1234", to.Int(is.In(int64(1234)))))
	assert.NoError(v.Validate(1.25, to.Int(is.In(int64(1)))))
	assert.EqualError(v.Validate("aaaa", to.Int(is.Any)), "(conversion) can not convert from string to int64")
}

func TestUint(t *testing.T) {
	assert := assert.New(t)

	// NOTE: to.Uint convert to uint64.
	assert.NoError(v.Validate("0", to.Uint(is.In(uint64(0)))))
	assert.EqualError(v.Validate("0", to.Uint(is.In(0))), "(inclusion) is not included in [0]")

	assert.NoError(v.Validate("1234", to.Uint(is.In(uint64(1234)))))
	assert.NoError(v.Validate(1.25, to.Uint(is.In(uint64(1)))))
	assert.EqualError(v.Validate(-1, to.Uint(is.Required)), "(conversion) can not convert from int to uint64")
}

func TestFloat(t *testing.T) {
	assert := assert.New(t)

	assert.NoError(v.Validate("1.25", to.Float(is.In(1.25))))
	assert.EqualError(v.Validate("0", to.Float(is.In(0))), "(inclusion) is not included in [0]")
	assert.EqualError(v.Validate("aaa", to.Float(is.Required)), "(conversion) can not convert from string to float64")
}

func TestMap(t *testing.T) {
}

func TestStringSlice(t *testing.T) {

}

func TestIntSlice(t *testing.T) {

}

func TestUintSlice(t *testing.T) {

}

func TestFloatSlice(t *testing.T) {

}
