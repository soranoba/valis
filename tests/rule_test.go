package tests

import (
	"errors"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	v *valis.Validator
)

func init() {
	v = valis.NewValidator()
	v.SetCommonRules()
}

type Person1 struct {
	Name string
}

func (p Person1) Validate() error {
	if p.Name == "" {
		return errors.New("name cannot be blank")
	}
	return nil
}

type Person2 struct {
	Name string
}

func (p *Person2) Validate() error {
	if p.Name == "" {
		return errors.New("name cannot be blank")
	}
	return nil
}

type Person3 struct {
	Name string
}

func (p Person3) Validate(validator *valis.Validator) {
	if p.Name == "" {
		field := valishelpers.GetField(&p, &p.Name)
		loc := validator.Location().FieldLocation(field)
		validator.ErrorCollector().Add(loc, &valis.ErrorDetail{
			Code:  "required",
			Value: p.Name,
			Err:   errors.New("cannot be blank"),
		})
	}
}

type Person4 struct {
	Name *string
}

func (p *Person4) Validate(validator *valis.Validator) {
	if p.Name == nil {
		field := valishelpers.GetField(&p, &p.Name)
		loc := validator.Location().FieldLocation(field)
		validator.ErrorCollector().Add(loc, &valis.ErrorDetail{
			Code:  "required",
			Value: p.Name,
			Err:   errors.New("cannot be blank"),
		})
	}
}

func TestValidatableRule(t *testing.T) {
	var _ valis.Validatable = Person1{}
	var _ valis.Validatable = &Person2{}
	var _ valis.ValidatableWithValidator = Person3{}
	var _ valis.ValidatableWithValidator = &Person4{}

	assert := assert.New(t)

	assert.EqualError(v.Validate(Person1{}, valis.ValidatableRule), "(validatable) name cannot be blank, but got tests.Person1{Name:\"\"}")
	assert.EqualError(v.Validate(&Person1{}, valis.ValidatableRule), "(validatable) name cannot be blank, but got &tests.Person1{Name:\"\"}")

	// NOTE: the argument should be a pointer.
	assert.NoError(v.Validate(Person2{}, valis.ValidatableRule))
	assert.EqualError(v.Validate(&Person2{}, valis.ValidatableRule), "(validatable) name cannot be blank, but got &tests.Person2{Name:\"\"}")

	assert.EqualError(v.Validate(Person3{}, valis.ValidatableRule), "(required) .Name cannot be blank, but got \"\"")
	assert.EqualError(v.Validate(&Person3{}, valis.ValidatableRule), "(required) .Name cannot be blank, but got \"\"")

	// NOTE: the argument should be a pointer.
	assert.NoError(v.Validate(Person4{}, valis.ValidatableRule))
	assert.EqualError(v.Validate(&Person4{}, valis.ValidatableRule), "(required) .Name cannot be blank, but got (*string)(nil)")

	// NOTE: it will succeed if it does not implement Validatable
	assert.NoError(v.Validate(struct{}{}, valis.ValidatableRule))
	assert.NoError(v.Validate("", valis.ValidatableRule))
}
