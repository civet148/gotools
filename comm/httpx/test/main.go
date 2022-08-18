package main

import (
	"github.com/civet148/gotools/comm/httpx"
	"github.com/civet148/gotools/log"
)

func main() {

	var strUrlGet = "https://www.baidu.com"
	log.Info("start http connection....")
	c := httpx.NewHttpClient(0)
	c.Header().SetAuthorization("123456").SetApplicationJson()
	resp, err := c.Get(strUrlGet, nil)
	if err != nil {
		log.Error("access [%v] error [%v]", strUrlGet, err.Error())
		return
	}
	c.Header().RemoveAll()
	log.Info("response [%+v]", resp)
}
