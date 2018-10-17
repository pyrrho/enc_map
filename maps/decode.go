package maps

import (
	"errors"
	"reflect"
)

func Unmarshal(src interface{}, v interface{}) error {
	err := defaultConfig.unmarshal(src, v)
	if err != nil {
		return err
	}
	return nil
}

type Unmarshaler interface {
	UnmarshalMapValue() error
}

var unmarshalerType = reflect.TypeOf(new(Unmarshaler)).Elem()

func (cfg *Config) Unmarshal(src interface{}, v interface{}) error {
	err := cfg.unmarshal(src, v)
	if err != nil {
		return err
	}
	return nil
}

func (cfg *Config) unmarshal(src interface{}, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return errors.New("encoding/maps: cannot unmarshal into non-pointer")
	} else if rv.IsNil() {
		return errors.New("encoding/maps: cannot unmarshal into nil pointer")
	}
	return errors.New("enc_map.Scan has not yet been implemented")
}
