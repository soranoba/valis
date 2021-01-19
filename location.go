package valis

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type (
	// LocationKind is the type of Location
	LocationKind int
	// Location is the interface of the private struct that indicates the location of the value.
	// Should not define your own Location struct.
	Location interface {
		// Kind returns a LocationKind of the Location.
		Kind() LocationKind
		// Parent returns a parent location, when it is not LocationKindRoot. Otherwise, it returns nil.
		Parent() Location

		// Field returns a reflect.StructField when it is LocationKindField. Otherwise, it panic.
		Field() reflect.StructField
		// Index returns a Index when it is LocationKindIndex. Otherwise, it panic.
		Index() int
		// Key returns a Key when it is LocationKindMapKey or LocationKindMapValue. Otherwise, it panic.
		Key() interface{}

		// FieldLocation returns a new Location.
		FieldLocation(field reflect.StructField) Location
		// IndexLocation returns a new Location.
		IndexLocation(index int) Location
		// MapKeyLocation returns a new Location.
		MapKeyLocation(key interface{}) Location
		// MapValueLocation returns a new Location.
		MapValueLocation(key interface{}) Location
	}
	// LocationNameResolver is an interface that creates a string corresponding to Location.
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
	// LocationKindRoot means the location of the root.
	LocationKindRoot LocationKind = iota
	// LocationKindField means the location of the field of the a struct.
	LocationKindField
	// LocationKindIndex means the location of the value in an array or slice.
	LocationKindIndex
	// LocationKindMapKey means the location of the key in a map.
	LocationKindMapKey
	// LocationKindMapValue means the location of the value in a map.
	LocationKindMapValue
)

var (
	// DefaultLocationNameResolver is a LocationNamResolver used by default
	DefaultLocationNameResolver LocationNameResolver = &defaultLocationNameResolver{}
	// JSONLocationNameResolver is a LocationNameResolver that creates LocationNames using the json tag
	JSONLocationNameResolver LocationNameResolver = &jsonLocationNameResolver{}
	// RequestLocationNameResolver is a LocationNameResolver that creates LocationNames using the json and query tag
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
		panic("Kind must be LocationKindRoot")
	}
	return loc.parent
}

func (loc *location) Field() reflect.StructField {
	if loc.kind != LocationKindField {
		panic("Kind must be LocationKindField")
	}
	return loc.value.(reflect.StructField)
}

func (loc *location) Index() int {
	if loc.kind != LocationKindIndex {
		panic("Kind must be LocationKindIndex")
	}
	return loc.value.(int)
}

func (loc *location) Key() interface{} {
	switch loc.kind {
	case LocationKindMapKey, LocationKindMapValue:
		return loc.value
	default:
		panic("Kind must be LocationKindMapKey or LocationKindMapValue")
	}
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

func (loc *location) MapKeyLocation(key interface{}) Location {
	return &location{
		parent: loc,
		kind:   LocationKindMapKey,
		value:  key,
	}
}

func (loc *location) MapValueLocation(key interface{}) Location {
	return &location{
		parent: loc,
		kind:   LocationKindMapValue,
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
	case LocationKindMapKey:
		return r.ResolveLocationName(loc.Parent()) + "[key: " + fmt.Sprintf("%v", loc.Key()) + "]"
	case LocationKindMapValue:
		return r.ResolveLocationName(loc.Parent()) + "[" + fmt.Sprintf("%v", loc.Key()) + "]"
	default:
		panic("invalid LocationKind")
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
	case LocationKindMapKey:
		return r.ResolveLocationName(loc.Parent()) + "#" + fmt.Sprintf("%v", loc.Key())
	case LocationKindMapValue:
		return r.ResolveLocationName(loc.Parent()) + "." + fmt.Sprintf("%v", loc.Key())
	default:
		panic("invalid LocationKind")
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
	case LocationKindMapKey:
		return r.ResolveLocationName(loc.Parent()) + "#" + fmt.Sprintf("%v", loc.Key())
	case LocationKindMapValue:
		return r.ResolveLocationName(loc.Parent()) + "." + fmt.Sprintf("%v", loc.Key())
	default:
		panic("invalid LocationKind")
	}
}
