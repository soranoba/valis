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
		// When InheritLocation is true, Clone keeps the Location.
		InheritLocation bool
		// When InheritErrorCollector is true, Clone keeps the ErrorCollector.
		InheritErrorCollector bool
		// When Location is not nil and InheritLocation is false, Clone set the Location to the new Validator.
		Location Location
		// When ErrorCollector is not nil and InheritErrorCollector is false, Clone set the ErrorCollector to the new Validator.
		ErrorCollector ErrorCollector
	}
)

// NewValidator returns a new Validator.
func NewValidator() *Validator {
	v := &Validator{
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
	v.commonRules = append(v.commonRules, rules...)
}

// SetErrorCollectorFactoryFunc is update ErrorCollectorFactoryFunc.
func (v *Validator) SetErrorCollectorFactoryFunc(f ErrorCollectorFactoryFunc) {
	v.errorCollectorFactoryFunc = f
}

// Clone returns a new Validator inheriting the settings.
func (v *Validator) Clone(opts *CloneOpts) *Validator {
	newValidator := *v
	if opts.InheritErrorCollector {
		newValidator.errorCollector = v.errorCollector
	} else {
		newValidator.errorCollector = opts.ErrorCollector
	}
	if !opts.InheritLocation {
		if opts.Location != nil {
			newValidator.loc = opts.Location
		} else {
			newValidator.loc = NewLocation()
		}
	}
	return &newValidator
}

// Location returns a current location.
func (v *Validator) Location() Location {
	return v.loc
}

func (v *Validator) DiveField(field *reflect.StructField, f func(v *Validator)) {
	loc := v.loc
	v.loc = loc.FieldLocation(field)
	f(v)
	v.loc = loc
}

func (v *Validator) DiveIndex(index int, f func(v *Validator)) {
	loc := v.loc
	v.loc = loc.IndexLocation(index)
	f(v)
	v.loc = loc
}

func (v *Validator) DiveMapKey(key interface{}, f func(v *Validator)) {
	loc := v.loc
	v.loc = loc.MapKeyLocation(key)
	f(v)
	v.loc = loc
}

func (v *Validator) DiveMapValue(key interface{}, f func(v *Validator)) {
	loc := v.loc
	v.loc = loc.MapValueLocation(key)
	f(v)
	v.loc = loc
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
	And(rules...).Validate(newValidator, value)

	if newValidator.ErrorCollector().HasError() {
		return newValidator.ErrorCollector().MakeError()
	}
	return nil
}
