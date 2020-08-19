package process

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetProcessName() (path string) {
	var strFile string
	strFile, _ = exec.LookPath(os.Args[0])
	path, _ = filepath.Abs(strFile)
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return strFile
	}
	return path[i+1:]
}

func GetProcessPath() (path string) {
	var strFile string
	strFile, _ = exec.LookPath(os.Args[0])
	path, _ = filepath.Abs(strFile)
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return //nothing
	}
	return path[0 : i+1]
}
