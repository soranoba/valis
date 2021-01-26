package valishelpers

import "fmt"

func ExampleGetField() {
	type User struct {
		Name string
		Age  int
	}
	u := User{}
	field := GetField(&u, &u.Name)
	fmt.Printf("Type of %s field is %s", field.Name, field.Type.String())

	// Output:
	// Type of Name field is string
}

func ExampleIsNumeric() {
	if IsNumeric(1.25) {
		fmt.Println("float is numeric")
	}
	if !IsNumeric("1.25") {
		fmt.Println("string is not numeric")
	}
	// Output:
	// float is numeric
	// string is not numeric
}

func ExampleIsNil() {
	if IsNil(nil) {
		fmt.Println("nil is nil")
	}
	if IsNil((*string)(nil)) {
		fmt.Println("(*string)(nil) is nil")
	}
	if !IsNil("") {
		fmt.Println("\"\" is not nil")
	}
	// Output:
	// nil is nil
	// (*string)(nil) is nil
	// "" is not nil
}
