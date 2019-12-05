package main

import (
	"time"

	log "github.com/civet148/gotools/log"
)

/**
* 1. 通过参数直接指定日志文件、输出级别(DEBUG,INFO,WARN,ERROR, FATAL)和属性
*
*	1.1. 直接输入文件名
*	Open("test.log")
*
*	1.2. 设置文件日志输出级别和分块大小(单位：MB)
*  	Open("file:///var/log/test.log?log_level=INFO&file_size=50")
*
*	1.3. 设置文件日志输出级别和分块大小(单位：MB)及邮件通知级别、邮件地址、邮件标题
*	Open("file:///var/log/test.log?log_level=INFO&file_size=50&email_level=FATAL&email=civet148@126.com&email_title=service-error-message")
*
* 2. 	通过指定json配置文件设置日志级别、日志文件及属性
*
*   2.1. 指定json配置文件
*   Open("json:///tmp/test.json")
*
*   test.json 范例
*   {
*      "file_path":"/tmp/test.log",
*      "log_level":"INFO",
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
	abc        int       //非导出字段(不处理会报panic错误)
	str        string 	 //非导出字段
	flt32      float32 	 //非导出字段
	flt64      float64 	 //非导出字段
	ui32       uint32 	 //非导出字段
	ui8        uint8 	 //非导出字段
	i8         int8      //非导出字段
	i64        int64     //非导出字段
	slice      []string  //非导出字段: 切片
	arr2       [5]byte   //非导出字段: 数组
	test       testSubSt //非导出字段: 结构体
	MySubSt    testSubSt
	MySubStPtr *testSubSt
}

func main() {

	strUrl := "test.log" //指定当前目录创建日志文件（Windows+linux通用）
	//strUrl := "file://e:/test.log" //指定日志文件但不指定属性（Windows）
	//strUrl := "file:///tmp/test.log" //指定日志文件但不指定属性(Linux)
	//strUrl := "json:///tmp/test.json" //json文件名(Linux)
	//strUrl := "json://test.json" //json文件名(Windows)
	//strUrl := "file:///var/log/test.log?log_level=INFO&file_size=50" //文件属性
	//strUrl := "file://e:/test.log?log_level=WARN&file_size=50"

	log.Open(strUrl)
	defer log.Close()

	//log.SetLevel(log.LEVEL_INFO)
	log.Debug("This is debug message")
	log.Info("This is info message")
	log.Warn("This is warn message")
	log.Error("This is error message")
	log.Fatal("This is fatal message")
	time.Sleep(5 * time.Millisecond)

	st1 := testSt{abc:10086, flt32: 0.58, flt64: 0.96666, ui8: 25, ui32: 10032, i8: 44, i64: 100000000000019, str: "ni hao", slice: []string{"str1", "str2"},
		MyInt: 1, MyFloat64: 2.00, MySubSt: testSubSt{SubInt: 1, SubStr: "MySubSt"}, MySubStPtr: &testSubSt{SubInt: 19, SubStr: "MySubStPtr"}}
	st2 := &testSt{MyInt: 2, MyFloat64: 4.00, abc:9999}
	log.Json(st1, st2)
	log.Json("hello , I'm a string object", 123456, []int{100, 200, 300, 400}, map[string]interface{}{"key1": 109, "key2": "hello"})
	log.Struct(st1, st2)
	log.Info("Program exit...")
}
