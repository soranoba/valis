// Package tagrule implements some valis.Rule related to field tag.
package tagrule

import (
	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/soranoba/valis/to"
	"github.com/soranoba/valis/when"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type (
	ValidateTagHandler struct{}
)

type (
	requiredTagHandler struct{}
	patternTagHandler  struct{}
	enumsTagHandler    struct{}
)

var (
	// Required is a `required` tag rule.
	// See also is.Required rule.
	//
	// For example,
	//   `required:"true"`
	//
	Required = valis.NewFieldTagRule("required", &requiredTagHandler{})
	// Format is a `pattern` tag rule.
	// See also is.MatchString rule.
	//
	// For example,
	//   `pattern:"^[0-9]+$"`
	//   `pattern:"^\\d+$"`
	//
	Pattern = valis.NewFieldTagRule("pattern", &patternTagHandler{})
	// Enums is a `enums` tag rule.
	// After converting the type, it applies the is.In rule.
	//
	// For example,
	//   `enums:"a,b"`
	//   `enums:"1,2"`
	Enums = valis.NewFieldTagRule("enums", &enumsTagHandler{})
	// Validate is a `validate` tag rule.
	Validate = valis.NewFieldTagRule("validate", &ValidateTagHandler{})
)

var (
	validateTagSubKeys = map[string]func(string) ([]valis.Rule, error){
		"required": func(v string) ([]valis.Rule, error) { // required
			return []valis.Rule{is.Required}, nil
		},
		"nonzero": func(v string) ([]valis.Rule, error) { // nonzero
			return []valis.Rule{is.NonZero}, nil
		},
		"zero": func(v string) ([]valis.Rule, error) { // zero
			return []valis.Rule{is.Zero}, nil
		},
		"lte": func(v string) ([]valis.Rule, error) { // lte=10
			var num float64
			if _, err := SplitAndParseTagValues(v, " ", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{when.IsNil().Else(is.LessThanOrEqualTo(num))}, nil
		},
		"lt": func(v string) ([]valis.Rule, error) { // lt=10
			var num float64
			if _, err := SplitAndParseTagValues(v, " ", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{when.IsNil().Else(is.LessThan(num))}, nil
		},
		"gte": func(v string) ([]valis.Rule, error) { // gte=10
			var num float64
			if _, err := SplitAndParseTagValues(v, " ", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{when.IsNil().Else(is.GreaterThanOrEqualTo(num))}, nil
		},
		"gt": func(v string) ([]valis.Rule, error) { // gt=10
			var num float64
			if _, err := SplitAndParseTagValues(v, " ", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{when.IsNil().Else(is.GreaterThan(num))}, nil
		},
		"min": func(v string) ([]valis.Rule, error) { // min=1
			var min int
			if _, err := SplitAndParseTagValues(v, " ", &min); err != nil {
				return nil, err
			}
			return []valis.Rule{
				when.IsNil().
					ElseWhen(when.IsNumeric(is.Min(min))).
					ElseWhen(when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.LengthBetween(min, math.MaxInt64))).
					Else(is.LenBetween(min, math.MaxInt64)),
			}, nil
		},
		"max": func(v string) ([]valis.Rule, error) { // max=10
			var max int
			if _, err := SplitAndParseTagValues(v, " ", &max); err != nil {
				return nil, err
			}
			return []valis.Rule{
				when.IsNil().
					ElseWhen(when.IsNumeric(is.Max(max))).
					ElseWhen(when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.LengthBetween(0, max))).
					Else(is.LenBetween(0, max)),
			}, nil
		},
		"oneof": func(v string) ([]valis.Rule, error) { // oneof=1 2
			if v == "" {
				return nil, errInsufficientNumberOfTagParameters
			}
			elems := henge.New(strings.Split(v, " ")).Slice().Value()
			return []valis.Rule{when.IsNil().Else(to.String(is.In(elems...)))}, nil
		},
		"url": func(v string) ([]valis.Rule, error) { // url=http https
			if v == "" {
				return []valis.Rule{is.URL()}, nil
			}
			return []valis.Rule{is.URL(strings.Split(v, " ")...)}, nil
		},
	}
)

func (h *requiredTagHandler) ParseTagValue(tagValue string) ([]valis.Rule, error) {
	ok, _ := strconv.ParseBool(tagValue)
	if ok {
		return []valis.Rule{is.Required}, nil
	}
	return []valis.Rule{}, nil
}

func (h *ValidateTagHandler) ParseTagValue(tagValue string) ([]valis.Rule, error) {
	elems := strings.Split(tagValue, ",")
	rules := make([]valis.Rule, 0)

	for _, elem := range elems {
		if elem == "-" {
			return []valis.Rule{}, nil
		}

		subKv := strings.SplitN(elem, "=", 2)
		if f, ok := validateTagSubKeys[subKv[0]]; ok {
			subKey := ""
			if len(subKv) == 2 {
				subKey = subKv[1]
			}
			r, err := f(subKey)
			if err != nil {
				return nil, err
			}
			rules = append(rules, r...)
		}
	}
	return rules, nil
}

func (h *patternTagHandler) ParseTagValue(tagValue string) ([]valis.Rule, error) {
	if tagValue == "" {
		return nil, errInsufficientNumberOfTagParameters
	}
	return []valis.Rule{when.IsNil().Else(is.MatchString(tagValue))}, nil
}

func (h *enumsTagHandler) ParseTagValue(tagValue string) ([]valis.Rule, error) {
	if tagValue == "" {
		return nil, errInsufficientNumberOfTagParameters
	}
	elems := henge.New(strings.Split(tagValue, ",")).Slice().Value()
	return []valis.Rule{when.IsNil().Else(to.String(is.In(elems...)))}, nil
}
