package numerics

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type IntType int64

func init() {

}

//for printf %s %v %+v
func (t IntType) String() (name string) {
	return fmt.Sprintf("%d", t)
}

//for printf %#v
func (t IntType) GoString() (name string) {
	return t.String()
}

func (t IntType) ToInt() int {
	return int(t)
}

func (t IntType) ToUint() uint {
	return uint(t)
}

func (t IntType) ToInt8() int8 {
	return int8(t)
}

func (t IntType) ToInt16() int16 {
	return int16(t)
}

func (t IntType) ToInt32() int32 {
	return int32(t)
}

func (t IntType) ToInt64() int64 {
	return int64(t)
}

func (t IntType) ToUnt8() uint8 {
	return uint8(t)
}

func (t IntType) ToUint16() uint16 {
	return uint16(t)
}

func (t IntType) ToUint32() uint32 {
	return uint32(t)
}

func (t IntType) ToUint64() uint64 {
	return uint64(t)
}

func (t IntType) ToString() string {
	return t.String()
}

func (t *IntType) FromInt(v int) {
	*t = IntType(v)
}

func (t *IntType) FromInt8(v int8) {
	*t = IntType(v)
}

func (t *IntType) FromInt16(v int16) {
	*t = IntType(v)
}

func (t *IntType) FromInt32(v int32) {
	*t = IntType(v)
}

func (t *IntType) FromInt64(v int64) {
	*t = IntType(v)
}

func (t *IntType) FromUint(v uint) {
	*t = IntType(v)
}

func (t *IntType) FromUint8(v uint8) {
	*t = IntType(v)
}

func (t *IntType) FromUint16(v uint16) {
	*t = IntType(v)
}

func (t *IntType) FromUint32(v uint32) {
	*t = IntType(v)
}

func (t *IntType) FromUint64(v uint64) {
	*t = IntType(v)
}

func (t *IntType) FromString(v string) {
	i64, _ := strconv.ParseInt(v, 10, 64)
	t.FromInt64(i64)
}

// Scan implements the sql.Scanner interface for database deserialization.
func (t *IntType) Scan(value interface{}) error {
	switch v := value.(type) {
	case int:
		t.FromInt(v)
	case int8:
		t.FromInt8(v)
	case int16:
		t.FromInt16(v)
	case int32:
		t.FromInt32(v)
	case int64:
		t.FromInt64(v)
	case uint:
		t.FromUint(v)
	case uint8:
		t.FromUint8(v)
	case uint16:
		t.FromUint16(v)
	case uint32:
		t.FromUint32(v)
	case uint64:
		t.FromUint64(v)
	case string:
		t.FromString(v)
	default:
		return fmt.Errorf("this type not support yet")
	}
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (t IntType) Value() (driver.Value, error) {
	return t.String(), nil
}
