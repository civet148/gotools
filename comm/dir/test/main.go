package main

import (
	"github.com/civet148/gotools/comm/dir"
	"github.com/civet148/gotools/log"
)

func main() {

	//测试获取指定文件或目录的上级目录
	log.Info("parent dir [%v]", dir.GetParentDir("/tmp/test.log"))

	//测试正则表达式查找/tmp目录下xxx.log文件以及 xxx.log.20200217_150326类似的日志备份文件
	entry, _ := dir.GetFilesAndDirs("/tmp", ".log", ".\\d{8}_\\d{6}")
	for _, v := range entry.GetFiles() {
		log.Json("file", v) //打印文件信息
	}
	//for _, v := range entry.GetDirs() {
	//	log.Json("dir", v) //打印目录信息
	//}
}
