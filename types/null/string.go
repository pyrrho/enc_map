package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// NullString is a wrapper around the database/sql NullString type that
// implements all of the encoding/type interfaces that sql.NullString doesn't
// implement out of the box.
//
// If the NullString is valid and contains the empty string, it will be
// considered non-nil, and zero.
type NullString struct {
	sql.NullString
}

// String creates a new NullString based on the type and value of the given
// interface. This function intentionally sacrafices compile-time safety for
// developer convenience.
//
// If the interface is nil or a nil *string, the new NullString will be null.
//
// If the interface is a string or a non-nil *string, the new NullString will be
// valid, and will be initialized with the (possibly dereferenced) value of the
// interface.
//
// If the interface is any other type this function will panic.
func String(i interface{}) NullString {
	switch v := i.(type) {
	case string:
		return StringFrom(v)
	case *string:
		return StringFromPtr(v)
	case nil:
		return NullString{}
	}
	panic(fmt.Errorf(
		"null.NullString: invalid constructor argument; %#v of type %T "+
			"is not of type string, *string, or nil", i, i))
}

// StringFrom creates a valid String from s.
func StringFrom(s string) NullString {
	return NullString{sql.NullString{
		String: s,
		Valid:  true,
	}}
}

// StringFromPtr creates a valid String from *s.
func StringFromPtr(s *string) NullString {
	if s == nil {
		return NullString{}
	}
	return StringFrom(*s)
}

// ValueOrZero returns the value of this NullString if it is valid; otherwise it
// returns the zero value for a string.
func (s NullString) ValueOrZero() string {
	if !s.Valid {
		return ""
	}
	return s.String
}

// Ptr returns a pointer to this NullString's value if it is valid; otherwise
// returns a nil pointer. The captured pointer will be able to modify the value
// of this NullString.
func (s *NullString) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// Set modifies the value stored in this NullString, and guarantees it is valid.
func (s *NullString) Set(v string) {
	s.String = v
	s.Valid = true
}

// Null marks this NullString as null with no meaningful value.
func (s *NullString) Null() {
	s.Valid = false
}

// Interface

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if this NullString is null.
func (s NullString) IsNil() bool {
	return !s.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if this NullString is null or if its value is 0.
func (s NullString) IsZero() bool {
	return !s.Valid || s.String == ""
}

// MarshalJSON implements the encoding/json Marshaler interface. It will return
// the value of this NullString if valid, otherwise 'null'.
func (s NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It will
// decode a given []byte into this NullString, so long as the provided []byte
// is a valid JSON string or a null.
//
// An empty string will result in a valid-but-empty NullString. The keyword
// 'null' will result in a null NullString. JSON objects in the form of
// '{"String":<string>,"Valid":<bool>`}' will decode directly into this
// NullString. The string '"null"' is considered to be a string -- not a keyword
// -- and will result in a valid NullString.
//
// If the decode fails, the value of this NullString will be unchanged.
func (s *NullString) UnmarshalJSON(data []byte) error {
	if s == nil {
		return fmt.Errorf("null.NullString: UnmarshalJSON called on nil pointer")
	}
	var j interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	switch val := j.(type) {
	case string:
		s.String = val
		s.Valid = true
		return nil
	case map[string]interface{}:
		return json.Unmarshal(data, &s.NullString)
	case nil:
		s.String = ""
		s.Valid = false
		return nil
	default:
		return fmt.Errorf("null.NullString: cannot unmarshal JSON of type %T (%v)",
			val, data)
	}
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode this NullString into an interface{} representation for use in a
// map[string]interface{} if valid, or return nil otherwise.
func (s NullString) MarshalMapValue() (interface{}, error) {
	if s.Valid {
		return s.String, nil
	}
	return nil, nil
}
