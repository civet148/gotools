package main

import (
	"fmt"
	"github.com/civet148/gotools/numerics"
)

func main() {
	var ints = []uint32{123, 124, 125, 126}
	var floats = []float32{123.1, 124.2, 125.3, 126.4}

	fmt.Printf("numeric int join [%+v] \n", numerics.Join(ints, ","))
	fmt.Printf("numeric float join [%+v] \n", numerics.Join(floats, ","))
}
