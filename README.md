valis
=====
[![CircleCI](https://circleci.com/gh/soranoba/henge.svg?style=svg&circle-token=3c8c20a0a57a6333fb949dd6b901c610656e9da6)](https://circleci.com/gh/soranoba/henge)
[![Go Report Card](https://goreportcard.com/badge/github.com/soranoba/valis)](https://goreportcard.com/report/github.com/soranoba/valis)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/soranoba/valis)](https://pkg.go.dev/github.com/soranoba/valis)

Valis is a validation framework for Go.

## Overviews

- 👌  Validation in various ways
  - struct tag, Validate methods, and various other ways.
- 🔧  Excellent customizability.
- 🌎  Support translations.

## Motivations

### Validation in various ways

[go-playground/validator](https://github.com/go-playground/validator) is a great validation framework, but it only supports the way using struct tag.<br>
[go-ozzo/ozzo-validation](https://github.com/go-ozzo/ozzo-validation) supports the way using Validate methods, but it does not have enough features to use struct tag.<br>
Valis supports to validate in various ways due to its excellent extensibility.<br>

### Create validation rules from multiple struct tags

When using various libraries, it is often the case that constraints with the same content are described by other struct tags.<br>
For example, `required:"true" validate:"required"`.<br>

What this !?!?<br>
Can you guarantee that both are always in agreement?<br>

Valis solves this by supporting to read all tags.

### Customizability is power

Performance is important, but it's faster not to use the library.<br>
Therefore, Valis take emphasis on customizability.<br>

So, requests for customizability are also welcomed.

## Installation

To install it, run:

```
go get -u github.com/soranoba/valis
```

## Usage

### Basic

```go
package main

import (
	"fmt"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
)

func main() {
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
}
```

### Struct tag

```go
package main

import (
	"fmt"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/tagrule"
	"github.com/soranoba/valis/when"
)

func main() {
	type User struct {
		Name    string `required:"true"`
		Age     int    `validate:"min=20"`
		Company struct {
			Location string `required:"true"`
		}
	}

	v := valis.NewValidator()
	// Use the CommonRule if you want to automatically search and validate all hierarchies.
	v.SetCommonRules(when.IsStruct(valis.EachFields(tagrule.Required, tagrule.Validate)))

	user := User{}
	if err := v.Validate(&user); err != nil {
		fmt.Println(err)
	}
}
```
Please refer to [documents](https://pkg.go.dev/github.com/soranoba/valis) for other usages.
