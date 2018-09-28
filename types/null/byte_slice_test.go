package null_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types/null"
)

// Helper Functions

func assertByteSlice(t *testing.T, expected []byte, bs null.NullByteSlice, fileLine string) {
	if !bs.Valid {
		t.Fatalf("%s: NullByteSlice is null, but should be valid", fileLine)
	}
	if !bytes.Equal(expected, bs.ByteSlice) {
		t.Fatalf("%s: %v (%s) ≠ %v (%s)",
			fileLine,
			expected, expected,
			bs.ByteSlice, bs.ByteSlice,
		)
	}
}

func assertNullByteSlice(t *testing.T, bs null.NullByteSlice, fileLine string) {
	if bs.Valid {
		t.Fatalf("%s: NullByteSlice is valid, but should be null", fileLine)
	}
}

// Tests

func TestByteSliceFrom(t *testing.T) {
	assertByteSlice(t, byteSliceValue, null.ByteSliceFrom(byteSliceValue), FileLine())
	assertByteSlice(t, []byte{}, null.ByteSliceFrom([]byte{}), FileLine())

	var nl []byte
	assertNullByteSlice(t, null.ByteSliceFrom(nil), FileLine())
	assertNullByteSlice(t, null.ByteSliceFrom(nl), FileLine())
}

func TestByteSliceFromPtr(t *testing.T) {
	assertByteSlice(t, byteSliceValue, null.ByteSliceFromPtr(&byteSliceValue), FileLine())

	var nlPtr *[]byte
	var nl []byte
	ptrToNl := &nl
	assertNullByteSlice(t, null.ByteSliceFromPtr(nil), FileLine())
	assertNullByteSlice(t, null.ByteSliceFromPtr(nlPtr), FileLine())
	assertNullByteSlice(t, null.ByteSliceFromPtr(ptrToNl), FileLine())
}

func TestByteSliceCtor(t *testing.T) {
	var nilPtr *[]byte

	assertByteSlice(t, byteSliceValue, null.ByteSlice(byteSliceValue), FileLine())
	assertByteSlice(t, []byte{}, null.ByteSlice([]byte{}), FileLine())
	assertByteSlice(t, byteSliceValue, null.ByteSlice(&byteSliceValue), FileLine())
	assertNullByteSlice(t, null.ByteSlice(nil), FileLine())
	assertNullByteSlice(t, null.ByteSlice(nilPtr), FileLine())
}

func TestFailureNewByteSliceFromInt(t *testing.T) {
	defer ShouldPanic(t, FileLine())
	_ = null.ByteSlice(2012)
}

func TestFailureNewByteSliceFromString(t *testing.T) {
	defer ShouldPanic(t, FileLine())
	_ = null.ByteSlice("DAICON V")
}

func TestByteSliceValueOrZero(t *testing.T) {
	valid := null.ByteSlice(byteSliceValue)
	if !bytes.Equal(byteSliceValue, valid.ValueOrZero()) {
		t.Fatalf("unexpected ValueOrZero(), %s ≠ %s ", byteSliceValue, valid.ValueOrZero())
	}

	nul := null.NullByteSlice{}
	if !bytes.Equal([]byte{}, nul.ValueOrZero()) {
		t.Fatalf("unexpected ValueOrZero(), %s ≠ %s", []byte{}, nul.ValueOrZero())
	}
}

func TestByteSlicePtr(t *testing.T) {
	bs := null.ByteSlice(byteSliceValue)
	ptr := bs.Ptr()
	if !bytes.Equal(*ptr, byteSliceValue) {
		t.Fatalf("bad %s byte slice: %#v ≠ %v\n", "pointer", ptr, byteSliceValue)
	}

	nul := null.NullByteSlice{}
	ptr = nul.Ptr()
	if ptr != nil {
		t.Fatalf("bad %s byte slice: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestByteSliceSet(t *testing.T) {
	bs := null.NullByteSlice{}
	assertNullByteSlice(t, bs, FileLine())
	bs.Set(byteSliceValue)
	assertByteSlice(t, byteSliceValue, bs, FileLine())
	bs.Set([]byte{})
	assertByteSlice(t, []byte{}, bs, FileLine())
	bs.Set(nil)
	assertNullByteSlice(t, bs, FileLine())
}

func TestByteSliceNull(t *testing.T) {
	bs := null.ByteSlice(byteSliceValue)
	assertByteSlice(t, byteSliceValue, bs, FileLine())
	bs.Null()
	assertNullByteSlice(t, bs, FileLine())
}

func TestByteSliceIsNil(t *testing.T) {
	bs := null.ByteSlice(byteSliceValue)
	if bs.IsNil() {
		t.Fatalf("IsNil() should be false")
	}
	empty := null.ByteSlice([]byte{})
	if empty.IsNil() {
		t.Fatalf("IsNil() should be false")
	}
	nul := null.NullByteSlice{}
	if !nul.IsNil() {
		t.Fatalf("IsNil() should be true")
	}
}

func TestByteSliceIsZero(t *testing.T) {
	bs := null.ByteSlice(byteSliceValue)
	if bs.IsZero() {
		t.Fatalf("IsZero() should be false")
	}
	empty := null.ByteSlice([]byte{})
	if !empty.IsZero() {
		t.Fatalf("IsZero() should be true")
	}
	nul := null.NullByteSlice{}
	if !nul.IsZero() {
		t.Fatalf("IsZero() should be true")
	}
}

func TestByteSliceSQLValue(t *testing.T) {
	bs := null.ByteSlice(byteSliceValue)
	val, err := bs.Value()
	fatalIf(t, err, FileLine())
	if !bytes.Equal(byteSliceBase64, val.([]byte)) {
		t.Fatalf("NullByteSlice{..., true}.Value() should return a valid driver.Value ([]byte)")
	}

	empty := null.ByteSlice([]byte{})
	val, err = empty.Value()
	fatalIf(t, err, FileLine())
	if !bytes.Equal([]byte{}, val.([]byte)) {
		t.Fatalf("NullByteSlice{..., true}.Value() should return a valid driver.Value ([]byte)")
	}

	nul := null.NullByteSlice{}
	val, err = nul.Value()
	fatalIf(t, err, FileLine())
	if nil != val {
		t.Fatalf("NullByteSlice{..., false}.Value() should return a nil driver.Value")
	}
}

func TestByteSliceSQLScan(t *testing.T) {
	var bs null.NullByteSlice
	err := bs.Scan(byteSliceBase64)
	fatalIf(t, err, FileLine())
	assertByteSlice(t, byteSliceValue, bs, FileLine())

	var str null.NullByteSlice
	err = str.Scan(string(byteSliceBase64))
	fatalIf(t, err, FileLine())
	assertByteSlice(t, byteSliceValue, str, FileLine())

	var empty null.NullByteSlice
	err = empty.Scan([]byte{})
	fatalIf(t, err, FileLine())
	assertByteSlice(t, []byte{}, empty, FileLine())

	var nul null.NullByteSlice
	err = nul.Scan(nil)
	fatalIf(t, err, FileLine())
	assertNullByteSlice(t, nul, FileLine())

	var wrong null.NullByteSlice
	err = wrong.Scan(int64(42))
	fatalUnless(t, err, FileLine())
}

func TestByteSliceMarshalJSON(t *testing.T) {
	bs := null.ByteSlice(byteSliceValue)
	data, err := json.Marshal(bs)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, string(byteSliceJSON), FileLine())
	data, err = json.Marshal(&bs)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, string(byteSliceJSON), FileLine())

	empty := null.ByteSlice([]byte{})
	data, err = json.Marshal(empty)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, `""`, FileLine())
	data, err = json.Marshal(&empty)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, `""`, FileLine())

	nul := null.NullByteSlice{}
	data, err = json.Marshal(nul)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "null", FileLine())
	data, err = json.Marshal(&nul)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "null", FileLine())
}

func TestByteSliceUnmarshalJSON(t *testing.T) {
	// Successful Valid Parses

	var bs null.NullByteSlice
	err := json.Unmarshal(byteSliceJSON, &bs)
	fatalIf(t, err, FileLine())
	assertByteSlice(t, byteSliceValue, bs, FileLine())

	var quotes null.NullByteSlice
	err = json.Unmarshal([]byte(`""`), &quotes)
	fatalIf(t, err, FileLine())
	assertByteSlice(t, []byte(""), quotes, FileLine())

	var nullStrQuoted null.NullByteSlice
	err = json.Unmarshal([]byte(`"null"`), &nullStrQuoted)
	fatalIf(t, err, FileLine())
	// Skip checking what this decoded to; it's garbage.

	// Successful Null Parses

	var nullStr null.NullByteSlice
	err = json.Unmarshal([]byte("null"), &nullStr)
	fatalIf(t, err, FileLine())
	assertNullByteSlice(t, nullStr, FileLine())

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var badType null.NullByteSlice
	// Ints are never byte slices.
	err = json.Unmarshal(intJSON, &badType)
	fatalUnless(t, err, FileLine())

	var invalid null.NullByteSlice
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf("expected json.SyntaxError, not %T", err)
	}
}

func TestByteSliceMarshalMapValue(t *testing.T) {
	wrapper := struct{ Slice null.NullByteSlice }{null.ByteSlice(byteSliceValue)}
	data, err := maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Slice": byteSliceValue}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Slice": byteSliceValue}, FileLine())

	wrapper = struct{ Slice null.NullByteSlice }{null.ByteSlice([]byte{})}
	data, err = maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Slice": []byte{}}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Slice": []byte{}}, FileLine())

	// Null NullByteSlices should be encoded as "nil"
	wrapper = struct{ Slice null.NullByteSlice }{null.NullByteSlice{}}
	data, err = maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Slice": nil}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Slice": nil}, FileLine())
}
