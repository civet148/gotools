package gomd5

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"
)

func Md5File(strFileName string) string {
	f, err := os.Open(strFileName)
	if err != nil {
		return ""
	}
	defer f.Close()

	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, f); err != nil {
		return ""
	}

	return fmt.Sprintf("%x", md5Hash.Sum(nil))
}

func Md5Sum(data []byte) string {
	md5Hash := md5.New()
	md5Hash.Write(data)
	return fmt.Sprintf("%x", md5Hash.Sum(nil))
}

func Md5SumUpper(data []byte) string {
	return strings.ToUpper(Md5Sum(data))
}

func Md5String(data string) string {
	return Md5Sum([]byte(data))
}

func Md5StringUpper(data string) string {
	return strings.ToUpper(Md5String(data))
}
