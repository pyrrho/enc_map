/*
Package maps provides helper functions and struct tags that facilitate the
conversion of structs to `map[string]interface{}`s, and vice-versa.

This package is primarily inspired by the encode/json and database/sql packages,
as well as other open-source alternatives.

Note that this package relies _heavily_ on the reflect package and, as such,
has severely weakened compile-time type-safety. Be sure to keep an eye on your
error returns.
*/
package maps
