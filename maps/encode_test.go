package maps_test

import (
	"testing"

	"github.com/pyrrho/encoding/maps"
	"github.com/stretchr/testify/require"
)

type SimpleStruct struct {
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

	s := &SimpleStruct{
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
	actual, err = maps.Marshal(s)

	require.NoError(err)
	require.Equal(expected, actual)
}

type SimpleStructWithTags struct {
	FieldOne   int        ``                  // undecorated
	FieldTwo   float64    `map:"-"`           // explicitly ignored
	FieldThree string     `map:"field_three"` // explicitly named
	fieldFour  complex128 `map:"field_four"`  // unexported with name (ignored)
	fieldFive  bool       ``                  // unexported sans name
}

func TestSimpleTaggedStruct(t *testing.T) {
	require := require.New(t)

	s := &SimpleStructWithTags{
		42,
		3.14,
		"Hello World",
		complex(1, 2),
		true,
	}
	expected := map[string]interface{}{
		"FieldOne":    42,
		"field_three": "Hello World",
	}
	actual, err := maps.Marshal(s)

	require.NoError(err)
	require.Equal(expected, actual)
}

type ParentStruct struct {
	AMap    map[int]int
	AStruct NestedStruct
}

type NestedStruct struct {
	AnInt  int
	AFloat float64
}

func TestNestedStructsAndMaps(t *testing.T) {
	require := require.New(t)

	s := &ParentStruct{
		map[int]int{
			1: 2,
			3: 4,
		},
		NestedStruct{5, 6.7},
	}
	expected := map[string]interface{}{
		"AMap": map[int]int{
			1: 2,
			3: 4,
		},
		"AStruct": map[string]interface{}{
			"AnInt":  5,
			"AFloat": 6.7,
		},
	}
	actual, err := maps.Marshal(s)

	require.NoError(err)
	require.Equal(expected, actual)
}

type TopLevelStruct struct {
	AnInt  int
	WeMust // embedded
}

type WeMust struct {
	Go // embedded
}

type Go struct {
	Deeper // embedded
}

type Deeper struct {
	unexported int
	Exported   int
}

func TestSimpleEmbeddedStructs(t *testing.T) {
	require := require.New(t)

	s := &TopLevelStruct{
		42,
		WeMust{Go{Deeper{1, 2}}},
	}
	expected := map[string]interface{}{
		"AnInt":    42,
		"Exported": 2,
	}
	actual, err := maps.Marshal(s)

	require.NoError(err)
	require.Equal(expected, actual)
}

type LevelOne struct {
	LevelTwoLeft  // embedded, with contentious field names
	LevelTwoRight // embedded, with contentious field names
}

type LevelTwoLeft struct {
	AnInt   int
	AString string
	AFloat  float64
}

type LevelTwoRight struct {
	// Accessing `LevelOne.AnInt` would cause an ambiguous selector compiler
	// error, but `map` tag means there won't be contention when marshalling.
	AnInt int `map:"AnInt"`

	LevelThree // embedded
}

type LevelThree struct {
	// `LevelThree.AString` will be shadowed by `LevelTwoLeft.AString`.
	AString string
	// `LevelThree.AFloat` will be shadowed by `LevelTwoLeft.AFloat`, despite
	// the `map` tag.
	AFloat float64 `map:"AFloat"`
}

func TestContendingEmbeddedStructs(t *testing.T) {
	require := require.New(t)

	s := &LevelOne{
		LevelTwoLeft{
			100,
			"foo",
			3.14,
		},
		LevelTwoRight{
			200,
			LevelThree{
				"bar",
				6.28,
			},
		},
	}
	require.Equal("foo", s.AString)
	require.Equal(3.14, s.AFloat)
	// Compile-time error
	// require.Equal(0, s.AnInt)
	expected := map[string]interface{}{
		"AnInt":   200,   // From LevelTwoRight
		"AString": "foo", // From LevelTwoLeft
		"AFloat":  3.14,  // From LevelTwoLeft
	}
	actual, err := maps.Marshal(s)

	require.NoError(err)
	require.Equal(expected, actual)
}

type MarahalerParent struct {
	AnInt            int
	AnArrayIshStruct MarshalerImplementor
}

type MarshalerImplementor struct {
	AnArray  [3]int
	Constant int
}

func (mi *MarshalerImplementor) MarshalMapValue() (interface{}, error) {
	return map[string]int{
		"Arr0": mi.AnArray[0] + mi.Constant,
		"Arr1": mi.AnArray[1] + mi.Constant,
		"Arr2": mi.AnArray[2] + mi.Constant,
	}, nil
}

func TestMarshalerInterface(t *testing.T) {
	require := require.New(t)

	s := &MarahalerParent{
		42,
		MarshalerImplementor{
			[3]int{1, 2, 3},
			10,
		},
	}
	expected := map[string]interface{}{
		"AnInt": 42,
		"AnArrayIshStruct": map[string]int{
			"Arr0": 11,
			"Arr1": 12,
			"Arr2": 13,
		},
	}
	actual, err := maps.Marshal(s)

	require.NoError(err)
	require.Equal(expected, actual)
}
