package maps_test

import (
	"testing"

	"github.com/pyrrho/encoding/maps"
	"github.com/stretchr/testify/require"
)

type A struct {
	FieldOne   int
	FieldTwo   float64
	FieldThree string
	FieldFour  complex128
}

func TestSimpleUntaggedStruct(t *testing.T) {
	require := require.New(t)

	var (
		err              error
		actual, expected map[string]interface{}
	)

	a := A{
		42,
		3.14,
		"Hello World",
		complex(1, 2),
	}
	expected = map[string]interface{}{
		"FieldOne":   42,
		"FieldTwo":   float64(3.14),
		"FieldThree": "Hello World",
		"FieldFour":  complex(1, 2),
	}
	actual, err = maps.Marshal(a)

	require.NoError(err)
	require.Equal(expected, actual)
}

type B struct {
	FieldOne   int        ``                 // undecorated
	fieldTwo   float64    ``                 // unexported
	FieldThree string     `map:"-"`          // explicitly ignored
	FieldFour  complex128 `map:"field_four"` // explicitly named
}

func TestSimpleTaggedStruct(t *testing.T) {
	require := require.New(t)

	a := B{
		42,
		3.14,
		"Hello World",
		complex(1, 2),
	}
	expected := map[string]interface{}{
		"FieldOne":   42,
		"field_four": complex(1, 2),
	}
	actual, err := maps.Marshal(a)

	require.NoError(err)
	require.Equal(expected, actual)
}

type C struct {
	AMap    map[int]int
	AStruct D
}

type D struct {
	A int
	B int
}

func TestNestedStructsAndMaps(t *testing.T) {
	require := require.New(t)

	a := C{
		map[int]int{
			1: 2,
			3: 4,
		},
		D{5, 6},
	}
	expected := map[string]interface{}{
		"AMap": map[int]int{
			1: 2,
			3: 4,
		},
		"AStruct": map[string]interface{}{
			"A": 5,
			"B": 6,
		},
	}
	actual, err := maps.Marshal(a)

	require.NoError(err)
	require.Equal(expected, actual)
}

type E struct {
	A int
	F // embedded
}

type F struct {
	G // embedded
}

type G struct {
	H // embedded
}

type H struct {
	unexported int
	Exported   int
}

func TestSimpleEmbeddedStructs(t *testing.T) {
	require := require.New(t)

	a := E{
		42,
		F{G{H{1, 2}}},
	}
	expected := map[string]interface{}{
		"A":        42,
		"Exported": 2,
	}
	actual, err := maps.Marshal(a)

	require.NoError(err)
	require.Equal(expected, actual)
}

type I struct {
	J // embedded, with contentious field names
	K // embedded, with contentious field names
}

type J struct {
	AInt    int
	AString string
	AFloat  float64
}

type K struct {
	// Accessing `I.AInt` would cause an ambiguous selector compiler error, but
	// `map` tag means there won't be contention when marshalling.
	AInt int `map:"AInt"`
	// embedded
	L
}

type L struct {
	// `L.AString` will be shadowed by `J.AString`.
	AString string
	// `L.AFloat` will be shadowed by `J.AFloat`, despite the `map` tag.
	AFloat float64
}

func TestContendingEmbeddedStructs(t *testing.T) {
	require := require.New(t)

	a := I{
		J{
			100,
			"foo",
			3.14,
		},
		K{
			200,
			L{
				"bar",
				6.28,
			},
		},
	}
	require.Equal("foo", a.AString)
	require.Equal(3.14, a.AFloat)
	// Compile-time error
	// require.Equal(0, a.AInt)
	expected := map[string]interface{}{
		"AInt":    200,   // From K
		"AString": "foo", // From J
		"AFloat":  3.14,  // From J
	}
	actual, err := maps.Marshal(a)

	require.NoError(err)
	require.Equal(expected, actual)
}
