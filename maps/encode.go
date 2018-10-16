package maps

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/pyrrho/encoding"
)

func Marshal(src interface{}) (map[string]interface{}, error) {
	ret, err := defaultConfig.marshal(src)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func MarshalSlice(src interface{}) ([]map[string]interface{}, error) {
	ret, err := defaultConfig.marshalSlice(src)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func MarshalWithConfig(src interface{}, cfg *Config) (map[string]interface{}, error) {
	ret, err := cfg.marshal(src)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func MarshalSliceWithConfig(src interface{}, cfg *Config) ([]map[string]interface{}, error) {
	ret, err := cfg.marshalSlice(src)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type Marshaler interface {
	MarshalMapValue() (interface{}, error)
}

var marshalerType = reflect.TypeOf(new(Marshaler)).Elem()

func (cfg *Config) marshal(src interface{}) (m map[string]interface{}, err error) {
	srcv := reflect.ValueOf(src)
	if srcv.Kind() == reflect.Ptr {
		srcv = srcv.Elem()
	}
	if srcv.Kind() != reflect.Struct {
		return nil, errors.New("src must be a struct, or pointer-to-struct")
	}

	// Any panics after this point should be converted to errors, and returned
	// normally. Unless it's a runtime error, it's a raw string, or it's not of
	// type `error`. In which case, do panic.
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			} else if s, ok := r.(string); ok {
				panic(s)
			} else if e, ok := r.(error); !ok {
				panic(r)
			} else {
				err = e
			}
		}
	}()

	ret := lookupEncodeFn(srcv.Type(), cfg)(srcv, cfg)
	return ret.(map[string]interface{}), nil
}

func (cfg *Config) marshalSlice(src interface{}) (ret []map[string]interface{}, err error) {
	srcv := reflect.ValueOf(src)
	if srcv.Kind() == reflect.Ptr {
		srcv = srcv.Elem()
	}
	if !(srcv.Kind() == reflect.Array || srcv.Kind() == reflect.Slice) {
		return nil, errors.New("src must be a slice, array, or pointer to either")
	}

	// Any panics after this point should be converted to errors, and returned
	// normally. Unless it's a runtime error, it's a raw string, or it's not of
	// type `error`. In which case, do panic.
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			} else if s, ok := r.(string); ok {
				panic(s)
			} else if e, ok := r.(error); !ok {
				panic(r)
			} else {
				err = e
			}
		}
	}()

	ret = make([]map[string]interface{}, srcv.Len())
	for i := 0; i < srcv.Len(); i++ {
		elemv := srcv.Index(i)
		if elemv.Kind() == reflect.Ptr || elemv.Kind() == reflect.Interface {
			elemv = elemv.Elem()
		}
		ret[i] = (lookupEncodeFn(elemv.Type(), cfg)(elemv, cfg)).(map[string]interface{})
	}
	return ret, nil
}

type encodeFn func(src reflect.Value, cfg *Config) interface{}

type encoderFnCacheKey struct {
	t reflect.Type
	c Config
}

// `encodeFnCache` is based on encode/json's encoderCache. It stores the given
// type's encodeFn s.t. the construction of new encodeFn wrappers need only
// happen once.
var encodeFnCache sync.Map // map[encoderFnCacheKey]encodeFn

func lookupEncodeFn(t reflect.Type, cfg *Config) encodeFn {
	key := encoderFnCacheKey{t, *cfg}
	// Early-out on quick cache-hits.
	if fn, ok := encodeFnCache.Load(key); ok {
		return fn.(encodeFn)
	}

	// From encoding/json/encode.go@typeEncoder, modified;
	// > To deal with recursive types, populate the map with an indirect func
	// > before we build [a real encoding function]. This [function] waits on
	// > the real func (f) to be ready and then calls it. This indirect func is
	// > only used for recursive types.
	var (
		wg sync.WaitGroup
		fn encodeFn
	)
	wg.Add(1)
	fi, loaded := encodeFnCache.LoadOrStore(
		key,
		encodeFn(func(src reflect.Value, cfg *Config) interface{} {
			wg.Wait()
			return fn(src, cfg)
		}),
	)
	if loaded {
		// This is *not* a new type and the correct encodeFn has already
		// been stored; return that.
		return fi.(encodeFn)
	}

	// This type does not have a correct encodeFn loaded into the cache;
	// find/construct the correct encoder and replace the indirect fn.
	fn = newEncodeValueFn(t, cfg, true)
	wg.Done()
	encodeFnCache.Store(key, fn)
	return fn
}

func newEncodeValueFn(t reflect.Type, cfg *Config, firstPass bool) encodeFn {
	if t.Implements(marshalerType) {
		return encodeMarshaller
	}
	if firstPass && t.Kind() != reflect.Ptr && reflect.PtrTo(t).Implements(marshalerType) {
		return newConditionalEncoder(
			func(src reflect.Value) bool {
				return src.CanAddr()
			},
			encodeAddrMarshaller,
			newEncodeValueFn(t, cfg, false),
		)
	}
	switch t.Kind() {
	case reflect.Struct:
		return newStructEncoder(t, cfg)
	default:
		// We assume that if the type is non-nilable, and not a struct, we can
		// just return an enclosing interface{}, and call it good.
		return encodeInterface
	}
}

func encodeInterface(src reflect.Value, cfg *Config) interface{} {
	if !src.CanInterface() {
		panic(errors.New("How did you get here with a non-interfaceable value?"))
	}
	return src.Interface()
}

func encodeMarshaller(src reflect.Value, cfg *Config) interface{} {
	if src.Kind() == reflect.Ptr && src.IsNil() {
		return nil
	}
	m, ok := src.Interface().(Marshaler)
	if !ok {
		panic(errors.New("How did you get here w/o an enc_map.Marshaler?"))
	}
	ret, err := m.MarshalMapValue()
	if err != nil {
		panic(err)
	}
	return ret
}

func encodeAddrMarshaller(src reflect.Value, cfg *Config) interface{} {
	srca := src.Addr()
	if srca.IsNil() {
		return nil
	}
	m, ok := srca.Interface().(Marshaler)
	if !ok {
		panic(errors.New("How did you get here w/o a pointer-to enc_map.Marshaler?"))
	}
	ret, err := m.MarshalMapValue()
	if err != nil {
		panic(err)
	}
	return ret
}

type condEncoder struct {
	cond func(src reflect.Value) bool
	tru  encodeFn
	fls  encodeFn
}

func (cm *condEncoder) marshalValue(src reflect.Value, cfg *Config) interface{} {
	if cm.cond(src) {
		return cm.tru(src, cfg)
	}
	return cm.fls(src, cfg)
}

func newConditionalEncoder(cond func(src reflect.Value) bool, tru encodeFn, fls encodeFn) encodeFn {
	cm := &condEncoder{cond, tru, fls}
	return cm.marshalValue
}

func fieldByIndex(v reflect.Value, index []int) reflect.Value {
	for _, i := range index {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return reflect.Value{}
			}
			v = v.Elem()
		}
		v = v.Field(i)
	}
	return v
}

func typeByIndex(t reflect.Type, index []int) reflect.Type {
	for _, i := range index {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		t = t.Field(i).Type
	}
	return t
}

type structEncoder struct {
	fields    []field
	fieldEncs []encodeFn
}

func (se *structEncoder) encode(src reflect.Value, cfg *Config) interface{} {
	ret := make(map[string]interface{}, len(se.fields))
	for i, f := range se.fields {
		fv := fieldByIndex(src, f.index)
		if !fv.IsValid() ||
			(f.options.Contains("omitZero") && encoding.IsValueZero(fv)) ||
			(f.options.Contains("omitNil") && encoding.IsValueNil(fv)) {
			continue
		}
		if !src.CanInterface() {
			panic(fmt.Errorf("How did you get here with a non-interfaceable value?"))
		}
		ret[f.name] = se.fieldEncs[i](fv, cfg)
	}
	return ret
}

func newStructEncoder(t reflect.Type, cfg *Config) encodeFn {
	fields := cachedTypeFields(t, cfg)
	se := structEncoder{
		fields:    fields,
		fieldEncs: make([]encodeFn, len(fields)),
	}
	for i, f := range fields {
		if f.options.Contains("value") {
			se.fieldEncs[i] = encodeInterface
		} else {
			se.fieldEncs[i] = lookupEncodeFn(typeByIndex(t, f.index), cfg)
		}
	}
	return se.encode
}
