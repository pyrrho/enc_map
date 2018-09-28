package null_test

import (
	"encoding/json"
	"testing"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types/null"
)

// Helper Functions

func assertBool(t *testing.T, expected bool, b null.NullBool, fileLine string) {
	if !b.Valid {
		t.Fatalf("%s: NullBool is null, but should be valid", fileLine)
	}
	if expected != b.Bool {
		t.Fatalf("%s: %v ≠ %v", fileLine, expected, b.Bool)
	}
}

func assertNullBool(t *testing.T, b null.NullBool, fileLine string) {
	if b.Valid {
		t.Fatalf("%s: NullBool is valid, but should be null", fileLine)
	}
}

// Tests

func TestBoolFrom(t *testing.T) {
	assertBool(t, true, null.BoolFrom(true), FileLine())
	assertBool(t, false, null.BoolFrom(false), FileLine())
}

func TestBoolFromPtr(t *testing.T) {
	tr := true
	fl := false

	assertBool(t, true, null.BoolFromPtr(&tr), FileLine())
	assertBool(t, false, null.BoolFromPtr(&fl), FileLine())

	assertNullBool(t, null.BoolFromPtr(nil), FileLine())
}

func TestBoolCtor(t *testing.T) {
	tr := true
	var nilPtr *bool

	assertBool(t, true, null.Bool(true), FileLine())
	assertBool(t, false, null.Bool(false), FileLine())
	assertBool(t, tr, null.Bool(tr), FileLine())
	assertBool(t, tr, null.Bool(&tr), FileLine())
	assertNullBool(t, null.Bool(nil), FileLine())
	assertNullBool(t, null.Bool(nilPtr), FileLine())
}

func TestFailureNewBoolFromInt(t *testing.T) {
	defer ShouldPanic(t, FileLine())
	_ = null.Bool(0)
}

func TestFailureNewBoolFromString(t *testing.T) {
	defer ShouldPanic(t, FileLine())
	_ = null.Bool("false")
}

func TestBoolValueOrZero(t *testing.T) {
	valid := null.Bool(true)
	if valid.ValueOrZero() != true {
		t.Fatalf("unexpected ValueOrZero, %v ≠ %v", true, valid.ValueOrZero())
	}

	nul := null.NullBool{}
	if nul.ValueOrZero() != false {
		t.Fatalf("unexpected ValueOrZero, %v ≠ %v", false, nul.ValueOrZero())
	}
}

func TestBoolPtr(t *testing.T) {
	b := null.Bool(true)
	ptr := b.Ptr()
	if *ptr != true {
		t.Fatalf("bad %s bool: %#v ≠ %v\n", "pointer", ptr, true)
	}

	nul := null.NullBool{}
	ptr = nul.Ptr()
	if ptr != nil {
		t.Fatalf("bad %s bool: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestBoolSet(t *testing.T) {
	b := null.NullBool{}
	assertNullBool(t, b, FileLine())
	b.Set(true)
	assertBool(t, true, b, FileLine())
	b.Set(false)
	assertBool(t, false, b, FileLine())
}

func TestBoolNull(t *testing.T) {
	b := null.Bool(true)
	assertBool(t, true, b, FileLine())
	b.Null()
	assertNullBool(t, b, FileLine())
}

func TestBoolIsNil(t *testing.T) {
	a := null.Bool(true)
	if a.IsNil() {
		t.Fatal("NullBool{true, true}.IsNil() should be false")
	}
	b := null.Bool(false)
	if b.IsNil() {
		t.Fatal("NullBool{false, true}.IsNil() should be false")
	}
	nul := null.NullBool{}
	if !nul.IsNil() {
		t.Fatal("NullBool{false, false}.IsNil() should be true")
	}
}

func TestBoolIsZero(t *testing.T) {
	a := null.Bool(true)
	if a.IsZero() {
		t.Fatal("NullBool{true, true}.IsZero() should be false")
	}
	b := null.Bool(false)
	if !b.IsZero() {
		t.Fatal("NullBool{false, true}.IsZero() should be true")
	}
	nul := null.NullBool{}
	if !nul.IsZero() {
		t.Fatal("NullBool{false, false}.IsZero() should be true")
	}
}

func TestBoolSQLValue(t *testing.T) {
	b := null.Bool(true)
	val, err := b.Value()
	fatalIf(t, err, FileLine())
	if true != val.(bool) {
		t.Fatalf("NullBool{true, true}.Value() should return a valid driver.Value (bool)")
	}

	nul := null.NullBool{}
	val, err = nul.Value()
	fatalIf(t, err, FileLine())
	if nil != val {
		t.Fatalf("NullBool{false, false}.Value() should return a nil driver.Value")
	}
}

func TestBoolSQLScan(t *testing.T) {
	var b null.NullBool
	err := b.Scan(true)
	fatalIf(t, err, FileLine())
	assertBool(t, true, b, FileLine())

	var nul null.NullBool
	err = nul.Scan(nil)
	fatalIf(t, err, FileLine())
	assertNullBool(t, nul, FileLine())

	var wrong null.NullBool
	err = wrong.Scan(int64(42))
	fatalUnless(t, err, FileLine())
}

func TestBoolMarshalJSON(t *testing.T) {
	b := null.Bool(true)
	data, err := json.Marshal(b)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "true", FileLine())
	data, err = json.Marshal(&b)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "true", FileLine())

	zero := null.Bool(false)
	data, err = json.Marshal(zero)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "false", FileLine())
	data, err = json.Marshal(&zero)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "false", FileLine())

	// Null NullBools should be encoded as "null"
	nul := null.NullBool{}
	data, err = json.Marshal(nul)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "null", FileLine())
	data, err = json.Marshal(&nul)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "null", FileLine())

	wrapper := struct {
		Foo null.NullBool
		Bar null.NullBool
	}{
		null.Bool(true),
		null.NullBool{},
	}
	data, err = json.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, `{"Foo":true,"Bar":null}`, FileLine())
	data, err = json.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, `{"Foo":true,"Bar":null}`, FileLine())
}

func TestBoolUnmarshalJSON(t *testing.T) {
	// Successful Valid Parses

	var b null.NullBool
	err := json.Unmarshal(boolTrueJSON, &b)
	fatalIf(t, err, FileLine())
	assertBool(t, true, b, FileLine())

	// Successful Null Parses

	var nul null.NullBool
	err = json.Unmarshal([]byte("null"), &nul)
	fatalIf(t, err, FileLine())
	assertNullBool(t, nul, FileLine())

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var str null.NullBool
	// Booleans wrapped in quotes aren't booleans.
	err = json.Unmarshal(boolStringJSON, &str)
	fatalUnless(t, err, FileLine())

	var empty null.NullBool
	// An empty string is not a boolean.
	err = json.Unmarshal([]byte(`""`), &empty)
	fatalUnless(t, err, FileLine())

	var badType null.NullBool
	// Ints are never booleans.
	err = json.Unmarshal([]byte("1"), &badType)
	fatalUnless(t, err, FileLine())

	var invalid null.NullBool
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf("expected json.SyntaxError, not %T", err)
	}
}

func TestBoolMarshalMapValue(t *testing.T) {
	wrapper := struct{ Bool null.NullBool }{null.Bool(true)}
	data, err := maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Bool": true}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Bool": true}, FileLine())

	wrapper = struct{ Bool null.NullBool }{null.Bool(false)}
	data, err = maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Bool": false}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Bool": false}, FileLine())

	// Null NullBools should be encoded as "nil"
	wrapper = struct{ Bool null.NullBool }{null.NullBool{}}
	data, err = maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Bool": nil}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Bool": nil}, FileLine())
}
