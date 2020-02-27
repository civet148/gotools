package main

import (
	"github.com/civet148/gotools/log"
	"sync"
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

type SubSubSt struct {
	Name string
}

type TestSubSt struct {
	SubInt int
	SubStr string
	sst    *SubSubSt
}

type testSt struct {
	MyPtr      *string
	MyInt      int
	MyFloat64  float64
	MyMap      map[string]string
	MyMapPtr   *map[string]string
	abc        int       //非导出字段(不处理会报panic错误)
	str        string    //非导出字段
	flt32      float32   //非导出字段
	flt64      float64   //非导出字段
	ui32       uint32    //非导出字段
	ui8        uint8     //非导出字段
	i8         int8      //非导出字段
	i64        int64     //非导出字段
	slice      []string  //非导出字段: 切片
	arr2       [5]byte   //非导出字段: 数组
	test       TestSubSt //非导出字段: 结构体
	ip         *int32    //非导出字段
	MySubSt    TestSubSt
	MySubStPtr *TestSubSt
}

func main() {

	log.Enter()

	//strUrl := "test.log" //指定当前目录创建日志文件（Windows+linux通用）
	//strUrl := "file://e:/test.log" //指定日志文件但不指定属性（Windows）
	//strUrl := "file:///tmp/test.log" //指定日志文件但不指定属性(Linux)
	//strUrl := "json:///tmp/test.json" //json文件名(Linux)
	//strUrl := "json://f:/test/test.json" //json文件名(Windows)
	//strUrl := "file:///var/log/test.log?log_level=INFO&file_size=50" //Linux/Unix文件带属性
	//strUrl := "file://e:/test.log?log_level=WARN&file_size=50" //Windows文件带属性

	//log.Open(strUrl)
	//defer log.Close()

	//log.SetLevel(log.LEVEL_INFO)
	log.Debug("This is debug message")
	log.Info("This is info message")
	log.Warn("This is warn message")
	log.Error("This is error message")
	log.Fatal("This is fatal message")

	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {

		wg.Add(1)
		//go PrintFuncExecuteTime(i, wg)
		PrintFuncExecuteTime(i, wg)
	}

	wg.Wait()
	log.Leave()

	//打印方法执行调用次数、总时间、平均时间和错误次数
	log.Info("report summary: %v", log.Report())
}

func PrintFuncExecuteTime(i int, wg *sync.WaitGroup) {

	log.Enter("index", i) //enter PrintLog function (start statistic)

	var ip int32 = 10
	st1 := testSt{
		flt32:      0.58,
		flt64:      0.96666,
		ui8:        25,
		ui32:       10032,
		i8:         44,
		i64:        100000000000019,
		str:        "hello...",
		slice:      []string{"str1", "str2"},
		MyInt:      1,
		MyFloat64:  2.00,
		ip:         &ip,
		MySubSt:    TestSubSt{SubInt: 1, SubStr: "My sub str"},
		MySubStPtr: &TestSubSt{SubInt: 19, SubStr: "MySubStPtr", sst: &SubSubSt{Name: "I'm subsubst object"}}}

	st2 := &testSt{MyInt: 2, MyFloat64: 4.00, abc: 9999}
	log.Json(st1, st2)
	//log.Struct(st1, st2)

	log.Leave("index", i) //leave PrintLog function (end statistic)
	wg.Done()
}
