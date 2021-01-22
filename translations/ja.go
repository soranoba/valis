package translations

import (
	"github.com/soranoba/valis/code"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
)

// DefaultJapanese is a CatalogRegistrationFunc for Japanese.
func DefaultJapanese(c *catalog.Builder) {
	tag := language.Japanese

	// type error
	c.Set(tag, code.StringOnly, catalog.String("must be a string"))
	c.Set(tag, code.StructOnly, catalog.String("must be any struct"))
	c.Set(tag, code.ArrayOnly, catalog.String("must be any array"))
	c.Set(tag, code.MapOnly, catalog.String("must be any map"))
	c.Set(tag, code.NotAssignable, catalog.String("can't assign to %[1]s"))

	// not found error
	c.Set(tag, code.NoKey, catalog.String("requires the value at the key (%[1]v)"))

	// convert error
	c.Set(tag, code.ConversionFailed, catalog.String("%[1]v"))

	// others
	c.Set(tag, code.Custom, catalog.String("%[1]v"))
	c.Set(tag, code.Invalid, catalog.String("は不正な値です"))
	c.Set(tag, code.Required, catalog.String("は必須です"))
	c.Set(tag, code.NilOrNonZero, catalog.String("を指定する場合は空白にすることはできません"))
	c.Set(tag, code.ZeroOnly, catalog.String("は指定できません"))

	c.Set(tag, code.TooLongLength, catalog.String("は%[1]d文字までです"))
	c.Set(tag, code.TooShortLength, catalog.String("は%[1]d文字以上必要です"))
	c.Set(tag, code.TooLongLen, catalog.String("は%[1]d要素までです"))
	c.Set(tag, code.TooShortLen, catalog.String("は%[1]d要素以上必要です"))

	c.Set(tag, code.Inclusion, catalog.String("は %[1]v のいずれかである必要があります"))
}
