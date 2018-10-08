package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pyrrho/encoding/types"
	"github.com/twpayne/go-geom"
)

// SFPoint is a wrapper around types.SFPoint that makes the type null-aware, in
// terms of both the JSON 'null' keyword, and SQL NULL values. It implements all
// of the pyrrho/encoding/types interfaces detailed in the package comments.
type SFPoint struct {
	Point types.SFPoint
	Valid bool
}

// Constructors

// NullSFPoint constructs and returns a new null SFPoint object.
func NullSFPoint() SFPoint {
	return SFPoint{
		Point: types.SFPoint{},
		Valid: false,
	}
}

// NewSFPoint constructs and returns a new SFPoint object based on the given
// types.SFPoint p. If p is of zero-value, the new SFPoint will be null.
// Otherwise a new, valid SFPoint will be initialized with a copy of p.
func NewSFPoint(p types.SFPoint) SFPoint {
	if p.IsZero() {
		return NullSFPoint()
	}
	return SFPoint{
		Point: types.NewSFPoint(p.Point),
		Valid: true,
	}
}

// NewSFPointXY constructs and returns a new SFPoint object based on the given
// longitude and latitude coordinates.
func NewSFPointXY(x float64, y float64) SFPoint {
	return SFPoint{
		Point: types.NewSFPointXY(x, y),
		Valid: true,
	}
}

// NewSFPointXY constructs and returns a new SFPoint object based on the given
// longitude, latitude, and altitude coordinates.
func NewSFPointXYZ(x float64, y float64, z float64) SFPoint {
	return SFPoint{
		Point: types.NewSFPointXYZ(x, y, z),
		Valid: true,
	}
}

// Getters and Setters

// ValueOrZero will return the value of p if it is valid, or a newly constructed
// zero-value types.SFPoint otherwise.
func (p SFPoint) ValueOrZero() types.SFPoint {
	if !p.Valid {
		return types.SFPoint{}
	}
	return p.Point
}

// Set copies the given types.SFPoint value into p. If the given value is nil,
// p will be nulled.
func (p *SFPoint) Set(v types.SFPoint) {
	if v.IsNil() {
		p.Point = types.SFPoint{} // Let the garbage collector have this types.SFPoint.
		p.Valid = false
		return
	}
	p.Point = v
	p.Valid = true
}

// SetXY will set the latitude and longitude of p.
func (p *SFPoint) SetXY(x float64, y float64) {
	p.Point.SetCoords(geom.Coord{x, y})
	p.Valid = true
}

// SetXYZ will set the latitude, longitude, and altitude of p.
func (p *SFPoint) SetXYZ(x float64, y float64, z float64) {
	p.Point.SetCoords(geom.Coord{x, y, z})
	p.Valid = true
}

// Null will set p to null; p.Valid will be false, and p.Point will contain no
// meaningful value.
func (p *SFPoint) Null() {
	p.Point = types.SFPoint{} // Let the garbage collector have this types.SFPoint.
	p.Valid = false
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if j is null.
func (p SFPoint) IsNil() bool {
	return !p.Valid
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if p is null or if the contained SFPoint is a zero value.
func (p SFPoint) IsZero() bool {
	if !p.Valid {
		return true
	}
	return p.Point.IsZero()
}

// Value implements the database/sql/driver Valuer interface. It will return the
// value of p as a driver.Value. If p is valid, this function will first
// validate the contained SFPoint returning either any encouted parsing errors,
// or a []byte as a driver.Value. If p is null, nil will be returned, and no
// validation will occur.
func (p SFPoint) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}
	return p.Point.Value()
}

// Scan implements the database/sql Scanner interface. It expects to receive a
// valid WKB encoded byte describing a Point as a []byte, or NULL as a nil from
// an SQL database. A zero-length []byte or a nil will be considered NULL, and p
// will be nulled. Otherwise, the value will be passed to types.SFPoint to be
// scanned and parsed as a WKB Point.
func (p *SFPoint) Scan(src interface{}) error {
	if p == nil {
		return fmt.Errorf("null.SFPoint: Scan called on nil pointer")
	}
	switch x := src.(type) {
	case nil:
		p.Point = types.SFPoint{}
		p.Valid = false
		return nil
	case []byte:
		if len(x) == 0 {
			p.Point = types.SFPoint{}
			p.Valid = false
			return nil
		}
		err := p.Point.Scan(x)
		if err != nil {
			p.Point = types.SFPoint{}
			p.Valid = false
			return err
		}
		p.Valid = true
		return nil
	default:
		return fmt.Errorf("null.SFPoint: cannot scan type %T (%v)", src, src)
	}
}

// MarshalJSON implements the encoding/json Marshaler interface. It will return
// the value of p as a JSON-encoded []byte. If p is valid, this function will
// first validate the contained JSON returning either any encouted parsing
// errors, or a []byte. If p is null, "null" will be returned, and no validation
// will occur.
func (p SFPoint) MarshalJSON() ([]byte, error) {
	if !p.Valid {
		return []byte("null"), nil
	}
	return p.Point.MarshalJSON()
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It expects
// to receive a valid JSON value, and will assign that value to this SFPoint. If
// the incoming JSON is the 'null' keyword, p will be nulled. UnmarshalJSON will
// validate the incoming JSON as part of the "Is this JSON null?" check.
func (p *SFPoint) UnmarshalJSON(data []byte) error {
	if p == nil {
		return fmt.Errorf("null.SFPoint: UnmarshalJSON called on nil pointer")
	}
	var k interface{}
	if err := json.Unmarshal(data, &k); err != nil {
		return err
	}
	if k == nil {
		p.Point = types.SFPoint{}
		p.Valid = false
		return nil
	}
	p.Point.UnmarshalJSON(data)
	p.Valid = true
	return nil
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will encode p into its interface{} representation for use in a
// map[string]interface{} by passing it through JSON.Unmarshal if valid, or the
// 'null' keyword otherwise.
func (p SFPoint) MarshalMapValue() (interface{}, error) {
	if !p.Valid {
		return []byte("null"), nil
	}
	return p.Point.MarshalMapValue()
}
