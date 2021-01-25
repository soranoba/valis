package tests

import (
	"errors"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TagValueHandler struct {
}

func (h *TagValueHandler) ParseTagValue(tagValue string) ([]valis.Rule, error) {
	if tagValue == "true" {
		return []valis.Rule{is.Required}, nil
	}
	return nil, errors.New("invalid required tag")
}

func TestFieldTagRule(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name string `required:"true"`
		Age  int    `required:"false"`
	}

	u := User{}
	requiredTagRule := valis.NewFieldTagRule("required", &TagValueHandler{})

	assert.EqualError(
		v.Validate(&u, requiredTagRule),
		"(not_struct) must be any struct",
	)
	assert.EqualError(
		v.Validate(&u, valis.Field(&u.Name, requiredTagRule)),
		"(required) .Name is required",
	)
	assert.Panics(func() {
		v.Validate(&u, valis.Field(&u.Age, requiredTagRule))
	})
}
