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
	// These are all OpenGIS Simple Feature representations of an XY Polygon with
	// X == 1.2 and Y == 2.3, converted between representations with
	// https://rodic.fr/blog/online-conversion-between-geometric-formats/
	testPolygonWKT     = []byte("POLYGON((30 10,40 40,20 40,10 20,30 10))")
	testPolygonGeoJSON = []byte(`{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}`)
	testPolygonWKB     = []byte{
		0x01, 0x03, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x3e, 0x40, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x24, 0x40, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x44, 0x40, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x44, 0x40, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x34, 0x40, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x44, 0x40, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x24, 0x40, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x34, 0x40, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x3e, 0x40, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x24, 0x40,
	}
	testPolygonCoords = [][]geom.Coord{{
		{30, 10},
		{40, 40},
		{20, 40},
		{10, 20},
		{30, 10},
	}}
	testPolygonGoGeom = *geom.NewPolygon(geom.XY).MustSetCoords(testPolygonCoords)
)

func TestSFPolygonCtors(t *testing.T) {
	require := require.New(t)

	// types.SFPolygon is a wrapper around go-geom's Polygon class. As such,
	// construction typically uses their conventions.
	pa := types.NewSFPolygon(*geom.NewPolygon(geom.XY).MustSetCoords(
		[][]geom.Coord{{
			{30, 10},
			{40, 40},
			{20, 40},
			{10, 20},
			{30, 10},
		}}))
	require.Equal(testPolygonCoords, pa.Coords())
	require.Equal(testPolygonGoGeom, pa.Polygon)

	// We have some helpers to make it easier, though.
	pb, err := types.NewSFPolygonXY([]geom.Coord{
		{30, 10},
		{40, 40},
		{20, 40},
		{10, 20},
		{30, 10},
	})
	require.NoError(err)
	pc, err := types.NewSFPolygonXYZ([]geom.Coord{
		{30, 10, 1},
		{40, 40, 2},
		{20, 40, 3},
		{10, 20, 4},
		{30, 10, 5},
	})
	require.NoError(err)
	pd, err := types.NewSFPolygonXYZM([]geom.Coord{
		{30, 10, 1, 42},
		{40, 40, 2, 42},
		{20, 40, 3, 42},
		{10, 20, 4, 42},
		{30, 10, 5, 42},
	})
	require.NoError(err)

	require.Equal(
		*geom.NewPolygon(geom.XY).MustSetCoords(
			[][]geom.Coord{{
				{30, 10},
				{40, 40},
				{20, 40},
				{10, 20},
				{30, 10}},
			}),
		pb.Polygon)
	require.Equal(
		*geom.NewPolygon(geom.XYZ).MustSetCoords(
			[][]geom.Coord{{
				{30, 10, 1},
				{40, 40, 2},
				{20, 40, 3},
				{10, 20, 4},
				{30, 10, 5},
			}}),
		pc.Polygon)
	require.Equal(
		*geom.NewPolygon(geom.XYZM).MustSetCoords(
			[][]geom.Coord{{
				{30, 10, 1, 42},
				{40, 40, 2, 42},
				{20, 40, 3, 42},
				{10, 20, 4, 42},
				{30, 10, 5, 42},
			}}),
		pd.Polygon)

	// If you do something bad with one of these constructors, you'll get the
	// error expected from go-geom.
	_, err = types.NewSFPolygonXY([]geom.Coord{
		{30, 10, 1, 42},
		{40, 40, 2, 42},
		{20, 40, 3, 42},
		{10, 20, 4, 42},
		{30, 10, 5, 42},
	})
	require.Error(err) // geom: stride mismatch, got 4, want 2
}

func TestSFPolygonIsNil(t *testing.T) {
	require := require.New(t)

	p, _ := types.NewSFPolygonXY([]geom.Coord{
		{30, 10},
		{40, 40},
		{20, 40},
		{10, 20},
		{30, 10}})
	require.False(p.IsNil())

	// NB. I consider this a defect, but it's how go-geom is implemented, so...
	zero, _ := types.NewSFPolygonXY([]geom.Coord{})
	require.True(zero.IsNil())

	empty := types.SFPolygon{}
	require.True(empty.IsNil())
}

func TestSFPolygonIsZero(t *testing.T) {
	require := require.New(t)

	p, _ := types.NewSFPolygonXY([]geom.Coord{
		{30, 10},
		{40, 40},
		{20, 40},
		{10, 20},
		{30, 10}})
	require.False(p.IsZero())

	zero, _ := types.NewSFPolygonXY([]geom.Coord{})
	require.True(zero.IsZero())

	empty := types.SFPolygon{}
	require.True(empty.IsZero())
}

func TestSFPolygonSQLValue(t *testing.T) {
	require := require.New(t)
	var val driver.Value
	var err error

	p, err := types.NewSFPolygonXY([]geom.Coord{
		{30, 10},
		{40, 40},
		{20, 40},
		{10, 20},
		{30, 10},
	})
	require.NoError(err)
	val, err = p.Value()
	require.NoError(err)
	require.EqualValues(testPolygonWKB, val)
}

func TestSFPolygonSQLScan(t *testing.T) {
	require := require.New(t)
	var err error

	var p types.SFPolygon
	err = p.Scan(driver.Value(testPolygonWKB))
	require.NoError(err)
	require.Equal(
		[][]geom.Coord{{
			{30, 10},
			{40, 40},
			{20, 40},
			{10, 20},
			{30, 10},
		}}, p.Coords())

	var bad types.SFPolygon
	err = bad.Scan(driver.Value(nil))
	require.Error(err)
}

func TestSFPolygonMarshalJSON(t *testing.T) {
	require := require.New(t)
	var data []byte
	var err error

	p, _ := types.NewSFPolygonXY([]geom.Coord{
		{30, 10},
		{40, 40},
		{20, 40},
		{10, 20},
		{30, 10},
	})
	data, err = json.Marshal(p)
	require.NoError(err)
	require.EqualValues(testPolygonGeoJSON, data)
	data, err = json.Marshal(&p)
	require.NoError(err)
	require.EqualValues(testPolygonGeoJSON, data)

	bad := types.SFPolygon{}
	_, err = json.Marshal(bad)
	require.Error(err)
	_, err = json.Marshal(&bad)
	require.Error(err)
}

func TestSFPolygonUnmarshalJSON(t *testing.T) {
	require := require.New(t)
	var err error

	var p types.SFPolygon
	err = json.Unmarshal(testPolygonGeoJSON, &p)
	require.NoError(err)
	require.Equal([][]geom.Coord{{
		{30, 10},
		{40, 40},
		{20, 40},
		{10, 20},
		{30, 10},
	}}, p.Coords())
}

func TestSFPolygonMarshsalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ Polygon types.SFPolygon }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	wrapper = Wrapper{*types.MustSFPolygon(types.NewSFPolygonXY([]geom.Coord{
		{30, 10},
		{40, 40},
		{20, 40},
		{10, 20},
		{30, 10},
	}))}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(*types.MustSFPolygon(types.NewSFPolygonXY([]geom.Coord{
		{30, 10},
		{40, 40},
		{20, 40},
		{10, 20},
		{30, 10},
	})), data["Polygon"])
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(*types.MustSFPolygon(types.NewSFPolygonXY([]geom.Coord{
		{30, 10},
		{40, 40},
		{20, 40},
		{10, 20},
		{30, 10},
	})), data["Polygon"])
}
