package hmacsha

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
)

func HmacSha1(strContent string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(strContent))
	sha := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

func HmacSha256(strContent string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strContent))
	sha := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

func HmacSha512(strContent string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha512.New, key)
	h.Write([]byte(strContent))
	sha := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}
