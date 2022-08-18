package hmacsha

import (
	"fmt"
	"testing"
)

func TestHmacSha256(t *testing.T) {
	strOut := HmacSha256Hex("symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000&timestamp=1499827319559",
		"NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j")
	fmt.Printf("%s\n", strOut)
}
