package pubcopy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type subStruct struct {
	Name   string
	Value  string
	hidden string
}

type Hidden struct {
	Value string
}

type structure struct {
	subStruct
	Hidden

	Another     subStruct
	AnotherPtr  *subStruct
	another     Hidden
	Value       int
	Pointer     *int
	hiddenValue int
	String      string

	Items []subStruct
	Map   map[string]subStruct
}

func TestCopy(t *testing.T) {
	intValue := 4
	anotherIntValue := 4
	input := structure{
		subStruct: subStruct{
			Name:   "name",
			Value:  "value",
			hidden: "hidden",
		},
		Hidden: Hidden{
			Value: "actually visible",
		},
		Another: subStruct{
			Name:   "name",
			Value:  "value",
			hidden: "hidden",
		},
		AnotherPtr: &subStruct{
			Name:   "name",
			Value:  "value",
			hidden: "hidden",
		},
		another: Hidden{
			Value: "value",
		},
		Value:       77,
		Pointer:     &intValue,
		hiddenValue: 12,
		String:      "string",
		Items: []subStruct{
			{
				Name:   "name",
				Value:  "value",
				hidden: "hidden",
			},
		},
		Map: map[string]subStruct{
			"key": {
				Name:   "name",
				Value:  "value",
				hidden: "hidden",
			},
		},
	}
	mustBeAfterCopy := structure{
		Hidden: Hidden{
			Value: "actually visible",
		},
		Another: subStruct{
			Name:  "name",
			Value: "value",
		},
		AnotherPtr: &subStruct{
			Name:  "name",
			Value: "value",
		},
		Value:   77,
		Pointer: &anotherIntValue,
		String:  "string",
		Items: []subStruct{
			{
				Name:  "name",
				Value: "value",
			},
		},
		Map: map[string]subStruct{
			"key": {
				Name:  "name",
				Value: "value",
			},
		},
	}
	var dest structure

	if err := Copy(input, &dest, PublicOnly); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, mustBeAfterCopy, dest)
}
