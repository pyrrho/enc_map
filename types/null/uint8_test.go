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

func TestUint8Ctors(t *testing.T) {
	require := require.New(t)

	// null.NullUint8() returns a new null null.Uint8.
	// This is equivalent to null.Uint8{}.
	nul := null.NullUint8()
	require.False(nul.Valid)

	empty := null.Uint8{}
	require.False(empty.Valid)

	// null.NewUint8 constructs a new, valid null.Uint8.
	i := null.NewUint8(123)
	require.True(i.Valid)
	require.Equal(uint8(123), i.Uint8)

	z := null.NewUint8(0)
	require.True(z.Valid)
	require.Equal(uint8(0), z.Uint8)
}

func TestUint8ValueOrZero(t *testing.T) {
	require := require.New(t)

	valid := null.NewUint8(123)
	require.Equal(uint8(123), valid.Uint8)

	nul := null.Uint8{}
	require.Equal(uint8(0), nul.Uint8)
}

func TestUint8Set(t *testing.T) {
	require := require.New(t)

	i := null.Uint8{}
	require.False(i.Valid)

	i.Set(123)
	require.True(i.Valid)
	require.Equal(uint8(123), i.Uint8)

	i.Set(0)
	require.True(i.Valid)
	require.Equal(uint8(0), i.Uint8)
}

func TestUint8Null(t *testing.T) {
	require := require.New(t)

	i := null.NewUint8(123)

	i.Null()
	require.False(i.Valid)
}

func TestUint8IsNil(t *testing.T) {
	require := require.New(t)

	i := null.NewUint8(123)
	require.False(i.IsNil())

	z := null.NewUint8(0)
	require.False(z.IsNil())

	nul := null.Uint8{}
	require.True(nul.IsNil())
}

func TestUint8IsZero(t *testing.T) {
	require := require.New(t)

	i := null.NewUint8(123)
	require.False(i.IsZero())

	z := null.NewUint8(0)
	require.True(z.IsZero())

	nul := null.Uint8{}
	require.True(nul.IsZero())
}

func TestUint8SQLValue(t *testing.T) {
	require := require.New(t)
	var val driver.Value
	var err error

	i := null.NewUint8(123)
	val, err = i.Value()
	require.NoError(err)
	require.Equal(int64(123), val)

	z := null.NewUint8(0)
	val, err = z.Value()
	require.NoError(err)
	require.Equal(int64(0), val)

	nul := null.Uint8{}
	val, err = nul.Value()
	require.NoError(err)
	require.Equal(nil, val)
}

func TestUint8SQLScan(t *testing.T) {
	require := require.New(t)
	var err error

	var i null.Uint8
	err = i.Scan(123)
	require.NoError(err)
	require.True(i.Valid)
	require.Equal(uint8(123), i.Uint8)

	var u8Str null.Uint8
	// NB. Scan will coerce strings, but UnmarshalJSON won't.
	err = u8Str.Scan("123")
	require.NoError(err)
	require.True(u8Str.Valid)
	require.Equal(uint8(123), u8Str.Uint8)

	var nul null.Uint8
	err = nul.Scan(nil)
	require.NoError(err)
	require.False(nul.Valid)

	var wrong null.Uint8
	err = wrong.Scan("hello world")
	require.Error(err)

	var overflow null.Uint8
	err = overflow.Scan(12345)
	require.Error(err)

	var negative null.Uint8
	err = negative.Scan(-123)
	require.Error(err)

	var f null.Uint8
	err = f.Scan(1.2345)
	require.Error(err)

	var b null.Uint8
	err = b.Scan(true)
	require.Error(err)
}

func TestUint8MarshalJSON(t *testing.T) {
	require := require.New(t)
	var data []byte
	var err error

	i := null.NewUint8(123)
	data, err = json.Marshal(i)
	require.NoError(err)
	require.EqualValues("123", data)
	data, err = json.Marshal(&i)
	require.NoError(err)
	require.EqualValues("123", data)

	z := null.NewUint8(0)
	data, err = json.Marshal(z)
	require.NoError(err)
	require.EqualValues("0", data)
	data, err = json.Marshal(&z)
	require.NoError(err)
	require.EqualValues("0", data)

	nul := null.Uint8{}
	data, err = json.Marshal(nul)
	require.NoError(err)
	require.EqualValues("null", data)
	data, err = json.Marshal(&nul)
	require.NoError(err)
	require.EqualValues("null", data)
}

func TestUint8UnmarshalJSON(t *testing.T) {
	require := require.New(t)
	var err error

	// Successful Valid Parses

	var i null.Uint8
	err = json.Unmarshal([]byte("123"), &i)
	require.NoError(err)
	require.True(i.Valid)
	require.Equal(uint8(123), i.Uint8)

	// Successful Null Parses

	var nul null.Uint8
	err = json.Unmarshal([]byte("null"), &nul)
	require.NoError(err)
	require.False(nul.Valid)

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var intStr null.Uint8
	// Ints wrapped in quotes aren't ints.
	err = json.Unmarshal([]byte(`"123"`), &intStr)
	require.Error(err)

	var empty null.Uint8
	err = json.Unmarshal([]byte(""), &empty)
	require.Error(err)

	var quotes null.Uint8
	err = json.Unmarshal([]byte(`""`), &quotes)
	require.Error(err)

	var f null.Uint8
	// Non-integer numbers should not be coerced to ints.
	err = json.Unmarshal([]byte("1.2345"), &f)
	require.Error(err)

	var overflow null.Uint8
	err = json.Unmarshal([]byte("12345"), &overflow)
	require.Error(err)

	var negative null.Uint8
	err = json.Unmarshal([]byte("-123"), &negative)
	require.Error(err)

	var invalid null.Uint8
	err = invalid.UnmarshalJSON([]byte(":->"))
	if _, ok := err.(*json.SyntaxError); !ok {
		require.FailNowf(
			"Unexpected Error Type",
			"expected *json.SyntaxError, not %T", err)
	}
}

func TestUint8UnmarshalJSONOverflow(t *testing.T) {
	require := require.New(t)
	var err error

	uint8Overflow := uint64(math.MaxUint8)

	// Max uint8 should decode successfully
	var i null.Uint8
	err = json.Unmarshal([]byte(strconv.FormatUint(uint8Overflow, 10)), &i)
	require.NoError(err)

	// Attempt to overflow
	uint8Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(uint8Overflow, 10)), &i)
	// Decoded values should overflow uint8
	require.Error(err)
}

func TestUint8MarshalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ Uint8 null.Uint8 }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	wrapper = Wrapper{null.NewUint8(123)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Uint8": uint8(123)}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Uint8": uint8(123)}, data)

	wrapper = Wrapper{null.NewUint8(0)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Uint8": uint8(0)}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Uint8": uint8(0)}, data)

	// Null NullUint8s should be encoded as "nil"
	wrapper = Wrapper{null.Uint8{}}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Uint8": nil}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Uint8": nil}, data)
}
