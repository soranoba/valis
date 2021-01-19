package valis

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

const (
	ConvertToCode = "convert_to"
)

// To returns a new rule that verifies the converted value met all rules.
func To(convertFunc ConvertFunc, rules ...Rule) Rule {
	return &toRule{convertFunc: convertFunc, rules: rules}
}

func (rule *toRule) Validate(validator *Validator, value interface{}) {
	newValue, err := rule.convertFunc(value)
	if err != nil {
		validator.ErrorCollector().Add(validator.Location(), &ErrorDetail{
			ConvertToCode,
			rule,
			value,
			nil,
			err,
		})
		return
	}

	errorCollector := newToRuleErrorCollector(validator.ErrorCollector(), validator.Location(), value)
	newValidator := validator.Clone(&CloneOpts{InheritLocation: true, ErrorCollector: errorCollector})
	And(rule.rules...).Validate(newValidator, newValue)
}

func newToRuleErrorCollector(errorCollector ErrorCollector, location Location, value interface{}) ErrorCollector {
	return &toRuleErrorCollector{ErrorCollector: errorCollector, loc: location, value: value}
}

func (c *toRuleErrorCollector) Add(loc Location, detail *ErrorDetail) {
	if loc == c.loc {
		newDetail := *detail
		newDetail.UnconvertedValue = c.value
		detail = &newDetail
	}
	c.ErrorCollector.Add(loc, detail)
}
