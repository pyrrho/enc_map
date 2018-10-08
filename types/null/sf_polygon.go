package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pyrrho/encoding/types"
)

// SFPolygon is a wrapper around types.SFPolygon that makes the type null-aware, in
// terms of both the JSON 'null' keyword, and SQL NULL values. It implements all
// of the pyrrho/encoding/types interfaces detailed in the package comments.
type SFPolygon struct {
	Polygon types.SFPolygon
	Valid   bool
}

// Constructors

// NullSFPolygon constructs and returns a new null SFPolygon object.
func NullSFPolygon() SFPolygon {
	return SFPolygon{
		Polygon: types.SFPolygon{},
		Valid:   false,
	}
}

// NewSFPolygon constructs and returns a new SFPolygon object based on the given
// types.SFPolygon p. If p is of zero-value, the new SFPolygon will be null.
// Otherwise a new, valid SFPolygon will be initialized with a copy of p.
func NewSFPolygon(p types.SFPolygon) SFPolygon {
	if p.IsZero() {
		return NullSFPolygon()
	}
	return SFPolygon{
		Polygon: types.NewSFPolygon(p.Polygon),
		Valid:   true,
	}
}

// NewSFPolygonXY constructs and returns a new SFPolygon object based on the given
// external and (optionally) internal shapes.
func NewSFPolygonXY(external [][2]float64, internals ...[][2]float64) SFPolygon {
	return SFPolygon{
		Polygon: types.NewSFPolygonXY(external, internals...),
		Valid:   true,
	}
}

// NewSFPolygonXY constructs and returns a new SFPolygon object based on the given
// external and (optionally) internal shapes.
func NewSFPolygonXYZ(external [][3]float64, internals ...[][3]float64) SFPolygon {
	return SFPolygon{
		Polygon: types.NewSFPolygonXYZ(external, internals...),
		Valid:   true,
	}
}

// Getters and Setters

// ValueOrZero will return the value of p if it is valid, or a newly constructed
// zero-value types.SFPolygon otherwise.
func (p SFPolygon) ValueOrZero() types.SFPolygon {
	if !p.Valid {
		return types.SFPolygon{}
	}
	return p.Polygon
}

// Set copies the given types.SFPolygon value into p. If the given value is nil,
// p will be nulled.
func (p *SFPolygon) Set(v types.SFPolygon) {
	if v.IsNil() {
		p.Polygon = types.SFPolygon{} // Let the garbage collector have this types.SFPolygon.
		p.Valid = false
		return
	}
	p.Polygon = v
	p.Valid = true
}

// Null will set p to null; p.Valid will be false, and p.Polygon will contain no
// meaningful value.
func (p *SFPolygon) Null() {
	p.Polygon = types.SFPolygon{} // Let the garbage collector have this types.SFPolygon.
	p.Valid = false
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if j is null.
func (p SFPolygon) IsNil() bool {
	return !p.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if p is null or if the contained SFPolygon is a zero value.
func (p SFPolygon) IsZero() bool {
	if !p.Valid {
		return true
	}
	return p.Polygon.IsZero()
}

// Value implements the database/sql/driver Valuer interface. It will return the
// value of p as a driver.Value. If p is valid, this function will first
// validate the contained SFPolygon returning either any encouted parsing errors,
// or a []byte as a driver.Value. If p is null, nil will be returned, and no
// validation will occur.
func (p SFPolygon) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}
	return p.Polygon.Value()
}

// Scan implements the database/sql Scanner interface. It expects to receive a
// valid WKB encoded byte describing a Polygon as a []byte, or NULL as a nil
// from an SQL database. A zero-length []byte or a nil will be considered NULL,
// and p will be nulled. Otherwise, the value will be passed to types.SFPolygon
// to be scanned and parsed as a WKB Polygon.
func (p *SFPolygon) Scan(src interface{}) error {
	if p == nil {
		return fmt.Errorf("null.SFPolygon: Scan called on nil pointer")
	}
	switch x := src.(type) {
	case nil:
		p.Polygon = types.SFPolygon{}
		p.Valid = false
		return nil
	case []byte:
		if len(x) == 0 {
			p.Polygon = types.SFPolygon{}
			p.Valid = false
			return nil
		}
		err := p.Polygon.Scan(x)
		if err != nil {
			p.Polygon = types.SFPolygon{}
			p.Valid = false
			return err
		}
		p.Valid = true
		return nil
	default:
		return fmt.Errorf("null.SFPolygon: cannot scan type %T (%v)", src, src)
	}
}

// MarshalJSON implements the encoding/json Marshaler interface. It will return
// the value of p as a JSON-encoded []byte. If p is valid, this function will
// first validate the contained JSON returning either any encouted parsing
// errors, or a []byte. If p is null, "null" will be returned, and no validation
// will occur.
func (p SFPolygon) MarshalJSON() ([]byte, error) {
	if !p.Valid {
		return []byte("null"), nil
	}
	return p.Polygon.MarshalJSON()
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It expects
// to receive a valid JSON value, and will assign that value to this SFPolygon. If
// the incoming JSON is the 'null' keyword, p will be nulled. UnmarshalJSON will
// validate the incoming JSON as part of the "Is this JSON null?" check.
func (p *SFPolygon) UnmarshalJSON(data []byte) error {
	if p == nil {
		return fmt.Errorf("null.SFPolygon: UnmarshalJSON called on nil pointer")
	}
	var k interface{}
	if err := json.Unmarshal(data, &k); err != nil {
		return err
	}
	if k == nil {
		p.Polygon = types.SFPolygon{}
		p.Valid = false
		return nil
	}
	p.Polygon.UnmarshalJSON(data)
	p.Valid = true
	return nil
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode p into its interface{} representation for use in a
// map[string]interface{} by passing it through JSON.Unmarshal if valid, or the
// 'null' keyword otherwise.
func (p SFPolygon) MarshalMapValue() (interface{}, error) {
	if !p.Valid {
		return []byte("null"), nil
	}
	return p.Polygon.MarshalMapValue()
}
