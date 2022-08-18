package main

import (
	"github.com/civet148/gotools/comm/codec"
	"github.com/civet148/gotools/log"
)

func main() {

	ParseSpecialUrl()
}

func ParseSpecialUrl() {
	//URL have some special characters in password
	//strUrl := "https://root:`~!@#$%^&*()-_=+@127.0.0.1:8082/my_path/abc?input=yes&console=no#golang"
	//strUrl := "mysql://root:123456@127.0.0.1:3306/mydb?"
	strUrl := "http://127.0.0.1:8080"
	u := codec.ParseUrl(strUrl)
	log.Debugf("url [%+v]", u)
}
