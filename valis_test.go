package valis_test

import (
	"fmt"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
)

func Example() {
	type User struct {
		Name string
		Age  int
	}

	u := &User{}
	if err := valis.Validate(
		&u,
		valis.Field(&u.Name, is.Required),
		valis.Field(&u.Age, is.Min(20)),
	); err != nil {
		fmt.Println(err)
	}

	u.Name = "Alice"
	u.Age = 20
	if err := valis.Validate(
		&u,
		valis.Field(&u.Name, is.Required),
		valis.Field(&u.Age, is.Min(20)),
	); err != nil {
		fmt.Println(err)
	}

	// Output:
	// (required) .Name is required
	// (gte) .Age must be greater than or equal to 20
}
