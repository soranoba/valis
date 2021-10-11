// Package code define error codes.
package code

// Type error codes.
const (
	NotString      = "not_string"
	NotStruct      = "not_struct"
	NotStructField = "not_struct_field"
	NotArray       = "not_array"
	NotMap         = "not_map"
	NotNumeric     = "not_numeric"
	NotIterable    = "not_iterable"
	NotAssignable  = "not_assignable" // %[1]s = TypeName
)

// Not found error codes.
const (
	NoKey      = "no_key"       // %[1]v = Key
	OutOfRange = "out_of_range" // %[1]d = Index, %[2]d = Length
)

// Conversion error codes.
const (
	ConversionFailed = "conversion" // %[1]w = Error
)

// Validation error codes.
const (
	Custom             = "custom" // %[1]w = Error
	Invalid            = "invalid"
	Required           = "required"
	NonZero            = "non_zero"
	NilOrNonZero       = "nil_or_non_zero"
	ZeroOnly           = "zero_only"
	TooLongLength      = "too_long_length"  // %[1]d = Count
	TooShortLength     = "too_short_length" // %[1]d = Count
	TooLongLen         = "too_long_len"     // %[1]d = Count
	TooShortLen        = "too_short_len"    // %[1]d = Count
	GreaterThan        = "gt"               // %[1]v = Number
	LessThan           = "lt"               // %[1]v = Number
	GreaterThanOrEqual = "gte"              // %[1]v = Number
	LessThanOrEqual    = "lte"              // %[1]v = Number
	Inclusion          = "inclusion"        // %[1]v = List
	RegexpMismatch     = "regexp"           // %[1]s = regexp
	InvalidURLFormat   = "invalid_url"      // %[1]v = Error
	InvalidScheme      = "invalid_scheme"   // %[1]v = List
	InvalidEmailFormat = "invalid_email"
)
