package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types/null"
	"github.com/stretchr/testify/require"
)

func TestStringCtors(t *testing.T) {
	require := require.New(t)

	// null.NullString() returns a new null null.String.
	// This is equivalent to null.String{}.
	nul := null.NullString()
	require.False(nul.Valid)

	empty := null.String{}
	require.False(empty.Valid)

	// null.NewString constructs a new, valid null.String.
	s := null.NewString("Hello World")
	require.True(s.Valid)
	require.Equal("Hello World", s.String)

	qs := null.NewString("")
	require.True(qs.Valid)
	require.Equal("", qs.String)
}

func TestStringValueOrZero(t *testing.T) {
	require := require.New(t)

	valid := null.NewString("test")
	require.Equal("test", valid.String)

	nul := null.String{}
	require.Equal("", nul.String)
}

func TestStringSet(t *testing.T) {
	require := require.New(t)

	s := null.String{}
	require.False(s.Valid)

	s.Set("test")
	require.True(s.Valid)
	require.Equal("test", s.String)

	s.Set("")
	require.True(s.Valid)
	require.Equal("", s.String)
}

func TestStringNull(t *testing.T) {
	require := require.New(t)

	s := null.NewString("test")

	s.Null()
	require.False(s.Valid)
}

func TestStringIsNil(t *testing.T) {
	require := require.New(t)

	s := null.NewString("test")
	require.False(s.IsNil())

	qs := null.NewString("")
	require.False(qs.IsNil())

	nul := null.String{}
	require.True(nul.IsNil())
}

func TestStringIsZero(t *testing.T) {
	require := require.New(t)

	s := null.NewString("test")
	require.False(s.IsZero())

	qs := null.NewString("")
	require.True(qs.IsZero())

	nul := null.String{}
	require.True(nul.IsZero())
}

func TestStringSQLValue(t *testing.T) {
	require := require.New(t)
	var val driver.Value
	var err error

	s := null.NewString("test")
	val, err = s.Value()
	require.NoError(err)
	require.Equal("test", val)

	nul := null.String{}
	val, err = nul.Value()
	require.NoError(err)
	require.Equal(nil, val)
}

func TestStringSQLScan(t *testing.T) {
	require := require.New(t)
	var err error

	var str null.String
	err = str.Scan("test")
	require.NoError(err)
	require.True(str.Valid)
	require.Equal("test", str.String)

	var empty null.String
	err = empty.Scan("")
	require.NoError(err)
	require.True(empty.Valid)
	require.Equal("", empty.String)

	var nul null.String
	err = nul.Scan(nil)
	require.NoError(err)
	require.False(nul.Valid)

	// NB. Scan is aggressive about converting values to strings. UnmarshalJSON
	// are less so.
	var i null.String
	err = i.Scan(12345)
	require.NoError(err)
	require.True(i.Valid)
	require.Equal("12345", i.String)

	var f null.String
	err = f.Scan(1.2345)
	require.NoError(err)
	require.True(f.Valid)
	require.Equal("1.2345", f.String)

	var b null.String
	err = b.Scan(true)
	require.NoError(err)
	require.True(b.Valid)
	require.Equal("true", b.String)
}

func TestStringMarshalJSON(t *testing.T) {
	require := require.New(t)
	var data []byte
	var err error

	str := null.NewString("test")
	data, err = json.Marshal(str)
	require.NoError(err)
	require.EqualValues(`"test"`, data)
	data, err = json.Marshal(&str)
	require.NoError(err)
	require.EqualValues(`"test"`, data)

	zero := null.NewString("")
	data, err = json.Marshal(zero)
	require.NoError(err)
	require.EqualValues(`""`, data)
	data, err = json.Marshal(&zero)
	require.NoError(err)
	require.EqualValues(`""`, data)

	null := null.String{}
	data, err = json.Marshal(null)
	require.NoError(err)
	require.EqualValues("null", data)
	data, err = json.Marshal(&null)
	require.NoError(err)
	require.EqualValues("null", data)
}

func TestStringMarshalJSONInStruct(t *testing.T) {
	require := require.New(t)
	var sj []byte
	var err error

	type stringTestStruct struct {
		NullString null.String `json:"null_string"`
		String     string      `json:"string"`
	}

	s := stringTestStruct{
		NullString: null.NewString("valid"),
		String:     "test",
	}
	sj, err = json.Marshal(s)
	require.NoError(err)
	require.EqualValues(`{"null_string":"valid","string":"test"}`, sj)

	s = stringTestStruct{
		NullString: null.String{},
		String:     "test",
	}
	sj, err = json.Marshal(s)
	require.NoError(err)
	require.EqualValues(`{"null_string":null,"string":"test"}`, sj)
}

func TestStringUnmarshalJSON(t *testing.T) {
	require := require.New(t)
	var err error

	// Successful Valid Parses

	var str null.String
	err = json.Unmarshal([]byte(`"test"`), &str)
	require.NoError(err)
	require.True(str.Valid)
	require.Equal("test", str.String)

	var quotes null.String
	err = json.Unmarshal([]byte(`""`), &quotes)
	require.NoError(err)
	require.True(quotes.Valid)
	require.Equal("", quotes.String)

	var nullStr null.String
	err = json.Unmarshal([]byte(`"null"`), &nullStr)
	require.NoError(err)
	require.True(nullStr.Valid)
	require.Equal("null", nullStr.String)

	// Successful Null Parses

	var nul null.String
	err = json.Unmarshal([]byte("null"), &nul)
	require.NoError(err)
	require.False(nul.Valid)

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var badType null.String
	// Ints are never string.
	err = json.Unmarshal([]byte("12345"), &badType)
	require.Error(err)

	var invalid null.String
	err = invalid.UnmarshalJSON([]byte(":->"))
	if _, ok := err.(*json.SyntaxError); !ok {
		require.FailNowf(
			"Unexpected Error Type",
			"expected *json.SyntaxError, not %T", err)
	}
}

func TestStringMarshalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ Slice null.String }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	wrapper = Wrapper{null.NewString("test")}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Slice": "test"}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Slice": "test"}, data)

	wrapper = Wrapper{null.NewString("")}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Slice": ""}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Slice": ""}, data)

	// Null NullStrings should be encoded as "nil"
	wrapper = Wrapper{null.String{}}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Slice": nil}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Slice": nil}, data)
}
