package types_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/pyrrho/encoding/maps"

	"github.com/pyrrho/encoding/types"
	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-geom"
)

var (
	// These are all OpenGIS Simple Feature representations of an XY Point with
	// X == 1.2 and Y == 2.3, converted between representations with
	// https://rodic.fr/blog/online-conversion-between-geometric-formats/
	testPointWKT     = []byte("POINT(1.2 2.3)")
	testPointGeoJSON = []byte(`{"type":"Point","coordinates":[1.2,2.3]}`)
	testPointWKB     = []byte{
		0x01, 0x01, 0x00, 0x00, 0x00, 0x33, 0x33, 0x33,
		0x33, 0x33, 0x33, 0xf3, 0x3f, 0x66, 0x66, 0x66,
		0x66, 0x66, 0x66, 0x02, 0x40,
	}
)

func TestSFPointCtors(t *testing.T) {
	require := require.New(t)

	// types.SFPoint is a wrapper around go-geom's Point class. As such,
	// construction typically uses their conventions.
	pa := types.NewSFPoint(
		*geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1.2, 2.3}))
	require.Equal(1.2, pa.X())
	require.Equal(2.3, pa.Y())
	require.Equal(
		*geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1.2, 2.3}),
		pa.Point)

	// We have some helpers to make it easier, though.
	pb, err := types.NewSFPointXY(1.2, 2.3)
	require.NoError(err)
	pc, err := types.NewSFPointXYZ(1.2, 2.3, 3.4)
	require.NoError(err)
	pd, err := types.NewSFPointXYZM(1.2, 2.3, 3.4, 4.5)
	require.NoError(err)

	require.Equal(
		*geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{1.2, 2.3}),
		pb.Point)
	require.Equal(
		*geom.NewPoint(geom.XYZ).MustSetCoords(geom.Coord{1.2, 2.3, 3.4}),
		pc.Point)
	require.Equal(
		*geom.NewPoint(geom.XYZM).MustSetCoords(geom.Coord{1.2, 2.3, 3.4, 4.5}),
		pd.Point)

	// If you do something bad with one of these constructors, you'll get the
	// error expected from go-geom.
	_, err = types.NewSFPointXY(1.2, 2.3, 3.4)
	require.Error(err) // geom: stride mismatch, got 3, want 2
}

func TestSFPointIsNil(t *testing.T) {
	require := require.New(t)

	p, _ := types.NewSFPointXY(1.2, 2.3)
	require.False(p.IsNil())

	zero, _ := types.NewSFPointXY(0.0, 0.0)
	require.False(zero.IsNil())

	empty := types.SFPoint{}
	require.True(empty.IsNil())
}

func TestSFPointIsZero(t *testing.T) {
	require := require.New(t)

	p, _ := types.NewSFPointXY(1.2, 2.3)
	require.False(p.IsZero())

	zero, _ := types.NewSFPointXY(0.0, 0.0)
	require.True(zero.IsZero())

	empty := types.SFPoint{}
	require.True(empty.IsZero())
}

func TestSFPointSQLValue(t *testing.T) {
	require := require.New(t)
	var val driver.Value
	var err error

	p, err := types.NewSFPointXY(1.2, 2.3)
	require.NoError(err)
	val, err = p.Value()
	require.NoError(err)
	require.EqualValues(testPointWKB, val)
}

func TestSFPointSQLScan(t *testing.T) {
	require := require.New(t)
	var err error

	var p types.SFPoint
	err = p.Scan(driver.Value(testPointWKB))
	require.NoError(err)
	require.Equal([]float64{1.2, 2.3}, p.FlatCoords())

	var bad types.SFPoint
	err = bad.Scan(driver.Value(nil))
	require.Error(err)
}

func TestSFPointMarshalJSON(t *testing.T) {
	require := require.New(t)
	var data []byte
	var err error

	p, _ := types.NewSFPointXY(1.2, 2.3)
	data, err = json.Marshal(p)
	require.NoError(err)
	require.EqualValues(testPointGeoJSON, data)
	data, err = json.Marshal(&p)
	require.NoError(err)
	require.EqualValues(testPointGeoJSON, data)

	bad := types.SFPoint{}
	_, err = json.Marshal(bad)
	require.Error(err)
	_, err = json.Marshal(&bad)
	require.Error(err)
}

func TestSFPointUnmarshalJSON(t *testing.T) {
	require := require.New(t)
	var err error

	var p types.SFPoint
	err = json.Unmarshal(testPointGeoJSON, &p)
	require.NoError(err)
	require.Equal(1.2, p.X())
	require.Equal(2.3, p.Y())
}

func TestSFPointMarshsalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ Point types.SFPoint }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	wrapper = Wrapper{*types.MustSFPoint(types.NewSFPointXY(1.2, 2.3))}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(*types.MustSFPoint(types.NewSFPointXY(1.2, 2.3)), data["Point"])
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(*types.MustSFPoint(types.NewSFPointXY(1.2, 2.3)), data["Point"])
}
