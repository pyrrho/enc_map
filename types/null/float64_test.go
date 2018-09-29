package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"math"
	"testing"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types/null"
	"github.com/stretchr/testify/require"
)

func TestFloat64Ctors(t *testing.T) {
	require := require.New(t)

	// null.NullFloat64() returns a new null null.Float64.
	// This is equivalent to null.Float64{}.
	nul := null.NullFloat64()
	require.False(nul.Valid)

	empty := null.Float64{}
	require.False(empty.Valid)

	// null.NewFloat64 constructs a new, valid null.Float64.
	f := null.NewFloat64(1.2345)
	require.True(f.Valid)
	require.Equal(1.2345, f.Float64)

	i := null.NewFloat64(12345)
	require.True(i.Valid)
	require.Equal(float64(12345), i.Float64)

	z := null.NewFloat64(0)
	require.True(z.Valid)
	require.Equal(0.0, z.Float64)
}

func TestFloat64ValueOrZero(t *testing.T) {
	require := require.New(t)

	f := null.NewFloat64(1.2345)
	require.Equal(1.2345, f.ValueOrZero())

	i := null.NewFloat64(12345)
	require.Equal(float64(12345), i.ValueOrZero())

	z := null.NewFloat64(0)
	require.Equal(0.0, z.ValueOrZero())

	nul := null.Float64{}
	require.Equal(0.0, nul.ValueOrZero())
}

func TestFloat64Set(t *testing.T) {
	require := require.New(t)

	f := null.Float64{}
	require.False(f.Valid)

	f.Set(1.2345)
	require.True(f.Valid)
	require.Equal(1.2345, f.Float64)

	f.Set(0.0)
	require.True(f.Valid)
	require.Equal(0.0, f.Float64)
}

func TestFloat64Null(t *testing.T) {
	require := require.New(t)

	f := null.NewFloat64(1.2345)

	f.Null()
	require.False(f.Valid)
}

func TestFloat64IsNil(t *testing.T) {
	require := require.New(t)

	f := null.NewFloat64(1.2345)
	require.False(f.IsNil())

	z := null.NewFloat64(0)
	require.False(z.IsNil())

	nul := null.Float64{}
	require.True(nul.IsNil())
}

func TestFloat64IsZero(t *testing.T) {
	require := require.New(t)

	f := null.NewFloat64(1.2345)
	require.False(f.IsZero())

	z := null.NewFloat64(0)
	require.True(z.IsZero())

	nul := null.Float64{}
	require.True(nul.IsZero())
}

func TestFloat64SQLValue(t *testing.T) {
	require := require.New(t)
	var val driver.Value
	var err error

	f := null.NewFloat64(1.2345)
	val, err = f.Value()
	require.NoError(err)
	require.Equal(1.2345, val)

	zero := null.NewFloat64(0)
	val, err = zero.Value()
	require.NoError(err)
	require.Equal(0.0, val)

	nul := null.Float64{}
	val, err = nul.Value()
	require.NoError(err)
	require.Equal(nil, val)
}

func TestFloat64SQLScan(t *testing.T) {
	require := require.New(t)

	var f null.Float64
	err := f.Scan(1.2345)
	require.NoError(err)
	require.True(f.Valid)
	require.Equal(1.2345, f.Float64)

	var i null.Float64
	err = i.Scan(12345)
	require.NoError(err)
	require.True(i.Valid)
	require.Equal(float64(12345), i.Float64)

	var f64Str null.Float64
	// NB. Scan will coerce strings, but UnmarshalJSON won't.
	err = f64Str.Scan("1.2345")
	require.NoError(err)
	require.True(f.Valid)
	require.Equal(1.2345, f.Float64)

	var nul null.Float64
	err = nul.Scan(nil)
	require.NoError(err)
	require.False(nul.Valid)

	var wrong null.Float64
	err = wrong.Scan("hello world")
	require.Error(err)
}

func TestFloat64MarshalJSON(t *testing.T) {
	require := require.New(t)
	var data []byte
	var err error

	f := null.NewFloat64(1.2345)
	data, err = json.Marshal(f)
	require.NoError(err)
	require.EqualValues("1.2345", data)
	data, err = json.Marshal(&f)
	require.NoError(err)
	require.EqualValues("1.2345", data)

	i := null.NewFloat64(12345)
	data, err = json.Marshal(i)
	require.NoError(err)
	require.EqualValues("12345", data)
	data, err = json.Marshal(&i)
	require.NoError(err)
	require.EqualValues("12345", data)

	zero := null.NewFloat64(0)
	data, err = json.Marshal(zero)
	require.NoError(err)
	require.EqualValues("0", data)
	data, err = json.Marshal(&zero)
	require.NoError(err)
	require.EqualValues("0", data)

	nul := null.Float64{}
	data, err = json.Marshal(nul)
	require.NoError(err)
	require.EqualValues("null", data)
	data, err = json.Marshal(&nul)
	require.NoError(err)
	require.EqualValues("null", data)

	nan := null.NewFloat64(math.NaN())
	data, err = json.Marshal(nan)
	require.Error(err)
	data, err = json.Marshal(&nan)
	require.Error(err)

	inf := null.NewFloat64(math.Inf(1))
	data, err = json.Marshal(inf)
	require.Error(err)
	data, err = json.Marshal(&inf)
	require.Error(err)
}

func TestFloat64UnmarshalJSON(t *testing.T) {
	require := require.New(t)
	var err error

	// Successful Valid Parses

	var f null.Float64
	err = json.Unmarshal([]byte("1.2345"), &f)
	require.NoError(err)
	require.True(f.Valid)
	require.Equal(1.2345, f.Float64)

	var i null.Float64
	err = json.Unmarshal([]byte("12345"), &i)
	require.NoError(err)
	require.True(i.Valid)
	require.Equal(float64(12345), i.Float64)

	// Successful Null Parses

	var nul null.Float64
	err = json.Unmarshal([]byte("null"), &nul)
	require.NoError(err)
	require.False(nul.Valid)

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var f64Str null.Float64
	// Floats wrapped in quotes aren't floats.
	err = json.Unmarshal([]byte(`"1.2345"`), &f64Str)
	require.Error(err)

	var empty null.Float64
	err = json.Unmarshal([]byte(""), &empty)
	require.Error(err)

	var quotes null.Float64
	err = json.Unmarshal([]byte(`""`), &quotes)
	require.Error(err)

	var badType null.Float64
	// Booleans are never floats.
	err = json.Unmarshal([]byte("true"), &badType)
	require.Error(err)

	// The JSON specification does not include NaN, INF, Infinity, NegInfinity
	// or any other common literal for the IEEE 754 floating point
	// not-really-number values. As such, un-marshaling them from JSON will
	// result in errors.
	var nan null.Float64
	err = json.Unmarshal([]byte("NaN"), &nan)
	require.Error(err)

	var inf null.Float64
	err = json.Unmarshal([]byte("INF"), &inf)
	require.Error(err)

	var invalid null.Float64
	err = invalid.UnmarshalJSON([]byte(":->"))
	if _, ok := err.(*json.SyntaxError); !ok {
		require.FailNowf(
			"Unexpected Error Type",
			"expected *json.SyntaxError, not %T", err)
	}
}

func TestFloat64MarshalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ Float64 null.Float64 }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	wrapper = Wrapper{null.NewFloat64(1.2345)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Float64": 1.2345}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Float64": 1.2345}, data)

	wrapper = Wrapper{null.NewFloat64(0)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Float64": 0.0}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Float64": 0.0}, data)

	// Null NullFloat64s should be encoded as "nil"
	wrapper = Wrapper{null.Float64{}}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Float64": nil}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Float64": nil}, data)
}
