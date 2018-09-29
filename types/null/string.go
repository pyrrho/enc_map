package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// String is a wrapper around the database/sql NullString type that implements
// all of the pyrrho/encoding/types interfaces detailed in the package comments
// that sql.NullString doesn't implement out of the box.
//
// If the String is valid and contains the empty string, it will be considered
// non-nil, and of zero value.
type String struct {
	sql.NullString
}

// Constructors

// NullString constructs and returns a new null String.
func NullString() String {
	return String{
		sql.NullString{
			String: "",
			Valid:  false,
		}}
}

// NewString constructs and returns a new, valid String initialized with the
// value of the given s.
func NewString(s string) String {
	return String{
		sql.NullString{
			String: s,
			Valid:  true,
		}}
}

// Getters and Setters

// ValueOrZero returns the value of s if it is valid; otherwise it returns the
// zero value for a string ("").
func (s String) ValueOrZero() string {
	if !s.Valid {
		return ""
	}
	return s.String
}

// Set modifies the value stored in s, and guarantees it is valid.
func (s *String) Set(v string) {
	s.String = v
	s.Valid = true
}

// Null marks s as null with no meaningful value.
func (s *String) Null() {
	s.Valid = false
}

// Interface

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if s is null.
func (s String) IsNil() bool {
	return !s.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if s is null or if its value is the empty string.
func (s String) IsZero() bool {
	return !s.Valid || s.String == ""
}

// MarshalJSON implements the encoding/json Marshaler interface. It will return
// the value of s if valid, otherwise 'null'.
func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It will
// decode a given []byte into s, so long as the provided []byte is a valid JSON
//string or a null.
//
// An empty string will result in a valid-but-empty String. The keyword 'null'
// will result in a null String. The string '"null"' is considered to be a
// string -- not a keyword -- and will result in a valid String.
//
// If the decode fails, the value of s will be unchanged.
func (s *String) UnmarshalJSON(data []byte) error {
	if s == nil {
		return fmt.Errorf("null.String: UnmarshalJSON called on nil pointer")
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
	case nil:
		s.String = ""
		s.Valid = false
		return nil
	default:
		return fmt.Errorf("null.String: cannot unmarshal JSON of type %T (%v)",
			val, data)
	}
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode s into an interface{} representation for use in a
// map[string]interface{} if valid, or return nil otherwise.
func (s String) MarshalMapValue() (interface{}, error) {
	if s.Valid {
		return s.String, nil
	}
	return nil, nil
}
