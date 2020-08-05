package ecb

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/civet148/gotools/cryptos/goaes"
)

type CryptoAES_ECB struct {
	key, iv []byte
}

func init() {
	goaes.Register(goaes.AES_Mode_ECB, NewCryptoAES_ECB)
}

//key 长度必须为16/24/32字节(128/192/256 bits)
//iv 向量可为nil
func NewCryptoAES_ECB(key, iv []byte) goaes.CryptoAES {
	goaes.AssertKey(key)

	return &CryptoAES_ECB{
		key: key,
		iv:  iv,
	}
}

//加密后返回二进制字节数据切片
func (c *CryptoAES_ECB) Encrypt(in []byte) (out []byte, err error) {
	var block cipher.Block
	if block, err = aes.NewCipher(c.key); err != nil {
		return
	}
	var data = make([]byte, len(in))
	copy(data, in)
	blockSize := block.BlockSize()
	//填充
	data = goaes.PKCS7Padding(data, blockSize)
	//分配返回数据切片空间
	out = make([]byte, len(data))
	//存储每次加密的数据
	tmpData := make([]byte, blockSize)
	//分组分块加密
	for index := 0; index < len(data); index += blockSize {
		offset := index + blockSize
		block.Encrypt(tmpData, data[index:offset])
		copy(out[index:], tmpData)
	}
	return
}

//加密后将密文做BASE64编码字符串
func (c *CryptoAES_ECB) EncryptBase64(in []byte) (out string, err error) {
	var enc []byte
	if enc, err = c.Encrypt(in); err != nil {
		return
	}
	out = base64.StdEncoding.EncodeToString(enc)
	return
}

//解密后返回二进制字节数据切片
func (c *CryptoAES_ECB) Decrypt(in []byte) (out []byte, err error) {
	var block cipher.Block
	if block, err = aes.NewCipher(c.key); err != nil {
		return
	}
	//分配解密数据存储空间
	out = make([]byte, len(in))
	//存储每次加密的数据
	tmpData := make([]byte, block.BlockSize())
	//分组分块加密
	for index := 0; index < len(in); index += block.BlockSize() {
		offset := index + block.BlockSize()
		block.Decrypt(tmpData, in[index:offset])
		copy(out[index:], tmpData)
	}
	out = goaes.PKCS7UnPadding(out)
	return
}

//解密BASE64编码字符串的密文后返回二进制切片
func (c *CryptoAES_ECB) DecryptBase64(in string) (out []byte, err error) {
	var data []byte
	if data, err = base64.StdEncoding.DecodeString(in); err != nil {
		return
	}
	return c.Decrypt(data)
}

//获取当前AES模式
func (c *CryptoAES_ECB) GetMode() goaes.AES_Mode {
	return goaes.AES_Mode_ECB
}
