package types

import (
	"bytes"
	"database/sql/driver"
	"fmt"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkb"
)

type SFPolygon struct {
	geom.Polygon
}

// Constructors

func NewSFPolygon(p geom.Polygon) *SFPolygon {
	return &SFPolygon{p}
}

func MustSFPolygon(p *SFPolygon, err error) *SFPolygon {
	if err != nil {
		panic(err)
	}
	return p
}

func NewSFPolygonXY(c ...[]geom.Coord) (*SFPolygon, error) {
	p, err := geom.NewPolygon(geom.XY).SetCoords(c)
	if err != nil {
		return nil, err
	}
	return &SFPolygon{*p}, err
}

func NewSFPolygonXYM(c ...[]geom.Coord) (*SFPolygon, error) {
	p, err := geom.NewPolygon(geom.XYM).SetCoords(c)
	if err != nil {
		return nil, err
	}
	return &SFPolygon{*p}, err
}

func NewSFPolygonXYZ(c ...[]geom.Coord) (*SFPolygon, error) {
	p, err := geom.NewPolygon(geom.XYZ).SetCoords(c)
	if err != nil {
		return nil, err
	}
	return &SFPolygon{*p}, err
}

func NewSFPolygonXYZM(c ...[]geom.Coord) (*SFPolygon, error) {
	p, err := geom.NewPolygon(geom.XYZM).SetCoords(c)
	if err != nil {
		return nil, err
	}
	return &SFPolygon{*p}, err
}

// Interfaces

func (p SFPolygon) IsNil() bool {
	return p.FlatCoords() == nil
}

func (p SFPolygon) IsZero() bool {
	return len(p.FlatCoords()) == 0
}

func (p SFPolygon) Value() (driver.Value, error) {
	b := &bytes.Buffer{}
	if err := wkb.Write(b, wkb.NDR, &p.Polygon); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

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

func (p SFPolygon) MarshalJSON() ([]byte, error) {
	if p.IsNil() {
		return nil, fmt.Errorf("types.SFPolygon: cannot unmarshal an uninitialized SFPolygon")
	}
	return geojson.Marshal(&p.Polygon)
}

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

func (p SFPolygon) MarshalMapValue() (interface{}, error) {
	return p, nil
}
