package main

import (
	"github.com/civet148/gotools/comm/dir"
	"github.com/civet148/gotools/log"
)

func main() {

	//测试正则表达式查找/tmp目录下xxx.log文件以及 xxx.log.20200217_150326类似的日志备份文件
	entry, _ := dir.GetFilesAndDirs("/tmp", ".log", ".\\d{8}_\\d{6}")
	for _, v := range entry.GetFiles() {
		log.Json("files", v)
	}
	for _, v := range entry.GetDirs() {
		log.Json("dirs", v)
	}
}
