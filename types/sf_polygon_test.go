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
	// These are all OpenGIS Simple Feature representations of the same
	// XY Polygon, converted between representations with
	// https://rodic.fr/blog/online-conversion-between-geometric-formats/
	testPolygonWKT     = []byte("POLYGON((30 10,40 40,20 40,10 20,30 10),(28 15,15 21,22 35,35 35,28 15))")
	testPolygonGeoJSON = []byte(`{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]],[[28,15],[15,21],[22,35],[35,35],[28,15]]]}`)
	testPolygonWKB     = []byte{
		0x01, 0x03, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
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
		0x00, 0x00, 0x00, 0x24, 0x40, 0x05, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3c,
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e,
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e,
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x35,
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x36,
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41,
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41,
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41,
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3c,
		0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e,
		0x40,
	}
	testPolygonExternal = [][2]float64{
		{30, 10},
		{40, 40},
		{20, 40},
		{10, 20},
		{30, 10},
	}
	testPolygonInternal = [][2]float64{
		{28, 15},
		{15, 21},
		{22, 35},
		{35, 35},
		{28, 15},
	}
	testPolygonCoords = [][]geom.Coord{
		{
			{30, 10},
			{40, 40},
			{20, 40},
			{10, 20},
			{30, 10},
		},
		{
			{28, 15},
			{15, 21},
			{22, 35},
			{35, 35},
			{28, 15},
		},
	}
	testPolygonGoGeom = *geom.NewPolygon(geom.XY).MustSetCoords(testPolygonCoords)
)

func TestSFPolygonCtors(t *testing.T) {
	require := require.New(t)

	// types.SFPolygon is a wrapper around go-geom's Polygon class. As such,
	// construction typically uses their conventions.
	pa := types.NewSFPolygon(*geom.NewPolygon(geom.XY).MustSetCoords(
		[][]geom.Coord{
			{
				{30, 10},
				{40, 40},
				{20, 40},
				{10, 20},
				{30, 10},
			},
			{
				{28, 15},
				{15, 21},
				{22, 35},
				{35, 35},
				{28, 15},
			},
		}))
	require.Equal(testPolygonCoords, pa.Coords())
	require.Equal(testPolygonGoGeom, pa.Polygon)
	// We have some helpers to make it easier, though.
	pb := types.NewSFPolygonXY(
		[][2]float64{
			{30, 10},
			{40, 40},
			{20, 40},
			{10, 20},
			{30, 10},
		},
		[][2]float64{
			{28, 15},
			{15, 21},
			{22, 35},
			{35, 35},
			{28, 15},
		})
	require.Equal(
		*geom.NewPolygon(geom.XY).MustSetCoords(
			[][]geom.Coord{
				{
					{30, 10},
					{40, 40},
					{20, 40},
					{10, 20},
					{30, 10},
				},
				{
					{28, 15},
					{15, 21},
					{22, 35},
					{35, 35},
					{28, 15},
				},
			}),
		pb.Polygon)
	pc := types.NewSFPolygonXYZ([][3]float64{
		{30, 10, 1},
		{40, 40, 2},
		{20, 40, 3},
		{10, 20, 4},
		{30, 10, 5},
	})
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
}

func TestSFPolygonIsNil(t *testing.T) {
	require := require.New(t)

	p := types.NewSFPolygonXY(testPolygonExternal, testPolygonInternal)
	require.False(p.IsNil())

	malformed := types.NewSFPolygonXY([][2]float64{
		{0.0, 0.0},
		{0.0, 0.0},
		{0.0, 0.0},
		{0.0, 0.0},
	})
	require.False(malformed.IsNil())

	zero := types.NewSFPolygonXY([][2]float64{})
	require.True(zero.IsNil())

	nul := types.NewSFPolygonXY(nil)
	require.True(nul.IsNil())

	empty := types.SFPolygon{}
	require.True(empty.IsNil())
}

func TestSFPolygonIsZero(t *testing.T) {
	require := require.New(t)

	p := types.NewSFPolygonXY(testPolygonExternal, testPolygonInternal)
	require.False(p.IsZero())

	malformed := types.NewSFPolygonXY([][2]float64{
		{0.0, 0.0},
		{0.0, 0.0},
		{0.0, 0.0},
		{0.0, 0.0},
	})
	require.True(malformed.IsZero())

	zero := types.NewSFPolygonXY([][2]float64{})
	require.True(zero.IsZero())

	nul := types.NewSFPolygonXY(nil)
	require.True(nul.IsZero())

	empty := types.SFPolygon{}
	require.True(empty.IsZero())
}

func TestSFPolygonSQLValue(t *testing.T) {
	require := require.New(t)
	var val driver.Value
	var err error

	p := types.NewSFPolygonXY(testPolygonExternal, testPolygonInternal)
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
	require.Equal(testPolygonCoords, p.Coords())

	var bad types.SFPolygon
	err = bad.Scan(driver.Value(nil))
	require.Error(err)
}

func TestSFPolygonMarshalJSON(t *testing.T) {
	require := require.New(t)
	var data []byte
	var err error

	p := types.NewSFPolygonXY(testPolygonExternal, testPolygonInternal)
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
	require.Equal(testPolygonCoords, p.Coords())
}

func TestSFPolygonMarshsalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ Polygon types.SFPolygon }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	wrapper = Wrapper{types.NewSFPolygonXY(testPolygonExternal, testPolygonInternal)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(types.NewSFPolygon(testPolygonGoGeom), data["Polygon"])
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(types.NewSFPolygon(testPolygonGoGeom), data["Polygon"])
}
