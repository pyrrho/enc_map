package null

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// ByteSlice is a nullable wrapper around the []byte type. It implements all of
// the pyrrho/encoding/types interfaces detailed in the package comments. This
// type makes a distinction between nil and valid-but-empty []bytes. It is valid
// to initialzed a ByteSlice with an empty []byte{}, but it is not valid to
// initialize a ByteSlice with a []byte(nil). This distinction holds true for
// the Set, Scan, an Unmarshal function groups.
//
// To maintain consistency with the encoding/json package -- and to ensure we
// never attempt to marshal non-ASCII characters -- this type will emit base64
// encoded strings from MarshalJSON, Value, and MarshalMapValue (when non-null),
// and expect to receive base64 encoded strings in UnmarshalJSON and Scan.
type ByteSlice struct {
	ByteSlice []byte
	Valid     bool
}

// Constructors

// NullByteSlice constructs and returns a new null ByteSlice.
func NullByteSlice() ByteSlice {
	return ByteSlice{
		ByteSlice: nil,
		Valid:     false,
	}
}

// NewByteSlice constructs and returns a new ByteSlice based on the given []byte
// b. If b is of zero length the new ByteSlice will be null. Otherwise a new,
// valid ByteSlice will be initialized with a copy of b.
func NewByteSlice(b []byte) ByteSlice {
	if b == nil {
		return NullByteSlice()
	}
	return ByteSlice{
		ByteSlice: append([]byte{}, b...),
		Valid:     true,
	}
}

// NewByteSliceStr constructs and returns a new ByteSlice object based on the
// given string s. If s is the empty string, the new ByteSlice will be null.
// Otherwise s will be cast to []byte and used to initialize the new, valid
// ByteSlice.
func NewByteSliceStr(s string) ByteSlice {
	return ByteSlice{
		ByteSlice: []byte(s),
		Valid:     true,
	}
}

// NewByteSliceFromBase64 constructs and returns a new ByteSlice object based on
// the given base64 encoded []byte b.
func NewByteSliceFromBase64(b []byte) (ByteSlice, error) {
	if b == nil {
		return NullByteSlice(), nil
	}

	tmp := make([]byte, base64.StdEncoding.DecodedLen(len(b)))
	n, err := base64.StdEncoding.Decode(tmp, b)
	if err != nil {
		return ByteSlice{}, err
	}
	return ByteSlice{
		ByteSlice: tmp[:n],
		Valid:     true,
	}, nil

}

// NewByteSliceFromBase64Str constructs and returns a new ByteSlice object based on
// the given base64 encoded []byte b.
func NewByteSliceFromBase64Str(s string) (ByteSlice, error) {
	tmp := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	n, err := base64.StdEncoding.Decode(tmp, []byte(s))
	if err != nil {
		return ByteSlice{}, err
	}
	return ByteSlice{
		ByteSlice: tmp[:n],
		Valid:     true,
	}, nil

}

// Getters and Setters

// ValueOrZero returns the value of b if it is valid; otherwise,it returns an
// empty []byte.
func (b ByteSlice) ValueOrZero() []byte {
	if !b.Valid {
		return []byte{}
	}
	return b.ByteSlice
}

// Set copies the given []byte v into b. If v is of length zero, b will be
// nulled.
func (b *ByteSlice) Set(v []byte) {
	if v == nil {
		b.ByteSlice = nil
		b.Valid = false
		return
	}
	// Make sure Set initializes null ByteSlices.
	// If b.ByteSlice is nil and len(v) == 0, b.ByteSlice would remain nil. We
	// don't want to that.
	if !b.Valid {
		b.ByteSlice = []byte{}
		b.Valid = true
	}
	b.ByteSlice = append(b.ByteSlice[0:0], v...)
	b.Valid = true
}

// SetStr copies the given string v into b. If v is of length zero, b will be
// nulled.
func (b *ByteSlice) SetStr(v string) {
	b.ByteSlice = append(b.ByteSlice[0:0], v...)
	b.Valid = true
}

// Null marks b as null with no meaningful value.
func (b *ByteSlice) Null() {
	b.ByteSlice = nil
	b.Valid = false
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if b is null.
func (b ByteSlice) IsNil() bool {
	return !b.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if b is null or if its value is false.
func (b ByteSlice) IsZero() bool {
	return !b.Valid || len(b.ByteSlice) == 0
}

// Value implements the database/sql/driver Valuer interface. It will base64
// encode valid values prior to returning them.
func (b ByteSlice) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	// TODO: For all base64.StdEncoding.Encode calls, consider performing an
	// optimization similar to the one implemented in encoding/json/encode.go's
	// encodeByteSlice function.
	enc := make([]byte, base64.StdEncoding.EncodedLen(len(b.ByteSlice)))
	base64.StdEncoding.Encode(enc, b.ByteSlice)
	return enc, nil
}

// Scan implements the database/sql Scanner interface. It will receive a value
// from an SQL database and assign it to b, so long as the provided data can be
// is a []byte, a string, or nil. Valid data is expected to be base64 encoded.
func (b *ByteSlice) Scan(src interface{}) error {
	if b == nil {
		return fmt.Errorf("null.ByteSlice: Scan called on nil pointer")
	}
	switch val := src.(type) {
	case nil:
		b.ByteSlice = nil
		b.Valid = false
		return nil
	case []byte:
		if val == nil {
			b.ByteSlice = nil
			b.Valid = false
			return nil
		}
		tmp := make([]byte, base64.StdEncoding.DecodedLen(len(val)))
		n, err := base64.StdEncoding.Decode(tmp, val)
		if err != nil {
			return err
		}
		b.ByteSlice = tmp[:n]
		b.Valid = true
		return nil
	case string:
		tmp, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return err
		}
		b.ByteSlice = tmp
		b.Valid = true
		return nil
	default:
		return fmt.Errorf("null.ByteSlice: cannot scan type %T (%v)",
			val, src)
	}
}

// MarshalJSON implements the encoding/json Marshaler interface. It will encode
// b into its base64 representation if valid, or 'null' otherwise.
func (b ByteSlice) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	// Because we're passing a []byte into json.Marshal, the json package will
	// handle any base64 decoding that needs to happen.
	return json.Marshal(b.ByteSlice)
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It will
// decode a given []byte into b, so long as the provided []byte is
// a valid base64 encoded string or a null.
//
// An empty string will result in a valid-but-empty ByteSlice. The keyword
// 'null' will result in a null ByteSlice. The string '"null"' is considered
// to be a string -- not a keyword -- and will result in base64 decoded garbage.
//
// If the decode fails, the value of b will be unchanged.
func (b *ByteSlice) UnmarshalJSON(data []byte) error {
	if b == nil {
		return fmt.Errorf("null.ByteSlice: UnmarshalJSON called on nil pointer")
	}
	var j interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	switch val := j.(type) {
	case nil:
		b.ByteSlice = nil
		b.Valid = false
		return nil
	case string:
		if len(val) == 0 {
			// We were passed something similar to a string that is an empty
			// string (`""`). This should result in an empty-but-valid slice.
			b.ByteSlice = []byte{}
			b.Valid = true
			return nil
		}
		// Call json.Unmarshal again, this time with a []byte as the dest. This
		// lets encoding/json package take care of the base64 decoding.
		var tmp []byte
		err := json.Unmarshal(data, &tmp)
		if err != nil {
			return err
		}
		b.ByteSlice = tmp
		b.Valid = true
		return nil
	default:
		return fmt.Errorf("null.ByteSlice: cannot unmarshal JSON of type %T (%v)",
			val, data)
	}
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode b into its base64 encoded interface{} representation for use in a
// map[string]interface{} if valid, or return nil otherwise.
func (b ByteSlice) MarshalMapValue() (interface{}, error) {
	if !b.Valid {
		return nil, nil
	}
	// TODO: For all base64.StdEncoding.Encode calls, consider performing an
	// optimization similar to the one implemented in encoding/json/encode.go's
	// encodeByteSlice function.
	enc := make([]byte, base64.StdEncoding.EncodedLen(len(b.ByteSlice)))
	base64.StdEncoding.Encode(enc, b.ByteSlice)
	return enc, nil
}
