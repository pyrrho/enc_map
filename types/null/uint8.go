package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// Uint8 is a nullable wrapper around the uint8 type that implementing all of
// the pyrrho/encoding/types interfaces detailed in the package comments.
//
// If the Uint8 is valid and contains 0, it will be considered non-nil, and of
// zero value.
type Uint8 struct {
	Uint8 uint8
	Valid bool
}

// Constructors

// NullUint8 constructs and returns a new null Uint8.
func NullUint8() Uint8 {
	return Uint8{
		Uint8: 0,
		Valid: false,
	}
}

// NewUint8 constructs and returns a new, valid Uint8 initialized with the value
// of the given i.
func NewUint8(i uint8) Uint8 {
	return Uint8{
		Uint8: i,
		Valid: true,
	}
}

// Getters and Setters

// ValueOrZero returns the value of i if it is valid; otherwise it returns the
// zero value for a uint8 (0).
func (i Uint8) ValueOrZero() uint8 {
	if !i.Valid {
		return 0
	}
	return i.Uint8
}

// Set modifies the value stored in i, and guarantees it is valid.
func (i *Uint8) Set(v uint8) {
	i.Uint8 = v
	i.Valid = true
}

// Null marks i as null with no meaningful value.
func (i *Uint8) Null() {
	i.Uint8 = 0
	i.Valid = false
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if i is null.
func (i Uint8) IsNil() bool {
	return !i.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if i is null or if its value is 0.
func (i Uint8) IsZero() bool {
	return !i.Valid || i.Uint8 == 0
}

// Value implements the database/sql/driver Valuer interface. Nil is a valid
// type to be stored in a driver.Value, but uint8 isn't, so if this Uint8 is
// valid it will cast its uint8 to an int64.
func (i Uint8) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return int64(i.Uint8), nil
}

// Scan implements the database/sql Scanner interface. It will receive a value
// from an SQL database and assign it to i, so long as the provided data is of
// type nil, uint8, string, or another integer or float type that doesn't
// overflow uint8. All other types will result in an error.
func (i *Uint8) Scan(src interface{}) error {
	// Helpers for repetitive int scans
	scanUint := func(src interface{}, vi uint64) error {
		if vi > math.MaxUint8 {
			return fmt.Errorf("null.Uint8: failed to scan type %T (%v): overflow", src, src)
		}
		i.Uint8 = uint8(vi)
		i.Valid = true
		return nil
	}
	scanInt := func(src interface{}, vi int64) error {
		if vi > math.MaxUint8 {
			return fmt.Errorf("null.Uint8: failed to scan type %T (%v): overflow", src, src)
		} else if vi < 0 {
			return fmt.Errorf("null.Uint8: failed to scan type %T (%v): negative value", src, src)
		}
		i.Uint8 = uint8(vi)
		i.Valid = true
		return nil
	}

	if src == nil {
		i.Uint8 = 0
		i.Valid = false
		return nil
	}
	if i == nil {
		return fmt.Errorf("null.Uint8: Scan called on nil pointer")
	}

	switch val := src.(type) {
	case uint8:
		i.Uint8 = val
		i.Valid = true
		return nil
	case uint:
		return scanUint(src, uint64(val))
	case uint16:
		return scanUint(src, uint64(val))
	case uint32:
		return scanUint(src, uint64(val))
	case uint64:
		return scanUint(src, uint64(val))
	case int:
		return scanInt(src, int64(val))
	case int8:
		return scanInt(src, int64(val))
	case int16:
		return scanInt(src, int64(val))
	case int32:
		return scanInt(src, int64(val))
	case int64:
		return scanInt(src, int64(val))
	case string:
		parsedUint, err := strconv.ParseUint(val, 10, 0)
		if err != nil {
			return fmt.Errorf("null.Uint8: failed to scan type %T (%v): %v", src, src, err)
		}
		i.Uint8 = uint8(parsedUint)
		i.Valid = true
		return nil
	case float64:
		// Use a string intermediate so we can generate an error on any loss of
		// precision.
		s := strconv.FormatFloat(val, 'g', -1, 64)
		parsedUint, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return fmt.Errorf("null.Uint8: failed to convert driver.Value type %T (%v): %v", src, s, err)
		}
		i.Uint8 = uint8(parsedUint)
		i.Valid = true
		return nil
	case float32:
		// Use a string intermediate so we can generate an error on any loss of
		// precision.
		s := strconv.FormatFloat(float64(val), 'g', -1, 32)
		parsedUint, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return fmt.Errorf("null.Uint8: failed to convert driver.Value type %T (%v): %v", src, s, err)
		}
		i.Uint8 = uint8(parsedUint)
		i.Valid = true
		return nil
	case nil:
		i.Uint8 = 0
		i.Valid = false
		return nil
	default:
		return fmt.Errorf("null.Uint8: cannot scan type %T (%v)", src, src)
	}
}

// MarshalJSON implements the encoding/json Marshaler interface. It will encode
// i into its JSON representation if valid, or 'null' otherwise.
func (i Uint8) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint8), 10)), nil
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It will
// decode a given []byte into i, so long as the provided []byte is a valid JSON
// representation of an int. The 'null' keyword will decode into a null Uint8.
//
// If the decode fails, the value of i will be unchanged.
func (i *Uint8) UnmarshalJSON(data []byte) error {
	if i == nil {
		return fmt.Errorf("null.Uint8: UnmarshalJSON called on nil pointer")
	}
	var j interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	switch val := j.(type) {
	case float64:
		// Perform a second unmarshal, this time into a uint8. This give the
		// JSON parse a chance to meaningfully fail (eg. if the conversion from
		// float to integer will result in a loss of precision).
		var tmp uint8
		err := json.Unmarshal(data, &tmp)
		if err != nil {
			return err
		}
		i.Uint8 = tmp
		i.Valid = true
		return nil
	case nil:
		i.Uint8 = 0
		i.Valid = false
		return nil
	default:
		return fmt.Errorf("null.Uint8: cannot unmarshal JSON of type %T (%v)",
			val, data)
	}
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode i into its interface{} representation for use in a
// map[string]interface{} if valid, or return nil otherwise.
func (i Uint8) MarshalMapValue() (interface{}, error) {
	if i.Valid {
		return i.Uint8, nil
	}
	return nil, nil
}
