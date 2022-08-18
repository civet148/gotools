package process

import (
	"os"
	"strings"
)

const (
	ENV_PATH_NAME = "PATH"
)

func GetEnv(strKey string) string {
	return os.Getenv(strKey)
}

func GetWorkDir() string {
	strWorkDir, _ := os.Getwd()
	return strWorkDir
}

func GetEnvPath() string {
	return os.Getenv(ENV_PATH_NAME)
}

func GetEnvPathSlice() []string {
	strEnv := GetEnvPath()
	return strings.Split(strEnv, string(os.PathListSeparator))
}
