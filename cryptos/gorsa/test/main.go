package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/civet148/gotools/cryptos/gorsa"
)

func init() {

}

var Pubkey = `-----BEGIN 公钥-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAk+89V7vpOj1rG6bTAKYM
56qmFLwNCBVDJ3MltVVtxVUUByqc5b6u909MmmrLBqS//PWC6zc3wZzU1+ayh8xb
UAEZuA3EjlPHIaFIVIz04RaW10+1xnby/RQE23tDqsv9a2jv/axjE/27b62nzvCW
eItu1kNQ3MGdcuqKjke+LKhQ7nWPRCOd/ffVqSuRvG0YfUEkOz/6UpsPr6vrI331
hWRB4DlYy8qFUmDsyvvExe4NjZWblXCqkEXRRAhi2SQRCl3teGuIHtDUxCskRIDi
aMD+Qt2Yp+Vvbz6hUiqIWSIH1BoHJer/JOq2/O6X3cmuppU4AdVNgy8Bq236iXvr
MQIDAQAB
-----END 公钥-----
`

var Pirvatekey = `-----BEGIN 私钥-----
MIIEpAIBAAKCAQEAk+89V7vpOj1rG6bTAKYM56qmFLwNCBVDJ3MltVVtxVUUByqc
5b6u909MmmrLBqS//PWC6zc3wZzU1+ayh8xbUAEZuA3EjlPHIaFIVIz04RaW10+1
xnby/RQE23tDqsv9a2jv/axjE/27b62nzvCWeItu1kNQ3MGdcuqKjke+LKhQ7nWP
RCOd/ffVqSuRvG0YfUEkOz/6UpsPr6vrI331hWRB4DlYy8qFUmDsyvvExe4NjZWb
lXCqkEXRRAhi2SQRCl3teGuIHtDUxCskRIDiaMD+Qt2Yp+Vvbz6hUiqIWSIH1BoH
Jer/JOq2/O6X3cmuppU4AdVNgy8Bq236iXvrMQIDAQABAoIBAQCCbxZvHMfvCeg+
YUD5+W63dMcq0QPMdLLZPbWpxMEclH8sMm5UQ2SRueGY5UBNg0WkC/R64BzRIS6p
jkcrZQu95rp+heUgeM3C4SmdIwtmyzwEa8uiSY7Fhbkiq/Rly6aN5eB0kmJpZfa1
6S9kTszdTFNVp9TMUAo7IIE6IheT1x0WcX7aOWVqp9MDXBHV5T0Tvt8vFrPTldFg
IuK45t3tr83tDcx53uC8cL5Ui8leWQjPh4BgdhJ3/MGTDWg+LW2vlAb4x+aLcDJM
CH6Rcb1b8hs9iLTDkdVw9KirYQH5mbACXZyDEaqj1I2KamJIU2qDuTnKxNoc96HY
2XMuSndhAoGBAMPwJuPuZqioJfNyS99x++ZTcVVwGRAbEvTvh6jPSGA0k3cYKgWR
NnssMkHBzZa0p3/NmSwWc7LiL8whEFUDAp2ntvfPVJ19Xvm71gNUyCQ/hojqIAXy
tsNT1gBUTCMtFZmAkUsjqdM/hUnJMM9zH+w4lt5QM2y/YkCThoI65BVbAoGBAMFI
GsIbnJDNhVap7HfWcYmGOlWgEEEchG6Uq6Lbai9T8c7xMSFc6DQiNMmQUAlgDaMV
b6izPK4KGQaXMFt5h7hekZgkbxCKBd9xsLM72bWhM/nd/HkZdHQqrNAPFhY6/S8C
IjRnRfdhsjBIA8K73yiUCsQlHAauGfPzdHET8ktjAoGAQdxeZi1DapuirhMUN9Zr
kr8nkE1uz0AafiRpmC+cp2Hk05pWvapTAtIXTo0jWu38g3QLcYtWdqGa6WWPxNOP
NIkkcmXJjmqO2yjtRg9gevazdSAlhXpRPpTWkSPEt+o2oXNa40PomK54UhYDhyeu
akuXQsD4mCw4jXZJN0suUZMCgYAgzpBcKjulCH19fFI69RdIdJQqPIUFyEViT7Hi
bsPTTLham+3u78oqLzQukmRDcx5ddCIDzIicMfKVf8whertivAqSfHytnf/pMW8A
vUPy5G3iF5/nHj76CNRUbHsfQtv+wqnzoyPpHZgVQeQBhcoXJSm+qV3cdGjLU6OM
HgqeaQKBgQCnmL5SX7GSAeB0rSNugPp2GezAQj0H4OCc8kNrHK8RUvXIU9B2zKA2
z/QUKFb1gIGcKxYr+LqQ25/+TGvINjuf6P3fVkHL0U8jOG0IqpPJXO3Vl9B8ewWL
cFQVB/nQfmaMa4ChK0QEUe+Mqi++MwgYbRHx1lIOXEfUJO+PXrMekw==
-----END 私钥-----
`

func main() {

	TestPriKeyEncryptAndPubKeyDecrypt()       //私钥加密、公钥解密
	TestPubKeyEncryptAndPriKeyDecrypt()       //公钥解密、私钥加密
	TestPriKeyEncryptAndPubKeyDecryptByFile() //私钥证书加密、公钥证书解密
	TestPubKeyEncryptAndPriKeyDecryptByFile() //公钥证书加密、私钥证书解密
}

func TestPriKeyEncryptAndPubKeyDecrypt() {
	var err error
	var prikey *rsa.PrivateKey
	var pubkey *rsa.PublicKey
	var plainText = "12345678900000你好,这真实一个明朗的天气，下面我准备测试一下通过标准公/私钥" +
		"进行RSA加解密，buffered data waiting to be encoded1111111112344444444444444444444444444444444444444" +
		"TGvINjuf6P3fVkHL0U8jOG0IqpPJXO3Vl9B8ewWLo2oXNa40PomK54UhYDhyeuHgqeaQKBgQCnmL5SX7GSAeB0rSNugPp2GezAQj0H4OCc8kNrHK8RUvXIU9B2zKA2"

	fmt.Println("-------- 通过标准私钥加密、标准公钥解密 ---------")
	//读取私钥证书文件test.pfx用于加密数据
	prikey, err = myrsa.GetPriKey([]byte(Pirvatekey))
	if err != nil {
		fmt.Println(err)
		return
	}

	//读取私钥证书文件test.cer用于密数据
	pubkey, err = myrsa.GetPubKey([]byte(Pubkey))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("加密前字符串 [", plainText, "] 长度", len(plainText))
	encData, errEnc := myrsa.PriKeyEncrypt(prikey, []byte(plainText))
	if errEnc != nil {
		fmt.Println(errEnc)
		return
	}

	decData, errDec := myrsa.PubKeyDecrypt(pubkey, encData)
	if errDec != nil {
		fmt.Println(errDec)
		return
	}

	fmt.Println("解密后字符串 [", string(decData), "] 长度", len(decData))
}

func TestPubKeyEncryptAndPriKeyDecrypt() {
	var err error
	var prikey *rsa.PrivateKey
	var pubkey *rsa.PublicKey
	var plainText = "12345678900000你好,这真实一个明朗的天气，下面我准备测试一下标准公/私钥" +
		"进行RSA加解密，buffered data waiting to be encoded1111111112344444444444444444444444444444444444444" +
		"TGvINjuf6P3fVkHL0U8jOG0IqpPJXO3Vl9B8ewWLo2oXNa40PomK54UhYDhyeuHgqeaQKBgQCnmL5SX7GSAeB0rSNugPp2GezAQj0H4OCc8kNrHK8RUvXIU9B2zKA2"

	fmt.Println("-------- 通过标准公钥加密、标准私钥解密 ---------")
	//读取私钥证书文件test.pfx用于加密数据
	prikey, err = myrsa.GetPriKey([]byte(Pirvatekey))
	if err != nil {
		fmt.Println(err)
		return
	}

	//读取私钥证书文件test.cer用于密数据
	pubkey, err = myrsa.GetPubKey([]byte(Pubkey))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("加密前字符串 [", plainText, "] 长度", len(plainText))
	encData, errEnc := myrsa.PubKeyEncrypt(pubkey, []byte(plainText))
	if errEnc != nil {
		fmt.Println(errEnc)
		return
	}

	decData, errDec := myrsa.PriKeyDecrypt(prikey, encData)
	if errDec != nil {
		fmt.Println(errDec)
		return
	}

	fmt.Println("解密后字符串 [", string(decData), "] 长度", len(decData))
}

//通过pfx私钥证书文件加密、cer公钥证书文件解密
func TestPriKeyEncryptAndPubKeyDecryptByFile() {

	var err error
	var prikey *rsa.PrivateKey
	var pubkey *rsa.PublicKey
	var plainText = "12345678900000你好,这真实一个明朗的天气，下面我准备测试一下通过证书" +
		"cer/pfx进行RSA加解密，buffered data waiting to be encoded1111111112344444444444444444444444444444444444444" +
		"TGvINjuf6P3fVkHL0U8jOG0IqpPJXO3Vl9B8ewWLo2oXNa40PomK54UhYDhyeuHgqeaQKBgQCnmL5SX7GSAeB0rSNugPp2GezAQj0H4OCc8kNrHK8RUvXIU9B2zKA2"

	fmt.Println("-------- 通过pfx私钥证书文件加密、cer公钥证书文件解密 ---------")
	//读取私钥证书文件test.pfx用于加密数据
	prikey, err = myrsa.GetPriKeyByPfxFile("test.pfx", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}

	//读取私钥证书文件test.cer用于密数据
	pubkey, err = myrsa.GetPubKeyByCerFile("test.cer")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("加密前字符串 [", plainText, "] 长度", len(plainText))
	encData, errEnc := myrsa.PriKeyEncrypt(prikey, []byte(plainText))
	if errEnc != nil {
		fmt.Println(errEnc)
		return
	}

	decData, errDec := myrsa.PubKeyDecrypt(pubkey, encData)
	if errDec != nil {
		fmt.Println(errDec)
		return
	}

	fmt.Println("解密后字符串 [", string(decData), "] 长度", len(decData))
}

//通过cer公钥证书文件加密、pfx私钥证书文件解密
func TestPubKeyEncryptAndPriKeyDecryptByFile() {

	var err error
	var prikey *rsa.PrivateKey
	var pubkey *rsa.PublicKey
	var plainText = "12345678900000你好,这真实一个明朗的天气，下面我准备测试一下通过证书" +
		"cer/pfx进行RSA加解密，buffered data waiting to be encoded1111111112344444444444444444444444444444444444444" +
		"TGvINjuf6P3fVkHL0U8jOG0IqpPJXO3Vl9B8ewWLo2oXNa40PomK54UhYDhyeuHgqeaQKBgQCnmL5SX7GSAeB0rSNugPp2GezAQj0H4OCc8kNrHK8RUvXIU9B2zKA2"

	fmt.Println("-------- 通过cer公钥证书文件加密、pfx私钥证书文件解密 ---------")

	//读取私钥证书文件test.pfx用于加密数据
	prikey, err = myrsa.GetPriKeyByPfxFile("test.pfx", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}

	//读取私钥证书文件test.cer用于密数据
	pubkey, err = myrsa.GetPubKeyByCerFile("test.cer")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("加密前字符串 [", plainText, "] 长度", len(plainText))
	encData, errEnc := myrsa.PubKeyEncrypt(pubkey, []byte(plainText))
	if errEnc != nil {
		fmt.Println(errEnc)
		return
	}
	//fmt.Println("加密后二进制数据 ", encData)

	decData, errDec := myrsa.PriKeyDecrypt(prikey, encData)
	if errDec != nil {
		fmt.Println(errDec)
		return
	}
	//fmt.Println("解密后二进制数据 ", decData)
	fmt.Println("解密后字符串 [", string(decData), "] 长度", len(decData))
}
