package valistag

import (
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/soranoba/valis/when"
	"math"
	"reflect"
	"strconv"
	"strings"
)

var (
	// Required is the `required` tag rule.
	//   - `required:"${flag}"`: When flag is true, the value must be not empty.
	Required = valis.NewFieldTagRule("required", requiredRules)
	//   - `required`: equiv to is.Required
	//   - `max=${max}`: equiv to is.LengthBetween(0, max) when string, otherwise equiv to is.LenBetween(0, max)
	//   - `min=${min}`:
	Validate = valis.NewFieldTagRule("validate", validateRules)
)

var (
	validateTagSubKeys = map[string]valis.FieldTagRuleFunc{
		"required": func(v string) ([]valis.Rule, error) { // required
			return []valis.Rule{is.Required}, nil
		},
		"lte": func(v string) ([]valis.Rule, error) { // lte=${num}
			var num float64
			if _, err := SplitAndParseTagValues(v, "=", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{is.LessThanOrEqualTo(num)}, nil
		},
		"lt": func(v string) ([]valis.Rule, error) { // lt=${num}
			var num float64
			if _, err := SplitAndParseTagValues(v, "=", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{is.LessThan(num)}, nil
		},
		"gte": func(v string) ([]valis.Rule, error) { // gte=${num}
			var num float64
			if _, err := SplitAndParseTagValues(v, "=", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{is.GreaterThanOrEqualTo(num)}, nil
		},
		"gt": func(v string) ([]valis.Rule, error) { // gt=${num}
			var num float64
			if _, err := SplitAndParseTagValues(v, "=", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{is.GreaterThan(num)}, nil
		},
		"min": func(v string) ([]valis.Rule, error) { // min=${min}
			var min int
			if _, err := SplitAndParseTagValues(v, "=", &min); err != nil {
				return nil, err
			}
			return []valis.Rule{
				when.IsNumeric(is.Min(min)).
					ElseWhen(when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.LengthBetween(min, math.MaxInt64))).
					Else(is.LenBetween(min, math.MaxInt64)),
			}, nil
		},
		"max": func(v string) ([]valis.Rule, error) { // max=${max}
			var max int
			if _, err := SplitAndParseTagValues(v, "=", &max); err != nil {
				return nil, err
			}
			return []valis.Rule{
				when.IsNumeric(is.Min(max)).
					ElseWhen(when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.LengthBetween(0, max))).
					Else(is.LenBetween(0, max)),
			}, nil
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
