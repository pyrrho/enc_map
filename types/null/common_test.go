package null_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// This file performs no meaningful tests; it's an aggregate for shared state
// and helper functions that are used in the package's test suite.

// Descriptive Tests
// -----------------
// These don't actually test anything, rather they demonstrate behavior of the
// language and other packages. They are included because they may help clarify
// why certain choices were made within this library.

func TestJSONUnmarshalErrors(t *testing.T) {
	require := require.New(t)
	// encoding/json considers and empty byte stream to be invalid JSON.
	// This is correct; _nothing_ is not a valid JSON _something_.
	var iface interface{}
	err := json.Unmarshal([]byte(""), &iface)
	require.Error(err)

	// A pair of double-quotes with nothing between them will be parsed as a
	// string or a []byte, but nothing else.
	var (
		s  string
		bs []byte
		i  int
		is []int
	)
	err = json.Unmarshal([]byte(`""`), &s)
	require.NoError(err)
	err = json.Unmarshal([]byte(`""`), &bs)
	require.NoError(err)
	err = json.Unmarshal([]byte(`""`), &i)
	require.Error(err)
	err = json.Unmarshal([]byte(`""`), &is)
	require.Error(err)
}
