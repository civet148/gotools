package myrsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/pkcs12"
	"io"
	"io/ioutil"
	"math/big"
)

var (
	ErrDataToLarge     = errors.New("message too long for RSA public key size")
	ErrDataLen         = errors.New("data length error")
	ErrDataBroken      = errors.New("data broken, first byte is not zero")
	ErrKeyPairDismatch = errors.New("data is not encrypted by the private key")
	ErrDecryption      = errors.New("decryption error")
	ErrPublicKey       = errors.New("get public key error")
	ErrPrivateKey      = errors.New("get private key error")
)

// 公钥加密
func PubKeyEncrypt(pubkey *rsa.PublicKey, input []byte) ([]byte, error) {
	if pubkey == nil {
		return []byte(""), errors.New(`Please set the public key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(pubkey, bytes.NewReader(input), output, true)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}

// 公钥解密
func PubKeyDecrypt(pubkey *rsa.PublicKey, input []byte) ([]byte, error) {
	if pubkey == nil {
		return []byte(""), errors.New(`Please set the public key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(pubkey, bytes.NewReader(input), output, false)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}

// 私钥加密
func PriKeyEncrypt(prikey *rsa.PrivateKey, input []byte) ([]byte, error) {
	if prikey == nil {
		return []byte(""), errors.New(`Please set the private key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(prikey, bytes.NewReader(input), output, true)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}

// 私钥解密
func PriKeyDecrypt(prikey *rsa.PrivateKey, input []byte) ([]byte, error) {
	if prikey == nil {
		return []byte(""), errors.New(`Please set the private key in advance`)
	}
	output := bytes.NewBuffer(nil)
	err := priKeyIO(prikey, bytes.NewReader(input), output, false)
	if err != nil {
		return []byte(""), err
	}

	return ioutil.ReadAll(output)
}

// 获得公钥(标准RSA PUBLIC KEY)
func GetPubKey(pubkey []byte) (*rsa.PublicKey, error) {
	// decode public key
	block, _ := pem.Decode(pubkey)
	if block == nil {
		return nil, errors.New("get public key error")
	}
	// x509 parse public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), err
}

// 获得私钥(标准RSA PRIVATE KEY)
func GetPriKey(prikey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(prikey)
	if block == nil {
		return nil, errors.New("get private key error")
	}
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return pri, nil
	}
	pri2, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pri2.(*rsa.PrivateKey), nil
}

//通过内存数组获得证书公钥
func GetPubKeyCer(cerKey []byte) (pubkey *rsa.PublicKey, err error) {
	pubkey, err = ParseCerPublicKey(cerKey)
	if err != nil {
		fmt.Println("parse certifacte key failed, err=", err)
		return
	}
	return
}

//通过指定文件获得证书公钥
func GetPubKeyByCerFile(strFile string) (pubkey *rsa.PublicKey, err error) {

	var data []byte
	if data, err = readFile(strFile); err != nil {
		fmt.Println("read cer/crt file error (%s)", err.Error())
		return nil, err
	}

	return GetPubKeyCer(data)
}

//通过内存数组获得证书私钥
func GetPriKeyPfx(pfxKey []byte, strPasswd string) (prikey *rsa.PrivateKey, err error) {

	private, _, errDec := pkcs12.Decode(pfxKey, strPasswd)
	if errDec != nil {
		fmt.Println("pkcs12 decode pfx error(%s)", errDec.Error())
		return nil, errDec
	}
	prikey = private.(*rsa.PrivateKey)
	return
}

//通过指定文件获得证书私钥
func GetPriKeyByPfxFile(strFile string, strPasswd string) (prikey *rsa.PrivateKey, err error) {

	var data []byte
	if data, err = readFile(strFile); err != nil {
		fmt.Println("read pfx file error (%s)", err.Error())
		return nil, err
	}
	return GetPriKeyPfx(data, strPasswd)
}

//读取pfx/cert（公/私钥文件），返回数组
func readFile(strFile string) (data []byte, err error) {

	data, err = ioutil.ReadFile(strFile)
	if err != nil {
		fmt.Println("read file[%s]error:[%s]", strFile, err.Error())
		return
	}
	return
}

//Certificate证书获取公钥对象
func ParseCerPublicKey(cerKey []byte) (*rsa.PublicKey, error) {

	//解密pem格式的公钥
	block, _ := pem.Decode(cerKey)
	if block == nil {
		return nil, fmt.Errorf("public key decode error")
	}
	//log.Debug("公钥block[%+v]", block)

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("parse certifacte key failed, err=", err)
		return nil, err
	}

	return cert.PublicKey.(*rsa.PublicKey), nil
}

// 公钥加密或解密byte
func pubKeyByte(pub *rsa.PublicKey, in []byte, isEncrytp bool) ([]byte, error) {
	k := (pub.N.BitLen() + 7) / 8
	if isEncrytp {
		k = k - 11
	}
	if len(in) <= k {
		if isEncrytp {
			return rsa.EncryptPKCS1v15(rand.Reader, pub, in)
		} else {
			return pubKeyDecrypt(pub, in)
		}
	} else {
		iv := make([]byte, k)
		out := bytes.NewBuffer(iv)
		if err := pubKeyIO(pub, bytes.NewReader(in), out, isEncrytp); err != nil {
			return nil, err
		}
		return ioutil.ReadAll(out)
	}
}

// 私钥加密或解密byte
func priKeyByte(pri *rsa.PrivateKey, in []byte, isEncrytp bool) ([]byte, error) {
	k := (pri.N.BitLen() + 7) / 8
	if isEncrytp {
		k = k - 11
	}
	if len(in) <= k {
		if isEncrytp {
			return priKeyEncrypt(rand.Reader, pri, in)
		} else {
			return rsa.DecryptPKCS1v15(rand.Reader, pri, in)
		}
	} else {
		iv := make([]byte, k)
		out := bytes.NewBuffer(iv)
		if err := priKeyIO(pri, bytes.NewReader(in), out, isEncrytp); err != nil {
			return nil, err
		}
		return ioutil.ReadAll(out)
	}
}

// 公钥加密或解密Reader
func pubKeyIO(pub *rsa.PublicKey, in io.Reader, out io.Writer, bEncrypt bool) (err error) {
	k := (pub.N.BitLen() + 7) / 8
	if bEncrypt {
		k = k - 11
	}
	buf := make([]byte, k)
	var b []byte
	size := 0
	for {
		size, err = in.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if size < k {
			b = buf[:size]
		} else {
			b = buf
		}
		if bEncrypt {
			b, err = rsa.EncryptPKCS1v15(rand.Reader, pub, b)
		} else {
			b, err = pubKeyDecrypt(pub, b)
		}
		if err != nil {
			return err
		}
		if _, err = out.Write(b); err != nil {
			return err
		}
	}
	return nil
}

// 私钥加密或解密Reader
func priKeyIO(pri *rsa.PrivateKey, r io.Reader, w io.Writer, bEncrypt bool) (err error) {
	k := (pri.N.BitLen() + 7) / 8
	if bEncrypt {
		k = k - 11
	}
	buf := make([]byte, k)
	var b []byte
	size := 0
	for {
		size, err = r.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if size < k {
			b = buf[:size]
		} else {
			b = buf
		}
		if bEncrypt {
			b, err = priKeyEncrypt(rand.Reader, pri, b)
		} else {
			b, err = rsa.DecryptPKCS1v15(rand.Reader, pri, b)
		}
		if err != nil {
			return err
		}
		if _, err = w.Write(b); err != nil {
			return err
		}
	}
	return nil
}

// 公钥解密
func pubKeyDecrypt(pub *rsa.PublicKey, data []byte) ([]byte, error) {
	k := (pub.N.BitLen() + 7) / 8
	if k != len(data) {
		return nil, ErrDataLen
	}
	m := new(big.Int).SetBytes(data)
	if m.Cmp(pub.N) > 0 {
		return nil, ErrDataToLarge
	}
	m.Exp(m, big.NewInt(int64(pub.E)), pub.N)
	d := leftPad(m.Bytes(), k)
	if d[0] != 0 {
		return nil, ErrDataBroken
	}
	if d[1] != 0 && d[1] != 1 {
		return nil, ErrKeyPairDismatch
	}
	var i = 2
	for ; i < len(d); i++ {
		if d[i] == 0 {
			break
		}
	}
	i++
	if i == len(d) {
		return nil, nil
	}
	return d[i:], nil
}

// 私钥加密
func priKeyEncrypt(rand io.Reader, priv *rsa.PrivateKey, hashed []byte) ([]byte, error) {
	tLen := len(hashed)
	k := (priv.N.BitLen() + 7) / 8
	if k < tLen+11 {
		return nil, ErrDataLen
	}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], hashed)
	m := new(big.Int).SetBytes(em)
	c, err := decrypt(rand, priv, m)
	if err != nil {
		return nil, err
	}
	copyWithLeftPad(em, c.Bytes())
	return em, nil
}

var bigZero = big.NewInt(0)
var bigOne = big.NewInt(1)

func encrypt(c *big.Int, pub *rsa.PublicKey, m *big.Int) *big.Int {
	e := big.NewInt(int64(pub.E))
	c.Exp(m, e, pub.N)
	return c
}

func decrypt(random io.Reader, priv *rsa.PrivateKey, c *big.Int) (m *big.Int, err error) {
	if c.Cmp(priv.N) > 0 {
		err = ErrDecryption
		return
	}
	var ir *big.Int
	if random != nil {
		var r *big.Int

		for {
			r, err = rand.Int(random, priv.N)
			if err != nil {
				return
			}
			if r.Cmp(bigZero) == 0 {
				r = bigOne
			}
			var ok bool
			ir, ok = modInverse(r, priv.N)
			if ok {
				break
			}
		}
		bigE := big.NewInt(int64(priv.E))
		rpowe := new(big.Int).Exp(r, bigE, priv.N)
		cCopy := new(big.Int).Set(c)
		cCopy.Mul(cCopy, rpowe)
		cCopy.Mod(cCopy, priv.N)
		c = cCopy
	}
	if priv.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, priv.D, priv.N)
	} else {
		m = new(big.Int).Exp(c, priv.Precomputed.Dp, priv.Primes[0])
		m2 := new(big.Int).Exp(c, priv.Precomputed.Dq, priv.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, priv.Primes[0])
		}
		m.Mul(m, priv.Precomputed.Qinv)
		m.Mod(m, priv.Primes[0])
		m.Mul(m, priv.Primes[1])
		m.Add(m, m2)

		for i, values := range priv.Precomputed.CRTValues {
			prime := priv.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}
	if ir != nil {
		m.Mul(m, ir)
		m.Mod(m, priv.N)
	}

	return
}

func copyWithLeftPad(dest, src []byte) {
	numPaddingBytes := len(dest) - len(src)
	for i := 0; i < numPaddingBytes; i++ {
		dest[i] = 0
	}
	copy(dest[numPaddingBytes:], src)
}

func nonZeroRandomBytes(s []byte, rand io.Reader) (err error) {
	_, err = io.ReadFull(rand, s)
	if err != nil {
		return
	}
	for i := 0; i < len(s); i++ {
		for s[i] == 0 {
			_, err = io.ReadFull(rand, s[i:i+1])
			if err != nil {
				return
			}
			s[i] ^= 0x42
		}
	}
	return
}

func leftPad(input []byte, size int) (out []byte) {
	n := len(input)
	if n > size {
		n = size
	}
	out = make([]byte, size)
	copy(out[len(out)-n:], input)
	return
}

func modInverse(a, n *big.Int) (ia *big.Int, ok bool) {
	g := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)
	g.GCD(x, y, a, n)
	if g.Cmp(bigOne) != 0 {
		return
	}
	if x.Cmp(bigOne) < 0 {
		x.Add(x, n)
	}
	return x, true
}
