package pubcopy

import (
	"reflect"
	"unicode"
)

// Filter a representation of filter to decide whether a field or object to be copied depending on field name and/or type
type Filter interface {
	Name(string) bool
	Type(p reflect.Type) bool
}

var _ Filter = oneFuncNameFilter(nil)

type oneFuncNameFilter func(string) bool

func (o oneFuncNameFilter) Name(name string) bool {
	return o(name)
}

func (o oneFuncNameFilter) Type(p reflect.Type) bool {
	if len(p.Name()) == 0 {
		// anonymous struct is OK
		return true
	}
	return o(p.Name())
}

// PublicOnly lets only public fields and types to go
var PublicOnly = oneFuncNameFilter(func(name string) bool {
	runes := []rune(name)
	if len(runes) == 0 {
		return false
	}

	return unicode.IsUpper(runes[0])
})
