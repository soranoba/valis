package valistag

import (
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"math"
	"strconv"
	"strings"
)

var (
	// Required is the `required` tag rule.
	//   - `required:"${flag}"`: When flag is true, the value must be not empty.
	Required = valis.NewTagRule("required", requiredRules)
	//   - `required`: equiv to is.Required
	//   - `max=${max}`: equiv to
	//   - `min=${min}`:
	Validate = valis.NewTagRule("validate", validateRules)
)

var (
	validateTagSubKeys = map[string]valis.TagRuleFunc{
		"required": func(v string) ([]valis.Rule, error) { // required
			return []valis.Rule{is.Required}, nil
		},
		"min": func(v string) ([]valis.Rule, error) { // min=${min}
			var min int
			if _, err := SplitAndParseTagValues(v, "=", &min); err != nil {
				return nil, err
			}
			return []valis.Rule{is.LengthRange(min, math.MaxInt64)}, nil
		},
		"max": func(v string) ([]valis.Rule, error) { // max=${max}
			var max int
			if _, err := SplitAndParseTagValues(v, "=", &max); err != nil {
				return nil, err
			}
			return []valis.Rule{is.LengthRange(0, max)}, nil
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
