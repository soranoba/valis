package when_test

import (
	"fmt"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/soranoba/valis/when"
)

func ExampleHasFieldTag() {
	type User struct {
		Name string `required:"true"`
	}

	v := valis.NewValidator()
	u := User{}
	if err := v.Validate(
		&u,
		valis.Field(&u.Name, when.HasFieldTag("required", is.Required)),
	); err != nil {
		fmt.Println(err)
	}

	// Output:
	// (required) .Name is required
}
