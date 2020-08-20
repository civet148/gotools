package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/numerics"
)

func main() {
	var ints = []uint32{123, 124, 125, 126}
	var floats = []float32{123.1, 124.2, 125.3, 126.4}

	log.Infof("numeric int join [%+v]", numerics.Join(ints, ","))
	log.Infof("numeric float join [%+v]", numerics.Join(floats, ","))
}
