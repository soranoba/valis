package valis

import "reflect"

type (
	TagRuleFunc func(tagValue string) []Rule
)

type (
	tagRule struct {
		key      string
		ruleFunc TagRuleFunc
	}
)

func TagRule(key string, ruleFunc TagRuleFunc) Rule {
	return &tagRule{key: key, ruleFunc: ruleFunc}
}

func (r *tagRule) Validate(validator *Validator, value interface{}) {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		if !fieldVal.IsValid() {
			continue
		}

		field := val.Type().Field(i)
		if tag := field.Tag.Get(r.key); tag != "" {
			newValidator := validator.WithLocation(validator.Location().FieldLocation(field))
			fieldValue := fieldVal.Interface()
			for _, rule := range r.ruleFunc(tag) {
				rule.Validate(newValidator, fieldValue)
			}
		}
	}
}
