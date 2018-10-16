/*
Package encoding is an extension to the Go Standard Library's encoding
(https://godoc.org/encoding) package. It defines interfaces and helper
functions shared by other packages (specifically the sub-packages of this
library).
*/
package encoding

import (
	"reflect"
)

// IsNiler is an interface implemented by an object with a nil value that may
// differ from Go's default nil value. This is used in encoding/map with the
// "omitnil" struct tag to give fields a chance to specify when they should be
// omitted due to containing a nil value.
type IsNiler interface {
	IsNil() bool
}

// IsNil returns true if the given value `v` is equivalent to nil,
// either because it is an `IsNiler` and `v.IsNil()` returns `true`, or it is a
// channel, function variable, interface, map, pointer, or slice whose value is
// the associated `nil`.
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}
	if n, ok := v.(IsNiler); ok {
		return n.IsNil()
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	}
	return false
}

// IsValueNil returns true if the given reflect.Value `v` is equivalent to nil,
// either because it is an `IsNiler` and `v.IsNil()` returns `true`, or it is a
// channel, function variable, interface, map, pointer, or slice whose value is
// the associated `nil`.
func IsValueNil(v reflect.Value) bool {
	if n, ok := v.Interface().(IsNiler); ok {
		return n.IsNil()
	}
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}

// IsZeroer is an interface implemented by an object with a zero value that may
// differ from Go's default zero value. This is used in encoding/map with the
// "omitzero" struct tag to give fields a chance to specify when they should be
// omitted due to containing a zero value.
type IsZeroer interface {
	IsZero() bool
}

// IsZero returns true if the given value `v` is that type's zero value, either
// because it is an `IsZeroer` and `v.IsZero()` returns `true`, or if it is
// equal to that type's default zero value.
func IsZero(v interface{}) bool {
	if v == nil {
		return true
	}
	if z, ok := v.(IsZeroer); ok {
		return z.IsZero()
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	// We consider an Array to be of Zero Value when all of its elements are
	// that type's Zero Value.
	case reflect.Array:
		ret := true
		for i := 0; i < rv.Len(); i++ {
			if !IsValueZero(rv.Index(i)) {
				ret = false
				break
			}
		}
		return ret
	case reflect.Struct:
		ret := true
		for i := 0; i < rv.NumField(); i++ {
			if !IsValueZero(rv.Field(i)) {
				ret = false
				break
			}
		}
		return ret
	case reflect.Map, reflect.Slice, reflect.String:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return rv == reflect.Zero(reflect.TypeOf(rv)).Interface()
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Ptr:
		return rv.IsNil()
	}
	return false
}

func IsValueZero(v reflect.Value) bool {
	if z, ok := v.Interface().(IsZeroer); ok {
		return z.IsZero()
	}
	switch v.Kind() {
	// We consider an Array to be of Zero Value when all of its elements are
	// that type's Zero Value.
	case reflect.Array:
		ret := true
		for i := 0; i < v.Len(); i++ {
			if !IsValueZero(v.Index(i)) {
				ret = false
				break
			}
		}
		return ret
	case reflect.Struct:
		ret := true
		for i := 0; i < v.NumField(); i++ {
			if !IsValueZero(v.Field(i)) {
				ret = false
				break
			}
		}
		return ret
	case reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v == reflect.Zero(reflect.TypeOf(v)).Interface()
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
