/*
Package types defines a number of helper types that bridge the packages
database/sql, encoding, encoding/json, and encoding/maps. Many are simple
wrappers around types from other libraries (json.RawMessage), but all implement
all of the following interfaces,
 - IsNiler         from pyrrho/encoding       --  IsNil() bool
 - IsZeroer        from pyrrho/encoding       --  IsZero() bool
 - Valuer          from database/sql/driver   --  Value() (driver.Value, error)
 - Scanner         from database/sql          --  Scan(value interface{}) error
 - Marshaler       from encoding/json         --  MarshalJSON() ([]byte, error)
 - Unmarshaler     from encoding/json         --  UnmarshalJSON(data []byte) error
 - Marshaler       from pyrrho/encoding/maps  --  MarshalMap() (map[string]interface{}, error)
 - Unmarshaler     from pyrrho/encoding/maps  --  [Pending maps.Unmarshal features]
*/
package types
