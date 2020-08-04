package main

import (
	. "github.com/civet148/gotools/cryptos/goaes"
	_ "github.com/civet148/gotools/cryptos/goaes/cbc"    //注册CBC加解密对象创建方法
	_ "github.com/civet148/gotools/cryptos/goaes/cfb"    //注册CFB加解密对象创建方法
	_ "github.com/civet148/gotools/cryptos/goaes/ctr128" //注册CTR128加解密对象创建方法
	_ "github.com/civet148/gotools/cryptos/goaes/ecb"    //注册ECB加解密对象创建方法
	_ "github.com/civet148/gotools/cryptos/goaes/ige256" //注册IGE256加解密对象创建方法
	_ "github.com/civet148/gotools/cryptos/goaes/ofb"    //注册OFB加解密对象创建方法
	"github.com/civet148/gotools/log"
)

var strKey16 = "1234567890123456"                           //加密KEY(16字节)
var strKey24 = "123456789012345678901234"                   //加密KEY(24字节)
var strKey32 = "12345678901234561234567890123456"           //加密KEY(32字节)
var strIV = "1234567890123456"                              //加密向量(固定16字节)
var strText = "wallet  RUNNING   pid 13027, uptime 0:00:15" //测试数据

func main() {
	AES_CBC()
	//AES_CFB()
	//AES_ECB()
	//AES_OFB()
	//AES_CTR128()
	//AES_IGE256()
}

func AES_CBC() {

	aes := NewCryptoAES(AES_Type_CBC, []byte(strKey32), []byte(strIV))
	enc, _ := aes.EncryptBase64([]byte(strText))
	log.Infof("AES CBC text [%v] encrypt -> [%v]", strText, enc)
	dec, _ := aes.DecryptBase64(enc)
	log.Infof("AES CBC base [%v] decrypt -> [%v]", enc, string(dec))
}

func AES_CFB() {
	aes := NewCryptoAES(AES_Type_CFB, []byte(strKey32), []byte(strIV))
	enc, _ := aes.EncryptBase64([]byte(strText))
	log.Infof("AES CFB text [%v] encrypt -> [%v]", strText, enc)
	dec, _ := aes.DecryptBase64(enc)
	log.Infof("AES CFB base [%v] decrypt -> [%v]", enc, string(dec))
}

func AES_ECB() {
	aes := NewCryptoAES(AES_Type_ECB, []byte(strKey32), []byte(strIV))
	enc, _ := aes.EncryptBase64([]byte(strText))
	log.Infof("AES ECB text [%v] encrypt -> [%v]", strText, enc)
	dec, _ := aes.DecryptBase64(enc)
	log.Infof("AES ECB base [%v] decrypt -> [%v]", enc, string(dec))
}

func AES_OFB() {

	aes := NewCryptoAES(AES_Type_OFB, []byte(strKey32), []byte(strIV))
	enc, _ := aes.EncryptBase64([]byte(strText))
	log.Infof("AES OFB text [%v] encrypt -> [%v]", strText, enc)
	dec, _ := aes.DecryptBase64(enc)
	log.Infof("AES OFB base [%v] decrypt -> [%v]", enc, string(dec))
}

func AES_CTR128() {

	aes := NewCryptoAES(AES_Type_CTR128, []byte(strKey16), []byte(strIV))
	enc, _ := aes.EncryptBase64([]byte(strText))
	log.Infof("AES CTR text [%v] encrypt -> [%v]", strText, enc)
	dec, _ := aes.DecryptBase64(enc)
	log.Infof("AES CTR base [%v] decrypt -> [%v]", enc, string(dec))
}

func AES_IGE256() {
	aes := NewCryptoAES(AES_Type_IGE256, []byte(strKey32), []byte(strIV))
	enc, _ := aes.EncryptBase64([]byte(strText))
	log.Infof("AES IGE text [%v] encrypt -> [%v]", strText, enc)
	dec, _ := aes.DecryptBase64(enc)
	log.Infof("AES IGE base [%v] decrypt -> [%v]", enc, string(dec))
}
