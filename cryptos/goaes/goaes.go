package goaes

import (
	"bytes"
	"fmt"
)

type AES_Type int

const (
	AES_Type_CBC    AES_Type = 1 //CBC模式(密钥长度16/24/32字节)
	AES_Type_CFB    AES_Type = 2 //CFB模式(密钥长度16/24/32字节)
	AES_Type_ECB    AES_Type = 3 //ECB模式(密钥长度16/24/32字节)
	AES_Type_OFB    AES_Type = 4 //OFB模式(密钥长度16/24/32字节)
	AES_Type_CTR128 AES_Type = 5 //CTR模式(密钥长度16/24/32字节)
	AES_Type_IGE256 AES_Type = 6 //IGE模式(密钥长度16/24/32字节)
)

func (t AES_Type) String() string {
	switch t {
	case AES_Type_CBC:
		return "AES_Type_CBC"
	case AES_Type_CFB:
		return "AES_Type_CFB"
	case AES_Type_ECB:
		return "AES_Type_ECB"
	case AES_Type_OFB:
		return "AES_Type_OFB"
	case AES_Type_CTR128:
		return "AES_Type_CTR128"
	case AES_Type_IGE256:
		return "AES_Type_IGE256"
	}
	return "AES_Type_Unknown"
}

func (t AES_Type) GoString() string {
	return t.String()
}

type CryptoAES interface {
	//加密后返回二进制字节数据切片
	Encrypt([]byte) ([]byte, error)
	//加密后将密文做BASE64编码字符串
	EncryptBase64([]byte) (string, error)
	//解密后返回二进制字节数据切片
	Decrypt([]byte) ([]byte, error)
	//解密BASE64编码字符串的密文后返回二进制切片
	DecryptBase64(string) ([]byte, error)
}

type instance func(key, iv []byte) CryptoAES

var mapInstances = make(map[AES_Type]instance, 1)

func Register(aesType AES_Type, inst instance) {
	mapInstances[aesType] = inst
}

//aesType AES加密模式
//key 长度必须为16/24/32字节(128/192/256 bits)
//iv  向量长度固定16字节(CBC模式可传nil)
func NewCryptoAES(aesType AES_Type, key, iv []byte) CryptoAES {

	fn, ok := mapInstances[aesType]
	if !ok {
		panic(fmt.Sprintf("AES type [%v] instance not registered", aesType))
	}

	return fn(key, iv)
}

//补码
func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, paddingText...)
}

//去码
func PKCS7UnPadding(cihperData []byte) []byte {
	length := len(cihperData)
	unpadding := int(cihperData[length-1])
	return cihperData[:(length - unpadding)]
}

func AssertKey(key []byte) {
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		panic("key length must be 16/24/32 bytes")
	}
}

func AssertIV(iv []byte) {
	keyLen := len(iv)
	if keyLen < 16 {
		panic("iv length must >= 16 bytes")
	}
}
