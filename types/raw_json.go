package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

// RawJSON is an alternative to the json.RawMessage type. RawJSON implements all
// of the pyrrho/encoding/types interfaces detailed in the package comments.
//
// This implementation should not be considered safe to use with NULL-able SQL
// columns; for that application please use the pyrrho/encoding/types/null
// package, specifically the null.RawJSON type.
type RawJSON []byte

// NewJSON will return a new RawJSON object that has been initialized with a
// copy of the contents of b.
func NewJSON(b []byte) RawJSON {
	ret := make(RawJSON, len(b))
	copy(ret, b)
	return ret
}

// NewJSONStr will return a new RawJSON object that has been initialized with
// the given string.
func NewJSONStr(s string) RawJSON {
	return RawJSON(s)
}

// Set will copy the contents of v into this RawJSON.
func (j *RawJSON) Set(v []byte) {
	// NB. This will re-use the array allocated for *j, if possible, filling it
	// with / sizing it to v.
	*j = append((*j)[0:0], v...)
}

// SetStr will copy the contents of v into j.
func (j *RawJSON) SetStr(v string) {
	// NB. This will re-use the array allocated for *j, if possible, filling it
	// with / sizing it to v.
	*j = append((*j)[0:0], v...)
}

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if j has a length of zero.
func (j RawJSON) IsNil() bool {
	return len(j) == 0
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if j has a length of zero, or if the contained JSON is a zero value. If the
// parsing the contined JSON results in an error, IsZero will return false.
func (j RawJSON) IsZero() bool {
	if len(j) == 0 {
		return true
	}
	b, _ := j.ValueIsZero()
	return b
}

// ValueIsZero will return true if the contained JSON is a zero value. If the
// contained JSON is invalid, ValueIsZero will return false and the resulting
// JSON parsing error.
func (j RawJSON) ValueIsZero() (bool, error) {
	if len(j) == 0 {
		return false, fmt.Errorf("types.RawJSON: invalid JSON, an empty string cannot be unmarshaled")
	}

	var k interface{}
	err := json.Unmarshal(j, &k)
	if err != nil {
		return false, err
	}

	v := reflect.ValueOf(k)
	switch v.Kind() {
	// json objects, arrays, strings ...
	case reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0, nil
	// ... numbers
	case reflect.Float64:
		return v.Float() == 0, nil
	// ... booleans
	case reflect.Bool:
		return v.Bool() == false, nil
	// ... and null
	case reflect.Invalid:
		return true, nil
	default:
		return false, fmt.Errorf("types.RawJSON: unexpected kind returned by json.Unmarshal; %s", v.Kind())
	}
}

// Value implements the database/sql/driver Valuer interface. It will return the
// value of j as a driver.Value; specifically a []byte. Before returning the
// value, this function will validate the contained JSON and return any parsing
// errors encountered.
func (j RawJSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		// An empty byte slice is not valid JSON. Return an error that's more
		// descriptive than the encoding/json message.
		return nil, fmt.Errorf("types.RawJSON: invalid JSON, an empty string cannot be unmarshaled")
	}
	// Unmarshal to look for errors, unmarshal into a RawJSON to keep
	// unnecessary allocations to a minimum.
	// NB. At the time of writing, what I _actually_ want to call is
	// json.checkValid, but that function is not exported. So.
	var tmp RawJSON
	if err := json.Unmarshal(j, &tmp); err != nil {
		return nil, err
	}
	return []byte(j), nil
}

// Scan implements the database/sql Scanner interface. It expects to receive a
// valid JSON string or []byte from an SQL database, and will assign that value
// to j. Scan will not validatet the incoming JSON.
func (j *RawJSON) Scan(src interface{}) error {
	if j == nil {
		return fmt.Errorf("types.RawJSON: Scan called on nil pointer")
	}
	switch x := src.(type) {
	case []byte:
		j.Set(x)
		return nil
	case string:
		j.SetStr(x)
		return nil
	default:
		return fmt.Errorf("types.RawJSON: cannot scan type %T (%v)", src, src)
	}
}

// MarshalJSON implements the encoding/json Marshaler interface. Before
// returning the encoded value, this function will validate the contained JSON
// and return any parsing errors encountered.
func (j RawJSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		// An empty byte slice is not valid JSON. Return an error that's more
		// descriptive than the encoding/json message.
		return nil, fmt.Errorf("types.RawJSON: invalid JSON, an empty string cannot be unmarshaled")
	}
	// Unmarshal to look for errors, unmarshal into a RawJSON to keep
	// unnecessary allocations to a minimum.
	// NB. At the time of writing, what I _actually_ want to call is
	// json.checkValid, but that function is not exported. So.
	var ret RawJSON
	if err := json.Unmarshal(j, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It expects
// to receive a valid JSON value, and will assign that value to j. UnmarshalJSON
// will not validate the incoming JSON.
func (j *RawJSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return fmt.Errorf("types.RawJSON: UnmarshalJSON called on nil pointer")
	}
	j.Set(data)
	return nil
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode j into its interface{} representation for use in a
// map[string]interface{} by passing it through json.Unmarshal.
func (j RawJSON) MarshalMapValue() (interface{}, error) {
	if len(j) == 0 {
		// An empty byte slice is not valid JSON. Return an error that's more
		// descriptive than the encoding/json message.
		return nil, fmt.Errorf("types.RawJSON: invalid JSON, an empty string cannot be unmarshaled")
	}
	var iface interface{}
	err := json.Unmarshal(j, &iface)
	if err != nil {
		return nil, err
	}
	return iface, nil
}
