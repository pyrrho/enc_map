package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Bool is a wrapper around the database/sql NullBool type that implements all
// of the pyrrho/encoding/types interfaces detailed in the package comments that
// sql.NullBool doesn't implement out of the box.
//
// If the Bool is valid and contains false, it will be considered non-nil, and
// of zero value.
type Bool struct {
	sql.NullBool
}

// Constructors

// NullBool constructs and returns a new null Bool.
func NullBool() Bool {
	return Bool{
		sql.NullBool{
			Bool:  false,
			Valid: false,
		}}
}

// NewBool constructs and returns a new, valid Bool initialized with the value
// of the given b.
func NewBool(b bool) Bool {
	return Bool{
		sql.NullBool{
			Bool:  b,
			Valid: true,
		}}
}

// Getters and Setters

// ValueOrZero returns the value of b if it is valid; otherwise, it returns the
// zero value for a bool (false).
func (b Bool) ValueOrZero() bool {
	if b.Valid {
		return b.Bool
	}
	return false
}

// Set modifies the value stored in b, and guarantees it is valid.
func (b *Bool) Set(v bool) {
	b.Bool = v
	b.Valid = true
}

// Null marks b as null with no meaningful value.
func (b *Bool) Null() {
	b.Bool = false
	b.Valid = false
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if b is null.
func (b Bool) IsNil() bool {
	return !b.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if b is null or if its value is false.
func (b Bool) IsZero() bool {
	return !b.Valid || !b.Bool
}

// MarshalJSON implements the encoding/json Marshaler interface. It will encode
// b into its JSON representation if valid, or 'null' otherwise.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It will
// decode a given []byte into b, so long as the provided []byte is a valid jSON
// representation of a bool or a null.
//
// The keyword 'null' will result in a null NullBool. The keywords 'true' and
// 'false' will result in a valid NullBool containing the value you would
// expect. The strings '"true"', '"false"', '"null"', and `""` are considered to
// be strings -- not keywords -- and will result in an error.
//
// If the decode fails, the value of b will be unchanged.
func (b *Bool) UnmarshalJSON(data []byte) error {
	if b == nil {
		return fmt.Errorf("null.Bool: UnmarshalJSON called on nil pointer")
	}
	var j interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	switch val := j.(type) {
	case bool:
		b.Bool = val
		b.Valid = true
		return nil
	case nil:
		b.Bool = false
		b.Valid = false
		return nil
	default:
		return fmt.Errorf("null.Bool: cannot unmarshal JSON of type %T (%v)",
			val, data)
	}
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode b into its interface{} representation for use in a
// map[string]interface{} if valid, or return nil otherwise.
func (b Bool) MarshalMapValue() (interface{}, error) {
	if b.Valid {
		return b.Bool, nil
	}
	return nil, nil
}
