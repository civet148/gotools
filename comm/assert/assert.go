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

// judgement: bool, string, struct, slice, map is nil or false?
func IsNilOrFalse(v interface{}) bool {
	switch v.(type) {
	case string:
		if v.(string) == "" {
			return true
		}
	case bool:
		return !v.(bool)
	default:
		{
			val := reflect.ValueOf(v)
			return val.IsNil()
		}
	}
	return false
}
