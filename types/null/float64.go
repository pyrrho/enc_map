package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// NullFloat64 is a wrapper around the database/sql NullFloat64 type that
// implements all of the encoding/type interfaces that sql.NullFloat64 doesn't
// implement out of the box.
//
// If the NullFloat64 is valid and contains 0, it will be considered non-nil,
// and zero.
type NullFloat64 struct {
	sql.NullFloat64
}

// Constructors

// Float64 creates a new NullFloat64 based on the type and value of the given
// interface. This function intentionally sacrafices compile-time safety for
// developer convenience.
//
// If the interface is nil, a nil *float64, or a nil *float32, the new
// NullFloat64 will be null.
//
// If the interface is a float64, a float32, an int, a non-nil *float64, or a
// non-nil *float32, the new NullFloat64 will be valid, and will be initialized
// with the (possibly dereferenced) value of the interface.
//
// If the interface is any other type, this function will panic.
func Float64(i interface{}) NullFloat64 {
	switch v := i.(type) {
	case float64:
		return Float64From(v)
	case *float64:
		return Float64FromPtr(v)
	case float32:
		return Float64From(float64(v))
	case *float32:
		if v == nil {
			return NullFloat64{}
		}
		return Float64From(float64(*v))
	case int:
		return Float64From(float64(v))
	case nil:
		return NullFloat64{}
	}
	panic(fmt.Errorf(
		"null.Float64: the given argument (%#v of type %T) was not of type "+
			"int, float64, *float64, float32, *float32, or nil", i, i))
}

// Float64From creates a valid NullFloat64 from f.
func Float64From(f float64) NullFloat64 {
	return NullFloat64{sql.NullFloat64{
		Float64: f,
		Valid:   true,
	}}
}

// Float64FromPtr creates a valid NullFloat64 from *f.
func Float64FromPtr(f *float64) NullFloat64 {
	if f == nil {
		return NullFloat64{}
	}
	return Float64From(*f)
}

// Getters and Setters

// ValueOrZero returns the value of this NullFloat64 if it is valid; otherwise
// it returns the zero value for a float64.
func (f NullFloat64) ValueOrZero() float64 {
	if !f.Valid {
		return 0
	}
	return f.Float64
}

// Ptr returns a pointer to this NullFloat64's value if it is valid; otherwise
// returns a nil pointer. The captured pointer will be able to modify the value
// of this NullFloat64.
func (f *NullFloat64) Ptr() *float64 {
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

// Set modifies the value stored in this NullFloat64, and guarantees it is
// valid.
func (f *NullFloat64) Set(v float64) {
	f.Float64 = v
	f.Valid = true
}

// Null marks this NullFloat64 as null with no meaningful value.
func (f *NullFloat64) Null() {
	f.Float64 = 0
	f.Valid = false
}

// Interface

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if this NullFloat64 is null.
func (f NullFloat64) IsNil() bool {
	return !f.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if this NullFloat64 is null or if its value is 0.
func (f NullFloat64) IsZero() bool {
	return !f.Valid || f.Float64 == 0.0
}

// MarshalText implements the encoding TextMarshaler interface. It will encode
// this NullFloat64 into its textual representation if valid, or an empty string
// otherwise.
func (f NullFloat64) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

// UnmarshalText implements the encoding TextUnmarshaler interface. It will
// decode a given []byte into this NullFloat64, so long as the provided string
// is a valid textual representation of a float or a null. Empty strings and
// "null" will decode into a null NullFloat64.
//
// If the decode fails, the value of this NullFloat64 will be unchanged.
func (f *NullFloat64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		f.Float64 = 0
		f.Valid = false
		return nil
	}
	tmp, err := strconv.ParseFloat(string(text), 64)
	if err != nil {
		return err
	}
	f.Float64 = tmp
	f.Valid = true
	return nil
}

// MarshalJSON implements the encoding/json Marshaler interface. It will attempt
// to encode this NullFloat64 into its JSON representation if valid. If the
// contained value is +/-INF or NaN, a json.UnsupportedValueError will be
// returned. If this NullFloat64 is not valid, it will encode to 'null'.
func (f NullFloat64) MarshalJSON() ([]byte, error) {
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
// decode a given []byte into this NullFloat64, so long as the provided []byte
// is a valid JSON representation of a float or a null.
//
// Empty strings and 'null' will both decode into a null NullFloat64. JSON
// objects in the form of '{"Float64":<float>,"Valid":<bool>}' will decode
// directly into this NullFloat64.
//
// If the decode fails, the value of this NullFloat64 will be unchanged.
func (f *NullFloat64) UnmarshalJSON(data []byte) error {
	var j interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	switch val := j.(type) {
	case float64:
		f.Float64 = val
		f.Valid = true
		return nil
	case map[string]interface{}:
		// If we've received a JSON object, try to decode it directly into our
		// sql.NullBool. Return any errors that occur.
		// TODO: Make sure this, if `data` is malformed, can't affect the value
		//       of this NullBool.
		return json.Unmarshal(data, &f.NullFloat64)
	case nil:
		f.Float64 = 0
		f.Valid = false
		return nil
	default:
		return fmt.Errorf(
			"null: cannot unmarshal %T (%#v) into Go value of type "+
				"null.NullFloat64",
			j, j,
		)
	}
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode this NullFloat64 into its interface{} representation for use in a
// map[string]interface{} if valid, or return nil otherwise.
func (f NullFloat64) MarshalMapValue() (interface{}, error) {
	if f.Valid {
		return f.Float64, nil
	}
	return nil, nil
}
