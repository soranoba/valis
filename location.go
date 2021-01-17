package valis

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type (
	LocationKind int
	Location     interface {
		Kind() LocationKind
		Parent() Location

		Field() reflect.StructField
		Index() int
		Key() interface{}

		FieldLocation(field reflect.StructField) Location
		IndexLocation(index int) Location
		KeyLocation(key interface{}) Location
	}
	LocationNameResolver interface {
		ResolveLocationName(loc Location) string
	}
)

type (
	location struct {
		parent Location
		kind   LocationKind
		value  interface{}
	}

	defaultLocationNameResolver struct {
	}
	jsonLocationNameResolver struct {
	}
	requestLocationNameResolver struct {
	}
)

const (
	LocationKindRoot LocationKind = iota
	LocationKindField
	LocationKindIndex
	LocationKindKey
)

var (
	DefaultLocationNameResolver LocationNameResolver = &defaultLocationNameResolver{}
	JSONLocationNameResolver    LocationNameResolver = &jsonLocationNameResolver{}
	RequestLocationNameResolver LocationNameResolver = &requestLocationNameResolver{}
)

func NewLocation() Location {
	return &location{}
}

func (loc *location) Kind() LocationKind {
	return loc.kind
}

func (loc *location) Parent() Location {
	if loc.kind == LocationKindRoot {
		panic("")
	}
	return loc.parent
}

func (loc *location) Field() reflect.StructField {
	if loc.kind != LocationKindField {
		panic("")
	}
	return loc.value.(reflect.StructField)
}

func (loc *location) Index() int {
	if loc.kind != LocationKindIndex {
		panic("")
	}
	return loc.value.(int)
}

func (loc *location) Key() interface{} {
	if loc.kind != LocationKindKey {
		panic("")
	}
	return loc.value
}

func (loc *location) FieldLocation(field reflect.StructField) Location {
	return &location{
		parent: loc,
		kind:   LocationKindField,
		value:  field,
	}
}

func (loc *location) IndexLocation(index int) Location {
	return &location{
		parent: loc,
		kind:   LocationKindIndex,
		value:  index,
	}
}

func (loc *location) KeyLocation(key interface{}) Location {
	return &location{
		parent: loc,
		kind:   LocationKindKey,
		value:  key,
	}
}

func (r *defaultLocationNameResolver) ResolveLocationName(loc Location) string {
	switch loc.Kind() {
	case LocationKindRoot:
		return ""
	case LocationKindField:
		return r.ResolveLocationName(loc.Parent()) + "." + loc.Field().Name
	case LocationKindIndex:
		return r.ResolveLocationName(loc.Parent()) + "[" + strconv.Itoa(loc.Index()) + "]"
	case LocationKindKey:
		return r.ResolveLocationName(loc.Parent()) + "[" + fmt.Sprintf("%v", loc.Key()) + "]"
	default:
		panic("")
	}
}

func (r *jsonLocationNameResolver) ResolveLocationName(loc Location) string {
	switch loc.Kind() {
	case LocationKindRoot:
		return ""
	case LocationKindField:
		field := loc.Field()
		name := field.Name
		if val := field.Tag.Get("json"); val != "" && val != "-" {
			attrs := strings.Split(val, ",")
			if len(attrs) > 0 {
				name = attrs[0]
			}
		}
		return r.ResolveLocationName(loc.Parent()) + "." + name
	case LocationKindIndex:
		return r.ResolveLocationName(loc.Parent()) + "[" + strconv.Itoa(loc.Index()) + "]"
	case LocationKindKey:
		return r.ResolveLocationName(loc.Parent()) + "." + fmt.Sprintf("%v", loc.Key())
	default:
		panic("")
	}
}

func (r *requestLocationNameResolver) ResolveLocationName(loc Location) string {
	switch loc.Kind() {
	case LocationKindRoot:
		return ""
	case LocationKindField:
		field := loc.Field()
		name := field.Name
		if val := field.Tag.Get("json"); val != "" && val != "-" {
			attrs := strings.Split(val, ",")
			if len(attrs) > 0 {
				name = attrs[0]
			}
		} else if val := field.Tag.Get("query"); val != "" && val != "-" {
			attrs := strings.Split(val, ",")
			if len(attrs) > 0 {
				name = attrs[0]
			}
		}
		return r.ResolveLocationName(loc.Parent()) + "." + name
	case LocationKindIndex:
		return r.ResolveLocationName(loc.Parent()) + "[" + strconv.Itoa(loc.Index()) + "]"
	case LocationKindKey:
		return r.ResolveLocationName(loc.Parent()) + "." + fmt.Sprintf("%v", loc.Key())
	default:
		panic("")
	}
}
