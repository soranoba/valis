package valis

import (
	"fmt"
	"sync"

	"github.com/soranoba/valis/code"
)

type (
	fieldTagRule struct {
		key             string
		ruleFactoryFunc func(tagValue string) ([]Rule, error)
		lock            *sync.RWMutex
		cache           map[string][]Rule
	}
)

// NewFieldTagRule returns a new rule related to the field tag.
// The rule verifies the value when it is a field value and has the specified tag.
func NewFieldTagRule(key string, ruleFactoryFunc func(tagValue string) ([]Rule, error)) *fieldTagRule {
	return &fieldTagRule{key: key, ruleFactoryFunc: ruleFactoryFunc, lock: &sync.RWMutex{}, cache: map[string][]Rule{}}
}

func (r *fieldTagRule) Validate(validator *Validator, value interface{}) {
	loc := validator.Location()
	if loc.Kind() != LocationKindField {
		validator.ErrorCollector().Add(loc, NewError(code.NotStruct, value))
		return
	}

	field := loc.Field()
	tag, ok := field.Tag.Lookup(r.key)
	if !ok {
		return
	}

	r.lock.RLock()
	rules, ok := r.cache[tag]
	r.lock.RUnlock()

	if !ok {
		var err error
		rules, err = r.ruleFactoryFunc(tag)
		if err != nil {
			panic(fmt.Sprintf("%s (key = %s, path = %s)", err.Error(), r.key, field.PkgPath))
		}

		r.lock.Lock()
		r.cache[tag] = rules
		r.lock.Unlock()
	}
	And(rules...).Validate(validator, value)
}
