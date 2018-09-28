package null_test

import (
	"encoding/json"
	"testing"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types/null"
	"github.com/stretchr/testify/require"
)

func base64ed(k string) []byte {
	var m = map[string][]byte{
		"DAICON V":   []byte("REFJQ09OIFY="),
		`"DAICON V"`: []byte(`"REFJQ09OIFY="`),
	}
	return m[k]
}

func TestByteSliceCtors(t *testing.T) {
	require := require.New(t)

	// null.NullByteSlice() retuns a new null null.ByteSlice.
	// This is equivalent to null.ByteSlice{}.
	nul := null.NullByteSlice()
	require.False(nul.Valid)

	empty := null.ByteSlice{}
	require.False(empty.Valid)

	// null.NewByteSlice constructs a new, valid, possibly zero-length
	// null.ByteSlice ...
	b := null.NewByteSlice([]byte("DAICON V"))
	require.True(b.Valid)
	require.Equal([]byte("DAICON V"), b.ByteSlice)

	b2 := null.NewByteSlice([]byte(""))
	require.True(b2.Valid)
	require.Equal([]byte{}, b2.ByteSlice)

	b3 := null.NewByteSlice([]byte{})
	require.True(b3.Valid)
	require.Equal([]byte{}, b3.ByteSlice)

	// ... unless a nil []byte is passed in.
	nul2 := null.NewByteSlice([]byte(nil))
	require.False(nul2.Valid)

	nul3 := null.NewByteSlice(nil)
	require.False(nul3.Valid)

	tmp := []byte(nil)
	nul4 := null.NewByteSlice(tmp)
	require.False(nul4.Valid)

	// null.NewByteSliceStr constructs a new, valid null.ByteSlice, but it takes
	// a string, rather than a []byte.
	s := null.NewByteSliceStr("DAICON V")
	require.True(s.Valid)
	require.EqualValues("DAICON V", s.ByteSlice)

	s2 := null.NewByteSliceStr("")
	require.True(s2.Valid)
	require.EqualValues([]byte{}, s2.ByteSlice)

	// You can also construct a null.ByteSlice with a base64 encoded string.
	// Note that this form may return an error.
	b3, err := null.NewByteSliceFromBase64([]byte("REFJQ09OIFY="))
	require.NoError(err)
	require.True(b3.Valid)
	require.EqualValues("DAICON V", b3.ByteSlice)

	s3, err := null.NewByteSliceFromBase64Str("REFJQ09OIFY=")
	require.NoError(err)
	require.True(s3.Valid)
	require.EqualValues("DAICON V", s3.ByteSlice)
}

func TestByteSliceValueOrZero(t *testing.T) {
	require := require.New(t)

	b := null.NewByteSliceStr("DAICON V")
	require.Equal([]byte("DAICON V"), b.ValueOrZero())

	n := null.ByteSlice{}
	require.Equal([]byte{}, n.ValueOrZero())
}

func TestByteSliceSet(t *testing.T) {
	require := require.New(t)

	bs := null.ByteSlice{}

	bs.Set([]byte("DAICON V"))
	require.True(bs.Valid)
	require.Equal([]byte("DAICON V"), bs.ByteSlice)

	bs.Set([]byte(nil))
	require.False(bs.Valid)

	bs.Set([]byte{})
	require.True(bs.Valid)
	require.Equal([]byte{}, bs.ByteSlice)

	bs.Set(nil)
	require.False(bs.Valid)

	bs.Set([]byte("Hello again!"))
	require.True(bs.Valid)
	require.Equal([]byte("Hello again!"), bs.ByteSlice)
}

func TestByteSliceSetStr(t *testing.T) {
	require := require.New(t)

	bs := null.ByteSlice{}

	bs.SetStr("DAICON V")
	require.True(bs.Valid)
	require.Equal([]byte("DAICON V"), bs.ByteSlice)

	bs.SetStr("")
	require.True(bs.Valid)
	require.Equal([]byte{}, bs.ByteSlice)

	bs.SetStr("Hello again!")
	require.True(bs.Valid)
	require.Equal([]byte("Hello again!"), bs.ByteSlice)
}

func TestByteSliceNull(t *testing.T) {
	require := require.New(t)

	bs := null.NewByteSliceStr("DAICON V")

	bs.Null()
	require.False(bs.Valid)
}

func TestByteSliceIsNil(t *testing.T) {
	require := require.New(t)

	bs := null.NewByteSliceStr("DAICON V")
	require.False(bs.IsNil())

	empty := null.NewByteSlice([]byte{})
	require.False(empty.IsNil())

	nul := null.ByteSlice{}
	require.True(nul.IsNil())
}

func TestByteSliceIsZero(t *testing.T) {
	require := require.New(t)

	bs := null.NewByteSliceStr("DAICON V")
	require.False(bs.IsZero())

	empty := null.NewByteSlice([]byte{})
	require.True(empty.IsZero())

	nul := null.ByteSlice{}
	require.True(nul.IsZero())
}

func TestByteSliceSQLValue(t *testing.T) {
	require := require.New(t)

	bs := null.NewByteSliceStr("DAICON V")
	val, err := bs.Value()
	require.NoError(err)
	require.Equal(base64ed("DAICON V"), val)

	empty := null.NewByteSlice([]byte{})
	val, err = empty.Value()
	require.NoError(err)
	require.Equal([]byte{}, val)

	nul := null.ByteSlice{}
	val, err = nul.Value()
	require.NoError(err)
	require.Equal(nil, val)
}

func TestByteSliceSQLScan(t *testing.T) {
	require := require.New(t)

	var bs null.ByteSlice
	err := bs.Scan(base64ed("DAICON V"))
	require.NoError(err)
	require.True(bs.Valid)
	require.Equal([]byte("DAICON V"), bs.ByteSlice)

	var str null.ByteSlice
	err = str.Scan(string(base64ed("DAICON V")))
	require.NoError(err)
	require.True(str.Valid)
	require.Equal([]byte("DAICON V"), str.ByteSlice)

	var empty null.ByteSlice
	err = empty.Scan([]byte{})
	require.NoError(err)
	require.True(empty.Valid)
	require.Equal([]byte{}, empty.ByteSlice)

	var nul null.ByteSlice
	err = nul.Scan(nil)
	require.NoError(err)
	require.False(nul.Valid)

	var wrong null.ByteSlice
	err = wrong.Scan(int64(42))
	require.Error(err)
}

func TestByteSliceMarshalJSON(t *testing.T) {
	require := require.New(t)

	bs := null.NewByteSliceStr("DAICON V")
	data, err := json.Marshal(bs)
	require.NoError(err)
	require.EqualValues(base64ed(`"DAICON V"`), data)
	data, err = json.Marshal(&bs)
	require.NoError(err)
	require.EqualValues(base64ed(`"DAICON V"`), data)

	empty := null.NewByteSlice([]byte{})
	data, err = json.Marshal(empty)
	require.NoError(err)
	require.EqualValues(`""`, data)
	data, err = json.Marshal(&empty)
	require.NoError(err)
	require.EqualValues(`""`, data)

	nul := null.ByteSlice{}
	data, err = json.Marshal(nul)
	require.NoError(err)
	require.EqualValues("null", data)
	data, err = json.Marshal(&nul)
	require.NoError(err)
	require.EqualValues("null", data)
}

func TestByteSliceUnmarshalJSON(t *testing.T) {
	require := require.New(t)
	var err error

	// Successful Valid Parses

	var bs null.ByteSlice
	err = json.Unmarshal(base64ed(`"DAICON V"`), &bs)
	require.NoError(err)
	require.True(bs.Valid)
	require.Equal([]byte("DAICON V"), bs.ByteSlice)

	var quotes null.ByteSlice
	err = json.Unmarshal([]byte(`""`), &quotes)
	require.NoError(err)
	require.True(quotes.Valid)
	require.Equal([]byte(""), quotes.ByteSlice)

	var nullStrQuoted null.ByteSlice
	err = json.Unmarshal([]byte(`"null"`), &nullStrQuoted)
	require.NoError(err)
	// Skip checking what this decoded to; it's garbage.

	// Successful Null Parses

	var nullStr null.ByteSlice
	err = json.Unmarshal([]byte("null"), &nullStr)
	require.NoError(err)
	require.False(nullStr.Valid)

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var badType null.ByteSlice
	// Ints are never byte slices.
	err = json.Unmarshal([]byte("12345"), &badType)
	require.Error(err)

	var invalid null.ByteSlice
	err = invalid.UnmarshalJSON([]byte(":->"))
	if _, ok := err.(*json.SyntaxError); !ok {
		require.FailNowf(
			"Unexpected Error Type",
			"expected *json.SyntaxError, not %T", err)
	}
}

func TestByteSliceMarshalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ Slice null.ByteSlice }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	wrapper = Wrapper{null.NewByteSlice([]byte("DAICON V"))}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Slice": base64ed("DAICON V")}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Slice": base64ed("DAICON V")}, data)

	wrapper = Wrapper{null.NewByteSlice([]byte{})}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Slice": []byte{}}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Slice": []byte{}}, data)

	// Null NullByteSlices should be encoded as "nil"
	wrapper = Wrapper{null.ByteSlice{}}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Slice": nil}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Slice": nil}, data)
}
