package null_test

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types/null"
)

// Helper Functions

func assertInt64(t *testing.T, expected int64, b null.NullInt64, fileLine string) {
	if !b.Valid {
		t.Fatalf("%s: NullInt64 is null, but should be valid", fileLine)
	}
	if expected != b.Int64 {
		t.Fatalf("%s: %v ≠ %v", fileLine, expected, b.Int64)
	}
}

func assertNullInt64(t *testing.T, b null.NullInt64, fileLine string) {
	if b.Valid {
		t.Fatalf("%s: NullInt64 is valid, but should be null", fileLine)
	}
}

// Tests

func TestInt64From(t *testing.T) {
	assertInt64(t, 12345, null.Int64From(12345), FileLine())
	assertInt64(t, 0, null.Int64From(0), FileLine())
}

func TestInt64FromPtr(t *testing.T) {
	i := int64(12345)
	z := int64(0)
	assertInt64(t, 12345, null.Int64FromPtr(&i), FileLine())
	assertInt64(t, 0, null.Int64FromPtr(&z), FileLine())

	assertNullInt64(t, null.Int64FromPtr(nil), FileLine())
}

func TestNewInt64(t *testing.T) {
	v := int64(12345)
	var nilPtr *int64

	assertInt64(t, 1, null.Int64(1), FileLine())
	assertInt64(t, 0, null.Int64(0), FileLine())
	assertInt64(t, v, null.Int64(v), FileLine())
	assertInt64(t, v, null.Int64(&v), FileLine())
	assertNullInt64(t, null.Int64(nil), FileLine())
	assertNullInt64(t, null.Int64(nilPtr), FileLine())
}

func TestFailureNewInt64FromBool(t *testing.T) {
	defer ShouldPanic(t, FileLine())
	_ = null.Int64(true)
}

func TestFailureNewInt64FromFloat(t *testing.T) {
	defer ShouldPanic(t, FileLine())
	_ = null.Int64(4.2)
}

func TestInt64ValueOrZero(t *testing.T) {
	valid := null.Int64(12345)
	if valid.ValueOrZero() != 12345 {
		t.Fatalf("unexpected ValueOrZero, %v ≠ %v", 12345, valid.ValueOrZero())
	}

	nul := null.NullInt64{}
	if nul.ValueOrZero() != 0 {
		t.Fatalf("unexpected ValueOrZero, %v ≠ %v", 0, nul.ValueOrZero())
	}
}

func TestInt64Ptr(t *testing.T) {
	i := null.Int64(12345)
	ptr := i.Ptr()
	if *ptr != 12345 {
		t.Fatalf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 12345)
	}

	nul := null.NullInt64{}
	ptr = nul.Ptr()
	if ptr != nil {
		t.Fatalf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestInt64Set(t *testing.T) {
	i := null.NullInt64{}
	assertNullInt64(t, i, FileLine())
	i.Set(12345)
	assertInt64(t, 12345, i, FileLine())
	i.Set(0)
	assertInt64(t, 0, i, FileLine())
}

func TestInt64Null(t *testing.T) {
	i := null.Int64(12345)
	assertInt64(t, 12345, i, FileLine())
	i.Null()
	assertNullInt64(t, i, FileLine())
}

func TestInt64IsNil(t *testing.T) {
	i := null.Int64(12345)
	if i.IsNil() {
		t.Fatalf("IsNil() should be false")
	}
	zero := null.Int64(0)
	if zero.IsNil() {
		t.Fatalf("IsNil() should be false")
	}
	nul := null.NullInt64{}
	if !nul.IsNil() {
		t.Fatalf("IsNil() should be true")
	}
}

func TestInt64IsZero(t *testing.T) {
	i := null.Int64(12345)
	if i.IsZero() {
		t.Fatalf("IsZero() should be false")
	}
	zero := null.Int64(0)
	if !zero.IsZero() {
		t.Fatalf("IsZero() should be true")
	}
	nul := null.NullInt64{}
	if !nul.IsZero() {
		t.Fatalf("IsZero() should be true")
	}
}

func TestInt64SQLValue(t *testing.T) {
	i := null.Int64(12345)
	val, err := i.Value()
	fatalIf(t, err, FileLine())
	if 12345 != val.(int64) {
		t.Fatalf("NullInt64{12345, true}.Value() should return a valid driver.Value (int64)")
	}

	zero := null.Int64(0)
	val, err = zero.Value()
	fatalIf(t, err, FileLine())
	if 0 != val.(int64) {
		t.Fatalf("NullInt64{0, true}.Value() should return a valid driver.Value (int64)")
	}

	nul := null.NullInt64{}
	val, err = nul.Value()
	fatalIf(t, err, FileLine())
	if nil != val {
		t.Fatalf("NullInt64{..., false}.Value() should return a nil driver.Value")
	}
}

func TestInt64SQLScan(t *testing.T) {
	var i null.NullInt64
	err := i.Scan(12345)
	fatalIf(t, err, FileLine())
	assertInt64(t, 12345, i, FileLine())

	var i64Str null.NullInt64
	// NB. Scan will coerce strings, but UnmarshalJSON won't.
	err = i64Str.Scan("12345")
	fatalIf(t, err, FileLine())
	assertInt64(t, 12345, i64Str, FileLine())

	var nul null.NullInt64
	err = nul.Scan(nil)
	fatalIf(t, err, FileLine())
	assertNullInt64(t, nul, FileLine())

	var wrong null.NullInt64
	err = wrong.Scan("hello world")
	fatalUnless(t, err, FileLine())

	var f null.NullInt64
	err = f.Scan(1.2345)
	fatalUnless(t, err, FileLine())

	var b null.NullInt64
	err = b.Scan(true)
	fatalUnless(t, err, FileLine())
}

func TestInt64MarshalJSON(t *testing.T) {
	i := null.Int64From(12345)
	data, err := json.Marshal(i)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "12345", FileLine())
	data, err = json.Marshal(&i)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "12345", FileLine())

	zero := null.Int64(0)
	data, err = json.Marshal(zero)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "0", FileLine())
	data, err = json.Marshal(&zero)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "0", FileLine())

	// Null Int64s should be encoded as 'null'
	nul := null.NullInt64{}
	data, err = json.Marshal(nul)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "null", FileLine())
	data, err = json.Marshal(&nul)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "null", FileLine())
}

func TestInt64UnmarshalJSON(t *testing.T) {
	// Successful Valid Parses

	var i null.NullInt64
	err := json.Unmarshal(intJSON, &i)
	fatalIf(t, err, FileLine())
	assertInt64(t, 12345, i, FileLine())

	var validObj null.NullInt64
	err = json.Unmarshal(validIntJSONObj, &validObj)
	fatalIf(t, err, FileLine())
	assertInt64(t, 12345, validObj, FileLine())

	// Successful Null Parses

	var nul null.NullInt64
	err = json.Unmarshal([]byte("null"), &nul)
	fatalIf(t, err, FileLine())
	assertNullInt64(t, nul, FileLine())

	var nullSQL null.NullInt64
	err = json.Unmarshal(nullIntJSONObj, &nullSQL)
	fatalIf(t, err, FileLine())
	assertNullInt64(t, nullSQL, FileLine())

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var intStr null.NullInt64
	// Ints wrapped in quotes aren't ints.
	err = json.Unmarshal(intStringJSON, &intStr)
	fatalIf(t, err, FileLine())

	var empty null.NullInt64
	err = json.Unmarshal([]byte(""), &empty)
	fatalUnless(t, err, FileLine())

	var quotes null.NullInt64
	err = json.Unmarshal([]byte(`""`), &quotes)
	fatalUnless(t, err, FileLine())

	var f null.NullInt64
	// Non-integer numbers should not be coerced to ints.
	err = json.Unmarshal(floatJSON, &f)
	fatalUnless(t, err, FileLine())

	var invalid null.NullInt64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf("expected json.SyntaxError, not %T", err)
	}
}

func TestInt64UnmarshalJSONOverflow(t *testing.T) {
	int64Overflow := uint64(math.MaxInt64)

	// Max int64 should decode successfully
	var i null.NullInt64
	err := json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	fatalIf(t, err, FileLine())

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	// Decoded values should overflow int64
	fatalUnless(t, err, FileLine())
}

func TestInt64MarshalMapValue(t *testing.T) {
	wrapper := struct{ Int64 null.NullInt64 }{null.Int64(12345)}
	data, err := maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Int64": int64(12345)}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Int64": int64(12345)}, FileLine())

	wrapper = struct{ Int64 null.NullInt64 }{null.Int64(0)}
	data, err = maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Int64": int64(0)}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Int64": int64(0)}, FileLine())

	// Null NullInt64s should be encoded as "nil"
	wrapper = struct{ Int64 null.NullInt64 }{null.NullInt64{}}
	data, err = maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Int64": nil}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Int64": nil}, FileLine())
}
