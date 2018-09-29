package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// Float64 is a wrapper around the database/sql NullFloat64 type that implements
// all of the pyrrho/encoding/types interfaces detailed in the package comments
// that sql.NullFloat64 doesn't implement out of the box.
//
// If the Float64 is valid and contains 0, it will be considered non-nil, and of
// zero value.
type Float64 struct {
	sql.NullFloat64
}

// Constructors

// NullFloat64 constructs and returns a new null Float64.
func NullFloat64() Float64 {
	return Float64{
		sql.NullFloat64{
			Float64: 0.0,
			Valid:   false,
		}}
}

// NewFloat64 constructs and returns a new, valid Float64 initialized with the
// value of the given f.
func NewFloat64(f float64) Float64 {
	return Float64{
		sql.NullFloat64{
			Float64: f,
			Valid:   true,
		}}
}

// Getters and Setters

// ValueOrZero returns the value of f if it is valid; otherwise it returns the
// zero value for a float64 (0.0).
func (f Float64) ValueOrZero() float64 {
	if !f.Valid {
		return 0.0
	}
	return f.Float64
}

// Set modifies the value stored in f, and guarantees it is valid.
func (f *Float64) Set(v float64) {
	f.Float64 = v
	f.Valid = true
}

// Null marks f as null with no meaningful value.
func (f *Float64) Null() {
	f.Float64 = 0.0
	f.Valid = false
}

// Interface

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if f is null.
func (f Float64) IsNil() bool {
	return !f.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if f is null or if its value is 0.
func (f Float64) IsZero() bool {
	return !f.Valid || f.Float64 == 0.0
}

// MarshalJSON implements the encoding/json Marshaler interface. It will attempt
// to encode f into its JSON representation if valid. If the contained value is
// +/-INF or NaN, a json.UnsupportedValueError will be returned. If f is not
// valid, it will encode to 'null'.
func (f Float64) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}
	if math.IsInf(f.Float64, 0) || math.IsNaN(f.Float64) {
		return nil, &json.UnsupportedValueError{
			Value: reflect.ValueOf(f.Float64),
			Str:   strconv.FormatFloat(f.Float64, 'g', -1, 64),
		}
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It will
// decode a given []byte into f, so long as the provided []byte is a valid JSON
// representation of a float or null. The 'null' keyword will decode into a null
// Float64.
//
// If the decode fails, the value of f will be unchanged.
func (f *Float64) UnmarshalJSON(data []byte) error {
	if f == nil {
		return fmt.Errorf("null.Float64: UnmarshalJSON called on nil pointer")
	}
	var j interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	switch val := j.(type) {
	case float64:
		f.Float64 = val
		f.Valid = true
		return nil
	case nil:
		f.Float64 = 0
		f.Valid = false
		return nil
	default:
		return fmt.Errorf("null.Float64: cannot unmarshal JSON of type %T (%v)",
			val, data)
	}
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode f into its interface{} representation for use in a
// map[string]interface{} if valid, or return nil otherwise.
func (f Float64) MarshalMapValue() (interface{}, error) {
	if f.Valid {
		return f.Float64, nil
	}
	return nil, nil
}
