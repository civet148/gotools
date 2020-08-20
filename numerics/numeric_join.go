package numerics

import (
	"fmt"
	"reflect"
	"strings"
)

//ints -> int8/int16/int32/int/int64/uint8/uint16/uint32/uint/uint64/float32/float64 slice
func Join(ints interface{}, sep string) (strJoin string) {

	var joins []string
	typ := reflect.TypeOf(ints)
	val := reflect.ValueOf(ints)

	kind := typ.Kind()

	if kind == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	switch typ.Kind() {
	case reflect.Slice:
		{
			for i := 0; i < val.Len(); i++ {
				joins = append(joins, fmt.Sprintf("%v", val.Index(i).Interface()))
			}
		}
	default:
		fmt.Println("WARN: Join function need a slice type")
	}

	return strings.Join(joins, sep)
}
