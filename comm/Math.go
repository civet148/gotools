package comm

import "math"

//截取浮点数 add by dragon
func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

const MIN = 0.000001
// MIN 为用户自定义的比较精度
func FloatIsEqual(f1, f2 float64) bool {
	return math.Dim(f1, f2) < MIN
}
