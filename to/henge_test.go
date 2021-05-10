package to

import (
	"fmt"
	"testing"

	"github.com/soranoba/henge/v2"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
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
	// (inclusion) is not included in [123]
}
