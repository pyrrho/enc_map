package types

import (
	"bytes"
	"database/sql/driver"
	"fmt"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkb"
)

// SFPolygon is a Simple Feature Polygon, named for the OpenGIS specification
// that backs WKB, WKT, and GeoJSON representations of geospatial data. An
// SFPolygon represents a series of [longitude, latitude] or [longitude,
// latitude, altitude] points in a given coordinate system that make up one
// external polygon bounding the given shape, and zero or more internal polygons
// bounding holes within the shape. These polygons should follow the right-hand
// rule for wrapping order; polygons with a positive area (the external,
// positive-space shape) should wrap counter-clockwise, and polygons with a
// negative area (the internal, negative-space holes) should wrap clockwise.
//
// This type is built on top of the go-geom geom.Polygon type, implementing all
// of the pyrrho/encoding/types interfaces detailed in the package comments.
// Database interactions (Value and Scan) will convert to and from a WKB (Well
// Known Binary) representation. JSON interactions (MarshalJSON and
// UnmarshalJSON) will convert to and from a GeoJSON representation.
type SFPolygon struct {
	geom.Polygon
}

// Constructors

// NewSFPolygon constructs and returns a new SFPolygon object initialized with
// the given geom.Polygon p.
func NewSFPolygon(p geom.Polygon) SFPolygon {
	return SFPolygon{p}
}

// NewSFPolygonXY constructs and returns a new SFPolygon object with longitude
// and latitude components initialized with the given external and (optionally)
// internal shapes.
func NewSFPolygonXY(external [][2]float64, internals ...[][2]float64) SFPolygon {
	l := len(internals) + 1
	polys := make([][]geom.Coord, l)

	polys[0] = make([]geom.Coord, len(external))
	for i := range external {
		polys[0][i] = append(geom.Coord(nil), external[i][:]...)
	}
	for j, internal := range internals {
		polys[j+1] = make([]geom.Coord, len(internal))
		for i := range internal {
			polys[j+1][i] = append(geom.Coord(nil), internal[i][:]...)
		}
	}

	p, err := geom.NewPolygon(geom.XY).SetCoords(polys)
	if err != nil {
		panic(err)
	}
	return SFPolygon{*p}
}

// NewSFPolygonXYZ constructs and returns a new SFPolygon object with longitude,
// latitude, and altitude components initialized with the given external and
// (optionally) internal shapes.
func NewSFPolygonXYZ(external [][3]float64, internals ...[][3]float64) SFPolygon {
	l := len(internals) + 1
	polys := make([][]geom.Coord, l)

	polys[0] = make([]geom.Coord, len(external))
	for i := range external {
		polys[0][i] = append(geom.Coord(nil), external[i][:]...)
	}
	for j, internal := range internals {
		polys[j+1] = make([]geom.Coord, len(internal))
		for i := range internal {
			polys[j+1][i] = append(geom.Coord(nil), internal[i][:]...)
		}
	}

	p, err := geom.NewPolygon(geom.XYZ).SetCoords(polys)
	if err != nil {
		panic(err)
	}
	return SFPolygon{*p}
}

// Interfaces

// IsNil implements the pyrrho/encoding IsNiler interface. It will return true
// if p contains no meaningful data. More specifically, if this if this SFPoint
// has been zero-initialized, or if it has been explicitly initialized with no
// layout;
//   var p types.SFPolygon
//   var p := types.SFPolygon{}
//   var p := types.NewSFPolygon(geom.Polygon{geom.NoLayout})
//   var p := types.NewSFPolygonXY(nil)
//   var p := types.NewSFPolygonXY([][2]float64{})
func (p SFPolygon) IsNil() bool {
	return p.FlatCoords() == nil || p.Layout() == geom.NoLayout
}

// IsZero implements the pyrrho/encoding IsZeroer interface. It will return true
// if p.IsNil() returns true, or if the contained data is of the zero-value.
func (p SFPolygon) IsZero() bool {
	for _, f := range p.FlatCoords() {
		if f != 0.0 {
			return false
		}
	}
	return true
}

// Value implements the database/sql/driver Valuer interface. It will return the
// value of p as a driver.Value; specifically a WKB encoded []byte.
func (p SFPolygon) Value() (driver.Value, error) {
	b := &bytes.Buffer{}
	if err := wkb.Write(b, wkb.NDR, &p.Polygon); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Scan implements the database/sql Scanner interface. It expects to receive a
// WKB encoded byte describing a Polygon from an SQL database, and will assign
// that value to p. If the incoming []byte is not a well formed WKB, or if that
// WKB value does not describe a Polygon, an error will be returned.
func (p *SFPolygon) Scan(src interface{}) error {
	if p == nil {
		return fmt.Errorf("types.SFPolygon: Scan called on nil SFLPolygoner")
	}
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("types.SFPolygon: cannot scan type %T (%v)", src, src)
	}
	g, err := wkb.Unmarshal(b)
	if err != nil {
		return err
	}
	t, ok := g.(*geom.Polygon)
	if !ok {
		return fmt.Errorf("types.SFPolygon: scan did not return a *geom.Polygon (got a %T)", t)
	}
	p.Polygon.Swap(t)
	return nil
}

// MarshalJSON implements the encoding/json Marshaler interface. It will return
// the GeoJSON encoded representation of p.
func (p SFPolygon) MarshalJSON() ([]byte, error) {
	if p.IsNil() {
		return nil, fmt.Errorf("types.SFPolygon: cannot unmarshal an uninitialized SFPolygon")
	}
	return geojson.Marshal(&p.Polygon)
}

// UnmarshalJSON implements the encoding/json Unmarshaler interface. It expects
// to receive a valid GeoJSON Geometry with of the type Polygon, and will assign
// the value of that data to p.
func (p *SFPolygon) UnmarshalJSON(data []byte) error {
	if p == nil {
		return fmt.Errorf("types.SFPolygon: UnmarshalJSON called on nil SFLPolygoner")
	}
	var gt geom.T
	if err := geojson.Unmarshal(data, &gt); err != nil {
		return err
	}
	p.Polygon.Swap(gt.(*geom.Polygon))
	return nil
}

// MarshalMapValue implements the pyrrho/encoding/maps Marshaler interface. It
// will return p wrapped in an interface{} for use in a map[string]interface{}.
func (p SFPolygon) MarshalMapValue() (interface{}, error) {
	return p, nil
}
