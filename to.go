package valis

import "github.com/soranoba/valis/code"

type (
	ConvertFunc func(value interface{}) (interface{}, error)
)

type (
	toRule struct {
		convertFunc ConvertFunc
		rules       []Rule
	}
	toRuleErrorCollector struct {
		ErrorCollector
		loc   Location
		value interface{}
	}
)

// To returns a new rule that verifies the converted value met all rules.
func To(convertFunc ConvertFunc, rules ...Rule) Rule {
	return &toRule{convertFunc: convertFunc, rules: rules}
}

func (rule *toRule) Validate(validator *Validator, value interface{}) {
	newValue, err := rule.convertFunc(value)
	if err != nil {
		validator.ErrorCollector().Add(validator.Location(), NewError(code.ConversionFailed, value, err))
		return
	}

	errorCollector := newToRuleErrorCollector(validator.ErrorCollector(), validator.Location(), value)
	newValidator := validator.Clone(&CloneOpts{InheritLocation: true, ErrorCollector: errorCollector})
	And(rule.rules...).Validate(newValidator, newValue)
}

func newToRuleErrorCollector(errorCollector ErrorCollector, location Location, value interface{}) ErrorCollector {
	return &toRuleErrorCollector{ErrorCollector: errorCollector, loc: location, value: value}
}

func (c *toRuleErrorCollector) Add(loc Location, err Error) {
	if loc == c.loc {
		newDetail := *err
		newDetail.valueBeforeConversion = c.value
		err = &newDetail
	}
	c.ErrorCollector.Add(loc, err)
}
