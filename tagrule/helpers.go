package tagrule

import (
	"errors"
	"github.com/soranoba/henge"
	"strings"
)

var (
	errInsufficientNumberOfTagParameters = errors.New("insufficient number of tag parameters")
)

func SplitAndParseTagValues(s string, sep string, outs ...interface{}) (count int, err error) {
	for _, elem := range strings.Split(s, sep) {
		if elem != "" {
			if len(outs) <= count {
				return count, errInsufficientNumberOfTagParameters
			}
			if err := henge.New(elem).Convert(outs[count]); err != nil {
				return count, err
			}
		}
		count += 1
	}
	return count, nil
}
