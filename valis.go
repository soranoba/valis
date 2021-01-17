package valis

var (
	standardValidator = NewValidator()
)

// Validate the value using the StandardValidator.
// See Validator.Validate
func Validate(value interface{}, rules ...Rule) error {
	return standardValidator.Clone(&CloneOpts{}).Validate(value, rules...)
}

// AddCommonRules add the rules to common rules of the StandardValidator.
// See Validator.AddCommonRules
func AddCommonRules(rules ...Rule) {
	standardValidator.AddCommonRules(rules...)
}

// SetErrorCollectorFactoryFunc is update ErrorCollectorFactoryFunc of the StandardValidator.
// See Validator.SetErrorCollectorFactoryFunc
func SetErrorCollectorFactoryFunc(f ErrorCollectorFactoryFunc) {
	standardValidator.SetErrorCollectorFactoryFunc(f)
}
