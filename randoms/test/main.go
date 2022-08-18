package main

import (
	"fmt"
	"github.com/civet148/gotools/randoms"
)

func main() {
	fmt.Printf("%s\n", randoms.RandomAlphaOrNumeric(64, true, true))
}
