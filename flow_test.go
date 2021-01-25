package valis_test

import (
	"fmt"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"reflect"
)

func ExampleWhen() {
	v := valis.NewValidator()

	isInt := func(ctx *valis.WhenContext) bool {
		return reflect.ValueOf(ctx.Value()).Kind() == reflect.Int
	}

	if err := v.Validate(0, valis.When(isInt, is.Required)); err != nil {
		fmt.Println(err)
	}
	if err := v.Validate("1", valis.When(isInt, is.Required).Else(is.In("a", "b", "c"))); err != nil {
		fmt.Println(err)
	}

	// Output:
	// (required) is required
	// (inclusion) is not included in [a b c]
}
