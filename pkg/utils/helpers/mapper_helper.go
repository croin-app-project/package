package helpers

import (
	"fmt"
	"reflect"
)

func MapStructs(src interface{}, dest interface{}) (interface{}, error) {
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest)

	if srcValue.Kind() != reflect.Struct || destValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Both src and dest should be structs")
	}

	if srcValue.Type() != destValue.Type() {
		return nil, fmt.Errorf("Source and destination structs should have the same type")
	}

	destPtr := reflect.New(destValue.Type()).Elem()
	destPtr.Set(srcValue)

	return destPtr.Interface(), nil
}
