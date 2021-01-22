package translations_test

import (
	"fmt"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/code"
	"github.com/soranoba/valis/translations"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"testing"
)

func TestCatalog(t *testing.T) {
	var _ catalog.Catalog = &translations.Catalog{}
}

func Example() {
	c := translations.NewCatalog(catalog.Fallback(language.English))
	if true /* When set all predefined translation data. */ {
		for _, f := range translations.AllPredefinedCatalogRegistrationFunc {
			c.Set(f)
		}
	} else /* When select and set by yourself. */ {
		c.Set(translations.DefaultEnglish)
		c.Set(translations.DefaultJapanese)
	}

	err := valis.NewError(code.TooShortLength, "value", 10)
	en := message.NewPrinter(language.English, message.Catalog(c))
	fmt.Printf("%s %s.\n", ".Name", en.Sprintf(err.Code(), err.Params()...))

	ja := message.NewPrinter(language.Japanese, message.Catalog(c))
	fmt.Printf("%s%s。\n", "名前", ja.Sprintf(err.Code(), err.Params()...))

	// Output:
	// .Name is too short length (minimum is 10 characters).
	// 名前は10文字以上必要です。
}
