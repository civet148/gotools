package ige256

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/civet148/gotools/cryptos/goaes"
)

type CryptoAES_IGE256 struct {
	key, iv []byte
}

func init() {
	goaes.Register(goaes.AES_Type_IGE256, NewCryptoAES_IGE256)
}

//key 长度必须为16/24/32字节(128/192/256 bits)
//iv 向量长度固定16字节
func NewCryptoAES_IGE256(key, iv []byte) goaes.CryptoAES {
	goaes.AssertKey(key)
	goaes.AssertIV(iv)

	return &CryptoAES_IGE256{
		key: key,
		iv:  iv,
	}
}

//加密后返回二进制字节数据切片
func (c *CryptoAES_IGE256) Encrypt(in []byte) (out []byte, err error) {
	var block cipher.Block
	if block, err = aes.NewCipher(c.key); err != nil {
		return
	}
	blockSize := block.BlockSize()
	goaes.PKCS7Padding(in, blockSize)
	return
}

//加密后将密文做BASE64编码字符串
func (c *CryptoAES_IGE256) EncryptBase64(in []byte) (out string, err error) {

	return
}

//解密后返回二进制字节数据切片
func (c *CryptoAES_IGE256) Decrypt(in []byte) (out []byte, err error) {

	return
}

//解密BASE64编码字符串的密文后返回二进制切片
func (c *CryptoAES_IGE256) DecryptBase64(in string) (out []byte, err error) {

	return
}
