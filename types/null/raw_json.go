package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pyrrho/encoding/types"
)

// RawJSON is a wrapper around types.RawJSON that makes the type null-aware, in
// terms of both the JSON 'null' keyword, and SQL NULL values. It implements all
// of the pyrrho/encoding/types interfaces detailed in the package comments.
type RawJSON struct {
	JSON  types.RawJSON
	Valid bool
}

// Constructors

// NullJSON constructs and returns a new null RawJSON object.
func NullJSON() RawJSON {
	return RawJSON{
		JSON:  nil,
		Valid: false,
	}
}

// NewJSON constructs and returns a new RawJSON object based on the given
// types.RawJSON j. If j is of zero-length the new RawJSON will be null.
// Otherwise a new, valid RawJSON will be initialized with a copy of j.
func NewJSON(j types.RawJSON) RawJSON {
	if len(j) == 0 {
		return NullJSON()
	}
	return RawJSON{
		JSON:  types.NewJSON(j),
		Valid: true,
	}
}

// NewJSONStr constructs and returns a new RawJSON object based on the given
// string s. If s is the empty string, the new RawJSON will be null. Otherwise s
// will be cast to types.RawJSON and used to initialize the new valid RawJSON.
func NewJSONStr(s string) RawJSON {
	if len(s) == 0 {
		return NullJSON()
	}
	return RawJSON{
		JSON:  types.NewJSONStr(s),
		Valid: true,
	}
}

// Getters and Setters

// ValueOrZero will return the value of j if it is valid, or a newly constructed
// zero-value types.RawJSON otherwise.
func (j RawJSON) ValueOrZero() types.RawJSON {
	if !j.Valid {
		return types.RawJSON{}
	}
	return j.JSON
}

// Set copies the given types.RawJSON value into j. If the given value is of
// length 0, j will be nulled.
func (j *RawJSON) Set(v types.RawJSON) {
	if len(v) == 0 {
		j.JSON = nil // Let the garbage collector have this types.RawJSON.
		j.Valid = false
		return
	}
	j.JSON.Set(v)
	j.Valid = true
}

// SetStr will copy the contents of v into j. If the given value is an empty
// string, j will be nulled.
func (j *RawJSON) SetStr(v string) {
	if len(v) == 0 {
		j.JSON = nil // Let the garbage collector have this types.RawJSON.
		j.Valid = false
		return
	}
	j.JSON.SetStr(v)
	j.Valid = true

}

// Null will set j to null; j.Valid will be false, and j.JSON will contain no
// meaningful value.
func (j *RawJSON) Null() {
	j.JSON = nil // Let the garbage collector have this types.RawJSON.
	j.Valid = false
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if j is null.
func (j RawJSON) IsNil() bool {
	return !j.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if j is null or if the contained JSON is a zero value.
func (j RawJSON) IsZero() bool {
	if !j.Valid {
		return true
	}
	return j.JSON.IsZero()
}

// Value implements the database/sql/driver Valuer interface. It will return the
// value of j as a driver.Value. If j is valid, this function will first
// validate the contained JSON returning either any encouted parsing errors, or
// a []byte as a driver.Value. If j is null, nil will be returned, and no
// validation will occur.
func (j RawJSON) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return j.JSON.Value()
}

// Scan implements the database/sql Scanner interface. It expects to receive a
// valid JSON value as a string or a []byte, or NULL as a nil from an SQL
// database. A zero-length string or []byte, or a nil will be considered NULL,
// and j will be nulled, otherwise the the value will be assigned to j. Scan
// will not validate the incoming JSON.
func (j *RawJSON) Scan(src interface{}) error {
	if j == nil {
		return fmt.Errorf("null.RawJSON: Scan called on nil pointer")
	}
	switch x := src.(type) {
	case nil:
		j.JSON = nil
		j.Valid = false
		return nil
	case []byte:
		if len(x) == 0 {
			j.JSON = nil
			j.Valid = false
			return nil
		}
		j.JSON.Set(x)
		j.Valid = true
		return nil
	case string:
		if len(x) == 0 {
			j.JSON = nil
			j.Valid = false
			return nil
		}
		j.JSON.SetStr(x)
		j.Valid = true
		return nil
	default:
		return fmt.Errorf("null.RawJSON: cannot scan type %T (%v)", src, src)
	}
}

// MarshalJSON implements the encoding/json Marshaler interface. It will return
// the value of j as a JSON-encoded []byte. If j is valid, this function will
// first validate the contained JSON returning either any encouted parsing
// errors, or a []byte. If j is null, "null" will be returned, and no validation
// will occur.
func (j RawJSON) MarshalJSON() ([]byte, error) {
	if !j.Valid {
		return []byte("null"), nil
	}
	return j.JSON.MarshalJSON()
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It expects
// to receive a valid JSON value, and will assign that value to this RawJSON. If
// the incoming JSON is the 'null' keyword, j will be nulled. UnmarshalJSON will
// validate the incoming JSON as part of the "Is this JSON null?" check.
func (j *RawJSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return fmt.Errorf("null.RawJSON: UnmarshalJSON called on nil pointer")
	}
	var k interface{}
	if err := json.Unmarshal(data, &k); err != nil {
		return err
	}
	if k == nil {
		j.JSON = nil
		j.Valid = false
		return nil
	}
	j.JSON.Set(data)
	j.Valid = true
	return nil
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode j into its interface{} representation for use in a
// map[string]interface{} by passing it through JSON.Unmarshal if valid, or the
// 'null' keyword otherwise.
func (j RawJSON) MarshalMapValue() (interface{}, error) {
	if !j.Valid {
		return []byte("null"), nil
	}
	return j.JSON.MarshalMapValue()
}
