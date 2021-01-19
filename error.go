package valis

import (
	"fmt"
	"strings"
)

type (
	// ValidationError is an error returned by Validator.Validate by default.
	ValidationError struct {
		details      []LocationAndErrorDetail
		nameResolver LocationNameResolver
		cacheMessage string
	}

	// ErrorDetail is the error content created by each Rule.
	ErrorDetail struct {
		// Code is the type code of the ErrorDetail.
		Code string
		// Rule is that created the ErrorDetail.
		// It is often nil.
		Rule Rule
		// Value is the verified value.
		Value interface{}
		// UnconvertedValue is a value before conversion when using To.
		// When not using To, it should be nil.
		// See: To
		UnconvertedValue interface{}
		// Err is an internal error.
		Err error
	}

	LocationAndErrorDetail struct {
		Location Location
		ErrorDetail *ErrorDetail
	}

	// ErrorCollector is an interface that receives some ErrorDetail of each rule and creates the error returned by Validator.Validate.
	ErrorCollector interface {
		HasError() bool
		Add(loc Location, detail *ErrorDetail)
		MakeError() error
	}
	ErrorCollectorFactoryFunc func() ErrorCollector
)

type (
	standardErrorCollector struct {
		nameResolver LocationNameResolver
		details     []LocationAndErrorDetail
	}

)

// NewValidationError returns a new ValidationError.
func NewValidationError(nameResolver LocationNameResolver, details []LocationAndErrorDetail) *ValidationError {
	return &ValidationError{details: details, nameResolver: nameResolver}
}

func (e *ValidationError) Details() []LocationAndErrorDetail {
	return e.details
}

func (e *ValidationError) Error() string {
	if e.cacheMessage != "" {
		return e.cacheMessage
	}

	messages := make([]string, 0)
	for _, locDetail := range e.details {
		var msg string

		name := e.nameResolver.ResolveLocationName(locDetail.Location)
		detail := locDetail.ErrorDetail
		if name == "" {
			msg = fmt.Sprintf("(%s) %s, but got %#v", detail.Code, detail.Err.Error(), detail.Value)
		} else {
			msg = fmt.Sprintf("(%s) %s %s, but got %#v", detail.Code, name, detail.Err.Error(), detail.Value)
		}
		if detail.UnconvertedValue != nil {
			msg += fmt.Sprintf(" (convert from %#v)", detail.UnconvertedValue)
		}
		messages = append(messages, msg)
	}
	e.cacheMessage = strings.Join(messages, ". ")
	return e.cacheMessage
}

// NewStandardErrorCollector returns an ErrorCollector used by default.
func NewStandardErrorCollector(nameResolver LocationNameResolver) ErrorCollector {
	return &standardErrorCollector{
		details:      make([]LocationAndErrorDetail, 0),
		nameResolver: nameResolver,
	}
}

func (c *standardErrorCollector) HasError() bool {
	return len(c.details) > 0
}

func (c *standardErrorCollector) Add(loc Location, detail *ErrorDetail) {
	c.details = append(c.details, LocationAndErrorDetail{loc, detail})
}

func (c *standardErrorCollector) MakeError() error {
	if c.HasError() {
		return NewValidationError(c.nameResolver, c.details)
	}
	return nil
}
