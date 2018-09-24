package null_test

import (
	"encoding/json"
	"testing"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types/null"
)

// Helper Functions

func assertString(t *testing.T, expected string, b null.NullString, fileLine string) {
	if !b.Valid {
		t.Fatalf("%s: NullString is null, but should be valid", fileLine)
	}
	if expected != b.String {
		t.Fatalf("%s: %v ≠ %v", fileLine, expected, b.String)
	}
}

func assertNullString(t *testing.T, b null.NullString, fileLine string) {
	if b.Valid {
		t.Fatalf("%s: NullString is valid, but should be null", fileLine)
	}
}

func TestStringFrom(t *testing.T) {
	assertString(t, "test", null.StringFrom("test"), FileLine())
	assertString(t, "", null.StringFrom(""), FileLine())
}

func TestStringFromPtr(t *testing.T) {
	s := "test"
	assertString(t, "test", null.StringFromPtr(&s), FileLine())
	assertNullString(t, null.StringFromPtr(nil), FileLine())
}

func TestStringCtor(t *testing.T) {
	v := "test"
	var nilPtr *string

	assertString(t, "true", null.String("true"), FileLine())
	assertString(t, "", null.String(""), FileLine())
	assertString(t, v, null.String(v), FileLine())
	assertString(t, v, null.String(&v), FileLine())
	assertNullString(t, null.String(nil), FileLine())
	assertNullString(t, null.String(nilPtr), FileLine())
}

func TestFailureNewStringFromInt(t *testing.T) {
	defer ShouldPanic(t, FileLine())
	_ = null.String(0)
}

func TestFailureNewStringFromBool(t *testing.T) {
	defer ShouldPanic(t, FileLine())
	_ = null.String(true)
}

func TestStringValueOrZero(t *testing.T) {
	valid := null.String("test")
	if valid.ValueOrZero() != "test" {
		t.Fatalf("unexpected ValueOrZero, %v ≠ %v", "test", valid.ValueOrZero())
	}

	nul := null.NullString{}
	if nul.ValueOrZero() != "" {
		t.Fatalf("unexpected ValueOrZero, %v ≠ %v", "", nul.ValueOrZero())
	}
}

func TestStringPtr(t *testing.T) {
	str := null.String("test")
	ptr := str.Ptr()
	if *ptr != "test" {
		t.Fatalf("bad %s string: %#v ≠ %v\n", "pointer", ptr, "test")
	}

	null := null.NullString{}
	ptr = null.Ptr()
	if ptr != nil {
		t.Fatalf("bad %s string: %#v ≠ %v\n", "nil pointer", ptr, "nil")
	}
}

func TestStringSet(t *testing.T) {
	s := null.NullString{}
	assertNullString(t, s, FileLine())
	s.Set("test")
	assertString(t, "test", s, FileLine())
	s.Set("")
	assertString(t, "", s, FileLine())
}

func TestStringNull(t *testing.T) {
	s := null.String("test")
	assertString(t, "test", s, FileLine())
	s.Null()
	assertNullString(t, s, FileLine())
}

func TestStringIsNil(t *testing.T) {
	a := null.String("test")
	if a.IsNil() {
		t.Fatal(`NullString{"test", true}.IsNil() should be false`)
	}
	b := null.String("")
	if b.IsNil() {
		t.Fatal(`NullString{"", true}.IsNil() should be false`)
	}
	nul := null.NullString{}
	if !nul.IsNil() {
		t.Fatal("NullString{..., false}.IsNil() should be true")
	}
}

func TestStringIsZero(t *testing.T) {
	a := null.String("test")
	if a.IsZero() {
		t.Fatal(`NullString{"test", true}.IsZero() should be false`)
	}
	b := null.String("")
	if !b.IsZero() {
		t.Fatal(`NullString{"", true}.IsZero() should be true`)
	}
	nul := null.NullString{}
	if !nul.IsZero() {
		t.Fatal("NullString{..., false}.IsZero() should be true")
	}
}

func TestStringSQLValue(t *testing.T) {
	s := null.String("test")
	val, err := s.Value()
	fatalIf(t, err, FileLine())
	if "test" != val.(string) {
		t.Fatalf(`NullString{"test", true}.Value() should return a valid driver.Value (string)`)
	}

	nul := null.NullString{}
	val, err = nul.Value()
	fatalIf(t, err, FileLine())
	if nil != val {
		t.Fatalf("NullString{..., false}.Value() should return a nil driver.Value")
	}
}

func TestStringSQLScan(t *testing.T) {
	var str null.NullString
	err := str.Scan("test")
	fatalIf(t, err, FileLine())
	assertString(t, "test", str, FileLine())

	var empty null.NullString
	err = empty.Scan("")
	fatalIf(t, err, FileLine())
	assertString(t, "", empty, FileLine())

	var nul null.NullString
	err = nul.Scan(nil)
	fatalIf(t, err, FileLine())
	assertNullString(t, nul, FileLine())

	// NB. Scan is aggressive about converting values to strings. UnmarshalJSON
	// are less so.
	var i null.NullString
	err = i.Scan(12345)
	fatalIf(t, err, FileLine())
	assertString(t, "12345", i, FileLine())

	var f null.NullString
	err = f.Scan(1.2345)
	fatalIf(t, err, FileLine())
	assertString(t, "1.2345", f, FileLine())

	var b null.NullString
	err = b.Scan(true)
	fatalIf(t, err, FileLine())
	assertString(t, "true", b, FileLine())
}

func TestStringMarshalJSON(t *testing.T) {
	str := null.String("test")
	data, err := json.Marshal(str)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, `"test"`, FileLine())
	data, err = json.Marshal(&str)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, `"test"`, FileLine())

	zero := null.String("")
	data, err = json.Marshal(zero)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, `""`, FileLine())
	data, err = json.Marshal(&zero)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, `""`, FileLine())

	null := null.NullString{}
	data, err = json.Marshal(null)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "null", FileLine())
	data, err = json.Marshal(&null)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, data, "null", FileLine())
}

func TestStringUnmarshalJSON(t *testing.T) {
	// Successful Valid Parses

	var str null.NullString
	err := json.Unmarshal(stringJSON, &str)
	fatalIf(t, err, FileLine())
	assertString(t, "test", str, FileLine())

	var quotes null.NullString
	err = json.Unmarshal([]byte(`""`), &quotes)
	fatalIf(t, err, FileLine())
	assertString(t, "", quotes, FileLine())

	var nullStr null.NullString
	err = json.Unmarshal([]byte(`"null"`), &nullStr)
	fatalIf(t, err, FileLine())
	assertString(t, "null", nullStr, FileLine())

	var validObj null.NullString
	err = json.Unmarshal(validStringJSONObj, &validObj)
	fatalIf(t, err, FileLine())
	assertString(t, "test", validObj, FileLine())

	var validButEmptyObj null.NullString
	err = json.Unmarshal(validButEmptyStringJSONObj, &validButEmptyObj)
	fatalIf(t, err, FileLine())
	assertString(t, "", validButEmptyObj, FileLine())

	// Successful Null Parses

	var nul null.NullString
	err = json.Unmarshal([]byte("null"), &nul)
	fatalIf(t, err, FileLine())
	assertNullString(t, nul, FileLine())

	var nullObj null.NullString
	err = json.Unmarshal(nullStringJSONObj, &nullObj)
	fatalIf(t, err, FileLine())
	assertNullString(t, nullObj, FileLine())

	// Unsuccessful Parses
	// TODO: make types for type mismatches on parsing, and check that the
	// correct error type is being returned here.

	var badType null.NullString
	// Ints are never string.
	err = json.Unmarshal(intJSON, &badType)
	fatalUnless(t, err, FileLine())

	var invalid null.NullString
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf("expected json.SyntaxError, not %T", err)
	}
}

func TestStringInStructMarshalJSON(t *testing.T) {
	type stringTestStruct struct {
		NullString null.NullString `json:"null_string"`
		String     string          `json:"string"`
	}

	s := stringTestStruct{
		NullString: null.String("valid"),
		String:     "test",
	}
	sj, err := json.Marshal(s)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, sj, `{"null_string":"valid","string":"test"}`, FileLine())

	s = stringTestStruct{
		NullString: null.NullString{},
		String:     "test",
	}
	sj, err = json.Marshal(s)
	fatalIf(t, err, FileLine())
	assertJSONEquals(t, sj, `{"null_string":null,"string":"test"}`, FileLine())
}

func TestStringMarshalMapValue(t *testing.T) {
	wrapper := struct{ Slice null.NullString }{null.String("test")}
	data, err := maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Slice": "test"}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Slice": "test"}, FileLine())

	wrapper = struct{ Slice null.NullString }{null.String("")}
	data, err = maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Slice": ""}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Slice": ""}, FileLine())

	// Null NullStrings should be encoded as "nil"
	wrapper = struct{ Slice null.NullString }{null.NullString{}}
	data, err = maps.Marshal(wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Slice": nil}, FileLine())
	data, err = maps.Marshal(&wrapper)
	fatalIf(t, err, FileLine())
	assertMapEquals(t, data, map[string]interface{}{"Slice": nil}, FileLine())
}
