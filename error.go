package valis

import (
	"fmt"
	"strings"
)

type (
	ValidationError struct {
		details      map[Location]*ErrorDetail
		nameResolver LocationNameResolver
		cacheMessage string
	}

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
		details      map[Location]*ErrorDetail
	}
)

func NewValidationError(nameResolver LocationNameResolver, details map[Location]*ErrorDetail) *ValidationError {
	return &ValidationError{details: details, nameResolver: nameResolver}
}

func (e *ValidationError) Details() map[Location]*ErrorDetail {
	return e.details
}

func (e *ValidationError) Error() string {
	if e.cacheMessage != "" {
		return e.cacheMessage
	}

	messages := make([]string, 0)
	for loc, detail := range e.details {
		var msg string

		name := e.nameResolver.ResolveLocationName(loc)
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

func NewStandardErrorCollector(nameResolver LocationNameResolver) ErrorCollector {
	return &standardErrorCollector{
		details:      make(map[Location]*ErrorDetail),
		nameResolver: nameResolver,
	}
}

func (c *standardErrorCollector) HasError() bool {
	return len(c.details) > 0
}

func (c *standardErrorCollector) Add(loc Location, detail *ErrorDetail) {
	c.details[loc] = detail
}

func (c *standardErrorCollector) MakeError() error {
	if c.HasError() {
		return NewValidationError(c.nameResolver, c.details)
	}
	return nil
}
