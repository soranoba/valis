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
	c.Set(tag, code.Invalid, catalog.String("は不正な値です"))
	c.Set(tag, code.Required, catalog.String("は必須です"))
	c.Set(tag, code.NonZero, catalog.String("を空白にすることはできません"))
	c.Set(tag, code.NilOrNonZero, catalog.String("を指定する場合は空白にすることはできません"))
	c.Set(tag, code.ZeroOnly, catalog.String("は指定できません"))

	c.Set(tag, code.TooLongLength, catalog.String("は%[1]d文字までです"))
	c.Set(tag, code.TooShortLength, catalog.String("は%[1]d文字以上必要です"))
	c.Set(tag, code.TooLongLen, catalog.String("は%[1]d要素までです"))
	c.Set(tag, code.TooShortLen, catalog.String("は%[1]d要素以上必要です"))
	c.Set(tag, code.GreaterThan, catalog.String("は%[1]v以上の値にする必要があります"))
	c.Set(tag, code.LessThan, catalog.String("は%[1]v以下の値にする必要があります"))
	c.Set(tag, code.GreaterThanOrEqual, catalog.String("は%[1]vより大きい値にする必要があります"))
	c.Set(tag, code.LessThanOrEqual, catalog.String("は%[1]vより小さい値にする必要があります"))

	c.Set(tag, code.Inclusion, catalog.String("は %[1]v のいずれかである必要があります"))
	c.Set(tag, code.RegexpMismatch, catalog.String("は正規表現 (%[1]s) に一致しません"))
	c.Set(tag, code.InvalidURLFormat, catalog.String("は不正なURLです"))
	c.Set(tag, code.InvalidScheme, catalog.String("のスキームは %[1]v のいずれかである必要があります"))
}
