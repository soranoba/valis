package valis_test

import (
	"fmt"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
)

func ExampleIndex() {
	v := valis.NewValidator()
	if err := v.Validate([]int{4, 2, 1, 3}, valis.Index(2, is.GreaterThan(2))); err != nil {
		fmt.Println(err)
	}
	if err := v.Validate([...]int{4, 2, 1, 3}, valis.Index(3, is.LessThan(2))); err != nil {
		fmt.Println(err)
	}
	if err := v.Validate([...]int{4, 2, 1, 3}, valis.Index(5, is.LessThan(2))); err != nil {
		fmt.Println(err)
	}

	// Output:
	// (gt) [2] must be greater than 2
	// (lt) [3] must be less than 2
	// (out_of_range) requires more than 5 elements
}

func ExampleIndexIfExist() {
	v := valis.NewValidator()
	if err := v.Validate([]int{4, 2, 1, 3}, valis.IndexIfExist(2, is.GreaterThan(2))); err != nil {
		fmt.Println(err)
	}
	if err := v.Validate([...]int{4, 2, 1, 3}, valis.IndexIfExist(3, is.LessThan(2))); err != nil {
		fmt.Println(err)
	}
	if err := v.Validate([...]int{4, 2, 1, 3}, valis.IndexIfExist(5, is.LessThan(2))); err != nil {
		fmt.Println(err)
	}

	// Output:
	// (gt) [2] must be greater than 2
	// (lt) [3] must be less than 2
}
