package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pyrrho/encoding/maps"
	"github.com/pyrrho/encoding/types"
	"github.com/pyrrho/encoding/types/null"
)

var (
	testSFPointXY  = types.NewSFPointXY(1.2, 2.3)
	testSFPointXYZ = types.NewSFPointXYZ(1.2, 2.3, 3.4)
	// WKB representation of the XY test point.
	testPointXYWKB = []byte{
		0x01, 0x01, 0x00, 0x00, 0x00, 0x33, 0x33, 0x33,
		0x33, 0x33, 0x33, 0xf3, 0x3f, 0x66, 0x66, 0x66,
		0x66, 0x66, 0x66, 0x02, 0x40,
	}
	// GeoJSON representation of the XY test point.
	testPointXYGeoJSON = []byte(`{"type":"Point","coordinates":[1.2,2.3]}`)
)

func TestSFPointCtors(t *testing.T) {
	require := require.New(t)

	// null.NullSFPoint returns a new null null.SFPoint.
	// This is equivalwent to null.SFPoint{}.
	pa := null.NullSFPoint()
	require.False(pa.Valid)

	// Passing a nil types.SFPoint to null.NewSFPoint does the same thing.
	pb := null.NewSFPoint(types.SFPoint{})
	require.False(pb.Valid)

	// Passing a non-nil types.SFPoint constructs a new, valid null.SFPoint.
	pc := null.NewSFPoint(testSFPointXY)
	require.True(pc.Valid)
	require.Equal(testSFPointXY, pc.Point)

	// You can also create null.SFPoints by passing coordinates.
	pd := null.NewSFPointXY(1.2, 2.3)
	require.Equal(testSFPointXY, pd.Point)

	pe := null.NewSFPointXYZ(1.2, 2.3, 3.4)
	require.Equal(testSFPointXYZ, pe.Point)
}

func TestSFPointValueOrZero(t *testing.T) {
	require := require.New(t)

	pa := null.NewSFPoint(testSFPointXY)
	require.EqualValues(testSFPointXY, pa.ValueOrZero())

	pb := null.NewSFPoint(testSFPointXYZ)
	require.EqualValues(testSFPointXYZ, pb.ValueOrZero())

	n := null.SFPoint{}
	require.EqualValues(types.SFPoint{}, n.ValueOrZero())
}

func TestSFPointSet(t *testing.T) {
	require := require.New(t)

	p := null.SFPoint{}

	p.Set(testSFPointXY)
	require.True(p.Valid)
	require.EqualValues(testSFPointXY, p.ValueOrZero())

	p.Set(types.SFPoint{})
	require.False(p.Valid)

	p.Set(testSFPointXYZ)
	require.True(p.Valid)
	require.EqualValues(testSFPointXYZ, p.ValueOrZero())
}

func TestSFPointNull(t *testing.T) {
	require := require.New(t)

	p := null.NewSFPointXY(1.2, 2.3)

	p.Null()
	require.False(p.Valid)
}

func TestSFPointIsNil(t *testing.T) {
	require := require.New(t)

	p := null.NewSFPointXY(1.2, 2.3)
	require.False(p.IsNil())

	zero := null.NewSFPointXY(0.0, 0.0)
	require.False(zero.IsNil())

	empty := null.SFPoint{}
	require.True(empty.IsNil())
}

func TestSFPointIsZero(t *testing.T) {
	require := require.New(t)

	p := null.NewSFPointXY(1.2, 2.3)
	require.False(p.IsZero())

	zero := null.NewSFPointXY(0.0, 0.0)
	require.True(zero.IsZero())

	empty := null.SFPoint{}
	require.True(empty.IsZero())
}

func TestSFPointSQLValue(t *testing.T) {
	require := require.New(t)
	var val driver.Value
	var err error

	p := null.NewSFPointXY(1.2, 2.3)
	val, err = p.Value()
	require.NoError(err)
	require.EqualValues(testPointXYWKB, val)
}

func TestSFPointSQLScan(t *testing.T) {
	require := require.New(t)
	var err error

	var p null.SFPoint
	err = p.Scan(driver.Value(testPointXYWKB))
	require.NoError(err)
	require.Equal(null.NewSFPoint(testSFPointXY), p)

	var n null.SFPoint
	err = n.Scan(driver.Value(nil))
	require.NoError(err)
	require.Equal(null.NullSFPoint(), n)
}

func TestSFPointMarshalJSON(t *testing.T) {
	require := require.New(t)
	var data []byte
	var err error

	p := null.NewSFPointXY(1.2, 2.3)
	data, err = json.Marshal(p)
	require.NoError(err)
	require.EqualValues(testPointXYGeoJSON, data)
	data, err = json.Marshal(&p)
	require.NoError(err)
	require.EqualValues(testPointXYGeoJSON, data)

	n := null.SFPoint{}
	data, err = json.Marshal(n)
	require.NoError(err)
	require.EqualValues("null", data)
	data, err = json.Marshal(&n)
	require.NoError(err)
	require.EqualValues("null", data)
}

func TestSFPointUnmarshalJSON(t *testing.T) {
	require := require.New(t)
	var err error

	var p null.SFPoint
	err = json.Unmarshal(testPointXYGeoJSON, &p)
	require.NoError(err)
	require.EqualValues(null.NewSFPoint(testSFPointXY), p)
}

func TestSFPointMarshsalMapValue(t *testing.T) {
	require := require.New(t)
	type Wrapper struct{ Point null.SFPoint }
	var wrapper Wrapper
	var data map[string]interface{}
	var err error

	wrapper = Wrapper{null.NewSFPointXY(1.2, 2.3)}
	data, err = maps.Marshal(wrapper)
	require.NoError(err)
	require.Equal(testSFPointXY, data["Point"])
	data, err = maps.Marshal(&wrapper)
	require.NoError(err)
	require.Equal(testSFPointXY, data["Point"])
}
