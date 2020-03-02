package reflector

func Struct(v interface{}) *Structure {

	return &Structure{
		value: v,
		dict:  make(map[string]string),
	}
}
