package types_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types"
	"github.com/stretchr/testify/require"
)

func TestSurprises(t *testing.T) {
	require := require.New(t)

	// If you cast an addressable []byte to RawJSON, modifying the original
	// []byte will modify the RawJSON. This is equivalent to assigning to a new
	// []byte rather than taking a copy.
	// NewJSON will handle making a copy.
	// NewJSONStr doesn't have this issue because strings don't have this issue.

	ba := []byte(`"Hello World"`)
	bb := ba
	bc := append([]byte(nil), ba...)
	ja := types.RawJSON(ba)                         // Don't do this.
	jb := types.RawJSON(append([]byte(nil), ba...)) // Don't do this.
	jc := types.NewJSON(ba)                         // Do this.
	jd := types.NewJSONStr(string(ba))              // This is extra work.

	ja[1] = 'h' // Title Case to lower case
	ja[7] = 'w' // Title Case to lower case

	require.EqualValues(`"hello world"`, ba) // indirectly modified [original]
	require.EqualValues(`"hello world"`, bb) // indirectly modified [assigned to ba]
	require.EqualValues(`"Hello World"`, bc) // not        modified [copied from ba]
	require.EqualValues(`"hello world"`, ja) // directly   modified [cast from ba]
	require.EqualValues(`"Hello World"`, jb) // not        modified [cast from a copy of ba]
	require.EqualValues(`"Hello World"`, jc) // not        modified [constructed]
	require.EqualValues(`"Hello World"`, jd) // not        modified [constructed]

	// Similar to NewJSON, .MarshalJSON will create and return new a object, so
	// older addressable objects will not be modified.

	jaJSON, err := ja.MarshalJSON() // returns ([]byte, error)
	require.NoError(err)

	jaJSON[1] = 'H' // lower case to Title Case
	jaJSON[7] = 'W' // lower case to Title Case

	require.EqualValues(`"hello world"`, ja)     // ja was unaffected by the ...
	require.EqualValues(`"Hello World"`, jaJSON) // jaJSON modifications.

	// .Set will attempt to reuse the underlying data in a RawJSON, so changes
	// _may_ be reflected in older objects.

	ja.Set([]byte(`"what?"`))

	require.EqualValues(`"what?"`, ja)       // This is correct.
	require.EqualValues(`"what?"world"`, ba) // This is.... not.

	// The same is true for .SetStr.

	ja.SetStr(`"Again?"`)

	require.EqualValues(`"Again?"`, ja)      // This is correct.
	require.EqualValues(`"Again?"orld"`, ba) // This is.... not.

	// If you trigger a reallocation in the .Set call, this won't occur.
	// Please do not rely on this behavior.

	ja.Set([]byte(`"And that's enough for today. Bye!`))

	require.EqualValues(`"And that's enough for today. Bye!`, ja)
	require.EqualValues(`"Again?"orld"`, ba) // Unmodified from earlier.

	// Note that .Set can be used to initialize a nil RawJSON.
	n := types.RawJSON(nil)
	require.Nil(n)

	n.Set([]byte(`"Here's lookin' at you, kid"`))
	require.NotNil(n)
	require.EqualValues(`"Here's lookin' at you, kid"`, n)

	// It can also be used to un-set an initialized RawJSON.
	m := types.NewJSONStr(`"Hello, once again."`)
	require.NotNil(m)

	m.Set(nil)
	require.True(len(m) == 0)
	require.NotNil(m) // Note that it won't nil an initialized RawJSON.
}

func TestRawJSONIsNil(t *testing.T) {
	require := require.New(t)

	str := types.RawJSON(`"Hello World"`)
	require.False(str.IsNil())
	emptyStr := types.RawJSON(`""`)
	require.False(emptyStr.IsNil())

	num := types.RawJSON("42.0")
	require.False(num.IsNil())
	emptyNum := types.RawJSON("0.0")
	require.False(emptyNum.IsNil())

	obj := types.RawJSON(`{"foo":42.0,bar:"baz"}`)
	require.False(obj.IsNil())
	emptyObj := types.RawJSON("{}")
	require.False(emptyObj.IsNil())

	arr := types.RawJSON("[1.0, 2.0, 3.0]")
	require.False(arr.IsNil())
	emptyArr := types.RawJSON("[]")
	require.False(emptyArr.IsNil())

	bol := types.RawJSON("true")
	require.False(bol.IsNil())
	emptyBol := types.RawJSON("false")
	require.False(emptyBol.IsNil())

	nul := types.RawJSON("null")
	require.False(nul.IsNil())

	empty := types.RawJSON{}
	require.True(empty.IsNil())

	nil_ := types.RawJSON(nil)
	require.True(nil_.IsNil())
}

func TestRawJSONIsZero(t *testing.T) {
	require := require.New(t)

	str := types.RawJSON(`"Hello World"`)
	require.False(str.IsZero())
	emptyStr := types.RawJSON(`""`)
	require.True(emptyStr.IsZero())

	num := types.RawJSON("42.0")
	require.False(num.IsZero())
	emptyNum := types.RawJSON("0.0")
	require.True(emptyNum.IsZero())

	obj := types.RawJSON(`{"foo":42.0,bar:"baz"}`)
	require.False(obj.IsZero())
	emptyObj := types.RawJSON("{}")
	require.True(emptyObj.IsZero())

	arr := types.RawJSON("[1.0, 2.0, 3.0]")
	require.False(arr.IsZero())
	emptyArr := types.RawJSON("[]")
	require.True(emptyArr.IsZero())

	bol := types.RawJSON("true")
	require.False(bol.IsZero())
	emptyBol := types.RawJSON("false")
	require.True(emptyBol.IsZero())

	nul := types.RawJSON("null")
	require.True(nul.IsZero())

	empty := types.RawJSON{}
	require.True(empty.IsZero())

	nil_ := types.RawJSON(nil)
	require.True(nil_.IsZero())
}

func TestRawJSONValueIsZero(t *testing.T) {
	require := require.New(t)
	var b bool
	var err error

	b, err = types.RawJSON(`"hello world"`).ValueIsZero()
	require.NoError(err)
	require.False(b)
	b, err = types.RawJSON(`""`).ValueIsZero()
	require.NoError(err)
	require.True(b)

	b, err = types.RawJSON("42.0").ValueIsZero()
	require.NoError(err)
	require.False(b)
	b, err = types.RawJSON("0.0").ValueIsZero()
	require.NoError(err)
	require.True(b)

	b, err = types.RawJSON(`{"foo":42.0,"bar":"baz"}`).ValueIsZero()
	require.NoError(err)
	require.False(b)
	b, err = types.RawJSON("{}").ValueIsZero()
	require.NoError(err)
	require.True(b)

	b, err = types.RawJSON("[1, 2, 3, 4]").ValueIsZero()
	require.NoError(err)
	require.False(b)
	b, err = types.RawJSON("[]").ValueIsZero()
	require.NoError(err)
	require.True(b)

	b, err = types.RawJSON("true").ValueIsZero()
	require.NoError(err)
	require.False(b)
	b, err = types.RawJSON("false").ValueIsZero()
	require.NoError(err)
	require.True(b)

	b, err = types.RawJSON("null").ValueIsZero()
	require.NoError(err)
	require.True(b)

	b, err = types.RawJSON("").ValueIsZero()
	require.Error(err)
	require.Contains(err.Error(), "RawJSON:") // err must come from RawJSON

	b, err = types.RawJSON(nil).ValueIsZero()
	require.Error(err)
	require.Contains(err.Error(), "RawJSON:") // err must come from RawJSON

	// `bar` should be quoted           ~~~
	b, err = types.RawJSON(`{"foo":42.0,bar:"baz"}`).ValueIsZero()
	require.Error(err)
	// This error should include information on the malformed object.
	require.Contains(err.Error(), "invalid character 'b'")
}

func TestRawJSONSQLValue(t *testing.T) {
	var val driver.Value
	var err error
	require := require.New(t)

	j := types.RawJSON(`{"foo":42.0,"bar":"baz"}`)
	val, err = j.Value()
	require.NoError(err)
	require.EqualValues(`{"foo":42.0,"bar":"baz"}`, val)

	empty := types.RawJSON{}
	val, err = empty.Value()
	require.Error(err)
	require.Contains(err.Error(), "RawJSON:") // err must come from RawJSON

	nil_ := types.RawJSON(nil)
	val, err = nil_.Value()
	require.Error(err)
	require.Contains(err.Error(), "RawJSON:") // err must come from RawJSON

	invalid := types.RawJSON(":->")
	val, err = invalid.Value()
	require.Error(err)
}

func TestRawJSONSQLScan(t *testing.T) {
	require := require.New(t)
	var err error

	var j types.RawJSON
	err = j.Scan(driver.Value([]byte(`{"foo":42.0,"bar":"baz"}`)))
	require.NoError(err)
	require.EqualValues(`{"foo":42.0,"bar":"baz"}`, j)

	var nul types.RawJSON
	err = nul.Scan(driver.Value([]byte("null")))
	require.NoError(err)
	require.EqualValues("null", nul)

	var empty1 types.RawJSON
	require.Nil(empty1) // empty1 starts nil ...
	err = empty1.Scan(driver.Value([]byte{}))
	require.NoError(err)
	require.Nil(empty1) // ... and stays nil after nothing is copied.

	var empty2 types.RawJSON
	require.Nil(empty2) // empty2 starts nil ...
	err = empty2.Scan(driver.Value([]byte(nil)))
	require.NoError(err)
	require.Nil(empty2) // ... and stays nil after nothing is copied.

	// NB. Scan assumes valid JSON is provided; it will not validate.
	var invalid types.RawJSON
	err = invalid.Scan(driver.Value([]byte(":->")))
	require.NoError(err)
	require.EqualValues(":->", invalid)
}

func TestRawJSONSQLScanStrs(t *testing.T) {
	require := require.New(t)
	var err error

	var j types.RawJSON
	err = j.Scan(driver.Value(`{"foo":42.0,"bar":"baz"}`))
	require.NoError(err)
	require.EqualValues(`{"foo":42.0,"bar":"baz"}`, j)

	var empty types.RawJSON
	require.Nil(empty) // empty starts nil ...
	err = empty.Scan(driver.Value(""))
	require.NoError(err)
	require.Nil(empty) // ... and stays nil after nothing is copied.

	var nul types.RawJSON
	err = nul.Scan(driver.Value("null"))
	require.NoError(err)
	require.EqualValues("null", nul)

	// NB. Scan assumes valid JSON is provided; it will not validate.
	var invalid types.RawJSON
	err = invalid.Scan(driver.Value(":->"))
	require.NoError(err)
	require.EqualValues(":->", invalid)
}

func TestRawJSONSQLScanNil(t *testing.T) {
	require := require.New(t)
	var err error

	var nil_ types.RawJSON
	err = nil_.Scan(driver.Value(nil))
	require.Error(err)
}

func TestRawJSONMarshalJSON(t *testing.T) {
	require := require.New(t)
	var err error
	var data []byte

	j := types.RawJSON(`{"foo":42.0,"bar":"baz"}`)
	data, err = json.Marshal(j)
	require.NoError(err)
	require.EqualValues(`{"foo":42.0,"bar":"baz"}`, data)
	data, err = json.Marshal(&j)
	require.NoError(err)
	require.EqualValues(`{"foo":42.0,"bar":"baz"}`, data)

	empty := types.RawJSON{}
	data, err = json.Marshal(empty)
	require.Error(err)
	require.Contains(err.Error(), "RawJSON:") // err must come from RawJSON
	data, err = json.Marshal(&empty)
	require.Error(err)
	require.Contains(err.Error(), "RawJSON:") // err must come from RawJSON

	// `bar` should be quoted                ~~~
	invalidObj := types.RawJSON(`{"foo":42.0,bar:"baz"}`)
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

	var j types.RawJSON
	err = json.Unmarshal([]byte(`{"foo":42.0,"bar":"baz"}`), &j)
	require.NoError(err)
	require.EqualValues(`{"foo":42.0,"bar":"baz"}`, j)

	var nul types.RawJSON
	err = json.Unmarshal([]byte("null"), &nul)
	require.NoError(err)
	require.EqualValues("null", nul)

	var quotes types.RawJSON
	err = json.Unmarshal([]byte(`""`), &quotes)
	require.NoError(err)
	require.EqualValues(`""`, quotes)

	var invalid types.RawJSON
	err = json.Unmarshal([]byte(`:->`), &invalid)
	require.Error(err)

	var invalidObj types.RawJSON
	// `bar` should be quoted                ~~~
	err = json.Unmarshal([]byte(`{"foo":42.0,bar:"baz"}`), &invalidObj)
	require.Error(err)
	// This error should include information on the malformed object.
	require.Contains(err.Error(), "invalid character 'b'")

	var empty types.RawJSON
	err = json.Unmarshal([]byte{}, &empty)
	require.Error(err)

	var nil_ types.RawJSON
	err = json.Unmarshal([]byte(nil), &nil_)
	require.Error(err)
}

func TestRawJSONMarshalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ JSONText types.RawJSON }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	wrapper = Wrapper{types.RawJSON(`{"foo":42.0,"bar":"baz"}`)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"foo": 42.0, "bar": "baz"}, data["JSONText"])
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"foo": 42.0, "bar": "baz"}, data["JSONText"])

	wrapper = Wrapper{types.RawJSON("true")}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(true, data["JSONText"])
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(true, data["JSONText"])

	wrapper = Wrapper{types.RawJSON("null")}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(nil, data["JSONText"])
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(nil, data["JSONText"])

	wrapper = Wrapper{types.RawJSON{}}
	data, err = maps.Marshal(wrapper)
	require.Error(err)
	require.Contains(err.Error(), "RawJSON:") // err must come from RawJSON
	data, err = maps.Marshal(&wrapper)
	require.Error(err)
	require.Contains(err.Error(), "RawJSON:") // err must come from RawJSON

	wrapper = Wrapper{types.RawJSON(nil)}
	data, err = maps.Marshal(wrapper)
	require.Error(err)
	require.Contains(err.Error(), "RawJSON:") // err must come from RawJSON
	data, err = maps.Marshal(&wrapper)
	require.Error(err)
	require.Contains(err.Error(), "RawJSON:") // err must come from RawJSON

	// `bar` should be quoted                                             ~~~
	wrapper = Wrapper{types.RawJSON(`{"foo":42.0,bar:"baz"}`)}
	data, err = maps.Marshal(wrapper)
	require.Error(err)
	// This error should include information on the malformed object.
	require.Contains(err.Error(), "invalid character 'b'")
	data, err = maps.Marshal(&wrapper)
	require.Error(err)
	// This error should include information on the malformed object.
	require.Contains(err.Error(), "invalid character 'b'")
}
