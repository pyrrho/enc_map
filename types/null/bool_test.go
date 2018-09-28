package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types/null"
	"github.com/stretchr/testify/require"
)

func TestBoolCtors(t *testing.T) {
	require := require.New(t)

	// null.NullBool() returns a new null null.Bool.
	// This is equivalent to null.Bool{}.
	nul := null.NullBool()
	require.False(nul.Valid)

	empty := null.Bool{}
	require.False(empty.Valid)

	// null.NewBool constructs a new, valid null.Bool.
	tr := null.NewBool(true)
	require.True(tr.Valid)
	require.Equal(true, tr.Bool)

	fl := null.NewBool(false)
	require.True(fl.Valid)
	require.Equal(false, fl.Bool)
}

func TestBoolValueOrZero(t *testing.T) {
	require := require.New(t)

	valid := null.NewBool(true)
	require.Equal(true, valid.ValueOrZero())

	nul := null.Bool{}
	require.Equal(false, nul.ValueOrZero())
}
func TestBoolSet(t *testing.T) {
	require := require.New(t)

	b := null.Bool{}
	require.False(b.Valid)

	b.Set(true)
	require.True(b.Valid)
	require.Equal(true, b.Bool)

	b.Set(false)
	require.True(b.Valid)
	require.Equal(false, b.Bool)
}

func TestBoolNull(t *testing.T) {
	require := require.New(t)

	b := null.NewBool(true)

	b.Null()
	require.False(b.Valid)
}

func TestBoolIsNil(t *testing.T) {
	require := require.New(t)

	tr := null.NewBool(true)
	require.False(tr.IsNil())

	fl := null.NewBool(false)
	require.False(fl.IsNil())

	nul := null.Bool{}
	require.True(nul.IsNil())
}

func TestBoolIsZero(t *testing.T) {
	require := require.New(t)

	tr := null.NewBool(true)
	require.False(tr.IsZero())

	fl := null.NewBool(false)
	require.True(fl.IsZero())

	nul := null.Bool{}
	require.True(nul.IsZero())
}

func TestBoolSQLValue(t *testing.T) {
	require := require.New(t)
	var val driver.Value
	var err error

	tr := null.NewBool(true)
	val, err = tr.Value()
	require.NoError(err)
	require.Equal(true, val)

	fl := null.NewBool(false)
	val, err = fl.Value()
	require.NoError(err)
	require.Equal(false, val)

	nul := null.Bool{}
	val, err = nul.Value()
	require.NoError(err)
	require.Equal(nil, val)
}

func TestBoolSQLScan(t *testing.T) {
	require := require.New(t)
	var err error

	var tr null.Bool
	err = tr.Scan(true)
	require.NoError(err)
	require.True(tr.Valid)
	require.Equal(true, tr.Bool)

	var fl null.Bool
	err = fl.Scan(false)
	require.NoError(err)
	require.True(fl.Valid)
	require.Equal(false, fl.Bool)

	var nul null.Bool
	err = nul.Scan(nil)
	require.NoError(err)
	require.False(nul.Valid)

	var wrong null.Bool
	err = wrong.Scan(int64(42))
	require.Error(err)
}

func TestBoolMarshalJSON(t *testing.T) {
	require := require.New(t)
	var data []byte
	var err error

	b := null.NewBool(true)
	data, err = json.Marshal(b)
	require.NoError(err)
	require.EqualValues("true", data)
	data, err = json.Marshal(&b)
	require.NoError(err)
	require.EqualValues("true", data)

	zero := null.NewBool(false)
	data, err = json.Marshal(zero)
	require.NoError(err)
	require.EqualValues("false", data)
	data, err = json.Marshal(&zero)
	require.NoError(err)
	require.EqualValues("false", data)

	// Null NullBools should be encoded as "null"
	nul := null.Bool{}
	data, err = json.Marshal(nul)
	require.NoError(err)
	require.EqualValues("null", data)
	data, err = json.Marshal(&nul)
	require.NoError(err)
	require.EqualValues("null", data)

	wrapper := struct {
		Foo null.Bool
		Bar null.Bool
	}{
		null.NewBool(true),
		null.Bool{},
	}
	data, err = json.Marshal(wrapper)
	require.NoError(err)
	require.EqualValues(`{"Foo":true,"Bar":null}`, data)
	data, err = json.Marshal(&wrapper)
	require.NoError(err)
	require.EqualValues(`{"Foo":true,"Bar":null}`, data)
}

func TestBoolUnmarshalJSON(t *testing.T) {
	require := require.New(t)
	var err error

	// Successful Valid Parses

	var tr null.Bool
	err = json.Unmarshal([]byte("true"), &tr)
	require.NoError(err)
	require.True(tr.Valid)
	require.Equal(true, tr.Bool)

	var fl null.Bool
	err = json.Unmarshal([]byte("false"), &fl)
	require.NoError(err)
	require.True(fl.Valid)
	require.Equal(false, fl.Bool)

	// Successful Null Parses

	var nul null.Bool
	err = json.Unmarshal([]byte("null"), &nul)
	require.NoError(err)
	require.False(nul.Valid)

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var str null.Bool
	// Booleans wrapped in quotes aren't booleans.
	err = json.Unmarshal([]byte(`"true"`), &str)
	require.Error(err)

	var empty null.Bool
	// An empty string is not a boolean.
	err = json.Unmarshal([]byte(`""`), &empty)
	require.Error(err)

	var badType null.Bool
	// Ints are never booleans.
	err = json.Unmarshal([]byte("1"), &badType)
	require.Error(err)

	var invalid null.Bool
	err = invalid.UnmarshalJSON([]byte(":->"))
	if _, ok := err.(*json.SyntaxError); !ok {
		require.FailNowf(
			"Unexpected Error Type",
			"expected *json.SyntaxError, not %T", err)
	}
}

func TestBoolMarshalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ Bool null.Bool }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	wrapper = Wrapper{null.NewBool(true)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Bool": true}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Bool": true}, data)

	wrapper = Wrapper{null.NewBool(false)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Bool": false}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Bool": false}, data)

	// Null NullBools should be encoded as "nil"
	wrapper = Wrapper{null.Bool{}}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Bool": nil}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Bool": nil}, data)
}
