package godh

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
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
	shareKey         [CRYPTO_DH_KEY_LENGTH]byte //共享密钥(用于双方加密/解密)
	privateKey       [CRYPTO_DH_KEY_LENGTH]byte //私钥
	publicKey        [CRYPTO_DH_KEY_LENGTH]byte //公钥
	shareKeyBase64   string                     //共享密钥BASE64字符串
	publicKeyBase64  string                     //公钥BASE64字符串
	privateKeyBase64 string                     //私钥BASE64字符串
}

func init() {
}

//创建DH加密对象并随机生成私钥和对应公钥(如果pri参数为空)
func NewCryptoDH(pri ...[32]byte) (dh *CryptoDH) {
	dh = &CryptoDH{}

	if len(pri) == 0 { //随机生成私钥
		dh.makeCurve25519PrivateKey()
	} else {
		dh.setPrivateKey(pri[0])
	}
	curve25519.ScalarBaseMult(&dh.publicKey, &dh.privateKey)
	dh.publicKeyBase64 = base64.StdEncoding.EncodeToString(dh.publicKey[:])
	dh.privateKeyBase64 = base64.StdEncoding.EncodeToString(dh.privateKey[:])
	return
}

//pub 对方的公钥(32字节byte数组)
//返回key：自己的私钥+对方公钥经DH算法计算出来的加密KEY
func (dh *CryptoDH) ScalarMult(pub [32]byte) [32]byte {
	curve25519.ScalarMult(&dh.shareKey, &dh.privateKey, &pub)
	dh.shareKeyBase64 = base64.StdEncoding.EncodeToString(dh.shareKey[:])
	return dh.shareKey
}

//base 对方的公钥(base64编码)
//返回key：自己的私钥+对方公钥经DH算法计算出来的加密KEY(base64编码)
func (dh *CryptoDH) ScalarMultBase64(base string) string {
	var pub [32]byte
	s, err := base64.StdEncoding.DecodeString(base)
	if err != nil {
		panic(fmt.Sprintf("parameter base64 [%v] illegal", base))
	}
	copy(pub[:], s)
	_ = dh.ScalarMult(pub)
	return dh.shareKeyBase64
}

func (dh *CryptoDH) GetPrivateKey() [32]byte {
	return dh.privateKey
}

func (dh *CryptoDH) GetPrivateKeyBase64() string {
	return dh.privateKeyBase64
}

func (dh *CryptoDH) GetPublicKey() [32]byte {
	return dh.publicKey
}

func (dh *CryptoDH) GetPublicKeyBase64() string {
	return dh.publicKeyBase64
}

func (dh *CryptoDH) GetShareKey() [32]byte {
	return dh.shareKey
}

func (dh *CryptoDH) GetShareKeyBase64() string {
	return dh.shareKeyBase64
}

func (dh *CryptoDH) GetCurve25519KeyPair() (pri, pub [32]byte) {
	return dh.privateKey, dh.publicKey
}

func (dh *CryptoDH) GetCurve25519KeyPairBase64() (priBase64, pubBase64 string) {
	return dh.privateKeyBase64, dh.publicKeyBase64
}

func (dh *CryptoDH) setPrivateKey(pri [32]byte) {
	dh.privateKey = pri
}

func (dh *CryptoDH) makeCurve25519PrivateKey() {
	if _, err := io.ReadFull(rand.Reader, dh.privateKey[:]); err != nil {
		panic("generate random DH private key data error")
	}
	return
}
