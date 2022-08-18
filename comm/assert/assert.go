package assert

import (
	"fmt"
	"reflect"
)

// assert bool, string or struct/slice/map nil or false, call panic
func Panic(v interface{}, strMsg string, args ...interface{}) {
	if IsNilOrFalse(v) {
		panic(fmt.Sprintf(strMsg, args...))
	}
}

// judgement: variant is a pointer type ?
func IsPtrType(v interface{}) bool {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		return true
	}
	return false
}

// judgement: bool, integer, string, struct, slice, map is nil or false?
func IsNilOrFalse(v interface{}) bool {
	switch v.(type) {
	case string:
		if v.(string) == "" {
			return true
		}
	case bool:
		return !v.(bool)
	case int8, int16, int, int32, int64, uint8, uint16, uint, uint32, uint64:
		{
			if fmt.Sprintf("%v", v) == "0" {
				return true
			}
		}
	default:
		{
			val := reflect.ValueOf(v)
			return val.IsNil()
		}
	}
	return false
}
