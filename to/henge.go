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
	String      = valis.To(HengeResultConvertFunc(toStringFunc))
	Int         = valis.To(HengeResultConvertFunc(toIntFunc))
	Uint        = valis.To(HengeResultConvertFunc(toUintFunc))
	Float       = valis.To(HengeResultConvertFunc(toFloatFunc))
	Map         = valis.To(HengeResultConvertFunc(toMapFunc))
	StringSlice = valis.To(HengeResultConvertFunc(toStringSliceFunc))
	IntSlice    = valis.To(HengeResultConvertFunc(toIntSliceFunc))
	UintSlice   = valis.To(HengeResultConvertFunc(toUintSliceFunc))
	FloatSlice  = valis.To(HengeResultConvertFunc(toFloatSliceFunc))
)

func (e *hengeError) Error() string {
	return fmt.Sprintf("can not convert from %s to %s", e.err.SrcType.String(), e.err.DstType.String())
}

func (e *hengeError) Unwrap() error {
	return e.err
}

// HengeResultConvertFunc convert to normalized error from *henge.ConvertError.
// For example:
//
//    valis.To(HengeResultConvertFunc(func(value interface{}) (interface{}, error) {
//        return henge.New(value).Int().Result()
//    }))
//
func HengeResultConvertFunc(f func(interface{}) (interface{}, error)) func(interface{}) (interface{}, error) {
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
