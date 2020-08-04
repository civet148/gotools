package godh

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/curve25519"
	"io"
)

/*
参考文档

原理说明 https://blog.csdn.net/andylau00j/article/details/82870841
代码示例 https://github.com/hwholiday/learning_tools/blob/master/encryption_algorithm/curve25519.go
*/

const (
	CRYPTO_DH_KEY_LENGTH = 32
)

type CryptoDH struct {
	Key              [CRYPTO_DH_KEY_LENGTH]byte
	PrivateKey       [CRYPTO_DH_KEY_LENGTH]byte
	PublicKey        [CRYPTO_DH_KEY_LENGTH]byte
	KeyBase64        string
	PublicKeyBase64  string
	PrivateKeyBase64 string
}

func init() {
}

//创建DH加密对象并随机生成私钥和对应公钥(如果pri参数为空)
func NewCryptoDH(pri ...[32]byte) (dh *CryptoDH) {
	dh = &CryptoDH{}

	if len(pri) == 0 {
		if _, err := io.ReadFull(rand.Reader, dh.PrivateKey[:]); err != nil {
			panic("generate random DH private key data error")
		}
	} else {
		dh.PrivateKey = pri[0]
	}

	curve25519.ScalarBaseMult(&dh.PublicKey, &dh.PrivateKey)
	dh.PublicKeyBase64 = base64.StdEncoding.EncodeToString(dh.PublicKey[:])
	dh.PrivateKeyBase64 = base64.StdEncoding.EncodeToString(dh.PrivateKey[:])
	return
}

//pub 对方的公钥(32字节byte数组)
//key=自己的私钥+对方公钥经DH算法计算出来的加密KEY
func (dh *CryptoDH) ScalarMult(pub [32]byte) [32]byte {
	curve25519.ScalarMult(&dh.Key, &dh.PrivateKey, &pub)
	dh.KeyBase64 = base64.StdEncoding.EncodeToString(dh.Key[:])
	return dh.Key
}

//pub 对方的公钥(32字节byte数组)
//key=自己的私钥+对方公钥经DH算法计算出来的加密KEY(base64编码)
func (dh *CryptoDH) ScalarMultBase64(pub [32]byte) string {
	_ = dh.ScalarMult(pub)
	return dh.KeyBase64
}

func (dh *CryptoDH) GetPrivateKey() [32]byte {
	return dh.PrivateKey
}

func (dh *CryptoDH) GetPrivateKeyBase64() string {
	return dh.PrivateKeyBase64
}

func (dh *CryptoDH) GetPublicKey() [32]byte {
	return dh.PublicKey
}

func (dh *CryptoDH) GetPublicKeyBase64() string {
	return dh.PublicKeyBase64
}

func (dh *CryptoDH) GetKey() [32]byte {
	return dh.Key
}

func (dh *CryptoDH) GetKeyBase64() string {
	return dh.KeyBase64
}
