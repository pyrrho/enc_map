package null

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// NullByteSlice is a nullable wrapper around the []byte type implementing all
// of the encoding/type interfaces.
//
// To maintain consistency with the encoding/json package -- and to ensure we
// never attempt to marshal non-ASCII characters -- this type will emit base64
// encoded strings from MarshalJSON,and Value (when non-null), and expect to
// receive base64 encoded strings in UnmarshalJSON and Scan. MarshalMapValue
// will produce an interface wrapping a []byte.
//
// We assume that if Valid is true, the contained ByteSlice will be non-nil. To
// manually construct a NullByteSlice where Valid is true and ByteSlice is a nil
// slice is considered user error, and will likely result in a panic.
type NullByteSlice struct {
	ByteSlice []byte
	Valid     bool
}

// Constructors

// ByteSlice creates a new NullByteSlice based on the type and value of the
// given interface. This function intentionally sacrafices compile-time safety
// for developer convenience.
//
// If the interface is nil, a nil []byte, a nil *[]byte, or a *[]byte that
// dereferences to a nil []byte the new NullByteSlice will be null.
//
// If the interface is a non-nil []byte, or a non-nil *[]byte that dereferences
// to a non-nil []byte, the new NullBytesSlice will be valid, and will be
// initialized with the (possibly dereferenced) value of the interface.
//
// If the interface is any other type this function will panic.
func ByteSlice(i interface{}) NullByteSlice {
	switch v := i.(type) {
	case []byte:
		return ByteSliceFrom(v)
	case *[]byte:
		return ByteSliceFromPtr(v)
	case nil:
		return NullByteSlice{}
	}
	panic(fmt.Errorf(
		"null.NullByteSlice: invalid constructor argument; %#v of type %T "+
			"is not of type []byte, *[]byte, or nil", i, i))
}

// ByteSliceFrom creates a valid Byteslice from b.
func ByteSliceFrom(b []byte) NullByteSlice {
	if b == nil {
		return NullByteSlice{}
	}
	return NullByteSlice{
		ByteSlice: b,
		Valid:     true,
	}
}

// ByteSliceFromPtr creates a valid ByteSlice from *b.
func ByteSliceFromPtr(b *[]byte) NullByteSlice {
	if b == nil {
		return NullByteSlice{}
	}
	return ByteSliceFrom(*b)
}

// Getters and Setters

// ValueOrZero returns the value of this NullByteSlice if it is valid;
// otherwise, it returns an empty but initialized []byte.
func (b NullByteSlice) ValueOrZero() []byte {
	if !b.Valid {
		tmp := make([]byte, 0)
		return tmp
	}
	return b.ByteSlice
}

// Ptr returns a pointer to this NullByteSlice's value if it is valid; otherwise
// returns a nil pointer. The captured pointer will be able to modify the value
// of this NullByteSlice.
func (b *NullByteSlice) Ptr() *[]byte {
	if !b.Valid {
		return nil
	}
	return &b.ByteSlice
}

// Set modifies the value stored in this NullByteSlice. If v is a nil []byte,
// this NullByteSlice will be marked null with no meaningful value, otherwise
// this NullByteSlice will wrap v.
func (b *NullByteSlice) Set(v []byte) {
	if v == nil {
		b.ByteSlice = nil // Let the garbage collector have this slice.
		b.Valid = false
		return
	}
	b.ByteSlice = v
	b.Valid = true
}

// Null marks this ByteSlice as null with no meaningful value.
func (b *NullByteSlice) Null() {
	b.ByteSlice = nil // Let the garbage collector have this slice.
	b.Valid = false
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if this NullByteSlice is null.
func (b NullByteSlice) IsNil() bool {
	return !b.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if this NullByteSlice is null or if its value is false.
func (b NullByteSlice) IsZero() bool {
	return !b.Valid || len(b.ByteSlice) == 0
}

// Value implements the database/sql/driver Valuer interface. It will base64
// encode valid values prior to returning them.
func (b NullByteSlice) Value() (driver.Value, error) {
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
// from an SQL database and assign it to this NullByteSlice, so long as the
// provided data can be coerced into a []byte or a nil. Valid data is expected
// to be base64 encoded.
func (b *NullByteSlice) Scan(src interface{}) error {
	if b == nil {
		return fmt.Errorf("null.NullByteSlice: Scan called on nil pointer")
	}
	switch val := src.(type) {
	case []byte:
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
	case nil:
		b.ByteSlice = nil
		b.Valid = false
		return nil
	default:
		return fmt.Errorf("null.NullByteSlice: cannot scan type %T (%v)",
			val, src)
	}
}

// MarshalJSON implements the encoding/json Marshaler interface. It will encode
// this NullByteSlice into its base64 representation if valid, or 'null'
// otherwise.
func (b NullByteSlice) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	// Because we're passing a []byte into json.Marshal, the json package will
	// handle any base64-decoding that needs to happen.
	return json.Marshal(b.ByteSlice)
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It will
// decode a given []byte into this NullByteSlice, so long as the provided []byte
// is a valid base64 encoded string or a null.
//
// An empty string will result in a valid-but-empty NullByteSlice. The keyword
// 'null' will result in a null NullByteSlice. The string '"null"' is considered
// to be a string -- not a keyword -- and will result in base64-decoded garbage.
//
// If the decode fails, the value of this NullByteSlice will be unchanged.
func (b *NullByteSlice) UnmarshalJSON(data []byte) error {
	if b == nil {
		return fmt.Errorf("null.NullByteSlice: UnmarshalJSON called on nil pointer")
	}
	var j interface{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	switch val := j.(type) {
	case string:
		if len(val) == 0 {
			// We were passed something similar to a string that is an empty
			// string (`""`). This should result in an empty-but-valid slice.
			b.ByteSlice = []byte{}
			b.Valid = true
			return nil
		}
		// Call json.Unmarshal again, this time with a []byte as the dest. This
		// lets encoding/json package take care of the base64-decoding.
		var tmp []byte
		err := json.Unmarshal(data, &tmp)
		if err != nil {
			return err
		}
		b.ByteSlice = tmp
		b.Valid = true
		return nil
	case nil:
		b.ByteSlice = nil
		b.Valid = false
		return nil
	default:
		return fmt.Errorf("null.NullByteSlice: cannot unmarshal JSON of type %T (%v)",
			val, data)
	}
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode this NullByteSlice into its interface{} representation for use in
// a map[string]interface{} if valid, or return nil otherwise.
func (b NullByteSlice) MarshalMapValue() (interface{}, error) {
	if b.Valid {
		return b.ByteSlice, nil
	}
	return nil, nil
}
