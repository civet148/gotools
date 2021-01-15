package dir

import (
	"os"
	"path/filepath"
)

//get current program abs path dir
func GetProgramAbsDir() (strPath string) {
	strPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	return
}
