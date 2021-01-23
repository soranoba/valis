package valis

import (
	"fmt"
	"reflect"
	"sync"
)

type (
	FieldTagRule     *fieldTagRule
	FieldTagRuleFunc func(tagValue string) ([]Rule, error)
)

type (
	fieldTagRule struct {
		key      string
		ruleFunc FieldTagRuleFunc
		lock     *sync.RWMutex
		cache    map[reflect.StructTag][]Rule
	}
)

// NewFieldTagRule returns a new rule, that call ruleFunc when it found any fields that have the tag.
func NewFieldTagRule(key string, ruleFunc FieldTagRuleFunc) *fieldTagRule {
	return &fieldTagRule{key: key, ruleFunc: ruleFunc, lock: &sync.RWMutex{}, cache: map[reflect.StructTag][]Rule{}}
}

func (r *fieldTagRule) Validate(validator *Validator, value interface{}) {
	loc := validator.Location()
	if loc.Kind() != LocationKindField {
		return
	}

	field := loc.Field()
	r.lock.RLock()
	rules, ok := r.cache[field.Tag]
	r.lock.RUnlock()

	if !ok {
		if tag, ok := field.Tag.Lookup(r.key); ok {
			var err error
			rules, err = r.ruleFunc(tag)
			if err != nil {
				panic(fmt.Sprintf("%s (key = %s, path = %s)", err.Error(), r.key, field.PkgPath))
			}

			r.lock.Lock()
			r.cache[field.Tag] = rules
			r.lock.Unlock()
		} else {
			return
		}
	}
	And(rules...).Validate(validator, value)
}
