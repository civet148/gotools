package gomd5

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func Md5File(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	defer f.Close()

	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5Hash.Sum(nil)), nil
}
