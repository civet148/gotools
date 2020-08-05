package main

import (
	. "github.com/civet148/gotools/cryptos/goaes"
	_ "github.com/civet148/gotools/cryptos/goaes/cbc" //注册CBC加解密对象创建方法
	_ "github.com/civet148/gotools/cryptos/goaes/cfb" //注册CFB加解密对象创建方法
	_ "github.com/civet148/gotools/cryptos/goaes/ctr" //注册CTR128加解密对象创建方法
	_ "github.com/civet148/gotools/cryptos/goaes/ecb" //注册ECB加解密对象创建方法
	_ "github.com/civet148/gotools/cryptos/goaes/ofb" //注册OFB加解密对象创建方法
	"github.com/civet148/gotools/log"
)

var strKey16 = "1234567890123456"                           //加密KEY(16字节)
var strKey24 = "123456789012345678901234"                   //加密KEY(24字节)
var strKey32 = "12345678901234561234567890123456"           //加密KEY(32字节)
var strIV = "1234567890123456"                              //加密向量(固定16字节)
var strText = "wallet  RUNNING   pid 13027, uptime 0:00:15" //测试数据

func main() {

	var modes = []AES_Mode{
		AES_Mode_CBC,
		AES_Mode_CFB,
		AES_Mode_ECB,
		AES_Mode_OFB,
		AES_Mode_CTR,
	}

	for _, v := range modes {

		aes := NewCryptoAES(v, []byte(strKey32), []byte(strIV))
		enc, err := aes.EncryptBase64([]byte(strText))
		if err != nil {
			log.Errorf("[%v] encrypt to base64 error [%v]", aes.GetMode(), err.Error())
			continue
		}
		log.Infof("[%v] text [%v] encrypt -> [%v]", aes.GetMode(), strText, enc)

		dec, err := aes.DecryptBase64(enc)
		if err != nil {
			log.Errorf("[%v] decrypt from base64 error [%v]", aes.GetMode(), err.Error())
			continue
		}
		log.Infof("[%v] base [%v] decrypt -> [%v]", aes.GetMode(), enc, string(dec))
	}
}
