package valistag

import (
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"strconv"
	"strings"
)

var (
	// Required is the `required` tag rule.
	//   - `required:"${flag}"`: When flag is true, the value must be not empty.
	Required = valis.NewTagRule("required", requiredRules)
	Validate = valis.NewTagRule("validate", validateRules)
)

func requiredRules(tagValue string) []valis.Rule {
	ok, _ := strconv.ParseBool(tagValue)
	if ok {
		return []valis.Rule{is.Required}
	}
	return []valis.Rule{}
}

func validateRules(tagValue string) []valis.Rule {
	elems := strings.Split(tagValue, ",")
	rules := make([]valis.Rule, 0)

	for _, elem := range elems {
		if elem == "-" {
			return []valis.Rule{}
		} else if elem == "required" {
			rules = append(rules, is.Required)
		}
	}
	return rules
}