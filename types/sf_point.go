package types

import (
	"bytes"
	"database/sql/driver"
	"fmt"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkb"
)

type SFPoint struct {
	geom.Point
}

// Constructors

func NewSFPoint(p geom.Point) *SFPoint {
	return &SFPoint{p}
}

func MustSFPoint(p *SFPoint, err error) *SFPoint {
	if err != nil {
		panic(err)
	}
	return p
}

func NewSFPointXY(c ...float64) (*SFPoint, error) {
	p, err := geom.NewPoint(geom.XY).SetCoords(c)
	if err != nil {
		return nil, err
	}
	return &SFPoint{*p}, err
}

func NewSFPointXYM(c ...float64) (*SFPoint, error) {
	p, err := geom.NewPoint(geom.XYM).SetCoords(c)
	if err != nil {
		return nil, err
	}
	return &SFPoint{*p}, err
}

func NewSFPointXYZ(c ...float64) (*SFPoint, error) {
	p, err := geom.NewPoint(geom.XYZ).SetCoords(c)
	if err != nil {
		return nil, err
	}
	return &SFPoint{*p}, err
}

func NewSFPointXYZM(c ...float64) (*SFPoint, error) {
	p, err := geom.NewPoint(geom.XYZM).SetCoords(c)
	if err != nil {
		return nil, err
	}
	return &SFPoint{*p}, err
}

// Interfaces

func (p SFPoint) IsNil() bool {
	return p.FlatCoords() == nil
}

func (p SFPoint) IsZero() bool {
	for _, f := range p.FlatCoords() {
		if f != 0.0 {
			return false
		}
	}
	return true
}

func (p SFPoint) Value() (driver.Value, error) {
	b := &bytes.Buffer{}
	if err := wkb.Write(b, wkb.NDR, &p.Point); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

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

func (p SFPoint) MarshalJSON() ([]byte, error) {
	if p.IsNil() {
		return nil, fmt.Errorf("types.SFPoint: cannot unmarshal an uninitialized SFPoint")
	}
	return geojson.Marshal(&p.Point)
}

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

func (p SFPoint) MarshalMapValue() (interface{}, error) {
	return p, nil
}
