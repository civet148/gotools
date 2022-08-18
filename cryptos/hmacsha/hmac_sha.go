package hmacsha

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
)

//HMAC+SHA1+HEX+BASE64
func HmacSha1(strContent string, strSecret string) string {
	sha := HmacSha1Hex(strContent, strSecret)
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

//HMAC+SHA256+HEX+BASE64
func HmacSha256(strContent string, strSecret string) string {
	sha := HmacSha1Hex(strContent, strSecret)
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

//HMAC+SHA512+HEX+BASE64
func HmacSha512(strContent string, strSecret string) string {
	sha := HmacSha1Hex(strContent, strSecret)
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

//HMAC+SHA1+HEX
func HmacSha1Hex(strContent string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(strContent))
	return hex.EncodeToString(h.Sum(nil))
}

//HMAC+SHA256+HEX
func HmacSha256Hex(strContent string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strContent))
	return hex.EncodeToString(h.Sum(nil))
}

//HMAC+SHA512+HEX
func HmacSha512Hex(strContent string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha512.New, key)
	h.Write([]byte(strContent))
	return hex.EncodeToString(h.Sum(nil))
}
