package translations_test

import (
	"errors"
	"fmt"
	"github.com/soranoba/valis/code"
	"github.com/soranoba/valis/translations"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"testing"
)

func TestTranslations(t *testing.T) {
	assert := assert.New(t)

	type Results map[language.Tag]string
	type Data struct {
		Code    string
		Params  []interface{}
		Results Results
	}

	f := func(code string, params ...interface{}) func(results Results) *Data {
		return func(results Results) *Data {
			return &Data{
				Code:    code,
				Params:  params,
				Results: results,
			}
		}
	}

	en := language.English
	ja := language.Japanese

	data := []*Data{
		f(code.StringOnly)(Results{
			en: "must be a string",
			ja: "must be a string",
		}),
		f(code.StructOnly)(Results{
			en: "must be any struct",
			ja: "must be any struct",
		}),
		f(code.ArrayOnly)(Results{
			en: "must be any array",
			ja: "must be any array",
		}),
		f(code.MapOnly)(Results{
			en: "must be any map",
			ja: "must be any map",
		}),
		f(code.NotAssignable, "string")(Results{
			en: "can't assign to string",
			ja: "can't assign to string",
		}),
		f(code.NoKey, "name")(Results{
			en: "requires the value at the key (name)",
			ja: "requires the value at the key (name)",
		}),
		f(code.ConversionFailed, errors.New("can't convert to string"))(Results{
			en: "can't convert to string",
			ja: "can't convert to string",
		}),
		f(code.Custom, errors.New("has error occurred"))(Results{
			en: "has error occurred",
			ja: "has error occurred",
		}),
		f(code.Invalid)(Results{
			en: "is invalid",
			ja: "は不正な値です",
		}),
		f(code.Required)(Results{
			en: "is required",
			ja: "は必須です",
		}),
		f(code.NilOrNonZero)(Results{
			en: "can't be blank (or zero) if specified",
			ja: "を指定する場合は空白にすることはできません",
		}),
		f(code.ZeroOnly)(Results{
			en: "must be blank",
			ja: "は指定できません",
		}),
		f(code.TooLongLength, 10)(Results{
			en: "is too long length (maximum is 10 characters)",
			ja: "は10文字までです",
		}),
		f(code.TooShortLength, 10)(Results{
			en: "is too short length (minimum is 10 characters)",
			ja: "は10文字以上必要です",
		}),
		f(code.TooLongLen, 10)(Results{
			en: "is too many elements (maximum is 10 elements)",
			ja: "は10要素までです",
		}),
		f(code.TooShortLen, 10)(Results{
			en: "is too few elements (minimum is 10 elements)",
			ja: "は10要素以上必要です",
		}),
		f(code.Inclusion, []interface{}{"male", "female"})(Results{
			en: "is not included in [male female]",
			ja: "は [male female] のいずれかである必要があります",
		}),
	}

	c := translations.NewCatalog()
	for _, registerFunc := range translations.AllPredefinedCatalogRegistrationFunc {
		c.Set(registerFunc)
	}

	for _, d := range data {
		for _, lang := range c.Languages() {
			p := message.NewPrinter(lang, message.Catalog(c))
			if !assert.Contains(d.Results, lang, fmt.Sprintf("%s does not have %s data", d.Code, lang.String())) {
				continue
			}
			assert.Equal(
				d.Results[lang],
				p.Sprintf(d.Code, d.Params...),
				fmt.Sprintf("lang: %s, code: %s", lang.String(), d.Code),
			)
		}
	}
}
