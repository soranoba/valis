package valis

import (
	"reflect"
	"sync"
)

type (
	TagRuleFunc func(tagValue string) ([]Rule, error)
)

type (
	tagRule struct {
		key      string
		ruleFunc TagRuleFunc
		lock     *sync.RWMutex
		cache    map[reflect.StructTag][]Rule
	}
)

const (
	InvalidStructTagCode = "invalid_struct_tag"
)

// NewTagRule returns a new rule, that call ruleFunc when it found any fields that have the tag.
func NewTagRule(key string, ruleFunc TagRuleFunc) Rule {
	return &tagRule{key: key, ruleFunc: ruleFunc, lock: &sync.RWMutex{}, cache: map[reflect.StructTag][]Rule{}}
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

		r.lock.RLock()
		rules, ok := r.cache[field.Tag]
		r.lock.RUnlock()

		if !ok {
			if tag, ok := field.Tag.Lookup(r.key); ok {
				var err error
				rules, err = r.ruleFunc(tag)
				if err != nil {
					validator.ErrorCollector().Add(validator.Location(), &ErrorDetail{
						InvalidStructTagCode,
						r,
						value,
						nil,
						err,
					})
					continue
				}

				r.lock.Lock()
				r.cache[field.Tag] = rules
				r.lock.Unlock()
			} else {
				continue
			}
		}

		fieldValue := fieldVal.Interface()
		validator.DiveField(&field, func(v *Validator) {
			for _, rule := range rules {
				rule.Validate(validator, fieldValue)
			}
		})
	}
}
