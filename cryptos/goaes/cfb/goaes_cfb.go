package cfb

import "github.com/civet148/gotools/cryptos/goaes"

type CryptoAES_CFB struct {
	key, iv []byte
}

func init() {
	goaes.Register(goaes.AES_Type_CFB, NewCryptoAES_CFB)
}

//key 长度必须为16/24/32字节(128/192/256 bits)
//iv 向量长度固定16字节
func NewCryptoAES_CFB(key, iv []byte) goaes.CryptoAES {

	goaes.AssertKey(key)
	goaes.AssertIV(iv)

	return &CryptoAES_CFB{
		key: key,
		iv:  iv,
	}
}

//加密后返回二进制字节数据切片
func (c *CryptoAES_CFB) Encrypt(in []byte) (out []byte, err error) {

	return
}

//加密后将密文做BASE64编码字符串
func (c *CryptoAES_CFB) EncryptBase64(in []byte) (out string, err error) {

	return
}

//解密后返回二进制字节数据切片
func (c *CryptoAES_CFB) Decrypt(in []byte) (out []byte, err error) {

	return
}

//解密BASE64编码字符串的密文后返回二进制切片
func (c *CryptoAES_CFB) DecryptBase64(in string) (out []byte, err error) {

	return
}
