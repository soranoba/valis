package valis_test

import (
	"encoding/json"
	"fmt"
	valistag "github.com/soranoba/valis/tag"
	"github.com/soranoba/valis/when"

	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/soranoba/valis/translations"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

func Example() {
	type User struct {
		Name string
		Age  int
	}

	u := &User{}
	if err := valis.Validate(
		&u,
		valis.Field(&u.Name, is.Required),
		valis.Field(&u.Age, is.Min(20)),
	); err != nil {
		fmt.Println(err)
	}

	u.Name = "Alice"
	u.Age = 20
	if err := valis.Validate(
		&u,
		valis.Field(&u.Name, is.Required),
		valis.Field(&u.Age, is.Min(20)),
	); err != nil {
		fmt.Println(err)
	}

	// Output:
	// (required) .Name is required
	// (gte) .Age must be greater than or equal to 20
}

func Example_customizeError() {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	v := valis.NewValidator()
	v.SetErrorCollectorFactoryFunc(func() valis.ErrorCollector {
		// You can create ErrorCollectors yourself that will generate your own error.
		// If you want to change only the attribute name, please change NameResolver.
		return valis.NewStandardErrorCollector(valis.JSONLocationNameResolver)
	})

	u := User{}
	if err := v.Validate(
		&u,
		valis.Field(&u.Name, is.Required),
		valis.Field(&u.Age, is.Min(20)),
	); err != nil {
		fmt.Println(err)
	}

	// Output:
	// (required) .name is required
	// (gte) .age must be greater than or equal to 20
}

func Example_translate() {
	type User struct {
		Name string
		Age  int
	}

	// It is recommended to have the catalog in a global variable etc. instead of creating it every time.
	// We welcome your contributions if it does not exist the language you use.
	c := translations.NewCatalog(catalog.Fallback(language.English))
	for _, f := range translations.AllPredefinedCatalogRegistrationFunc {
		c.Set(f)
	}

	u := User{}
	if err := valis.Validate(
		&u,
		valis.Field(&u.Name, is.Required),
		valis.Field(&u.Age, is.Min(20)),
	); err != nil {
		for _, lang := range []language.Tag{language.English, language.Japanese} {
			p := message.NewPrinter(lang, message.Catalog(c))
			// When you change the ErrorCollector and create errors other than ValidationError, you need an alternative.
			m := err.(*valis.ValidationError).Translate(p)
			b, _ := json.MarshalIndent(m, "", "  ")
			fmt.Printf("%s\n", b)
		}
	}

	// Output:
	// {
	//   ".Age": [
	//     "must be greater than or equal to 20"
	//   ],
	//   ".Name": [
	//     "is required"
	//   ]
	// }
	// {
	//   ".Age": [
	//     "は20より大きい値にする必要があります"
	//   ],
	//   ".Name": [
	//     "は必須です"
	//   ]
	// }
}

func Example_structTag() {
	type User struct {
		Name string `required:"true"`
		Age  int    `validate:"min=20"`
	}

	v := valis.NewValidator()
	v.SetCommonRules(when.IsStruct(valis.EachFields(valistag.Required, valistag.Validate)))

	u := User{}
	if err := v.Validate(&u); err != nil {
		fmt.Println(err)
	}

	// Output:
	// (required) .Name is required
	// (gte) .Age must be greater than or equal to 20
}
