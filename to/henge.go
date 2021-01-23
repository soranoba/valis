// The package implements some valis.CombinationRule that verifies the converted value.
package to

import (
	"errors"
	"fmt"

	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
)

type (
	hengeError struct {
		err *henge.ConvertError
	}
)

var (
	// String is a valis.CombinationRule that validates after the value convert to string.
	String valis.CombinationRule = NewCombinationRule(WrapHengeResultFunc(toStringFunc))
	// Int is a valis.CombinationRule that validates after the value convert to int64.
	Int valis.CombinationRule = NewCombinationRule(WrapHengeResultFunc(toIntFunc))
	// Uint is a valis.CombinationRule that validates after the value convert to uint64.
	Uint valis.CombinationRule = NewCombinationRule(WrapHengeResultFunc(toUintFunc))
	// Float is a valis.CombinationRule that validates after the value convert to float64.
	Float valis.CombinationRule = NewCombinationRule(WrapHengeResultFunc(toFloatFunc))
	// Map is a valis.CombinationRule that validates after the value convert to map[interface{}]interface{}.
	Map valis.CombinationRule = NewCombinationRule(WrapHengeResultFunc(toMapFunc))
)

func (e *hengeError) Error() string {
	from, to := "", ""
	if e.err.SrcType != nil {
		from = fmt.Sprintf(" from %s", e.err.SrcType.String())
	}
	if e.err.DstType != nil {
		to = fmt.Sprintf(" to %s", e.err.DstType.String())
	}
	return fmt.Sprintf("can not convert%s%s", from, to)
}

func (e *hengeError) Unwrap() error {
	return e.err
}

func NewCombinationRule(convertFunc valis.ConvertFunc) valis.CombinationRule {
	return func(rules ...valis.Rule) valis.Rule {
		return valis.To(convertFunc, rules...)
	}
}

// WrapHengeResultFunc is a higher-order function that wraps errors of henge framework to readable errors on valis.
func WrapHengeResultFunc(f valis.ConvertFunc) valis.ConvertFunc {
	return func(value interface{}) (interface{}, error) {
		val, err := f(value)
		var convertErr *henge.ConvertError
		if errors.As(err, &convertErr) {
			return val, &hengeError{err: convertErr}
		}
		return val, err
	}
}

func toStringFunc(value interface{}) (interface{}, error) {
	return henge.New(value).String().Result()
}

func toIntFunc(value interface{}) (interface{}, error) {
	return henge.New(value).Int().Result()
}

func toUintFunc(value interface{}) (interface{}, error) {
	return henge.New(value).Uint().Result()
}

func toFloatFunc(value interface{}) (interface{}, error) {
	return henge.New(value).Float().Result()
}

func toMapFunc(value interface{}) (interface{}, error) {
	return henge.New(value).Map().Result()
}
