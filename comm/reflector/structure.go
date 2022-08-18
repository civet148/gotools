package reflector

import (
	"fmt"
	"github.com/civet148/gotools/comm/assert"
	"reflect"
)

type Structure struct {
	value interface{}       //value
	dict  map[string]string //dictionary of structure tag and value
}

// parse struct tag and value to map
func (s *Structure) ToMap(tagName string) (m map[string]string) {

	typ := reflect.TypeOf(s.value)
	val := reflect.ValueOf(s.value)

	if typ.Kind() == reflect.Ptr { // pointer type
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() == reflect.Struct { // struct type

		s.parseStructField(typ, val, tagName)
	} else {
		assert.Panic(false, "not a struct object")
	}
	return s.dict
}

// map[string]string values assign to struct
func (s *Structure) AssignFromMap(v interface{}) (err error) {
	typ := reflect.TypeOf(s.value)
	val := reflect.ValueOf(s.value)

	if typ.Kind() == reflect.Ptr { // pointer type
		typ = typ.Elem()
		val = val.Elem()
	}

	switch typ.Kind() {
	case reflect.Map:
		{
			s.assignFromMapStringString(v.(map[string]string)) // warning: v must be a map[string]string type object, otherwise it will crash
		}
	default:
		err = fmt.Errorf("unsupport type [%v]", typ.Name())
		break
	}
	return
}

func (s *Structure) getTag(sf reflect.StructField, tagName string) string {

	return sf.Tag.Get(tagName)
}

// parse struct fields
func (s *Structure) parseStructField(typ reflect.Type, val reflect.Value, tagName string) {

	kind := typ.Kind()
	if kind == reflect.Struct {
		NumField := val.NumField()
		//log.Debug("Struct [%s] NumField = [%d]", typ.Name(), NumField)
		for i := 0; i < NumField; i++ {
			typField := typ.Field(i)
			valField := val.Field(i)

			if typField.Type.Kind() == reflect.Ptr {
				typField.Type = typField.Type.Elem()
				valField = valField.Elem()
			}
			if !valField.IsValid() || !valField.CanInterface() {
				//log.Debug("error: filed [%v] tag(%v) is not valid ", typField.Type.Name(), s.getTag(typField, tagName))
				continue
			}
			if typField.Type.Kind() == reflect.Struct {

				s.parseStructField(typField.Type, valField, tagName) //结构体需要递归调用
			} else {
				s.setValueByField(typField, valField, tagName) //按标签名赋值给结构体成员变量
			}
		}
	}
}

func (s *Structure) setValueByField(field reflect.StructField, val reflect.Value, tagName string) {

	tag := s.getTag(field, tagName)
	if tag != "" {

		s.dict[field.Name] = fmt.Sprintf("%v", val.Interface())
		//log.Debug("struct tag [%s] value [%v] save to map ok", tag, s.dict[field.Name])
	}
}

func (s *Structure) assignFromMapStringString(vm map[string]string) {
	//TODO @lory assignFromMapStringString
}
