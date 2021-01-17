package valis

import "reflect"

type (
	// Validator provides validation methods.
	// And, each rule uses a validator to save the error details.
	Validator struct {
		commonRules               []Rule
		errorCollectorFactoryFunc ErrorCollectorFactoryFunc

		loc            Location
		errorCollector ErrorCollector
	}
	// CloneOpts is an option of Clone.
	CloneOpts struct {
		// When KeepLocation is true, Clone keeps the Location.
		KeepLocation bool
		// When KeepErrorCollector is true, Clone keeps the ErrorCollector.
		KeepErrorCollector bool
		// When ErrorCollector is not nil and KeepErrorCollector is false, Clone set the ErrorCollector to the new Validator.
		ErrorCollector ErrorCollector
	}
)

// NewValidator returns a new Validator.
func NewValidator() *Validator {
	v := &Validator{
		commonRules: StandardRules[:],
		errorCollectorFactoryFunc: func() ErrorCollector {
			return NewStandardErrorCollector(DefaultLocationNameResolver)
		},
		loc: NewLocation(),
	}
	return v
}

// SetCommonRules is update common rules.
func (v *Validator) SetCommonRules(rules ...Rule) {
	v.commonRules = rules
}

// AddCommonRules add the rules to common rules.
func (v *Validator) AddCommonRules(rules ...Rule) {
	commonRules := make([]Rule, len(v.commonRules)+len(rules))
	copy(commonRules, v.commonRules)
	copy(commonRules[len(v.commonRules):], rules)
	v.commonRules = commonRules
}

// SetErrorCollectorFactoryFunc is update ErrorCollectorFactoryFunc.
func (v *Validator) SetErrorCollectorFactoryFunc(f ErrorCollectorFactoryFunc) {
	v.errorCollectorFactoryFunc = f
}

// Clone returns a new Validator inheriting the settings.
func (v *Validator) Clone(opts *CloneOpts) *Validator {
	newValidator := *v
	if !opts.KeepErrorCollector {
		if opts.ErrorCollector != nil {
			newValidator.errorCollector = opts.ErrorCollector
		} else {
			newValidator.errorCollector = nil
		}
	}
	if !opts.KeepLocation {
		newValidator.loc = NewLocation()
	}
	return &newValidator
}

// Location returns a current location.
func (v *Validator) Location() Location {
	return v.loc
}

// WithLocation returns a new Validator with the location.
func (v *Validator) WithLocation(loc Location) *Validator {
	newValidator := *v
	newValidator.loc = loc
	return &newValidator
}

// WithField is equiv to v.WithLocation(v.Location().FieldLocation(field))
func (v *Validator) WithField(field reflect.StructField) *Validator {
	return v.WithLocation(v.Location().FieldLocation(field))
}

// WithIndex is equiv to v.WithLocation(v.Location().IndexLocation(index))
func (v *Validator) WithIndex(index int) *Validator {
	return v.WithLocation(v.Location().IndexLocation(index))
}

// WithKey is equiv to v.WithLocation(v.Location().KeyLocation(key))
func (v *Validator) WithKey(key interface{}) *Validator {
	return v.WithLocation(v.Location().KeyLocation(key))
}

// ErrorCollector returns an ErrorCollector.
func (v *Validator) ErrorCollector() ErrorCollector {
	if v.errorCollector == nil {
		v.errorCollector = v.errorCollectorFactoryFunc()

		if v.errorCollector == nil {
			panic("failed to create an ErrorCollector")
		}
	}
	return v.errorCollector
}

// Validate the value.
// It returns an error if any rules are not met.
func (v *Validator) Validate(value interface{}, rules ...Rule) error {
	newValidator := v.Clone(&CloneOpts{})
	for _, rule := range append(v.commonRules, rules...) {
		rule.Validate(newValidator, value)
	}
	if newValidator.ErrorCollector().HasError() {
		return newValidator.ErrorCollector().MakeError()
	}
	return nil
}
