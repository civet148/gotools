package goaes

import (
	"bytes"
	"fmt"
)

type AES_Mode int

const (
	AES_Mode_CBC AES_Mode = 1 //CBC模式(密钥长度16/24/32字节)
	AES_Mode_CFB AES_Mode = 2 //CFB模式(密钥长度16/24/32字节)
	AES_Mode_ECB AES_Mode = 3 //ECB模式(密钥长度16/24/32字节)
	AES_Mode_OFB AES_Mode = 4 //OFB模式(密钥长度16/24/32字节)
	AES_Mode_CTR AES_Mode = 5 //CTR模式(密钥长度16/24/32字节)
)

func (t AES_Mode) String() string {
	switch t {
	case AES_Mode_CBC:
		return "AES_Mode_CBC"
	case AES_Mode_CFB:
		return "AES_Mode_CFB"
	case AES_Mode_ECB:
		return "AES_Mode_ECB"
	case AES_Mode_OFB:
		return "AES_Mode_OFB"
	case AES_Mode_CTR:
		return "AES_Mode_CTR"
	}
	return "AES_Mode_Unknown"
}

func (t AES_Mode) GoString() string {
	return t.String()
}

type CryptoAES interface {
	//获取当前AES模式
	GetMode() AES_Mode
	//加密后返回二进制字节数据切片
	Encrypt([]byte) ([]byte, error)
	//加密后将密文做BASE64编码字符串
	EncryptBase64([]byte) (string, error)
	//加密后将密文做HEX编码字符串
	EncryptHex([]byte) (string, error)
	//解密后返回二进制字节数据切片
	Decrypt([]byte) ([]byte, error)
	//解密BASE64编码字符串的密文后返回二进制切片
	DecryptBase64(string) ([]byte, error)
	//HEX编码字符串的密文后返回二进制切片
	DecryptHex(string) ([]byte, error)
}

type instance func(key, iv []byte) CryptoAES

var mapInstances = make(map[AES_Mode]instance, 1)

func Register(aesType AES_Mode, inst instance) {
	mapInstances[aesType] = inst
}

//aesType AES加密模式
//key 长度必须为16/24/32字节(128/192/256 bits)
//iv  向量长度固定16字节(ECB模式可传nil)
func NewCryptoAES(aesType AES_Mode, key, iv []byte) CryptoAES {

	fn, ok := mapInstances[aesType]
	if !ok {
		panic(fmt.Sprintf("AES type [%v] instance not registered", aesType))
	}

	return fn(key, iv)
}

//补码(AES算法PKCS#7和PKCS#5一致)
func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, paddingText...)
}

//去码(AES算法PKCS#7和PKCS#5一致)
func PKCS7UnPadding(cihperData []byte) []byte {
	length := len(cihperData)
	unpadding := int(cihperData[length-1])
	return cihperData[:(length - unpadding)]
}

//判断AES密钥长度是否合法
func AssertKey(key []byte) {
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		panic("key length must be 16/24/32 bytes")
	}
}

//判断AES向量长度是否合法
func AssertIV(iv []byte) {
	keyLen := len(iv)
	if keyLen < 16 {
		panic("iv length must >= 16 bytes")
	}
}
