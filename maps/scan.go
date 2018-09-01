package maps

import (
	"errors"
	"reflect"
)

func Scan(src map[string]interface{}, dest interface{}) error {
	return errors.New("enc_map.Scan has not yet been implemented")
}

type Scanner interface {
	ScanMap(src interface{}) error
}

var (
	scannerType = reflect.TypeOf(new(Scanner)).Elem()
)
