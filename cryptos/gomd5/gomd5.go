package gomd5

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
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
