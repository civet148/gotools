package structs

import "reflect"

func GetStructName(v interface{}) (strName string) {

	typ := reflect.TypeOf(v)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ.Name()
}
