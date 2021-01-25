// The package implements some valis.Rule related to field tag.
package valistag

import (
	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/code"
	"github.com/soranoba/valis/is"
	"github.com/soranoba/valis/when"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type (
	oneofRule struct {
		elems []string
	}
)

var (
	// Required is a `required` tag rule.
	Required = valis.NewFieldTagRule("required", requiredRules)
	// Validate is a `validate` tag rule.
	Validate = valis.NewFieldTagRule("validate", validateRules)
)

var (
	regexpCache     = map[string]*regexp.Regexp{}
	regexpCacheLock sync.RWMutex

	validateTagSubKeys = map[string]func(string) ([]valis.Rule, error){
		"required": func(v string) ([]valis.Rule, error) { // required
			return []valis.Rule{is.Required}, nil
		},
		"lte": func(v string) ([]valis.Rule, error) { // lte=10
			var num float64
			if _, err := SplitAndParseTagValues(v, " ", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{is.LessThanOrEqualTo(num)}, nil
		},
		"lt": func(v string) ([]valis.Rule, error) { // lt=10
			var num float64
			if _, err := SplitAndParseTagValues(v, " ", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{is.LessThan(num)}, nil
		},
		"gte": func(v string) ([]valis.Rule, error) { // gte=10
			var num float64
			if _, err := SplitAndParseTagValues(v, " ", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{is.GreaterThanOrEqualTo(num)}, nil
		},
		"gt": func(v string) ([]valis.Rule, error) { // gt=10
			var num float64
			if _, err := SplitAndParseTagValues(v, " ", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{is.GreaterThan(num)}, nil
		},
		"min": func(v string) ([]valis.Rule, error) { // min=1
			var min int
			if _, err := SplitAndParseTagValues(v, " ", &min); err != nil {
				return nil, err
			}
			return []valis.Rule{
				when.IsNumeric(is.Min(min)).
					ElseWhen(when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.LengthBetween(min, math.MaxInt64))).
					Else(is.LenBetween(min, math.MaxInt64)),
			}, nil
		},
		"max": func(v string) ([]valis.Rule, error) { // max=10
			var max int
			if _, err := SplitAndParseTagValues(v, " ", &max); err != nil {
				return nil, err
			}
			return []valis.Rule{
				when.IsNumeric(is.Min(max)).
					ElseWhen(when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.LengthBetween(0, max))).
					Else(is.LenBetween(0, max)),
			}, nil
		},
		"oneof": func(v string) ([]valis.Rule, error) { // oneof=1 2
			if v == "" {
				return nil, errInsufficientNumberOfTagParameters
			}
			elems := strings.Split(v, " ")
			return []valis.Rule{&oneofRule{elems: elems}}, nil
		},
		"regexp": func(v string) ([]valis.Rule, error) { // regexp=^[0-9]+$
			if v == "" {
				return nil, errInsufficientNumberOfTagParameters
			}

			regexpCacheLock.RLock()
			re, ok := regexpCache[v]
			regexpCacheLock.RUnlock()

			if !ok {
				re = regexp.MustCompile(v)
				regexpCacheLock.Lock()
				regexpCache[v] = re
				regexpCacheLock.Unlock()
			}
			return []valis.Rule{is.Match(re)}, nil
		},
	}
)

func requiredRules(tagValue string) ([]valis.Rule, error) {
	ok, _ := strconv.ParseBool(tagValue)
	if ok {
		return []valis.Rule{is.Required}, nil
	}
	return []valis.Rule{}, nil
}

func validateRules(tagValue string) ([]valis.Rule, error) {
	elems := strings.Split(tagValue, ",")
	rules := make([]valis.Rule, 0)

	for _, elem := range elems {
		if elem == "-" {
			return []valis.Rule{}, nil
		}

		subKv := strings.SplitN(elem, "=", 2)
		if f, ok := validateTagSubKeys[subKv[0]]; ok {
			subKey := ""
			if len(subKv) == 2 {
				subKey = subKv[1]
			}
			r, err := f(subKey)
			if err != nil {
				return nil, err
			}
			rules = append(rules, r...)
		}
	}
	return rules, nil
}

func (rule *oneofRule) Validate(validator *valis.Validator, value interface{}) {
	val, err := henge.New(value).String().Result()
	if err != nil {
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.ConversionFailed, value, err))
	}

	for _, elem := range rule.elems {
		if elem == val {
			return
		}
	}
	validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.Inclusion, value, rule.elems))
}
