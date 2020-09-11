package main

import (
	"fmt"
	"github.com/civet148/gotools/cryptos/gomd5"
	"os"
)

func main() {

	var strFile string
	if len(os.Args) != 2 {
		fmt.Println("please input file path")
		return
	}

	strFile = os.Args[1]
	strMd5 := gomd5.Md5File(strFile)
	fmt.Printf("Md5File -> [%v]\n", strMd5)
}
