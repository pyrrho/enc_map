package types_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types/null"
	"github.com/stretchr/testify/require"
)

var (
	timeString = "2012-12-21T21:21:21Z"
	timeJSON   = []byte(`"2012-12-21T21:21:21Z"`)
	timeValue  = time.Date(
		2012, time.December, 21,
		21, 21, 21, 0,
		time.UTC,
	)
	zeroTimeString = "0001-01-01T00:00:00Z"
	zeroTimeJSON   = []byte(`"0001-01-01T00:00:00Z"`)
)

func TestTimeCtors(t *testing.T) {
	require := require.New(t)

	// null.NullTime() returns a new null null.Time.
	// This is equivalent to null.Time{}.
	nul := null.NullTime()
	require.False(nul.Valid)

	empty := null.Time{}
	require.False(empty.Valid)

	// null.NewTime constructs a new, valid null.Time from a time.Time.
	ti := null.NewTime(time.Date(
		2012, time.December, 21,
		21, 21, 21, 0,
		time.UTC,
	))
	require.True(ti.Valid)
	require.Equal(timeValue, ti.Time)

	zero := null.NewTime(time.Time{})
	require.True(zero.Valid)
	require.Equal(time.Time{}, zero.Time)

	// Valid RFC 3339 strings can also be used to generate new Time objects.
	tis, err := null.NewTimeStr(timeString)
	require.NoError(err)
	require.True(tis.Valid)
	require.Equal(timeValue, tis.Time)

	zeros, err := null.NewTimeStr(zeroTimeString)
	require.NoError(err)
	require.True(zeros.Valid)
	require.Equal(time.Time{}, zeros.Time)

	nuls, err := null.NewTimeStr("")
	require.NoError(err)
	require.False(nuls.Valid)

	badString, err := null.NewTimeStr("December 12th, 12:02")
	require.Error(err)
	require.False(badString.Valid)
}

func TestTimeValueOrZero(t *testing.T) {
	require := require.New(t)

	ti := null.NewTime(timeValue)
	require.Equal(timeValue, ti.ValueOrZero())

	nul := null.Time{}
	require.Equal(time.Time{}, nul.ValueOrZero())
}

func TestTimeSet(t *testing.T) {
	require := require.New(t)

	ti := null.Time{}

	ti.Set(timeValue)
	require.True(ti.Valid)
	require.Equal(timeValue, ti.Time)

	ti.Set(time.Time{})
	require.True(ti.Valid)
	require.Equal(time.Time{}, ti.Time)
}

func TestTimeNull(t *testing.T) {
	require := require.New(t)

	ti := null.NewTime(timeValue)

	ti.Null()
	require.False(ti.Valid)
}

func TestTimeIsNil(t *testing.T) {
	require := require.New(t)

	ti := null.NewTime(timeValue)
	require.False(ti.IsNil())

	zero := null.NewTime(time.Time{})
	require.False(zero.IsNil())

	nul := null.Time{}
	require.True(nul.IsNil())
}

func TestTimeIsZero(t *testing.T) {
	require := require.New(t)

	ti := null.NewTime(timeValue)
	require.False(ti.IsZero())

	zero := null.NewTime(time.Time{})
	require.True(zero.IsZero())

	nul := null.Time{}
	require.True(nul.IsZero())
}

func TestTimeSQLValue(t *testing.T) {
	require := require.New(t)
	var val driver.Value
	var err error

	ti := null.NewTime(timeValue)
	val, err = ti.Value()
	require.NoError(err)
	require.Equal(timeValue, val)

	zero := null.NewTime(time.Time{})
	val, err = zero.Value()
	require.NoError(err)
	require.Equal(time.Time{}, val)

	nul := null.Time{}
	val, err = nul.Value()
	require.NoError(err)
	require.Equal(nil, val)
}

func TestTimeSQLScan(t *testing.T) {
	require := require.New(t)
	var err error

	var ti null.Time
	err = ti.Scan(timeValue)
	require.NoError(err)
	require.True(ti.Valid)
	require.Equal(timeValue, ti.Time)

	var zero null.Time
	err = zero.Scan(time.Time{})
	require.NoError(err)
	require.True(zero.Valid)
	require.Equal(time.Time{}, zero.Time)

	var nul null.Time
	err = nul.Scan(nil)
	require.NoError(err)
	require.False(nul.Valid)

	var wrong null.Time
	err = wrong.Scan("null")
	require.Error(err)
}

func TestTimeMarshalJSON(t *testing.T) {
	require := require.New(t)
	var data []byte
	var err error

	ti := null.NewTime(timeValue)
	data, err = json.Marshal(ti)
	require.NoError(err)
	require.Equal(timeJSON, data)
	data, err = json.Marshal(&ti)
	require.NoError(err)
	require.Equal(timeJSON, data)

	zero := null.NewTime(time.Time{})
	data, err = json.Marshal(zero)
	require.NoError(err)
	require.Equal(zeroTimeJSON, data)
	data, err = json.Marshal(&zero)
	require.NoError(err)
	require.Equal(zeroTimeJSON, data)

	nul := null.Time{}
	data, err = json.Marshal(nul)
	require.NoError(err)
	require.EqualValues("null", data)
	data, err = json.Marshal(&nul)
	require.NoError(err)
	require.EqualValues("null", data)
}

func TestTimeUnmarshalJSON(t *testing.T) {
	require := require.New(t)

	// Successful Valid Parses

	var ti null.Time
	err := json.Unmarshal(timeJSON, &ti)
	require.NoError(err)
	require.True(ti.Valid)
	require.Equal(timeValue, ti.Time)

	var zero null.Time
	err = json.Unmarshal(zeroTimeJSON, &zero)
	require.NoError(err)
	require.True(zero.Valid)
	require.Equal(time.Time{}, zero.Time)

	// Successful Null Parses

	var nul null.Time
	err = json.Unmarshal([]byte("null"), &nul)
	require.NoError(err)
	require.False(nul.Valid)

	var quotes null.Time
	err = json.Unmarshal([]byte(`""`), &quotes)
	require.NoError(err)
	require.False(nul.Valid)

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var badType null.Time
	err = json.Unmarshal([]byte("12345"), &badType)
	require.Error(err)

	var empty null.Time
	err = json.Unmarshal([]byte(""), &empty)
	require.Error(err)

	var invalid null.Time
	err = invalid.UnmarshalJSON([]byte(":->"))
	if _, ok := err.(*json.SyntaxError); !ok {
		require.FailNowf(
			"Unexpected Error Type",
			"expected *json.SyntaxError, not %T", err)
	}
}

func TestTimeMarshalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ Time null.Time }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	wrapper = Wrapper{null.NewTime(timeValue)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Time": timeValue}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Time": timeValue}, data)

	wrapper = Wrapper{null.NewTime(time.Time{})}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Time": time.Time{}}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Time": time.Time{}}, data)

	// Null NullTimes should be encoded as "nil"
	wrapper = struct{ Time null.Time }{null.Time{}}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Time": nil}, data)
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(map[string]interface{}{"Time": nil}, data)
}
