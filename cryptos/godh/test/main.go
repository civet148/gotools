package main

import (
	"fmt"
	"github.com/civet148/gotools/cryptos/godh"
)

func main() {

	dhUserA := godh.NewCryptoDH()
	dhUserB := godh.NewCryptoDH()

	dhShareKeyA := dhUserA.ScalarMultBase64(dhUserB.GetPublicKey())
	dhShareKeyB := dhUserB.ScalarMultBase64(dhUserA.GetPublicKey())
	fmt.Printf("User A share key [%v] \n", dhShareKeyA)
	fmt.Printf("User B share key [%v] \n", dhShareKeyB)
}
