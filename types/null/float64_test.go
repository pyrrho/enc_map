package null_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types/null"
)

// Helper Functions

func assertFloat64(t *testing.T, expected float64, f null.NullFloat64, fileLine string) {
	if !f.Valid {
		t.Fatalf("%s: NullFloat64 is null, but should be valid", fileLine)
	}
	if math.IsNaN(expected) {
		if !math.IsNaN(f.Float64) {
			t.Fatalf("%s: Expected NaN, received %v", fileLine, f.Float64)
		}
	} else if expected != f.Float64 {
		t.Fatalf("%s: %v ≠ %v", fileLine, expected, f.Float64)
	}
}

func assertNullFloat64(t *testing.T, f null.NullFloat64, fileLine string) {
	if f.Valid {
		t.Fatalf("%s: NullFloat64 is valid, but should be null", fileLine)
	}
}

// Tests

func TestFloat64From(t *testing.T) {
	assertFloat64(t, 1.2345, null.Float64From(1.2345), FileLine())
	assertFloat64(t, 4, null.Float64From(4), FileLine())
	assertFloat64(t, 0, null.Float64From(0), FileLine())
	assertFloat64(t, math.NaN(), null.Float64From(math.NaN()), FileLine())
	assertFloat64(t, math.Inf(1), null.Float64From(math.Inf(1)), FileLine())
}

func TestFloat64FromPtr(t *testing.T) {
	a := float64(1.2345)
	b := 0.0
	assertFloat64(t, 1.2345, null.Float64FromPtr(&a), FileLine())
	assertFloat64(t, 0.0, null.Float64FromPtr(&b), FileLine())

	assertNullFloat64(t, null.Float64FromPtr(nil), FileLine())
}

func TestFloat64Ctor(t *testing.T) {
	v := float64(1.2345)
	var nilPtr *float64

	assertFloat64(t, 1.2345, null.Float64(1.2345), FileLine())
	assertFloat64(t, 0.0, null.Float64(0.0), FileLine())
	assertFloat64(t, 7.0, null.Float64(7), FileLine())
	assertFloat64(t, v, null.Float64(v), FileLine())
	assertFloat64(t, v, null.Float64(&v), FileLine())
	assertFloat64(t, math.NaN(), null.Float64(math.NaN()), FileLine())
	assertFloat64(t, math.Inf(1), null.Float64(math.Inf(1)), FileLine())
	assertNullFloat64(t, null.Float64(nil), FileLine())
	assertNullFloat64(t, null.Float64(nilPtr), FileLine())
}

func TestFailureNewFloat64FromBool(t *testing.T) {
	defer ShouldPanic(t, FileLine())
	_ = null.Float64(true)
}

func TestFailureNewFloat64FromString(t *testing.T) {
	defer ShouldPanic(t, FileLine())
	_ = null.Float64("0")
}

func TestFloat64ValueOrZero(t *testing.T) {
	valid := null.Float64(1.2345)
	if valid.ValueOrZero() != 1.2345 {
		t.Fatalf("unexpected ValueOrZero, %v ≠ %v", 1.2345, valid.ValueOrZero())
	}

	nul := null.NullFloat64{}
	if nul.ValueOrZero() != 0 {
		t.Fatalf("unexpected ValueOrZero, %v ≠ %v", 0, nul.ValueOrZero())
	}
}

func TestFloat64Ptr(t *testing.T) {
	f := null.Float64(1.2345)
	ptr := f.Ptr()
	if *ptr != 1.2345 {
		t.Fatalf("bad %s float64: %#v ≠ %v\n", "pointer", ptr, 1.2345)
	}
	*ptr = 5.4321
	if f.Float64 != 5.4321 {
		t.Fatalf("bad %s float64: %#v ≠ %v\n", "pointer dereference", f.Float64, 5.4321)
	}

	nul := null.NullFloat64{}
	ptr = nul.Ptr()
	if ptr != nil {
		t.Fatalf("bad %s float64: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestFloat64Set(t *testing.T) {
	f := null.NullFloat64{}
	assertNullFloat64(t, f, FileLine())
	f.Set(1.2345)
	assertFloat64(t, 1.2345, f, FileLine())
	f.Set(0.0)
	assertFloat64(t, 0.0, f, FileLine())
}

func TestFloat64Null(t *testing.T) {
	f := null.Float64(1.2345)
	assertFloat64(t, 1.2345, f, FileLine())
	f.Null()
	assertNullFloat64(t, f, FileLine())
}

func TestFloat64IsNil(t *testing.T) {
	f := null.Float64(1.2345)
	if f.IsNil() {
		t.Fatalf("IsNil() should be false")
	}
	zero := null.Float64(0)
	if zero.IsNil() {
		t.Fatalf("IsNil() should be false")
	}
	nul := null.NullFloat64{}
	if !nul.IsNil() {
		t.Fatalf("IsNil() should be true")
	}
}

func TestFloat64IsZero(t *testing.T) {
	f := null.Float64(1.2345)
	if f.IsZero() {
		t.Fatalf("IsZero() should be false")
	}
	zero := null.Float64(0)
	if !zero.IsZero() {
		t.Fatalf("IsZero() should be true")
	}
	nul := null.NullFloat64{}
	if !nul.IsZero() {
		t.Fatalf("IsZero() should be true")
	}
}

func TestFloat64SQLValue(t *testing.T) {
	f := null.Float64(1.2345)
	val, err := f.Value()
	fatalIf(t, err, FileLine())
	if 1.2345 != val.(float64) {
		t.Fatalf("NullFloat64{1.2345, true}.Value() should return a valid driver.Value (float64)")
	}

	zero := null.Float64(0)
	val, err = zero.Value()
	fatalIf(t, err, FileLine())
	if 0 != val.(float64) {
		t.Fatalf("NullFloat64{0, true}.Value() should return a valid driver.Value (float64)")
	}

	nul := null.NullFloat64{}
	val, err = nul.Value()
	fatalIf(t, err, FileLine())
	if nil != val {
		t.Fatalf("NullFloat64{..., false}.Value() should return a nil driver.Value")
	}
}

func TestFloat64SQLScan(t *testing.T) {
	var f null.NullFloat64
	err := f.Scan(1.2345)
	fatalIf(t, err, FileLine())
	assertFloat64(t, 1.2345, f, FileLine())

	var i null.NullFloat64
	err = i.Scan(12345)
	fatalIf(t, err, FileLine())
	assertFloat64(t, 12345, i, FileLine())

	var f64Str null.NullFloat64
	// NB. Scan will coerce strings, but UnmarshalJSON won't.
	err = f64Str.Scan("1.2345")
	fatalIf(t, err, FileLine())
	assertFloat64(t, 1.2345, f, FileLine())

	var nul null.NullFloat64
	err = nul.Scan(nil)
	fatalIf(t, err, FileLine())
	assertNullFloat64(t, nul, FileLine())

	var wrong null.NullFloat64
	err = wrong.Scan("hello world")
	fatalUnless(t, err, FileLine())
}

func TestFloat64MarshalJSON(t *testing.T) {
	f := null.Float64(1.2345)
	data, err := json.Marshal(f)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "1.2345", FileLine())
	data, err = json.Marshal(&f)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "1.2345", FileLine())

	i := null.Float64(12345)
	data, err = json.Marshal(i)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "12345", FileLine())
	data, err = json.Marshal(&i)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "12345", FileLine())

	zero := null.Float64(0)
	data, err = json.Marshal(zero)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "0", FileLine())
	data, err = json.Marshal(&zero)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "0", FileLine())

	nul := null.NullFloat64{}
	data, err = json.Marshal(nul)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "null", FileLine())
	data, err = json.Marshal(&nul)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "null", FileLine())

	nan := null.Float64(math.NaN())
	data, err = json.Marshal(nan)
	fatalUnless(t, err, FileLine())
	data, err = json.Marshal(&nan)
	fatalUnless(t, err, FileLine())

	inf := null.Float64(math.Inf(1))
	data, err = json.Marshal(inf)
	fatalUnless(t, err, FileLine())
	data, err = json.Marshal(&inf)
	fatalUnless(t, err, FileLine())
}

func TestFloat64UnmarshalJSON(t *testing.T) {
	// Successful Valid Parses

	var f null.NullFloat64
	err := json.Unmarshal(floatJSON, &f)
	fatalIf(t, err, FileLine())
	assertFloat64(t, 1.2345, f, FileLine())

	var i null.NullFloat64
	err = json.Unmarshal(intJSON, &i)
	fatalIf(t, err, FileLine())
	assertFloat64(t, 12345, i, FileLine())

	var validObj null.NullFloat64
	err = json.Unmarshal(validFloatJSONObj, &validObj)
	fatalIf(t, err, FileLine())
	assertFloat64(t, 1.2345, validObj, FileLine())

	// Successful Null Parses

	var nul null.NullFloat64
	err = json.Unmarshal([]byte("null"), &nul)
	fatalIf(t, err, FileLine())
	assertNullFloat64(t, nul, FileLine())

	var nullObj null.NullFloat64
	err = json.Unmarshal(nullFloatJSONObj, &nullObj)
	fatalIf(t, err, FileLine())
	assertNullFloat64(t, nullObj, FileLine())

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var f64Str null.NullFloat64
	// Floats wrapped in quotes aren't floats.
	err = json.Unmarshal(floatStringJSON, &f64Str)
	fatalUnless(t, err, FileLine())

	var empty null.NullFloat64
	err = json.Unmarshal([]byte(""), &empty)
	fatalUnless(t, err, FileLine())

	var quotes null.NullFloat64
	err = json.Unmarshal([]byte(`""`), &quotes)
	fatalUnless(t, err, FileLine())

	var badType null.NullFloat64
	// Booleans are never floats.
	err = json.Unmarshal(boolTrueJSON, &badType)
	fatalUnless(t, err, FileLine())

	// The JSON specification does not include NaN, INF, Infinity, NegInfinity
	// or any other common literal for the IEEE 754 floating point
	// not-really-number values. As such, un-marshaling them from JSON will
	// result in errors.
	var nan null.NullFloat64
	err = json.Unmarshal([]byte("NaN"), &nan)
	fatalUnless(t, err, FileLine())

	var inf null.NullFloat64
	err = json.Unmarshal([]byte("INF"), &inf)
	fatalUnless(t, err, FileLine())

	var invalid null.NullFloat64
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf("expected json.SyntaxError, not %T", err)
	}
}

func TestFloat64MarshalMapValue(t *testing.T) {
	wrapper := struct{ Float64 null.NullFloat64 }{null.Float64(1.2345)}
	data, err := maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Float64": 1.2345}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Float64": 1.2345}, FileLine())

	wrapper = struct{ Float64 null.NullFloat64 }{null.Float64(0)}
	data, err = maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Float64": 0.0}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Float64": 0.0}, FileLine())

	// Null NullFloat64s should be encoded as "nil"
	wrapper = struct{ Float64 null.NullFloat64 }{null.NullFloat64{}}
	data, err = maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Float64": nil}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Float64": nil}, FileLine())
}
