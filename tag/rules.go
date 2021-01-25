// The package implements some valis.Rule related to field tag.
package valistag

import (
	"github.com/soranoba/henge"
	"github.com/soranoba/valis"
	"github.com/soranoba/valis/code"
	"github.com/soranoba/valis/is"
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
	formatTagHandler   struct{}
	oneofRule          struct {
		elems []string
	}
)

var (
	// Required is a `required` tag rule.
	// See also is.Required rule.
	//
	// For example,
	//   `required:"true"`
	//
	Required = valis.NewFieldTagRule("required", &requiredTagHandler{})
	// Format is a `format` tag rule.
	// See also is.MatchString rule.
	//
	// For example,
	//   `format:"^[0-9]+$"`
	//   `format:"^\\d+$"`
	//
	Format = valis.NewFieldTagRule("format", &formatTagHandler{})
	// Validate is a `validate` tag rule.
	Validate = valis.NewFieldTagRule("validate", &ValidateTagHandler{})
)

var (
	validateTagSubKeys = map[string]func(string) ([]valis.Rule, error){
		"required": func(v string) ([]valis.Rule, error) { // required
			return []valis.Rule{is.Required}, nil
		},
		"lte": func(v string) ([]valis.Rule, error) { // lte=10
			var num float64
			if _, err := SplitAndParseTagValues(v, " ", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{is.LessThanOrEqualTo(num)}, nil
		},
		"lt": func(v string) ([]valis.Rule, error) { // lt=10
			var num float64
			if _, err := SplitAndParseTagValues(v, " ", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{is.LessThan(num)}, nil
		},
		"gte": func(v string) ([]valis.Rule, error) { // gte=10
			var num float64
			if _, err := SplitAndParseTagValues(v, " ", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{is.GreaterThanOrEqualTo(num)}, nil
		},
		"gt": func(v string) ([]valis.Rule, error) { // gt=10
			var num float64
			if _, err := SplitAndParseTagValues(v, " ", &num); err != nil {
				return nil, err
			}
			return []valis.Rule{is.GreaterThan(num)}, nil
		},
		"min": func(v string) ([]valis.Rule, error) { // min=1
			var min int
			if _, err := SplitAndParseTagValues(v, " ", &min); err != nil {
				return nil, err
			}
			return []valis.Rule{
				when.IsNumeric(is.Min(min)).
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
				when.IsNumeric(is.Min(max)).
					ElseWhen(when.IsTypeOrElem(reflect.TypeOf((*string)(nil)), is.LengthBetween(0, max))).
					Else(is.LenBetween(0, max)),
			}, nil
		},
		"oneof": func(v string) ([]valis.Rule, error) { // oneof=1 2
			if v == "" {
				return nil, errInsufficientNumberOfTagParameters
			}
			elems := strings.Split(v, " ")
			return []valis.Rule{&oneofRule{elems: elems}}, nil
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

func (h *formatTagHandler) ParseTagValue(tagValue string) ([]valis.Rule, error) {
	if tagValue == "" {
		return nil, errInsufficientNumberOfTagParameters
	}
	return []valis.Rule{is.MatchString(tagValue)}, nil
}

func (rule *oneofRule) Validate(validator *valis.Validator, value interface{}) {
	val, err := henge.New(value).String().Result()
	if err != nil {
		validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.ConversionFailed, value, err))
	}

	for _, elem := range rule.elems {
		if elem == val {
			return
		}
	}
	validator.ErrorCollector().Add(validator.Location(), valis.NewError(code.Inclusion, value, rule.elems))
}
