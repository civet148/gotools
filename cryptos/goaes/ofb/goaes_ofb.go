package ofb

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/civet148/gotools/cryptos/goaes"
)

type CryptoAES_OFB struct {
	key, iv []byte
}

func init() {
	goaes.Register(goaes.AES_Mode_OFB, NewCryptoAES_OFB)
}

//key 长度必须为16/24/32字节(128/192/256 bits)
//iv 向量长度固定16字节
func NewCryptoAES_OFB(key, iv []byte) goaes.CryptoAES {
	goaes.AssertKey(key)
	goaes.AssertIV(iv)

	return &CryptoAES_OFB{
		key: key,
		iv:  iv,
	}
}

//加密后返回二进制字节数据切片
func (c *CryptoAES_OFB) Encrypt(in []byte) (out []byte, err error) {
	var block cipher.Block
	if block, err = aes.NewCipher(c.key); err != nil {
		return
	}
	var data = make([]byte, len(in))
	copy(data, in)
	blockMode := cipher.NewOFB(block, c.iv)
	out = make([]byte, len(data))
	blockMode.XORKeyStream(out, data)
	return
}

//加密后将密文做BASE64编码字符串
func (c *CryptoAES_OFB) EncryptBase64(in []byte) (out string, err error) {
	var enc []byte
	if enc, err = c.Encrypt(in); err != nil {
		return
	}
	out = base64.StdEncoding.EncodeToString(enc)
	return
}

//解密后返回二进制字节数据切片
func (c *CryptoAES_OFB) Decrypt(in []byte) (out []byte, err error) {
	var block cipher.Block
	if block, err = aes.NewCipher(c.key); err != nil {
		return
	}
	blockMode := cipher.NewOFB(block, c.iv)
	out = make([]byte, len(in))
	blockMode.XORKeyStream(out, in)
	return
}

//解密BASE64编码字符串的密文后返回二进制切片
func (c *CryptoAES_OFB) DecryptBase64(in string) (out []byte, err error) {
	var data []byte
	if data, err = base64.StdEncoding.DecodeString(in); err != nil {
		return
	}
	return c.Decrypt(data)
}

//获取当前AES模式
func (c *CryptoAES_OFB) GetMode() goaes.AES_Mode {
	return goaes.AES_Mode_OFB
}
