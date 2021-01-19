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
	// StringSlice is a valis.CombinationRule that validates after the value convert to []string.
	StringSlice valis.CombinationRule = NewCombinationRule(WrapHengeResultFunc(toStringSliceFunc))
	// IntSlice is a valis.CombinationRule that validates after the value convert to []int64.
	IntSlice valis.CombinationRule = NewCombinationRule(WrapHengeResultFunc(toIntSliceFunc))
	// UintSlice is a valis.CombinationRule that validates after the value convert to []uint64.
	UintSlice valis.CombinationRule = NewCombinationRule(WrapHengeResultFunc(toUintSliceFunc))
	// FloatSlice is a valis.CombinationRule that validates after the value convert to []float64.
	FloatSlice valis.CombinationRule = NewCombinationRule(WrapHengeResultFunc(toFloatSliceFunc))
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
func WrapHengeResultFunc(f func(interface{}) (interface{}, error)) func(interface{}) (interface{}, error) {
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

func toStringSliceFunc(value interface{}) (interface{}, error) {
	return henge.New(value).StringSlice().Result()
}

func toIntSliceFunc(value interface{}) (interface{}, error) {
	return henge.New(value).IntSlice().Result()
}

func toUintSliceFunc(value interface{}) (interface{}, error) {
	return henge.New(value).UintSlice().Result()
}

func toFloatSliceFunc(value interface{}) (interface{}, error) {
	return henge.New(value).FloatSlice().Result()
}
