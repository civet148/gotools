package gox3dh

import (
	"crypto/sha256"
	"github.com/civet148/gotools/cryptos/godh"
	"golang.org/x/crypto/hkdf"
	"hash"
)

/*
参考文档

原理说明 https://blog.csdn.net/andylau00j/article/details/82870841
代码示例 https://github.com/hwholiday/learning_tools/blob/master/encryption_algorithm/x3curve25519.go
*/

type CryptoX3DH struct {
	identityKeyPair  *godh.CryptoDH //身份密钥对(IPK)
	signedKeyPair    *godh.CryptoDH //已签名的预共享密钥(SPK)
	oneTimeKeyPair   *godh.CryptoDH //一次性预共享密钥(OPK)
	ephemeralKeyPair *godh.CryptoDH //一个临时密钥对(EPK)
	dH1              [32]byte
	dH2              [32]byte
	dH3              [32]byte
	dH4              [32]byte
	kdfPrefix        []byte
}

func (c *CryptoX3DH) NewCryptoX3DH() (x3dh *CryptoX3DH) {

	x3dh = &CryptoX3DH{
		identityKeyPair:  godh.NewCryptoDH(),
		signedKeyPair:    godh.NewCryptoDH(),
		oneTimeKeyPair:   godh.NewCryptoDH(),
		ephemeralKeyPair: godh.NewCryptoDH(),
	}

	return
}

//info 自定义明文字符串信息(可以为nil)
func (c *CryptoX3DH) kdf(data []byte, info string) []byte {
	// create reader
	r := hkdf.New(
		func() hash.Hash {
			return sha256.New()
		},
		data,
		make([]byte, 32), []byte(info),
	)
	var secret [32]byte
	_, err := r.Read(secret[:])
	if err != nil {
		panic(err)
	}
	return secret[:]
}

//内部产生下次KDF前缀和密钥(返回值为消息密钥)
func (c *CryptoX3DH) kdfSecret(data []byte, salt [32]byte) []byte {
	// create reader
	r := hkdf.New(
		func() hash.Hash {
			return sha256.New()
		},
		append(c.kdfPrefix[:], data...),
		salt[:], nil,
	)
	var secret [64]byte
	_, err := r.Read(secret[:])
	if err != nil {
		panic(err)
	}
	head := secret[:32]
	c.kdfPrefix = head
	tail := secret[32:]
	return tail
}
