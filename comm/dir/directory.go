package dir

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func GetSuffix(strFileName string) string {

	idx := strings.LastIndex(strFileName, ".")
	if idx == -1 {
		return ""
	}
	return strFileName[idx:]
}

func hasSuffix(strFileName string, suffixs []string) bool {

	for _, v := range suffixs {
		if strings.HasSuffix(strFileName, v) {
			return true
		}
	}
	return false
}

func hasSuffixRegexp(strFileName string, suffixs []string) bool {
	var err error
	strSuffix := GetSuffix(strFileName)

	for _, v := range suffixs {

		var r *regexp.Regexp
		if r, err = regexp.Compile(v); err != nil {
			//fmt.Println(fmt.Sprintf("regexp.Compile(%v) error [%v]", v, err.Error()))
			continue
		}

		if ok := r.Match([]byte(strSuffix)); ok {
			//fmt.Println(fmt.Sprintf("strFileName [%v] regexp suffix [%v] match ok", strFileName, v))
			return true
		}
	}
	return false
}

//获取指定文件或目录的上级目录
func GetParentDir(strFileName string) (dirParent string) {
	pathSep := string(os.PathSeparator)
	strFileName = strings.TrimSpace(strFileName)
	idx := strings.LastIndex(strFileName, pathSep)
	if idx > 0 {
		dirParent = strFileName[:idx]
	}
	return
}

//获取指定目录下的所有文件和目录(不递归）
func GetFilesAndDirs(dirPath string, suffixs ...string) (d DirEntry, err error) {
	pathSep := string(os.PathSeparator)

	dirPath = strings.TrimSpace(dirPath)
	idx := strings.LastIndex(dirPath, pathSep)
	if idx > 0 {
		dirPath = dirPath[:idx]
	}
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}

	d.Name = dirPath

	for _, f := range dir {
		if f.IsDir() { // 目录, 递归遍历

			d.dirs = append(d.dirs, newDirEntry(dirPath, pathSep, f))
		} else {

			var ok bool
			if len(suffixs) == 0 { //没有指定文件后缀
				ok = true
			} else {
				ok = hasSuffixRegexp(f.Name(), suffixs)
			}
			if ok {
				d.files = append(d.files, newFileEntry(dirPath, pathSep, f))
			}
		}
	}
	return
}
