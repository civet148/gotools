package comm


func AppendSlice(slice []interface{}, args...interface{}) []interface{} {
	for _, v := range args {
		slice = append(slice, v)
	}
	return slice
}