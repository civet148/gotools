package main

import (
	"fmt"
	"github.com/civet148/gotools/cryptos/gocrc"
)

var data = []byte{0x01, 0x02, 0x0e, 0x0f, 0x80, 0x11, 0x55}

func main() {
	cs32 := gocrc.CheckSum32(data)
	cs64ECMA := gocrc.CheckSum64ECMA(data)
	cs64ISO := gocrc.CheckSum64ISO(data)
	fmt.Printf("CRC32: [%v] CRC64: ECMA [%v] ISO [%v]\n", cs32, cs64ECMA, cs64ISO)
}
