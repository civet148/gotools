package endian

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

//对数据结构进行序列化(网络字节序)
func EncodeEndian(bLittleEndian bool, v interface{}) (data []byte, err error) {

	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	var buffer = bytes.NewBuffer(nil)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	switch typ.Kind() {
	case reflect.Struct: //结构体
		err = encodeEndianStruct(bLittleEndian, buffer, typ, val)
	default: //其他类型
		err = encodeEndianBaseType(bLittleEndian, buffer, typ, val)
	}
	return buffer.Bytes(), err
}

//将网络数据反序列化到结构体(网络字节序)
func DecodeEndian(bLittleEndian bool, data []byte, v interface{}) (err error) {

	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	var buffer = bytes.NewBuffer(data)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	switch typ.Kind() {
	case reflect.Struct: //结构体
		err = decodeEndianStruct(bLittleEndian, buffer, typ, val)
	default: //其他类型
		err = decodeEndianBaseType(bLittleEndian, buffer, typ, val)
	}
	return err
}

//遍历结构体，按大端/小端构造网络字节序数组
func encodeEndianStruct(bLittleEndian bool, buf io.Writer, typ reflect.Type, val reflect.Value) (err error) {

	kind := typ.Kind()
	if kind == reflect.Struct {

		Fields := val.NumField()
		for i := 0; i < Fields; i++ {

			typField := typ.Field(i)
			valField := val.Field(i)

			if typField.Type.Kind() == reflect.Ptr {
				typField.Type = typField.Type.Elem()
				valField = valField.Elem()
			}

			if typField.Type.Kind() == reflect.Struct {
				err = encodeEndianStruct(bLittleEndian, buf, typField.Type, valField)
			} else {
				err = encodeEndianBaseType(bLittleEndian, buf, typField.Type, valField)
			}
		}
	}

	return
}

//遍历结构体，将网络数据反序列化到结构体
func decodeEndianStruct(bLittleEndian bool, buf io.Reader, typ reflect.Type, val reflect.Value) (err error) {

	kind := typ.Kind()
	if kind == reflect.Struct {

		Fields := val.NumField()
		for i := 0; i < Fields; i++ {

			typField := typ.Field(i)
			valField := val.Field(i)

			if typField.Type.Kind() == reflect.Ptr {
				typField.Type = typField.Type.Elem()
				valField = valField.Elem()
			}

			if typField.Type.Kind() == reflect.Struct {
				err = decodeEndianStruct(bLittleEndian, buf, typField.Type, valField)
			} else {
				err = decodeEndianBaseType(bLittleEndian, buf, typField.Type, valField)
			}
		}
	}

	return
}

//遍历结构体成员基本类型，按大端/小端构造网络字节序数组
func encodeEndianBaseType(bLittleEndian bool, buf io.Writer, typ reflect.Type, val reflect.Value) (err error) {

	var bo binary.ByteOrder = binary.BigEndian
	//fmt.Printf("Struct field type [%v] kind [%v]\n", typ.Name(), typ.Kind())
	if bLittleEndian {
		bo = binary.LittleEndian
	}
	switch typ.Kind() {

	case reflect.Uint, reflect.Int:
		{
			err = fmt.Errorf("不支持的字段类型type [%v] kind [int/uint] 请使用int32/uint32类型\n", typ.Name())
			panic(err)
		}
	case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
		reflect.Int32, reflect.Uint32, reflect.Float32, reflect.Int64, reflect.Uint64,
		reflect.Float64, reflect.Complex64, reflect.Complex128: //数值类型
		{
			binary.Write(buf, bo, val.Interface())
		}
	case reflect.Array: //数组(单独处理)
		{
			var bs []byte
			for i := 0; i < val.Len(); i++ {
				v := val.Index(i).Interface()
				bs = append(bs, v.(uint8))
			}
			buf.Write(bs)
		}
	default: //不支持的类型
		{
			err = fmt.Errorf("不支持的字段类型type [%v] kind [%v] ", typ.Name(), typ.Kind())
			panic(err)
		}
	}
	return
}

//遍历结构体成员基本类型，按大端/小端构造网络字节序数组
func decodeEndianBaseType(bLittleEndian bool, buf io.Reader, typ reflect.Type, val reflect.Value) (err error) {

	var bo binary.ByteOrder = binary.BigEndian

	if bLittleEndian {
		bo = binary.LittleEndian
	}
	switch typ.Kind() {

	case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
		reflect.Int, reflect.Int32, reflect.Uint32, reflect.Float32, reflect.Int64, reflect.Uint64,
		reflect.Float64, reflect.Complex64, reflect.Complex128: //数值类型
		{
			binary.Read(buf, bo, val.Addr().Interface())
		}
	case reflect.Array: //数组(单独处理)
		{
			var bs []byte
			bs = make([]byte, val.Len())
			buf.Read(bs) //从网络数据读取bs切片长度的内容
			for i := 0; i < val.Len(); i++ {
				v := reflect.ValueOf(bs[i])
				val.Index(i).Set(v)
			}
		}
	default: //不支持的类型
		{
			return fmt.Errorf("invalid type [%v] to convert to endian binary", typ.Name())
		}
	}

	return
}
