package to

import (
	"fmt"
	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"testing"
)

func TestHengeError(t *testing.T) {
	var _ error = &hengeError{}
}

func ExampleWrapHengeResultFunc() {
	toString := NewCombinationRule(WrapHengeResultFunc(func(i interface{}) (interface{}, error) {
		return henge.New(i).String().Result()
	}))
	// Valid
	if err := valis.Validate(123, toString(is.In("123"))); err != nil {
		fmt.Println(err)
	}
	// Invalid
	if err := valis.Validate(123, toString(is.In(123))); err != nil {
		fmt.Println(err)
	}

	// Output:
	// (inclusion) is not included in [123], but got "123" (convert from 123)
}
