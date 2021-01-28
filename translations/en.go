package translations

import (
	"github.com/soranoba/valis/code"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

// DefaultEnglish is a CatalogRegistrationFunc for English.
func DefaultEnglish(c *catalog.Builder) {
	tag := language.English

	// type error
	c.Set(tag, code.NotString, catalog.String("must be any string"))
	c.Set(tag, code.NotStruct, catalog.String("must be any struct"))
	c.Set(tag, code.NotStructField, catalog.String("must be any struct field"))
	c.Set(tag, code.NotArray, catalog.String("must be any array"))
	c.Set(tag, code.NotMap, catalog.String("must be any map"))
	c.Set(tag, code.NotAssignable, catalog.String("can't assign to %[1]s"))

	// not found error
	c.Set(tag, code.NoKey, catalog.String("requires the value at the key (%[1]v)"))

	// convert error
	c.Set(tag, code.ConversionFailed, catalog.String("%[1]v"))

	// others
	c.Set(tag, code.Custom, catalog.String("%[1]v"))
	c.Set(tag, code.Invalid, catalog.String("is invalid"))
	c.Set(tag, code.Required, catalog.String("is required"))
	c.Set(tag, code.NilOrNonZero, catalog.String("can't be blank (or zero) if specified"))
	c.Set(tag, code.ZeroOnly, catalog.String("must be blank"))

	c.Set(tag, code.TooLongLength,
		catalog.Var("characters", plural.Selectf(1, "", plural.One, "character", plural.Other, "characters")),
		catalog.String("is too long length (maximum is %[1]d ${characters})"),
	)
	c.Set(tag, code.TooShortLength,
		catalog.Var("characters", plural.Selectf(1, "", plural.One, "character", plural.Other, "characters")),
		catalog.String("is too short length (minimum is %[1]d ${characters})"),
	)
	c.Set(tag, code.TooLongLen,
		catalog.Var("elements", plural.Selectf(1, "", plural.One, "element", plural.Other, "elements")),
		catalog.String("is too many elements (maximum is %[1]d ${elements})"),
	)
	c.Set(tag, code.TooShortLen,
		catalog.Var("elements", plural.Selectf(1, "", plural.One, "element", plural.Other, "elements")),
		catalog.String("is too few elements (minimum is %[1]d ${elements})"),
	)
	c.Set(tag, code.GreaterThan, catalog.String("must be greater than %[1]v"))
	c.Set(tag, code.LessThan, catalog.String("must be less than %[1]v"))
	c.Set(tag, code.GreaterThanOrEqual, catalog.String("must be greater than or equal to %[1]v"))
	c.Set(tag, code.LessThanOrEqual, catalog.String("must be less than or equal to %[1]v"))

	c.Set(tag, code.Inclusion, catalog.String("is not included in %[1]v"))
	c.Set(tag, code.RegexpMismatch, catalog.String("is a mismatch with the regular expression. (%[1]s)"))
	c.Set(tag, code.InvalidURLFormat, catalog.String("is an invalid url format"))
	c.Set(tag, code.InvalidScheme, catalog.String("which scheme is not included in %[1]v"))
}
