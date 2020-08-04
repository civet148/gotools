package cbc

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/civet148/gotools/cryptos/goaes"
)

type CryptoAES_CBC struct {
	key, iv []byte
}

func init() {
	goaes.Register(goaes.AES_Mode_CBC, NewCryptoAES_CBC)
}

//key 长度必须为16/24/32字节(128/192/256 bits)
//iv 向量长度固定16字节(CBC模式可传nil)
func NewCryptoAES_CBC(key, iv []byte) goaes.CryptoAES {

	goaes.AssertKey(key)
	return &CryptoAES_CBC{
		key: key,
		iv:  iv,
	}
}

//加密后返回二进制字节数据切片
func (c *CryptoAES_CBC) Encrypt(in []byte) (out []byte, err error) {
	var block cipher.Block
	if block, err = aes.NewCipher(c.key); err != nil {
		return
	}
	blockSize := block.BlockSize()
	goaes.PKCS7Padding(in, blockSize)
	return
}

//加密后将密文做BASE64编码字符串
func (c *CryptoAES_CBC) EncryptBase64(in []byte) (out string, err error) {

	return
}

//解密后返回二进制字节数据切片
func (c *CryptoAES_CBC) Decrypt(in []byte) (out []byte, err error) {

	return
}

//解密BASE64编码字符串的密文后返回二进制切片
func (c *CryptoAES_CBC) DecryptBase64(in string) (out []byte, err error) {

	return
}

//获取当前AES模式
func (c *CryptoAES_CBC) GetMode() goaes.AES_Mode {
	return goaes.AES_Mode_CBC
}
