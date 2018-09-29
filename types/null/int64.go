package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
)

// Int64 is a wrapper around the database/sql NullInt64 type that implements all
// of the pyrrho/encoding/types interfaces detailed in the package comments that
// sql.NullInt64 doesn't implement out of the box.
//
// If the Int64 is valid and contains 0, it will be considered non-nil, and of
// zero value.
type Int64 struct {
	sql.NullInt64
}

// Constructors

// NullInt64 constructs and returns a new null Int64.
func NullInt64() Int64 {
	return Int64{
		sql.NullInt64{
			Int64: 0,
			Valid: false,
		}}
}

// NewInt64 constructs and returns a new, valid Int64 initialized with the value
// of the given i.
func NewInt64(i int64) Int64 {
	return Int64{
		sql.NullInt64{
			Int64: i,
			Valid: true,
		}}
}

// Getters and Setters

// ValueOrZero returns the value of i if it is valid; otherwise it returns the
// zero value for a int64 (0).
func (i Int64) ValueOrZero() int64 {
	if !i.Valid {
		return 0
	}
	return i.Int64
}

// Set modifies the value stored in i, and guarantees it is valid.
func (i *Int64) Set(v int64) {
	i.Int64 = v
	i.Valid = true
}

// Null marks i as null with no meaningful value.
func (i *Int64) Null() {
	i.Int64 = 0
	i.Valid = false
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if i is null.
func (i Int64) IsNil() bool {
	return !i.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if i is null or if its value is 0.
func (i Int64) IsZero() bool {
	return !i.Valid || i.Int64 == 0
}

// MarshalJSON implements the encoding/json Marshaler interface. It will encode
// i into its JSON representation if valid, or 'null' otherwise.
func (i Int64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It will
// decode a given []byte into i, so long as the provided []byte is a valid JSON
// representation of an int. The 'null' keyword will decode into a null Int64.
//
// If the decode fails, the value of i will be unchanged.
func (i *Int64) UnmarshalJSON(data []byte) error {
	if i == nil {
		return fmt.Errorf("null.Int64: UnmarshalJSON called on nil pointer")
	}
	var j interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	switch val := j.(type) {
	case float64:
		// Perform a second unmarshal, this time into an int64. This give the
		// JSON parse a chance to meaningfully fail (eg. if the conversion from
		// float to integer will result in a loss of precision).
		var tmp int64
		err := json.Unmarshal(data, &tmp)
		if err != nil {
			return err
		}
		i.Int64 = tmp
		i.Valid = true
		return nil
	case nil:
		i.Int64 = 0
		i.Valid = false
		return nil
	default:
		return fmt.Errorf("null.Int64: cannot unmarshal JSON of type %T (%v)",
			val, data)
	}
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode i into its interface{} representation for use in a
// map[string]interface{} if valid, or return nil otherwise.
func (i Int64) MarshalMapValue() (interface{}, error) {
	if i.Valid {
		return i.Int64, nil
	}
	return nil, nil
}
