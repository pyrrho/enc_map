package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Time is a nullable wrapper around the time.Time type implementing all of the
// the pyrrho/encoding/types interfaces detailed in the package comments.
//
// If the Time is valid and contains the zero time instant, it will be
// considered non-null, and of zero value.
type Time struct {
	Time  time.Time
	Valid bool
}

// Constructors

// NullTime constructs and returns a new null Time.
func NullTime() Time {
	return Time{
		Time:  time.Time{},
		Valid: false,
	}
}

// NewTime constructs and returns a new, valid Time initialized with the value
// of the given t.
func NewTime(t time.Time) Time {
	return Time{
		Time:  t,
		Valid: true,
	}
}

// NewTimeStr parses a given string, s, as an RFC 3339 and returns
func NewTimeStr(s string) (Time, error) {
	if len(s) == 0 {
		return Time{}, nil
	}

	var tmp time.Time
	tmp, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return Time{}, err
	}
	return Time{
		Time:  tmp,
		Valid: true,
	}, nil
}

// Getters and Setters

// ValueOrZero returns the value of t if it is valid; otherwise it returns the
// zero value for a time.Time.
func (t Time) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// Set modifies the value stored in t, and guarantees it is valid.
func (t *Time) Set(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Null marks t as null with no meaningful value.
func (t *Time) Null() {
	t.Time = time.Time{}
	t.Valid = false
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if t is null.
func (t Time) IsNil() bool {
	return !t.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if t is null or if its value is the zero time instant.
func (t Time) IsZero() bool {
	return !t.Valid || t.Time == time.Time{}
}

// Value implements the database/sql/driver Valuer interface. As time.Time and
// nil are both valid types to be stored in a driver.Value, it will return this
// NullTime's value if valid, or nil otherwise.
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// Scan implements the database/sql Scanner interface. It will receive a value
// from an SQL database and assign it to t, so long as the provided data is of
// type nil or time.Time. All other types will result in an error.
func (t *Time) Scan(src interface{}) error {
	if t == nil {
		return fmt.Errorf("null.Time: Scan called on nil pointer")
	}
	switch val := src.(type) {
	case time.Time:
		t.Time = val
		t.Valid = true
		return nil
	case nil:
		t.Time = time.Time{}
		t.Valid = false
		return nil
	default:
		return fmt.Errorf("null.Time: cannot scan type %T (%v)", val, src)
	}
}

// MarshalJSON implements the encoding/json Marshaler interface. It will encode
// t into its JSON RFC 3339 string representation if valid, or
// 'null' otherwise.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalJSON()
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It will
// decode a given []byte into t so long as the provided []byte
// is a valid JSON representation of an RFC 3339 string. Empty strings and
// the 'null' keyword will both decode into a null NullTime.
//
// If the decode fails, the value of t will be unchanged.
func (t *Time) UnmarshalJSON(data []byte) error {
	if t == nil {
		return fmt.Errorf("null.Time: UnmarshalJSON called on nil pointer")
	}
	var j interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	switch val := j.(type) {
	case string:
		if len(val) == 0 {
			t.Time = time.Time{}
			t.Valid = false
			return nil
		}
		// TODO: If time.Time.UnmarshalJSON doesn't change the value of the
		// receiver we could skip the temporary object by calling .UnmarshalJSON
		// on t.Time.
		tmp, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return err
		}
		t.Time = tmp
		t.Valid = true
		return nil
	case nil:
		t.Time = time.Time{}
		t.Valid = false
		return nil
	default:
		return fmt.Errorf("null.Time: cannot unmarshal JSON of type %T (%v)",
			val, data)
	}
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode t into an interface{} representation for use in a
// map[Time]interface{} if valid, or return nil otherwise.
func (t Time) MarshalMapValue() (interface{}, error) {
	if t.Valid {
		return t.Time, nil
	}
	return nil, nil
}
