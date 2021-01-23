package valis_test

import (
	"errors"
	"fmt"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
)

func ExampleNewFieldTagRule() {
	type User struct {
		Name string `required:"true"`
		Age  int    `required:"true"`
	}

	v := valis.NewValidator()
	u := User{}
	requiredTagRule := valis.NewFieldTagRule("required", func(tagValue string) ([]valis.Rule, error) {
		if tagValue == "true" {
			return []valis.Rule{is.Required}, nil
		}
		return nil, errors.New("invalid required tag")
	})

	if err := v.Validate(
		&u,
		valis.Field(&u.Name, requiredTagRule),
		valis.Field(&u.Age, requiredTagRule),
	); err != nil {
		fmt.Println(err)
	}

	// Output:
	// (required) .Name is required
	// (required) .Age is required
}
