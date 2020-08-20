package main

import (
	"fmt"
	"github.com/civet148/gotools/cryptos/gomd5"
	"github.com/civet148/gotools/process"
)

var data = "12312431243543645666666666666666666666354345"

func main() {
	gopath := process.GetEnv("GOPATH")
	if gopath != "" {
		strMd5 := gomd5.Md5File(fmt.Sprintf("%s/src/github.com/civet148/gotools/cryptos/gomd5/gomd5.go", gopath))
		fmt.Printf("Md5File -> [%v]\n", strMd5)
	}
	strSum := gomd5.Md5Sum([]byte(data))
	fmt.Printf("Md5Sum  -> [%v]\n", strSum)
}
