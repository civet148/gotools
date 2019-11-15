package comm


import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)
//生成随机字符串
func genRandStr(length int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//生成32位MD5
func convToMD5(text string) string{
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

//生成随机字符串并MD5成32位

func GenRandStrMD5(length int) string {

	strRand := genRandStr(length)

	return convToMD5(strRand)
}
