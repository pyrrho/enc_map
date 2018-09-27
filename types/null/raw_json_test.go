package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types"
	"github.com/pyrrho/encoding/types/null"
)

func TestRawJSONCtors(t *testing.T) {
	require := require.New(t)

	// null.NullJSON() returns a new null null.RawJSON.
	// This is equivalent to null.RawJSON{}.
	nul := null.NullJSON()
	require.False(nul.Valid)

	empty := null.RawJSON{}
	require.False(empty.Valid)

	// null.NewJSON constructs a new, valid null.RawJSON.
	j := null.NewJSON(types.RawJSON(`"Hello World"`))
	require.True(j.Valid)
	require.EqualValues(`"Hello World"`, j.JSON)

	// If you want to, you can pass a []byte into null.NewJSON, because they are
	// silently converted into types.RawJSON on (function argument) copy.
	b := null.NewJSON([]byte(`"Hello World"`))
	require.True(b.Valid)
	require.EqualValues(`"Hello World"`, b.JSON)

	// Even easier is null.NewJSONStr, which takes a string.
	s := null.NewJSONStr(`"Hello World"`)
	require.True(s.Valid)
	require.EqualValues(`"Hello World"`, s.JSON)

	// If you give null.NewJSON or null.NewJSONStr a nil or zero-length
	// argument, you'll get a null null.RawJSON back.
	nil_ := null.NewJSON(nil)
	require.False(nil_.Valid)

	e1 := null.NewJSON(types.RawJSON{})
	require.False(e1.Valid)

	e2 := null.NewJSON(types.RawJSON(""))
	require.False(e2.Valid)

	e3 := null.NewJSON(types.RawJSON(nil))
	require.False(e3.Valid)

	e4 := null.NewJSON([]byte{})
	require.False(e4.Valid)

	e5 := null.NewJSON([]byte(""))
	require.False(e5.Valid)

	e6 := null.NewJSON([]byte(nil))
	require.False(e6.Valid)

	e7 := null.NewJSONStr("")
	require.False(e7.Valid)
}

func TestRawJSONValueOrZero(t *testing.T) {
	require := require.New(t)

	j := null.NewJSONStr(`"Hello World"`)
	require.EqualValues(`"Hello World"`, j.ValueOrZero())

	n := null.RawJSON{}
	require.EqualValues(types.RawJSON{}, n.ValueOrZero())
}

func TestRawJSONSet(t *testing.T) {
	require := require.New(t)

	j := null.RawJSON{}

	j.Set(types.RawJSON(`"Hello World"`))
	require.True(j.Valid)
	require.EqualValues(`"Hello World"`, j.JSON)

	j.Set(types.RawJSON{})
	require.False(j.Valid)

	j.Set([]byte(`"Hello again!"`))
	require.True(j.Valid)
	require.EqualValues(`"Hello again!"`, j.JSON)
}

func TestRawJSONSetStr(t *testing.T) {
	require := require.New(t)

	j := null.RawJSON{}

	j.SetStr(`"Hello World"`)
	require.True(j.Valid)
	require.EqualValues(`"Hello World"`, j.JSON)

	j.SetStr("")
	require.False(j.Valid)
}

func TestRawJSONNull(t *testing.T) {
	require := require.New(t)

	j := null.NewJSONStr(`"Hello World"`)

	j.Null()
	require.False(j.Valid)
}

func TestRawJSONIsNil(t *testing.T) {
	require := require.New(t)

	str := null.NewJSONStr(`"Hello World"`)
	require.False(str.IsNil())
	emptyStr := null.NewJSONStr(`""`)
	require.False(emptyStr.IsNil())

	num := null.NewJSONStr("42.0")
	require.False(num.IsNil())
	emptyNum := null.NewJSONStr("0.0")
	require.False(emptyNum.IsNil())

	obj := null.NewJSONStr(`{"foo":42.0,bar:"baz"}`)
	require.False(obj.IsNil())
	emptyObj := null.NewJSONStr("{}")
	require.False(emptyObj.IsNil())

	arr := null.NewJSONStr("[1.0, 2.0, 3.0]")
	require.False(arr.IsNil())
	emptyArr := null.NewJSONStr("[]")
	require.False(emptyArr.IsNil())

	bol := null.NewJSONStr("true")
	require.False(bol.IsNil())
	emptyBol := null.NewJSONStr("false")
	require.False(emptyBol.IsNil())

	nul := null.NewJSONStr("null")
	require.False(nul.IsNil())

	nil_ := null.NullJSON()
	require.True(nil_.IsNil())
}

func TestRawJSONIsZero(t *testing.T) {
	require := require.New(t)

	str := null.NewJSONStr(`"Hello World"`)
	require.False(str.IsZero())
	emptyStr := null.NewJSONStr(`""`)
	require.True(emptyStr.IsZero())

	num := null.NewJSONStr("42.0")
	require.False(num.IsZero())
	emptyNum := null.NewJSONStr("0.0")
	require.True(emptyNum.IsZero())

	obj := null.NewJSONStr(`{"foo":42.0,bar:"baz"}`)
	require.False(obj.IsZero())
	emptyObj := null.NewJSONStr("{}")
	require.True(emptyObj.IsZero())

	arr := null.NewJSONStr("[1.0, 2.0, 3.0]")
	require.False(arr.IsZero())
	emptyArr := null.NewJSONStr("[]")
	require.True(emptyArr.IsZero())

	bol := null.NewJSONStr("true")
	require.False(bol.IsZero())
	emptyBol := null.NewJSONStr("false")
	require.True(emptyBol.IsZero())

	nul := null.NewJSONStr("null")
	require.True(nul.IsZero())

	nil_ := null.NullJSON()
	require.True(nil_.IsZero())
}

func TestRawJSONSQLValue(t *testing.T) {
	var val driver.Value
	var err error
	require := require.New(t)

	// Empty null.RawJSON objects are handled correctly ...

	nul := null.NullJSON()
	val, err = nul.Value()
	require.NoError(err)
	require.EqualValues(nil, val)

	empty := null.RawJSON{}
	val, err = empty.Value()
	require.NoError(err)
	require.EqualValues(nil, val)

	// .. and other behavior is consistent with types.RawJSON.

	j := null.NewJSONStr(`{"foo":42.0,"bar":"baz"}`)
	val, err = j.Value()
	require.NoError(err)
	require.EqualValues(`{"foo":42.0,"bar":"baz"}`, val)

	invalid := null.NewJSONStr(":->")
	val, err = invalid.Value()
	require.Error(err)
}

func TestRawJSONSQLScan(t *testing.T) {
	require := require.New(t)
	var err error

	// nil values are handled correctly ...

	nil_ := null.NewJSONStr(`"Hello World"`)
	err = nil_.Scan(driver.Value(nil))
	require.NoError(err)
	require.False(nil_.Valid)

	empty1 := null.NewJSONStr(`"Hello World"`)
	err = empty1.Scan(driver.Value([]byte{}))
	require.NoError(err)
	require.False(empty1.Valid)

	empty2 := null.NewJSONStr(`"Hello World"`)
	err = empty2.Scan(driver.Value([]byte(nil)))
	require.NoError(err)
	require.False(empty2.Valid)

	empty3 := null.NewJSONStr(`"Hello World"`)
	err = empty3.Scan(driver.Value(""))
	require.NoError(err)
	require.False(empty3.Valid)

	// .. and other behavior is consistent with types.RawJSON.

	var jb null.RawJSON
	err = jb.Scan(driver.Value([]byte(`{"foo":42.0,"bar":"baz"}`)))
	require.NoError(err)
	require.True(jb.Valid)
	require.EqualValues(`{"foo":42.0,"bar":"baz"}`, jb.JSON)

	var js null.RawJSON
	err = js.Scan(driver.Value(`{"foo":42.0,"bar":"baz"}`))
	require.NoError(err)
	require.True(js.Valid)
	require.EqualValues(`{"foo":42.0,"bar":"baz"}`, js.JSON)

	var nulb null.RawJSON
	err = nulb.Scan(driver.Value([]byte("null")))
	require.NoError(err)
	require.True(nulb.Valid)
	require.EqualValues("null", nulb.JSON)

	var nuls null.RawJSON
	err = nuls.Scan(driver.Value("null"))
	require.NoError(err)
	require.True(nuls.Valid)
	require.EqualValues("null", nuls.JSON)

	// NB. Scan assumes valid JSON is provided; it will not validate.
	var invalidb null.RawJSON
	err = invalidb.Scan(driver.Value([]byte(":->")))
	require.NoError(err)
	require.True(invalidb.Valid)
	require.EqualValues(":->", invalidb.JSON)

	var invalids null.RawJSON
	err = invalids.Scan(driver.Value(":->"))
	require.NoError(err)
	require.True(invalids.Valid)
	require.EqualValues(":->", invalids.JSON)
}

func TestRawJSONMarshalJSON(t *testing.T) {
	require := require.New(t)
	var err error
	var data []byte

	// nil values are handled correctly ...

	empty := null.RawJSON{}
	data, err = json.Marshal(empty)
	require.NoError(err)
	require.EqualValues("null", data)
	data, err = json.Marshal(&empty)
	require.NoError(err)
	require.EqualValues("null", data)

	nul := null.NullJSON()
	data, err = json.Marshal(nul)
	require.NoError(err)
	require.EqualValues("null", data)
	data, err = json.Marshal(&nul)
	require.NoError(err)
	require.EqualValues("null", data)

	// .. and other behavior is consistent with types.RawJSON.

	j := null.NewJSONStr(`{"foo":42.0,"bar":"baz"}`)
	data, err = json.Marshal(j)
	require.NoError(err)
	require.EqualValues(`{"foo":42.0,"bar":"baz"}`, data)
	data, err = json.Marshal(&j)
	require.NoError(err)
	require.EqualValues(`{"foo":42.0,"bar":"baz"}`, data)

	// `bar` should be quoted                ~~~
	invalidObj := null.NewJSONStr(`{"foo":42.0,bar:"baz"}`)
	data, err = json.Marshal(invalidObj)
	require.Error(err)
	// This error should include information on the malformed object.
	require.Contains(err.Error(), "invalid character 'b'")
	data, err = json.Marshal(&invalidObj)
	require.Error(err)
	// This error should include information on the malformed object.
	require.Contains(err.Error(), "invalid character 'b'")
}

func TestRawJSONUnmarshalJSON(t *testing.T) {
	require := require.New(t)
	var err error

	// nil values are handled correctly ...

	var nul null.RawJSON
	err = json.Unmarshal([]byte("null"), &nul)
	require.NoError(err)
	require.False(nul.Valid)

	// ... zero-length JSON is still invalid, and generates errors ...

	var empty null.RawJSON
	err = json.Unmarshal([]byte{}, &empty)
	require.Error(err)

	var nilByte null.RawJSON
	err = json.Unmarshal([]byte(nil), &nilByte)
	require.Error(err)

	var nil_ null.RawJSON
	err = json.Unmarshal(nil, &nil_)
	require.Error(err)

	// .. and other behavior is consistent with types.RawJSON.

	var j null.RawJSON
	err = json.Unmarshal([]byte(`{"foo":42.0,"bar":"baz"}`), &j)
	require.NoError(err)
	require.True(j.Valid)
	require.EqualValues(`{"foo":42.0,"bar":"baz"}`, j.JSON)

	var quotes null.RawJSON
	err = json.Unmarshal([]byte(`""`), &quotes)
	require.NoError(err)
	require.True(quotes.Valid)
	require.EqualValues(`""`, quotes.JSON)

	var emptyObj null.RawJSON
	err = json.Unmarshal([]byte(`{}`), &emptyObj)
	require.NoError(err)
	require.True(emptyObj.Valid)
	require.EqualValues(`{}`, emptyObj.JSON)

	var invalid null.RawJSON
	err = json.Unmarshal([]byte(`:->`), &invalid)
	require.Error(err)

	var invalidObj null.RawJSON
	// `bar` should be quoted                ~~~
	err = json.Unmarshal([]byte(`{"foo":42.0,bar:"baz"}`), &invalidObj)
	require.Error(err)
	// This error should include information on the malformed object.
	require.Contains(err.Error(), "invalid character 'b'")
}

func TestRawJSONMarshalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ JSONText null.RawJSON }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	// nil values are handled correctly ...

	wrapper = Wrapper{null.RawJSON{}}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal([]byte("null"), data["JSONText"])
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal([]byte("null"), data["JSONText"])

	wrapper = Wrapper{null.NullJSON()}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal([]byte("null"), data["JSONText"])
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal([]byte("null"), data["JSONText"])

	// .. and other behavior is consistent with types.RawJSON.

	wrapper = Wrapper{null.NewJSONStr(`{"foo":42.0,"bar":"baz"}`)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"foo": 42.0, "bar": "baz"}, data["JSONText"])
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"foo": 42.0, "bar": "baz"}, data["JSONText"])

	wrapper = Wrapper{null.NewJSONStr("true")}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(true, data["JSONText"])
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(true, data["JSONText"])

	wrapper = Wrapper{null.NewJSONStr("null")}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(nil, data["JSONText"])
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(nil, data["JSONText"])

	// `bar` should be quoted                                             ~~~
	wrapper = Wrapper{null.NewJSONStr(`{"foo":42.0,bar:"baz"}`)}
	data, err = maps.Marshal(wrapper)
	require.Error(err)
	// This error should include information on the malformed object.
	require.Contains(err.Error(), "invalid character 'b'")
	data, err = maps.Marshal(&wrapper)
	require.Error(err)
	// This error should include information on the malformed object.
	require.Contains(err.Error(), "invalid character 'b'")
}
