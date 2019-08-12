package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/relvacode/iso8601"
	"time"
)

// Time is a wrapper around the time.Time type implementing all of the
// pyrrho/encoding/types interfaces detailed in the package comments.
//
// It offers more flexibility in parsing time strings.
type Time struct {
	Time time.Time
}

// Constructors

// NewTime constructs and returns a new, valid Time initialized with the value
// of the given t.
func NewTime(t time.Time) Time {
	return Time{
		Time: t,
	}
}

// NewTimeStr parses a given string, s, as ISO 8601 and returns
func NewTimeStr(s string) (Time, error) {
	if len(s) == 0 {
		return Time{}, nil
	}

	var tmp time.Time
	tmp, err := iso8601.Parse([]byte(s))
	if err != nil {
		return Time{}, err
	}
	return Time{
		Time: tmp,
	}, nil
}

// Setters

// Set modifies the value stored in t.
func (t Time) Set(v time.Time) {
	t.Time = v
}

// SetStr parses the given string. If the string is empty, it replaces t with a
// zero-value time. If parsing fails, it doesn't change t and returns an error.
// If string parsing succeeds, the value stored in t is replaced.
func (t Time) SetStr(s string) error {
	if len(s) == 0 {
		t.Time = time.Time{}
		return nil
	}

	var tmp time.Time
	tmp, err := iso8601.Parse([]byte(s))
	if err != nil {
		return err
	}

	t.Time = tmp
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if t's value has been zero-initialized. For a more meaningful distinction
// between nil and zero, consider using `null.Time`.
func (t Time) IsNil() bool {
	return t.Time == time.Time{}
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if t's value is the zero time instant.
func (t Time) IsZero() bool {
	return t.Time.IsZero()
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
		return fmt.Errorf("types.Time: Scan called on nil pointer")
	}
	switch val := src.(type) {
	case time.Time:
		t.Time = val
		return nil
	case nil:
		t.Time = time.Time{}
		return nil
	default:
		return fmt.Errorf("types.Time: cannot scan type %T (%v)", val, src)
	}
}

// MarshalJSON implements the encoding/json Marshaler interface. It simply wraps
// the underlying time.Time's MarshalJSON method.
func (t Time) MarshalJSON() ([]byte, error) {
	return t.Time.MarshalJSON()
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It will
// decode a given []byte into t so long as the provided []byte
// is a valid JSON representation of an ISO 8601 string. Empty strings and
// the 'null' keyword will both decode into a zero-value Time.
//
// If the decode fails, the value of t will be unchanged.
func (t *Time) UnmarshalJSON(data []byte) error {
	if t == nil {
		return fmt.Errorf("types.Time: UnmarshalJSON called on nil pointer")
	}
	var j interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	switch val := j.(type) {
	case string:
		if len(val) == 0 {
			t.Time = time.Time{}
			return nil
		}
		// TODO: If time.Time.UnmarshalJSON doesn't change the value of the
		// receiver we could skip the temporary object by calling .UnmarshalJSON
		// on t.Time.
		tmp, err := iso8601.Parse([]byte(val))
		if err != nil {
			return err
		}
		t.Time = tmp
		return nil
	case nil:
		t.Time = time.Time{}
		return nil
	default:
		return fmt.Errorf("types.Time: cannot unmarshal JSON of type %T (%v)",
			val, data)
	}
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will return t as an interface{} for use in a map[string]interface{}.
func (t Time) MarshalMapValue() (interface{}, error) {
	return t.Time, nil
}
