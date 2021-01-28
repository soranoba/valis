package translations_test

import (
	"errors"
	"fmt"
	"github.com/soranoba/valis/code"
	"github.com/soranoba/valis/translations"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"regexp"
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
		f(code.NotString)(Results{
			en: "must be any string",
			ja: "must be any string",
		}),
		f(code.NotStruct)(Results{
			en: "must be any struct",
			ja: "must be any struct",
		}),
		f(code.NotStructField)(Results{
			en: "must be any struct field",
			ja: "must be any struct field",
		}),
		f(code.NotArray)(Results{
			en: "must be any array",
			ja: "must be any array",
		}),
		f(code.NotMap)(Results{
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
		f(code.GreaterThan, 10)(Results{
			en: "must be greater than 10",
			ja: "は10以上の値にする必要があります",
		}),
		f(code.LessThan, 10)(Results{
			en: "must be less than 10",
			ja: "は10以下の値にする必要があります",
		}),
		f(code.GreaterThanOrEqual, 10)(Results{
			en: "must be greater than or equal to 10",
			ja: "は10より大きい値にする必要があります",
		}),
		f(code.LessThanOrEqual, 10)(Results{
			en: "must be less than or equal to 10",
			ja: "は10より小さい値にする必要があります",
		}),
		f(code.Inclusion, []interface{}{"male", "female"})(Results{
			en: "is not included in [male female]",
			ja: "は [male female] のいずれかである必要があります",
		}),
		f(code.RegexpMismatch, regexp.MustCompile("^[0-9]+$").String())(Results{
			en: "is a mismatch with the regular expression. (^[0-9]+$)",
			ja: "は正規表現 (^[0-9]+$) に一致しません",
		}),
		f(code.InvalidURLFormat, errors.New("hoge"))(Results{
			en: "is an invalid url format",
			ja: "は不正なURLです",
		}),
		f(code.InvalidScheme, []string{"http", "https"})(Results{
			en: "which scheme is not included in [http https]",
			ja: "のスキームは [http https] のいずれかである必要があります",
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
