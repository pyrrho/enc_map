package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// NullTime is a nullable wrapper around the time.Time type implementing all of
// the encoding/type interfaces. The zero time instant is considered non-null.
// if null.
type NullTime struct {
	Time  time.Time
	Valid bool
}

// Constructors

// Time creates a new NullTime based on the type and value of the given
// interface. This function intentionally sacrafices compile-time safety for
// developer convenience.
//
// If the interface is nil or a nil *time.Time, the new NullTime will be null.
//
// If the interface is a time.Time or a non-nil *time.Time, the new NullTime
// will be valid and will be initialized with the (possibly dereferenced) value
// of the interface.
//
// If the interface is any other type, this function will panic.
func Time(i interface{}) NullTime {
	switch v := i.(type) {
	case time.Time:
		return TimeFrom(v)
	case *time.Time:
		return TimeFromPtr(v)
	case nil:
		return NullTime{}
	}
	panic(fmt.Errorf(
		"null.Time: the given argument (%#v of type %T) was not of type "+
			"time.Time, *time.Time, or nil", i, i))
}

// TimeFrom creates a valid Time from t.
func TimeFrom(t time.Time) NullTime {
	return NullTime{
		Time:  t,
		Valid: true,
	}
}

// TimeFromPtr creates a valid Time from *t.
func TimeFromPtr(t *time.Time) NullTime {
	if t == nil {
		return NullTime{}
	}
	return TimeFrom(*t)
}

// Getters and Setters

// ValueOrZero returns the value of this NullTime if it is valid; otherwise
// it returns the zero value for a time.Time.
func (t NullTime) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// Ptr returns a pointer to this NullTime's value if it is valid; otherwise
// returns a nil pointer. The captured pointer will be able to modify the value
// of this NullTime.
func (t *NullTime) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// Set modifies the value stored in this NullTime, and guarantees it is
// valid.
func (t *NullTime) Set(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Null marks this NullTime as null with no meaningful value.
func (t *NullTime) Null() {
	t.Valid = false
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if this NullTime is null.
func (t NullTime) IsNil() bool {
	return !t.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if this NullTime is null or if its value is the zero time instant.
func (t NullTime) IsZero() bool {
	return !t.Valid || t.Time == time.Time{}
}

// Value implements the database/sql/driver Valuer interface. As time.Time and
// nil are both valid types to be stored in a driver.Value, it will return this
// NullTime's value if valid, or nil otherwise.
func (t NullTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// Scan implements the database/sql Scanner interface. It will receive a value
// from an SQL database and assign it to this NullTime, so long as the provided
// data is of type nil or time.Time. All other types will result in an error.
func (t *NullTime) Scan(value interface{}) error {
	switch x := value.(type) {
	case time.Time:
		t.Time = x
		t.Valid = true
		return nil
	case nil:
		t.Time = time.Time{}
		t.Valid = false
		return nil
	default:
		return fmt.Errorf("null: cannot scan type %T into null.Time: %v",
			value, value)
	}
}

// MarshalText implements the encoding TextMarshaler interface. It return the
// textual representation of this NullTime's value (by calling into the native
// time.Time's MarshalText()) if valid, or an empty string otherwise.
func (t NullTime) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte(""), nil
	}
	return t.Time.MarshalText()
}

// UnmarshalText implements the encoding TextUnmarshaler interface. It will
// decode a given []byte into this NullTime, so long as the provided string
// is a valid textual representation of a time.Time or a null.
//
// If the decode fails, the value of this NullTime will be unchanged.
func (t *NullTime) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Time = time.Time{}
		t.Valid = false
		return nil
	}
	// TODO: If time.Time.UnmarshalText doesn't change the value of the receiver
	// we could skip the temporary object by calling .UnmarshalText on t.Time.
	tmp := time.Time{}
	err := tmp.UnmarshalText(text)
	if err != nil {
		return err
	}
	t.Time = tmp
	t.Valid = true
	return nil
}

// MarshalJSON implements the encoding/json Marshaler interface. It will encode
// this NullTime into its JSON RFC 3339 string representation if valid, or
// 'null' otherwise.
func (t NullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalJSON()
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It will
// decode a given []byte into this NullTime so long as the provided []byte
// is a valid JSON representation of an RFC 3339 string or a null.
//
// Empty strings and 'null' will both decode into a null NullTime. JSON objects
// objects in the form of '{"Time":<RFC 3339 string>,"Valid":<bool>}' will
// decode directly into this NullTime.
//
// If the decode fails, the value of this NullTime will be unchanged.
func (t *NullTime) UnmarshalJSON(data []byte) error {
	var j interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	switch val := j.(type) {
	case string:
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
	case map[string]interface{}:
		ti, tiOK := val["Time"].(string)
		valid, validOK := val["Valid"].(bool)
		if !tiOK || !validOK {
			return fmt.Errorf(
				`null: unmarshalling object into Go value of type null.Time `+
					`requires key "Time" to be of type string and key "Valid" `+
					`to be of type bool; found %T and %T, respectively`,
				val["Time"], val["Valid"],
			)
		}
		// TODO: If time.Time.UnmarshalJSON doesn't change the value of the
		// receiver we could skip the temporary object by calling .UnmarshalJSON
		// on t.Time.
		tmp, err := time.Parse(time.RFC3339, ti)
		if err != nil {
			return err
		}
		t.Time = tmp
		t.Valid = valid
		return nil
	case nil:
		t.Time = time.Time{}
		t.Valid = false
		return nil
	default:
		return fmt.Errorf(
			"null: cannot unmarshal %T (%#v) into Go value of type "+
				"null.NullTime",
			j, j,
		)
	}
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode this NullTime into an interface{} representation for use in a
// map[Time]interface{} if valid, or return nil otherwise.
func (t NullTime) MarshalMapValue() (interface{}, error) {
	if t.Valid {
		return t.Time, nil
	}
	return nil, nil
}
