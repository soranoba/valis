package valis

import (
	"bytes"
	"fmt"
	"github.com/soranoba/valis/translations"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"sync"
)

type (
	// ValidationError is an error returned by Validator.Validate by default.
	ValidationError struct {
		errors       []*LocationError
		nameResolver LocationNameResolver
	}

	Error = *errorDetail

	LocationError struct {
		Error
		Location *Location
	}

	// ErrorCollector is an interface that receives some Error of each rule and creates the error returned by Validator.Validate.
	ErrorCollector interface {
		HasError() bool
		Add(loc *Location, err Error)
		MakeError() error
	}
	ErrorCollectorFactoryFunc func() ErrorCollector
)

type (
	standardErrorCollector struct {
		nameResolver LocationNameResolver
		errors       []*LocationError
	}
	errorDetail struct {
		code                  string
		params                []interface{}
		value                 interface{}
		valueBeforeConversion interface{}
	}
)

var (
	enCatalogOnce sync.Once
	enCatalog     *translations.Catalog
)

// NewValidationError returns a new ValidationError.
func NewValidationError(nameResolver LocationNameResolver, errors []*LocationError) *ValidationError {
	return &ValidationError{errors: errors, nameResolver: nameResolver}
}

func (e *ValidationError) Details() []*LocationError {
	return e.errors
}

func (e *ValidationError) Translate(p *message.Printer) map[string][]string {
	trans := make(map[string][]string)
	for _, locErr := range e.errors {
		key := e.nameResolver.ResolveLocationName(locErr.Location)
		trans[key] = append(trans[key], p.Sprintf(locErr.Error.Code(), locErr.Error.Params()...))
	}
	return trans
}

func (e *ValidationError) Error() string {
	enCatalogOnce.Do(func() {
		enCatalog = translations.NewCatalog()
		enCatalog.Set(translations.DefaultEnglish)
	})
	p := message.NewPrinter(language.English, message.Catalog(enCatalog))

	buf := bytes.NewBuffer(nil)
	for i, locErr := range e.errors {
		name := e.nameResolver.ResolveLocationName(locErr.Location)
		buf.WriteString(fmt.Sprintf("(%s) ", locErr.Code()))
		if name != "" {
			buf.WriteString(name)
			buf.WriteString(" ")
		}
		buf.WriteString(p.Sprintf(locErr.Error.Code(), locErr.Error.Params()...))
		if len(e.errors)-1 != i {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

// NewStandardErrorCollector returns an ErrorCollector used by default.
func NewStandardErrorCollector(nameResolver LocationNameResolver) ErrorCollector {
	return &standardErrorCollector{
		errors:       make([]*LocationError, 0),
		nameResolver: nameResolver,
	}
}

func (c *standardErrorCollector) HasError() bool {
	return len(c.errors) > 0
}

func (c *standardErrorCollector) Add(loc *Location, err Error) {
	c.errors = append(c.errors, &LocationError{
		Location: loc,
		Error:    err,
	})
}

func (c *standardErrorCollector) MakeError() error {
	if c.HasError() {
		return NewValidationError(c.nameResolver, c.errors)
	}
	return nil
}

func NewError(code string, value interface{}, params ...interface{}) Error {
	return &errorDetail{
		code:   code,
		params: params,
		value:  value,
	}
}

// Code returns an error code.
// See also code sub-package.
func (e *errorDetail) Code() string {
	return e.code
}

// Params returns translation parameters.
// See translations sub-package for details.
func (e *errorDetail) Params() []interface{} {
	return e.params
}

// Value returns a value received by Rule that created the error.
func (e *errorDetail) Value() interface{} {
	return e.value
}

// ValueBeforeConversion returns the original value at that location.
// If To is not used, the same value as Value is returned.
func (e *errorDetail) ValueBeforeConversion() interface{} {
	if e.valueBeforeConversion == nil {
		return e.value
	}
	return e.valueBeforeConversion
}

func (e *errorDetail) Error() string {
	return e.code
}
