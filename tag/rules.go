package valistag

import (
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"strconv"
)

var (
	Required = valis.NewTagRule("required", requiredRules)
)

func requiredRules(tagValue string) []valis.Rule {
	ok, _ := strconv.ParseBool(tagValue)
	if ok {
		return []valis.Rule{is.Required}
	}
	return []valis.Rule{}
}
