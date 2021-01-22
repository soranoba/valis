package code

// type error.
const (
	StringOnly    = "string_only"
	StructOnly    = "struct_only"
	ArrayOnly     = "array_only"
	MapOnly       = "map_only"
	NotIterable   = "not_iterable"
	NotAssignable = "not_assignable" // %[1]s = TypeName
)

// not found error.
const (
	NoKey = "no_key" // %[1]v = Key
)

// convert error.
const (
	ConversionFailed = "conversion" // %[1]w = Error
)

const (
	Custom         = "custom" // %[1]w = Error
	Invalid        = "invalid"
	Required       = "required"
	NilOrNonZero   = "nil_or_non_zero"
	ZeroOnly       = "zero_only"
	TooLongLength  = "too_long_length"  // %[1]d = Count
	TooShortLength = "too_short_length" // %[1]d = Count
	TooLongLen     = "too_long_len"     // %[1]d = Count
	TooShortLen    = "too_short_len"    // %[1]d = Count
	Inclusion      = "inclusion"        // %[1]v = List
)
