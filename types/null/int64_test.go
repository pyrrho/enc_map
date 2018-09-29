package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"math"
	"strconv"
	"testing"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types/null"
	"github.com/stretchr/testify/require"
)

func TestInt64Ctors(t *testing.T) {
	require := require.New(t)

	// null.NullInt64() returns a new null null.Int64.
	// This is equivalent to null.Int64{}.
	nul := null.NullInt64()
	require.False(nul.Valid)

	empty := null.Int64{}
	require.False(empty.Valid)

	// null.NewInt64 constructs a new, valid null.Int64.
	i := null.NewInt64(12345)
	require.True(i.Valid)
	require.Equal(int64(12345), i.Int64)

	z := null.NewInt64(0)
	require.True(z.Valid)
	require.Equal(int64(0), z.Int64)
}

func TestInt64ValueOrZero(t *testing.T) {
	require := require.New(t)

	valid := null.NewInt64(12345)
	require.Equal(int64(12345), valid.Int64)

	nul := null.Int64{}
	require.Equal(int64(0), nul.Int64)
}

func TestInt64Set(t *testing.T) {
	require := require.New(t)

	i := null.Int64{}
	require.False(i.Valid)

	i.Set(12345)
	require.True(i.Valid)
	require.Equal(int64(12345), i.Int64)

	i.Set(0)
	require.True(i.Valid)
	require.Equal(int64(0), i.Int64)
}

func TestInt64Null(t *testing.T) {
	require := require.New(t)

	i := null.NewInt64(12345)

	i.Null()
	require.False(i.Valid)
}

func TestInt64IsNil(t *testing.T) {
	require := require.New(t)

	i := null.NewInt64(12345)
	require.False(i.IsNil())

	z := null.NewInt64(0)
	require.False(z.IsNil())

	nul := null.Int64{}
	require.True(nul.IsNil())
}

func TestInt64IsZero(t *testing.T) {
	require := require.New(t)

	i := null.NewInt64(12345)
	require.False(i.IsZero())

	z := null.NewInt64(0)
	require.True(z.IsZero())

	nul := null.Int64{}
	require.True(nul.IsZero())
}

func TestInt64SQLValue(t *testing.T) {
	require := require.New(t)
	var val driver.Value
	var err error

	i := null.NewInt64(12345)
	val, err = i.Value()
	require.NoError(err)
	require.Equal(int64(12345), val)

	z := null.NewInt64(0)
	val, err = z.Value()
	require.NoError(err)
	require.Equal(int64(0), val)

	nul := null.Int64{}
	val, err = nul.Value()
	require.NoError(err)
	require.Equal(nil, val)
}

func TestInt64SQLScan(t *testing.T) {
	require := require.New(t)
	var err error

	var i null.Int64
	err = i.Scan(12345)
	require.NoError(err)
	require.True(i.Valid)
	require.Equal(int64(12345), i.Int64)

	var i64Str null.Int64
	// NB. Scan will coerce strings, but UnmarshalJSON won't.
	err = i64Str.Scan("12345")
	require.NoError(err)
	require.True(i64Str.Valid)
	require.Equal(int64(12345), i64Str.Int64)

	var nul null.Int64
	err = nul.Scan(nil)
	require.NoError(err)
	require.False(nul.Valid)

	var wrong null.Int64
	err = wrong.Scan("hello world")
	require.Error(err)

	var f null.Int64
	err = f.Scan(1.2345)
	require.Error(err)

	var b null.Int64
	err = b.Scan(true)
	require.Error(err)
}

func TestInt64MarshalJSON(t *testing.T) {
	require := require.New(t)
	var data []byte
	var err error

	i := null.NewInt64(12345)
	data, err = json.Marshal(i)
	require.NoError(err)
	require.EqualValues("12345", data)
	data, err = json.Marshal(&i)
	require.NoError(err)
	require.EqualValues("12345", data)

	z := null.NewInt64(0)
	data, err = json.Marshal(z)
	require.NoError(err)
	require.EqualValues("0", data)
	data, err = json.Marshal(&z)
	require.NoError(err)
	require.EqualValues("0", data)

	nul := null.Int64{}
	data, err = json.Marshal(nul)
	require.NoError(err)
	require.EqualValues("null", data)
	data, err = json.Marshal(&nul)
	require.NoError(err)
	require.EqualValues("null", data)
}

func TestInt64UnmarshalJSON(t *testing.T) {
	require := require.New(t)
	var err error

	// Successful Valid Parses

	var i null.Int64
	err = json.Unmarshal([]byte("12345"), &i)
	require.NoError(err)
	require.True(i.Valid)
	require.Equal(int64(12345), i.Int64)

	// Successful Null Parses

	var nul null.Int64
	err = json.Unmarshal([]byte("null"), &nul)
	require.NoError(err)
	require.False(nul.Valid)

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var intStr null.Int64
	// Ints wrapped in quotes aren't ints.
	err = json.Unmarshal([]byte(`"12345"`), &intStr)
	require.Error(err)

	var empty null.Int64
	err = json.Unmarshal([]byte(""), &empty)
	require.Error(err)

	var quotes null.Int64
	err = json.Unmarshal([]byte(`""`), &quotes)
	require.Error(err)

	var f null.Int64
	// Non-integer numbers should not be coerced to ints.
	err = json.Unmarshal([]byte("1.2345"), &f)
	require.Error(err)

	var invalid null.Int64
	err = invalid.UnmarshalJSON([]byte(":->"))
	if _, ok := err.(*json.SyntaxError); !ok {
		require.FailNowf(
			"Unexpected Error Type",
			"expected *json.SyntaxError, not %T", err)
	}
}

func TestInt64UnmarshalJSONOverflow(t *testing.T) {
	require := require.New(t)
	var err error

	int64Overflow := uint64(math.MaxInt64)

	// Max int64 should decode successfully
	var i null.Int64
	err = json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	require.NoError(err)

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	// Decoded values should overflow int64
	require.Error(err)
}

func TestInt64MarshalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ Int64 null.Int64 }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	wrapper = Wrapper{null.NewInt64(12345)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Int64": int64(12345)}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Int64": int64(12345)}, data)

	wrapper = Wrapper{null.NewInt64(0)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Int64": int64(0)}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Int64": int64(0)}, data)

	// Null NullInt64s should be encoded as "nil"
	wrapper = Wrapper{null.Int64{}}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Int64": nil}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Int64": nil}, data)
}
