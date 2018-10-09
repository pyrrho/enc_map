package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-geom"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types"
	"github.com/pyrrho/encoding/types/null"
)

var (
	// These are all OpenGIS Simple Feature representations of the XY test
	// Polygon, converted between representations with
	// https://rodic.fr/blog/online-conversion-between-geometric-formats/
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
	testSFPolygonXY = types.NewSFPolygonXY(
		testPolygonExternal,
		testPolygonInternal)
	// A different polygon to test the third dimension.
	testSFPolygonXYZ = types.NewSFPolygonXYZ([][3]float64{
		{30, 10, 1},
		{40, 40, 2},
		{20, 40, 3},
		{10, 20, 4},
		{30, 10, 5},
	})
	// A malformed test polygon.
	testMalformedPolygon = types.NewSFPolygonXY([][2]float64{
		{0.0, 0.0},
		{0.0, 0.0},
		{0.0, 0.0},
		{0.0, 0.0},
	})
)

func TestSFPolygonCtors(t *testing.T) {
	require := require.New(t)

	// null.NullSFPolygon returns a new null null.SFPolygon.
	// This is equivalent to null.SFPolygon{}.
	na := null.NullSFPolygon()
	require.False(na.Valid)

	// Passing a nil types.SFPolygon to null.NewSFPolygon does the same thing.
	nb := null.NewSFPolygon(types.SFPolygon{})
	require.False(nb.Valid)

	pa := null.NewSFPolygonXY(testPolygonExternal, testPolygonInternal)
	require.True(pa.Valid)
	require.Equal(testSFPolygonXY, pa.Polygon)

	pb := null.NewSFPolygonXYZ([][3]float64{
		{30, 10, 1},
		{40, 40, 2},
		{20, 40, 3},
		{10, 20, 4},
		{30, 10, 5},
	})
	require.True(pb.Valid)
	require.Equal(testSFPolygonXYZ, pb.Polygon)
}

func TestSFPolygonValueOrZero(t *testing.T) {
	require := require.New(t)

	p := null.NewSFPolygonXY(testPolygonExternal, testPolygonInternal)
	require.EqualValues(testSFPolygonXY, p.ValueOrZero())

	n := null.SFPolygon{}
	require.EqualValues(types.SFPolygon{}, n.ValueOrZero())
}

func TestSFPolygonSet(t *testing.T) {
	require := require.New(t)

	p := null.SFPolygon{}

	p.Set(testSFPolygonXY)
	require.True(p.Valid)
	require.EqualValues(testSFPolygonXY, p.ValueOrZero())

	p.Set(types.SFPolygon{})
	require.False(p.Valid)

	p.Set(testMalformedPolygon)
	require.True(p.Valid)
	require.EqualValues(testMalformedPolygon, p.ValueOrZero())
}

func TestSFPolygonNull(t *testing.T) {
	require := require.New(t)

	p := null.NewSFPolygon(testSFPolygonXY)

	p.Null()
	require.False(p.Valid)
}

func TestSFPolygonIsNil(t *testing.T) {
	require := require.New(t)

	p := null.NewSFPolygonXY(testPolygonExternal, testPolygonInternal)
	require.False(p.IsNil())

	malformed := null.NewSFPolygon(testMalformedPolygon)
	require.False(malformed.IsNil())

	zero := null.NewSFPolygonXY([][2]float64{})
	require.False(zero.IsNil())

	nul := null.NewSFPolygonXY(nil)
	require.False(nul.IsNil())

	empty := null.SFPolygon{}
	require.True(empty.IsNil())
}

func TestSFPolygonIsZero(t *testing.T) {
	require := require.New(t)

	p := null.NewSFPolygonXY(testPolygonExternal, testPolygonInternal)
	require.False(p.IsZero())

	malformed := null.NewSFPolygon(testMalformedPolygon)
	require.True(malformed.IsZero())

	zero := null.NewSFPolygonXY([][2]float64{})
	require.True(zero.IsZero())

	nul := null.NewSFPolygonXY(nil)
	require.True(nul.IsZero())

	empty := null.SFPolygon{}
	require.True(empty.IsZero())
}

func TestSFPolygonSQLValue(t *testing.T) {
	require := require.New(t)
	var val driver.Value
	var err error

	p := null.NewSFPolygonXY(testPolygonExternal, testPolygonInternal)
	val, err = p.Value()
	require.NoError(err)
	require.EqualValues(testPolygonWKB, val)
}

func TestSFPolygonSQLScan(t *testing.T) {
	require := require.New(t)
	var err error

	var p null.SFPolygon
	err = p.Scan(driver.Value(testPolygonWKB))
	require.NoError(err)
	require.Equal(null.NewSFPolygon(testSFPolygonXY), p)

	var n null.SFPolygon
	err = n.Scan(driver.Value(nil))
	require.NoError(err)
	require.Equal(null.NullSFPolygon(), n)
}

func TestSFPolygonMarshalJSON(t *testing.T) {
	require := require.New(t)
	var data []byte
	var err error

	p := null.NewSFPolygonXY(testPolygonExternal, testPolygonInternal)
	data, err = json.Marshal(p)
	require.NoError(err)
	require.EqualValues(testPolygonGeoJSON, data)
	data, err = json.Marshal(&p)
	require.NoError(err)
	require.EqualValues(testPolygonGeoJSON, data)

	n := null.SFPolygon{}
	data, err = json.Marshal(n)
	require.NoError(err)
	require.EqualValues("null", data)
	data, err = json.Marshal(&n)
	require.NoError(err)
	require.EqualValues("null", data)
}

func TestSFPolygonUnmarshalJSON(t *testing.T) {
	require := require.New(t)
	var err error

	var p null.SFPolygon
	err = json.Unmarshal(testPolygonGeoJSON, &p)
	require.NoError(err)
	require.Equal(null.NewSFPolygon(testSFPolygonXY), p)
}

func TestSFPolygonMarshsalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ Polygon null.SFPolygon }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	wrapper = Wrapper{
		null.NewSFPolygonXY(testPolygonExternal, testPolygonInternal)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(testSFPolygonXY, data["Polygon"])
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(testSFPolygonXY, data["Polygon"])
}
