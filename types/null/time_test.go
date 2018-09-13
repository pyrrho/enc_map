package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types/null"
)

// Helper Functions

func assertTime(t *testing.T, expected time.Time, actual null.NullTime, fileLine string) {
	if !actual.Valid {
		t.Fatalf("%s: NullTime is null, but should be valid", fileLine)
	}
	if expected != actual.Time {
		t.Fatalf("%s: %v ≠ %v", fileLine, expected, actual.Time)
	}
}

func assertNullTime(t *testing.T, actual null.NullTime, fileLine string) {
	if actual.Valid {
		t.Fatalf("%s: NullTime is valid, but should be null", fileLine)
	}
}

// Tests

func TestTimeFrom(t *testing.T) {
	assertTime(t, timeValue, null.TimeFrom(timeValue), FileLine())
	assertTime(t, time.Time{}, null.TimeFrom(zeroTimeValue), FileLine())
}

func TestTimeFromPtr(t *testing.T) {
	assertTime(t, timeValue, null.TimeFromPtr(&timeValue), FileLine())
	assertTime(t, time.Time{}, null.TimeFromPtr(&zeroTimeValue), FileLine())
	assertNullTime(t, null.TimeFromPtr(nil), FileLine())
}

func TestTimeCtor(t *testing.T) {
	var nilPtr *time.Time

	assertTime(t, timeValue, null.Time(timeValue), FileLine())
	assertTime(t, time.Time{}, null.Time(zeroTimeValue), FileLine())
	assertTime(t, timeValue, null.Time(&timeValue), FileLine())
	assertTime(t, time.Time{}, null.Time(&zeroTimeValue), FileLine())
	assertNullTime(t, null.Time(nil), FileLine())
	assertNullTime(t, null.Time(nilPtr), FileLine())
}

func TestFailureNewTimeFromInt(t *testing.T) {
	defer ShouldPanic(t, FileLine())
	_ = null.Time(2012)
}

func TestFailureNewTimeFromString(t *testing.T) {
	defer ShouldPanic(t, FileLine())
	_ = null.Time(timeString)
}

func TestTimeValueOrZero(t *testing.T) {
	ti := null.Time(timeValue)
	if ti.ValueOrZero().IsZero() || ti.ValueOrZero() != timeValue {
		t.Fatalf("unexpected ValueOrZero(), %s ≠ %s ", timeValue, ti.ValueOrZero())
	}

	nul := null.NullTime{}
	if !nul.ValueOrZero().IsZero() {
		t.Fatalf("unexpected ValueOrZero(), %s ≠ %s", zeroTimeValue, nul.ValueOrZero())
	}
}

func TestTimePtr(t *testing.T) {
	ti := null.Time(timeValue)
	ptr := ti.Ptr()
	if *ptr != timeValue {
		t.Fatalf("bad %s time: %#v ≠ %v\n", "pointer", ptr, timeValue)
	}

	nul := null.NullTime{}
	ptr = nul.Ptr()
	if ptr != nil {
		t.Fatalf("bad %s time: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestTimeSet(t *testing.T) {
	ti := null.NullTime{}
	assertNullTime(t, ti, FileLine())
	ti.Set(timeValue)
	assertTime(t, timeValue, ti, FileLine())
	ti.Set(time.Time{})
	assertTime(t, time.Time{}, ti, FileLine())
}

func TestTimeNull(t *testing.T) {
	ti := null.Time(timeValue)
	assertTime(t, timeValue, ti, FileLine())
	ti.Null()
	assertNullTime(t, ti, FileLine())
}

func TestTimeIsNil(t *testing.T) {
	ti := null.Time(timeValue)
	if ti.IsNil() {
		t.Fatalf("IsNil() should be false")
	}
	empty := null.Time(time.Time{})
	if empty.IsNil() {
		t.Fatalf("IsNil() should be false")
	}
	nul := null.NullTime{}
	if !nul.IsNil() {
		t.Fatalf("IsNil() should be true")
	}
}

func TestTimeIsZero(t *testing.T) {
	ti := null.Time(timeValue)
	if ti.IsZero() {
		t.Fatalf("IsZero() should be false")
	}
	empty := null.Time(time.Time{})
	if !empty.IsZero() {
		t.Fatalf("IsZero() should be true")
	}
	nul := null.NullByteSlice{}
	if !nul.IsZero() {
		t.Fatalf("IsZero() should be true")
	}
}

func TestTimeValue(t *testing.T) {
	var val driver.Value
	var err error

	ti := null.Time(timeValue)
	val, err = ti.Value()
	fatalIf(t, err, FileLine())
	if timeValue != val.(time.Time) {
		t.Fatalf("ti.Value() should return a valid driver.Value")
	}

	zero := null.Time(zeroTimeValue)
	val, err = zero.Value()
	fatalIf(t, err, FileLine())
	if zeroTimeValue != val.(time.Time) {
		t.Fatalf("zero.Value() should return a valid driver.Value")
	}

	nul := null.NullTime{}
	val, err = nul.Value()
	fatalIf(t, err, FileLine())
	if val != nil {
		t.Fatalf("nul.Value() should return a nul driver.Value")
	}
}

func TestTimeSQLValue(t *testing.T) {
	ti := null.Time(timeValue)
	val, err := ti.Value()
	fatalIf(t, err, FileLine())
	if timeValue != val.(time.Time) {
		t.Fatalf("NullTime{..., true}.Value() should return a valid driver.Value (time.Time)")
	}

	empty := null.Time(time.Time{})
	val, err = empty.Value()
	fatalIf(t, err, FileLine())
	if zeroTimeValue != val.(time.Time) {
		t.Fatalf("NullTime{..., true}.Value() should return a valid driver.Value (time.Time)")
	}

	nul := null.NullTime{}
	val, err = nul.Value()
	fatalIf(t, err, FileLine())
	if nil != val {
		t.Fatalf("NullTime{..., false}.Value() should return a nil driver.Value")
	}
}

func TestTimeSQLScan(t *testing.T) {
	var ti null.NullTime
	err := ti.Scan(timeValue)
	fatalIf(t, err, FileLine())
	assertTime(t, timeValue, ti, FileLine())

	var zero null.NullTime
	err = zero.Scan(zeroTimeValue)
	fatalIf(t, err, FileLine())
	assertTime(t, zeroTimeValue, zero, FileLine())

	var nul null.NullTime
	err = nul.Scan(nil)
	fatalIf(t, err, FileLine())
	assertNullTime(t, nul, FileLine())

	var wrong null.NullTime
	err = wrong.Scan("null")
	fatalUnless(t, err, FileLine())
}

func TestTimeMarshalText(t *testing.T) {
	ti := null.Time(timeValue)
	txt, err := ti.MarshalText()
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, txt, timeString, FileLine())
	txt, err = (&ti).MarshalText()
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, txt, timeString, FileLine())

	zero := null.Time(zeroTimeValue)
	txt, err = zero.MarshalText()
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, txt, zeroTimeString, FileLine())
	txt, err = (&zero).MarshalText()
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, txt, zeroTimeString, FileLine())

	nul := null.NullTime{}
	txt, err = nul.MarshalText()
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, txt, "", FileLine())
	txt, err = (&nul).MarshalText()
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, txt, "", FileLine())
}

func TestTimeUnmarshalText(t *testing.T) {
	// Successful Valid Parses

	var ti null.NullTime
	err := ti.UnmarshalText([]byte(timeString))
	fatalIf(t, err, FileLine())
	assertTime(t, timeValue, ti, FileLine())

	var zero null.NullTime
	err = zero.UnmarshalText([]byte(zeroTimeString))
	fatalIf(t, err, FileLine())
	assertTime(t, zeroTimeValue, zero, FileLine())

	// Successful Null Parses

	var nulStr null.NullTime
	err = nulStr.UnmarshalText([]byte("null"))
	fatalIf(t, err, FileLine())
	assertNullTime(t, nulStr, FileLine())

	var empty null.NullTime
	err = empty.UnmarshalText([]byte(""))
	fatalIf(t, err, FileLine())
	assertNullTime(t, nulStr, FileLine())

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var quotes null.NullTime
	err = quotes.UnmarshalText([]byte(`""`))
	fatalUnless(t, err, FileLine())

	var invalid null.NullTime
	err = invalid.UnmarshalText([]byte("hello world"))
	fatalUnless(t, err, FileLine())
}

func TestTimeMarshalJSON(t *testing.T) {
	ti := null.Time(timeValue)
	data, err := json.Marshal(ti)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, string(timeJSON), FileLine())
	data, err = json.Marshal(&ti)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, string(timeJSON), FileLine())

	zero := null.Time(zeroTimeValue)
	data, err = json.Marshal(zero)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, string(zeroTimeJSON), FileLine())
	data, err = json.Marshal(&zero)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, string(zeroTimeJSON), FileLine())

	nul := null.NullTime{}
	data, err = json.Marshal(nul)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "null", FileLine())
	data, err = json.Marshal(&nul)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "null", FileLine())
}

func TestTimeUnmarshalJSON(t *testing.T) {
	// Successful Valid Parses

	var ti null.NullTime
	err := json.Unmarshal(timeJSON, &ti)
	fatalIf(t, err, FileLine())
	assertTime(t, timeValue, ti, FileLine())

	var zero null.NullTime
	err = json.Unmarshal(zeroTimeJSON, &zero)
	fatalIf(t, err, FileLine())
	assertTime(t, zeroTimeValue, zero, FileLine())

	var validObj null.NullTime
	err = json.Unmarshal(validTimeJSONObj, &validObj)
	fatalIf(t, err, FileLine())
	assertTime(t, timeValue, validObj, FileLine())

	// Successful Null Parses

	var nul null.NullTime
	err = json.Unmarshal([]byte("null"), &nul)
	fatalIf(t, err, FileLine())
	assertNullTime(t, nul, FileLine())

	var nullObj null.NullTime
	err = json.Unmarshal(nullTimeJSONObj, &nullObj)
	fatalIf(t, err, FileLine())
	assertNullTime(t, nullObj, FileLine())

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var badType null.NullTime
	err = json.Unmarshal(intJSON, &badType)
	fatalUnless(t, err, FileLine())

	var empty null.NullTime
	err = json.Unmarshal([]byte(""), &empty)
	fatalUnless(t, err, FileLine())

	var quotes null.NullTime
	err = json.Unmarshal([]byte(`""`), &quotes)
	fatalUnless(t, err, FileLine())

	var invalid null.NullTime
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf("expected json.SyntaxError, not %T", err)
	}
	assertNullTime(t, invalid, FileLine())
}

func TestTimeMarshalMapValue(t *testing.T) {
	wrapper := struct{ Time null.NullTime }{null.Time(timeValue)}
	data, err := maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Time": timeValue}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Time": timeValue}, FileLine())

	wrapper = struct{ Time null.NullTime }{null.Time(zeroTimeValue)}
	data, err = maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Time": zeroTimeValue}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Time": zeroTimeValue}, FileLine())

	// Null NullTimes should be encoded as "nil"
	wrapper = struct{ Time null.NullTime }{null.NullTime{}}
	data, err = maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Time": nil}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Time": nil}, FileLine())
}
