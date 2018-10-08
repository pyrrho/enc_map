package types

import (
	"bytes"
	"database/sql/driver"
	"fmt"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkb"
)

// SFPoint is a Simple Feature Point, named for the OpenGIS specification that
// backs WKB, WKT, and GeoJSON representations of geospatial data. An SFPoint
// represents a single [longitude, latitude] or [longitude, latitude, altitude]
// point in a given coordinate system.
//
// This type is built on top of the go-geom geom.Point type, implementing all of
// the pyrrho/encoding/types interfaces detailed in the package comments.
// Database interactions (Value and Scan) will convert to and from a WKB (Well
// Known Binary) representation. JSON interactions (MarshalJSON and
// UnmarshalJSON) will convert to and from a GeoJSON representation.
type SFPoint struct {
	geom.Point
}

// Constructors

// NewSFPoint constructs and returns a new SFPoint object initialized with the
// given geom.Point p.
func NewSFPoint(p geom.Point) SFPoint {
	return SFPoint{p}
}

// NewSFPointXY constructs and returns a new SFPoint with longitude and
// latitude components.
func NewSFPointXY(x float64, y float64) SFPoint {
	p, err := geom.NewPoint(geom.XY).SetCoords(geom.Coord{x, y})
	if err != nil {
		panic(err)
	}
	return SFPoint{*p}
}

// NewSFPointXYZ constructs and returns a new SFPoint with longitude, latitude,
// and altitude components.
func NewSFPointXYZ(x float64, y float64, z float64) SFPoint {
	p, err := geom.NewPoint(geom.XYZ).SetCoords(geom.Coord{x, y, z})
	if err != nil {
		panic(err)
	}
	return SFPoint{*p}
}

// Getters

// Lng returns the longitude (northing, first) component of this SFPoint.
func (p SFPoint) Lng() float64 {
	return p.X()
}

// Lat returns the latitude (easting, second) component of this SFPoint.
func (p SFPoint) Lat() float64 {
	return p.Y()
}

// Alt returns the altitude (third) component of this SFPoint, or 0 if the
// SFPoint has no altitude component (is an XY SFPoint).
func (p SFPoint) Alt() float64 {
	return p.Z()
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if p contains no meaningful data. More specifically, if this if this SFPoint
// has been zero-initialized, or if it has been explicitly initialized with no
// layout;
//   var p types.SFPoint
//   var p := types.SFPoint{}
//   var p := types.NewSFPoint(geom.Point{geom.NoLayout})
func (p SFPoint) IsNil() bool {
	// return p.FlatCoords() == nil || p.Layout() == geom.NoLayout
	return len(p.FlatCoords()) == 0
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if p.IsNil() returns true, or if the contained data is of the zero-value.
func (p SFPoint) IsZero() bool {
	for _, f := range p.FlatCoords() {
		if f != 0.0 {
			return false
		}
	}
	return true
}

// Value implements the database/sql/driver Valuer interface. It will return the
// value of p as a driver.Value; specifically a WKB encoded []byte.
func (p SFPoint) Value() (driver.Value, error) {
	b := &bytes.Buffer{}
	if err := wkb.Write(b, wkb.NDR, &p.Point); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Scan implements the database/sql Scanner interface. It expects to receive a
// WKB encoded []byte describing a Point from an SQL database, and will assign
// that value to p. If the incoming []byte is not a well formed WKB, or if that
// WKB value does not describe a Point, an error will be returned.
func (p *SFPoint) Scan(src interface{}) error {
	if p == nil {
		return fmt.Errorf("types.SFPoint: Scan called on nil SFLpointer")
	}
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("types.SFPoint: cannot scan type %T (%v)", src, src)
	}
	g, err := wkb.Unmarshal(b)
	if err != nil {
		return err
	}
	t, ok := g.(*geom.Point)
	if !ok {
		return fmt.Errorf("types.SFPoint: scan did not return a *geom.Point (got a %T)", t)
	}
	p.Point.Swap(t)
	return nil
}

// MarshalJSON implements the encoding/json Marshaler interface. It will return
// the GeoJSON encoded representation of p.
func (p SFPoint) MarshalJSON() ([]byte, error) {
	if p.IsNil() {
		return nil, fmt.Errorf("types.SFPoint: cannot unmarshal an uninitialized SFPoint")
	}
	return geojson.Marshal(&p.Point)
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It expects
// to receive a valid GeoJSON Geometry of the type Point, and will assign
// the value of that data to p.
func (p *SFPoint) UnmarshalJSON(data []byte) error {
	if p == nil {
		return fmt.Errorf("types.SFPoint: UnmarshalJSON called on nil SFLpointer")
	}
	var gt geom.T
	if err := geojson.Unmarshal(data, &gt); err != nil {
		return err
	}
	p.Point.Swap(gt.(*geom.Point))
	return nil
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will return p wrapped in an interface{} for use in a map[string]interface{}.
func (p SFPoint) MarshalMapValue() (interface{}, error) {
	return p, nil
}
