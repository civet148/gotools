package main

import (
	"time"

	log "github.com/civet148/gotools.git/log"
)

/**
* 1. 通过参数直接指定日志文件、输出级别(DEBUG,INFO,WARN,ERROR, FATAL)和属性
*
*	1.1. 直接输入文件名
*	Open("test.log")
*
*	1.2. 设置文件日志输出级别和分块大小(单位：MB)
*  	Open("file:///var/log/test.log?file_level=INFO&file_size=50")
*
*	1.3. 设置文件日志输出级别和分块大小(单位：MB)及邮件通知级别、邮件地址、邮件标题
*	Open("file:///var/log/test.log?file_level=INFO&file_size=50&email_level=FATAL&email=civet148@126.com&email_title=service-error-message")
*
* 2. 	通过指定json配置文件设置日志级别、日志文件及属性
*
*   2.1. 指定json配置文件
*   Open("json:///tmp/test.json")
*
*   test.json 范例
*   {
*      "file_path":"/tmp/test.log",
*      "file_level":"INFO",
*      "file_size":"50",
*      "email_level":"FATAL",
*      "email_addr":"civet126@126.com",
*      "email_title":"error message title"
*   }
 */

type testSubSt struct {
	SubInt int
	SubStr string
}

type testSt struct {
	MyPtr      *string
	MyInt      int
	MyFloat64  float64
	MyMap      map[string]string
	MyMapPtr   *map[string]string
	MySubSt    testSubSt
	MySubStPtr *testSubSt
	abc        int //非导出字段(不处理会报panic错误)
}

func main() {

	//strUrl := "test.log" //指定当前目录创建日志文件（Windows+linux通用）
	//strUrl := "file://e:/test.log" //指定日志文件但不指定属性（Windows）
	//strUrl := "file:///tmp/test.log" //指定日志文件但不指定属性(Linux)
	//strUrl := "json:///tmp/test.json" //json文件名(Linux)
	strUrl := "json://test.json" //json文件名(Windows)
	//strUrl := "file:///var/log/test.log?file_level=INFO&file_size=50" //文件属性
	//strUrl := "file://e:/test.log?file_level=WARN&file_size=50"
	log.Open(strUrl)
	defer log.Close()

	for i := 0; i < 1; i++ {
		log.Debug("This is debug message")
		log.Info("This is info message")
		log.Warn("This is warn message")
		log.Error("This is error message")
		log.Fatal("This is fatal message")
		time.Sleep(5 * time.Millisecond)

		st1 := testSt{MyInt: 1, MyFloat64: 2.00, MySubSt: testSubSt{SubInt: 1, SubStr: "MySubSt"}, MySubStPtr: &testSubSt{SubInt: 19, SubStr: "MySubStPtr"}}
		st2 := &testSt{MyInt: 2, MyFloat64: 4.00}
		log.Struct("打印结构体", st1, st2)
		time.Sleep(5 * time.Second)
	}

	log.Info("Program exit...")
}
