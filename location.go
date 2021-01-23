package valis

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type (
	// LocationKind is the type of Location.
	LocationKind int
	// Location indicates the location of the validating value.
	Location struct {
		parent *Location
		kind   LocationKind
		value  interface{}
	}
	// LocationNameResolver is an interface that creates a string corresponding to Location.
	LocationNameResolver interface {
		ResolveLocationName(loc *Location) string
	}
)

type (
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

func newRootLocation() *Location {
	return &Location{}
}

// Kind returns a LocationKind of the Location.
func (loc *Location) Kind() LocationKind {
	return loc.kind
}

// Parent returns a parent location, when the Kind is not LocationKindRoot. Otherwise, occur panics.
func (loc *Location) Parent() *Location {
	if loc.kind == LocationKindRoot {
		panic("Kind must be LocationKindRoot")
	}
	return loc.parent
}

// Field returns a reflect.StructField when the Kind is LocationKindField. Otherwise, occur panics.
func (loc *Location) Field() *reflect.StructField {
	if loc.kind != LocationKindField {
		panic("Kind must be LocationKindField")
	}
	return loc.value.(*reflect.StructField)
}

// Index returns a Index when the Kind is LocationKindIndex. Otherwise, occur panics.
func (loc *Location) Index() int {
	if loc.kind != LocationKindIndex {
		panic("Kind must be LocationKindIndex")
	}
	return loc.value.(int)
}

// Key returns a Key when the Kind is LocationKindMapKey or LocationKindMapValue. Otherwise, occur panics.
func (loc *Location) Key() interface{} {
	switch loc.kind {
	case LocationKindMapKey, LocationKindMapValue:
		return loc.value
	default:
		panic("Kind must be LocationKindMapKey or LocationKindMapValue")
	}
}

// FieldLocation returns a new Location that indicates the value at the field in the struct.
func (loc *Location) FieldLocation(field *reflect.StructField) *Location {
	return &Location{
		parent: loc,
		kind:   LocationKindField,
		value:  field,
	}
}

// IndexLocation returns a new Location that indicates the value at the index in the array or slice.
func (loc *Location) IndexLocation(index int) *Location {
	return &Location{
		parent: loc,
		kind:   LocationKindIndex,
		value:  index,
	}
}

// MapKeyLocation returns a new Location that indicates the key.
func (loc *Location) MapKeyLocation(key interface{}) *Location {
	return &Location{
		parent: loc,
		kind:   LocationKindMapKey,
		value:  key,
	}
}

// MapValueLocation returns a new Location that indicates the value of the key.
func (loc *Location) MapValueLocation(key interface{}) *Location {
	return &Location{
		parent: loc,
		kind:   LocationKindMapValue,
		value:  key,
	}
}

func (r *defaultLocationNameResolver) ResolveLocationName(loc *Location) string {
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

func (r *jsonLocationNameResolver) ResolveLocationName(loc *Location) string {
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

func (r *requestLocationNameResolver) ResolveLocationName(loc *Location) string {
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
