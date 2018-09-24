package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// NullBool is a wrapper around the database/sql NullBool type that implements
// all of the encoding/type interfaces that sql.NullBool doesn't implement out
// of the box.
//
// If the NullBool is valid and contains false, it will be considered non-nil,
// and zero.
type NullBool struct {
	sql.NullBool
}

// Constructors

// Bool creates a new NullBool based on the type and value of the given
// interface. This function intentionally sacrafices compile-time safety for
// developer convenience.
//
// If the interface is nil or a nil *bool, the new NullBool will be null.
//
// If the interface is a bool or a non-nil *bool, the new NullBool will be
// valid, and will be initialized with the (possibly dereferenced) value of the
// interface.
//
// If the interface is any other type this function will panic.
func Bool(i interface{}) NullBool {
	switch v := i.(type) {
	case bool:
		return BoolFrom(v)
	case *bool:
		return BoolFromPtr(v)
	case nil:
		return NullBool{}
	}
	panic(fmt.Errorf(
		"null.Bool: the given argument (%#v of type %T) was not of type "+
			"bool, *bool, or nil", i, i))
}

// BoolFrom creates a valid NullBool from b.
func BoolFrom(b bool) NullBool {
	return NullBool{sql.NullBool{
		Bool:  b,
		Valid: true,
	}}
}

// BoolFromPtr creates a null NullBool if *b is a nil pointer, or a valid
// NullBool from *b otherwise.
func BoolFromPtr(b *bool) NullBool {
	if b == nil {
		return NullBool{}
	}
	return BoolFrom(*b)
}

// Getters and Setters

// ValueOrZero returns the value of this NullBool if it is valid; otherwise, it
// returns the zero value for a bool.
func (b NullBool) ValueOrZero() bool {
	if b.Valid {
		return b.Bool
	}
	return false
}

// Ptr returns a pointer to this NullBool's value if it is valid; otherwise
// returns a nil pointer. The captured pointer will be able to modify the value
// of this NullBool.
func (b *NullBool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// Set modifies the value stored in this NullBool, and guarantees it is valid.
func (b *NullBool) Set(v bool) {
	b.Bool = v
	b.Valid = true
}

// Null marks this NullBool as null with no meaningful value.
func (b *NullBool) Null() {
	b.Bool = false
	b.Valid = false
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if this NullBool is null.
func (b NullBool) IsNil() bool {
	return !b.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if this NullBool is null or if its value is false.
func (b NullBool) IsZero() bool {
	return !b.Valid || !b.Bool
}

// MarshalJSON implements the encoding/json Marshaler interface. It will encode
// this NullBool into its JSON representation if valid, or 'null' otherwise.
func (b NullBool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It will
// decode a given []byte into this NullBool, so long as the provided []byte is a
// valid JSON representation of a bool or a null.
//
// The keyword 'null' will result in a null NullBool. The keywords 'true' and
// 'false' will result in a valid NullBool containing the value you would
// expect. JSON objects in the form of '{"Bool":<bool>,"Valid":<bool>}' will
// decode directly into this NullBool. The strings '"true"', '"false"',
// '"null"', and `""` are considered to be strings -- not keywords -- and will
// result in an error.
//
// If the decode fails, the value of this NullBool will be unchanged.
func (b *NullBool) UnmarshalJSON(data []byte) error {
	var j interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	switch val := j.(type) {
	case bool:
		b.Bool = val
		b.Valid = true
		return nil
	case map[string]interface{}:
		// If we've received a JSON object, try to decode it directly into our
		// sql.NullBool. Return any errors that occur.
		// TODO: Make sure this, if `data` is malformed, can't affect the value
		//       of this NullBool.
		return json.Unmarshal(data, &b.NullBool)
	case nil:
		b.Bool = false
		b.Valid = false
		return nil
	default:
		return fmt.Errorf("null: cannot unmarshal %T (%#v) into Go value of type null.NullBool", j, j)
	}
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode this NullBool into its interface{} representation for use in a
// map[string]interface{} if valid, or return nil otherwise.
func (b NullBool) MarshalMapValue() (interface{}, error) {
	if b.Valid {
		return b.Bool, nil
	}
	return nil, nil
}
