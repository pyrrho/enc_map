package null_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// This file performs no meaningful tests; it's an aggregate for shared state
// and helper functions that are used in the package's test suite.

var (
	timeString = "2012-12-21T21:21:21Z"
	timeValue  = time.Date(
		2012, time.December, 21,
		21, 21, 21, 0,
		time.UTC,
	)
	zeroTimeString = "0001-01-01T00:00:00Z"
	zeroTimeValue  = time.Time{}

	boolFalseJSON   = []byte(`false`)
	boolTrueJSON    = []byte(`true`)
	boolStringJSON  = []byte(`"true"`)
	floatJSON       = []byte(`1.2345`)
	floatStringJSON = []byte(`"1.2345"`)
	intJSON         = []byte(`12345`)
	intStringJSON   = []byte(`"12345"`)
	invalidJSON     = []byte(`:)`)
	stringJSON      = []byte(`"test"`)
	timeJSON        = []byte(`"2012-12-21T21:21:21Z"`)
	zeroTimeJSON    = []byte(`"0001-01-01T00:00:00Z"`)
)

func FileLine() string {
	_, fileName, fileLine, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("\n%s:%d", fileName, fileLine)
	}
	return ""
}

func ShouldPanic(t *testing.T, fileLine string) {
	if r := recover(); r == nil {
		t.Fatalf("%s: This test did not panic", fileLine)
	}
}

func fatalIf(t *testing.T, err error, fileLine string) {
	if err != nil {
		t.Fatalf("%s: unexpected error %v", fileLine, err)
	}
}

func fatalUnless(t *testing.T, err error, fileLine string) {
	if err == nil {
		t.Fatalf("%s: error must not be nil", fileLine)
	}
}

func assertJSONEquals(t *testing.T, actual []byte, expected string, fileLine string) {
	if !bytes.Equal(actual, []byte(expected)) {
		t.Fatalf("%s: %s ≠ %s", fileLine, actual, expected)
	}
}

func assertMapEquals(t *testing.T, actual map[string]interface{}, expected map[string]interface{}, fileLine string) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("%s: %s ≠ %s", fileLine, actual, expected)
	}
}

// Descriptive Tests
// -----------------
// These don't actually test anything, rather they demonstrate behavior of the
// language and other packages. They are included because they may help clarify
// why certain choices were made within this library.

func TestJSONUnmarshalErrors(t *testing.T) {
	// encoding/json considers and empty byte stream to be invalid JSON.
	// This is correct; _nothing_ is not a valid JSON _something_.
	var iface interface{}
	err := json.Unmarshal([]byte(""), &iface)
	fatalUnless(t, err, FileLine())

	// A pair of double-quotes with nothing between them will be parsed as a
	// string or a []byte, but nothing else.
	var (
		s  string
		bs []byte
		i  int
		is []int
	)
	err = json.Unmarshal([]byte(`""`), &s)
	fatalIf(t, err, FileLine())
	err = json.Unmarshal([]byte(`""`), &bs)
	fatalIf(t, err, FileLine())
	err = json.Unmarshal([]byte(`""`), &i)
	fatalUnless(t, err, FileLine())
	err = json.Unmarshal([]byte(`""`), &is)
	fatalUnless(t, err, FileLine())
}
