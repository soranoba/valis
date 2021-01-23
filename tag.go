package valis

import (
	"fmt"
	"github.com/soranoba/valis/code"
	"reflect"
	"sync"
)

type (
	fieldTagRule struct {
		key             string
		ruleFactoryFunc func(tagValue string) ([]Rule, error)
		lock            *sync.RWMutex
		cache           map[reflect.StructTag][]Rule
	}
)

// NewFieldTagRule returns a new rule related to the field tag.
// The rule verifies the value when it is a field value and has the specified tag.
func NewFieldTagRule(key string, ruleFactoryFunc func(tagValue string) ([]Rule, error)) *fieldTagRule {
	return &fieldTagRule{key: key, ruleFactoryFunc: ruleFactoryFunc, lock: &sync.RWMutex{}, cache: map[reflect.StructTag][]Rule{}}
}

func (r *fieldTagRule) Validate(validator *Validator, value interface{}) {
	loc := validator.Location()
	if loc.Kind() != LocationKindField {
		validator.ErrorCollector().Add(loc, NewError(code.NotStruct, value))
		return
	}

	field := loc.Field()
	r.lock.RLock()
	rules, ok := r.cache[field.Tag]
	r.lock.RUnlock()

	if !ok {
		if tag, ok := field.Tag.Lookup(r.key); ok {
			var err error
			rules, err = r.ruleFactoryFunc(tag)
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
