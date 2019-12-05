package log

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var LevelName = []string{"[DEBUG]", "[INFO]", "[WARN]", "[ERROR]", "[FATAL]", "[JSON]", "[STRUCT]"}

const (
	LEVEL_DEBUG  = 0
	LEVEL_INFO   = 1
	LEVEL_WARN   = 2
	LEVEL_ERROR  = 3
	LEVEL_FATAL  = 4
	LEVEL_JSON   = 5
	LEVEL_STRUCT = 6
)

type LogContent struct {
	FilePath string `json:"file_path"`
	LogLevel string `json:"log_level"`
	FileSize int    `json:"file_size"`
	Console  bool   `json:"console"`
}

type LogJson struct {
	LogCon LogContent `json:"log"`
}

type LogUrl struct {
	Path     string //文件路径
	Host     string //主机名（文件路径）
	Fragment string //无用
	Rawpath  string //无用
	Rawquery string //无用
	Scheme   string //协议名
	User     string //用户名
	Password string //密码
}

//全局变量
var logfile *os.File   //日志文件对象
var logger *log.Logger //日志输出对象
var logurl LogUrl      //URL解析对象
var loglevel int       //文件日志输出级别
var filepath string    //文件日志路径
var filesize int       //文件日志分割大小(MB)
var console = true     //开启/关闭终端屏幕输出

/**  打开日志
* 1. 通过参数直接指定日志文件、输出级别(DEBUG,INFO,WARN,ERROR, FATAL)和属性
*
*	1.1. 直接输入文件名
*	Open("test.log")
*
*	1.2. 设置文件日志输出级别和分块大小(单位：MB)
*  	Open("file:///var/log/test.log?log_level=INFO&file_size=50")
*
* 2. 通过指定json配置文件设置日志级别、日志文件及属性
*
*   2.1. 指定json配置文件
*   Open("json:///etc/test.json")
*
*   JSON范例：
*   {
*      "file_path":"/tmp/test.log",
*      "log_level":"INFO",
*      "file_size":"50",
*      "email_level":"FATAL",
*      "email_addr":"civet126@126.com",
*      "email_title":"error message title"
*   }
 */
//JSON配置文件例子
var jsonExample = `{
      "file_path":"/tmp/test.log",
      "log_level":"INFO",
      "file_size":"50"
   }`

func init() {
	filesize = 50 //MB
}

func Open(strUrl string) bool {

	if strUrl == "" {

		Error("Open url is nil")
		return false
	}

	err := logurl.ParseUrl(strUrl)
	if err != nil {
		Error("%s", err)
		return false
	}

	if logurl.Scheme == "json" { //以 'json://' 开头的URL

		err = logurl.ReadFromJson() //从JSON配置文件读取
		if err != nil {
			Error("%s", err)
			return false
		}

	} else if logurl.Scheme == "file" || logurl.Scheme == "" { //以 'file://' 开头的URL或者没有协议名

		return logurl.CreateFile() //创建文件
	} else {
		Error("Unknown scheme [%s]", logurl.Scheme)
	}

	return true
}

//关闭日志
func Close() {
	if logfile != nil {

		err := logfile.Close()
		if err != nil {
			Error("%s", err)
			return
		}
		logfile = nil
	}
}

//解析Url
func (lu *LogUrl) ParseUrl(strUrl string) (err error) {

	var querys url.Values
	u, err := url.Parse(strUrl)
	if err != nil {
		return
	}

	lu.Path = u.Path
	lu.Host = u.Host
	lu.Scheme = u.Scheme
	if u.User != nil {
		lu.User = u.User.Username()
		lu.Password, _ = u.User.Password()
	}
	querys, err = url.ParseQuery(u.RawQuery)
	if err != nil {
		return
	}

	filepath = lu.Host + lu.Path

	//Info("scheme [%s] host [%s] path [%s] querys [%s]", lu.Scheme, lu.Host, lu.Path, querys)
	for k, v := range querys {
		//Info("key = [%s] v = [%s]", k, v[0])
		switch k {
		case "file_size":
			filesize, _ = strconv.Atoi(v[0])
		case "log_level":
			loglevel = GetLevel(v[0])
		}
	}

	// Debug("filelevel [%d] filepath [%s]  filesize [%d] emaillevel [%d] emailaddr  [%s] emailtitle  [%s]",
	// 	filelevel, filepath, filesize, emaillevel, emailaddr, emailtitle)
	return
}

//创建日志文件
func (lu *LogUrl) CreateFile() bool {
	var err error
	logfile, err = os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		Error("Open log file ", filepath, " failed ", err)
		return false
	}

	logger = log.New(logfile, "", log.Ldate|log.Ltime|log.LstdFlags)

	return true
}

//从json文件加载配置
func (lu *LogUrl) ReadFromJson() (err error) {

	strJsonFile := filepath

	go func() { //定时读取JSON文件更新配置信息

		for {
			var logjson LogJson //json文件序列化对象
			//fmt.Println("JSON ", strJsonFile)

			file, err := ioutil.ReadFile(strJsonFile)
			if err != nil {
				fmt.Println("JSON [", strJsonFile, "] read error [", err, "]")
				time.Sleep(5 * time.Second)
				continue
			}

			err = json.Unmarshal(file, &logjson)
			if err != nil {
				fmt.Println("JSON [", strJsonFile, "] parse error [", err, "]")
				time.Sleep(5 * time.Second)
				continue
			}

			filepath = logjson.LogCon.FilePath
			loglevel = GetLevel(logjson.LogCon.LogLevel)
			filesize = logjson.LogCon.FileSize
			console = logjson.LogCon.Console
			if !console {
				fmt.Println(logjson)
				fmt.Println("Console output closed by ", strJsonFile)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	return
}

//截取函数名称
func GetFuncName(pc uintptr) (name string) {

	n := runtime.FuncForPC(pc).Name()
	ns := strings.Split(n, ".")
	name = ns[len(ns)-1]
	return
}

//设置日志级别(0=DEBUG 1=INFO 2=WARN 3=ERROR 4=FATAL)
func SetLevel(nLevel int) {

	loglevel = nLevel
}

//通过级别名称获取索引
func GetLevel(name string) (idx int) {

	name = "[" + name + "]"
	switch name {

	case LevelName[LEVEL_DEBUG]:
		idx = LEVEL_DEBUG
	case LevelName[LEVEL_INFO]:
		idx = LEVEL_INFO
	case LevelName[LEVEL_WARN]:
		idx = LEVEL_WARN
	case LevelName[LEVEL_ERROR]:
		idx = LEVEL_ERROR
	case LevelName[LEVEL_FATAL]:
		idx = LEVEL_FATAL
	default:
		idx = LEVEL_INFO
	}

	//Debug("Name [%s] level [%d]", name, idx)
	return
}

//内部格式化输出函数
func Output(level int, fmtstr string, args ...interface{}) {
	var inf, code string
	var colorName string

	if level < loglevel {
		return
	}

	Name := LevelName[level]
	switch level {
	case LEVEL_DEBUG:
		colorName = fmt.Sprintf("\033[34m%v %s", time.Now().Format("2006-01-02 15:04:05"), Name)
	case LEVEL_INFO:
		colorName = fmt.Sprintf("\033[32m%v %s", time.Now().Format("2006-01-02 15:04:05"), Name)
	case LEVEL_WARN:
		colorName = fmt.Sprintf("\033[33m%v %s", time.Now().Format("2006-01-02 15:04:05"), Name)
	case LEVEL_ERROR:
		colorName = fmt.Sprintf("\033[31m%v %s", time.Now().Format("2006-01-02 15:04:05"), Name)
	case LEVEL_FATAL:
		colorName = fmt.Sprintf("\033[35m%v %s", time.Now().Format("2006-01-02 15:04:05"), Name)
	default:
		colorName = fmt.Sprintf("\033[34m%v %s", time.Now().Format("2006-01-02 15:04:05"), Name)
	}

	if fmtstr != "" {
		inf = fmt.Sprintf(fmtstr, args...)
	} else {
		inf = fmt.Sprint(args...)
	}

	pc, file, line, ok := runtime.Caller(2)
	if ok {
		code = "<" + path.Base(file) + ":" + strconv.Itoa(line) + " " + GetFuncName(pc) + "()" + ">"
	}

	var output string

	switch runtime.GOOS {
	case "windows": //Windows终端不支持颜色显示
		output = time.Now().Format("2006-01-02 15:04:05") + " " + Name + " " + code + " " + inf
	default: //Unix类终端支持颜色显示
		output = "\033[1m" + colorName + " " + code + "\033[0m " + inf
	}

	//打印到终端屏幕
	if console {
		fmt.Println(output)
	}

	//输出到文件（如果Open函数传入了正确的文件路径）
	if logger != nil {
		fi, e := os.Stat(filepath)
		if e == nil {
			fs := fi.Size()
			if fs > int64(filesize*1024*1024) {

				logfile.Close()
				newpath := filepath + "." + time.Now().Format("2006-01-02_15_04_05")
				e = os.Rename(filepath, newpath) //将文件备份
				if e != nil {
					Error("%s", e)
					return
				} else {
					logurl.CreateFile() //重新创建文件
				}
			}
		}

		logger.Println(Name + " " + code + " " + inf)
	}
}

//输出调试级别信息
func Debug(fmtstr string, args ...interface{}) {
	Output(LEVEL_DEBUG, fmtstr, args...)
}

//输出运行级别信息
func Info(fmtstr string, args ...interface{}) {
	Output(LEVEL_INFO, fmtstr, args...)
}

//输出警告级别信息
func Warn(fmtstr string, args ...interface{}) {
	Output(LEVEL_WARN, fmtstr, args...)
}

//输出警告级别信息
func Warning(fmtstr string, args ...interface{}) {
	Output(LEVEL_WARN, fmtstr, args...)
}

//输出错误级别信息
func Error(fmtstr string, args ...interface{}) {
	Output(LEVEL_ERROR, fmtstr, args...)
}

//输出危险级别信息
func Fatal(fmtstr string, args ...interface{}) {
	Output(LEVEL_FATAL, fmtstr, args...)
}

//输出到空设备
func Null(fmtstr string, args ...interface{}) {

}

//打印结构体JSON
func Json(args ...interface{}) {

	if !console {
		return
	}

	for i := range args {
		arg := args[i]
		typ := reflect.TypeOf(arg)
		val := reflect.ValueOf(arg)
		if typ.Kind() == reflect.Ptr { //如果是指针类型则先转为对象

			typ = typ.Elem()
			val = val.Elem()
		}
		data, _ := json.MarshalIndent(arg, "", "\t")
		strTypeName := ""
		if typ.Name() == "" {
			strTypeName = "slice/base type"
		}
		Output(LEVEL_JSON, fmt.Sprintf("(%v) %v ", strTypeName, string(data)))
	}
}

//打印结构体
func Struct(args ...interface{}) {

	if !console {
		return
	}

	var strLog string
	for i := range args {
		arg := args[i]
		typ := reflect.TypeOf(arg)
		val := reflect.ValueOf(arg)
		if typ.Kind() == reflect.Ptr { //如果是指针类型则先转为对象

			typ = typ.Elem()
			val = val.Elem()
		}

		var nDeep int
		switch typ.Kind() {

		case reflect.Struct:
			strLog = fmtStruct(nDeep, typ, val) //遍历结构体成员标签和值存到map[string]string中
		default:
			fmt.Printf("<%v> = <%+v> \n", typ.Name(), val.Interface())
		}

		Output(LEVEL_STRUCT, strLog)
	}
}

//将字段值存到其他类型的变量中
func fmtStruct(deep int, typ reflect.Type, val reflect.Value, args...interface{}) (strLog string) {

	kind := typ.Kind()
	nCurDeep := deep

	var bPointer bool
	var strParentName string
	if len(args) > 0 {
		bPointer = args[0].(bool)
		strParentName = args[1].(string)
	}

	if !val.IsValid() {
		if bPointer { //this variant is a struct pointer
			strLog = fmt.Sprintf("%v%v (*%v) = <nil>\n", fmtDeep(deep), strParentName, typ.String())
		} else {
			strLog = fmt.Sprintf("%v%v (%v) = <nil>\n", fmtDeep(deep), strParentName, typ.String())
		}
		return
	}

	if bPointer {//this variant is a struct pointer
		//strLog = fmt.Sprintf("%v%v (*%v) {\n", fmtDeep(deep) , typ.Kind().String(), typ.String())
		strLog = fmt.Sprintf("%v%v (*%v) {\n", fmtDeep(deep) , strParentName, typ.String())
	} else {
		strLog = fmt.Sprintf("%v%v (%v) {\n", fmtDeep(deep) , strParentName, typ.String())
	}


	if kind == reflect.Struct {
		deep++
		NumField := val.NumField()
		for i := 0; i < NumField; i++ {

			var isPointer bool
			typField := typ.Field(i)
			valField := val.Field(i)
			if typField.Type.Kind() == reflect.Ptr { //如果是指针类型则先转为对象

				typField.Type = typField.Type.Elem()
				valField = valField.Elem()
				isPointer = true
			}

			if typField.Type.Kind() == reflect.Struct {

				strLog += fmtStruct(deep, typField.Type, valField, isPointer, typField.Name) //结构体需要递归调用
			} else {
				//var strLog string
				if !valField.IsValid() { //字段为空指针
					strLog += fmtDeep(deep) + fmt.Sprintf("%s = <nil> \n", typField.Name)
				} else if !valField.CanInterface() { //非导出字段
					strLog += fmtDeep(deep) + fmt.Sprintf("%s = <%+v> \n", typField.Name, valField)
				} else {

					switch typField.Type.Kind() {
					case reflect.String:
						strLog += fmtDeep(deep) + fmt.Sprintf("%s = <%+v> \n", typField.Name, valField.Interface())
					default:
						strLog += fmtDeep(deep) + fmt.Sprintf("%s = <%+v> \n", typField.Name, valField.Interface())
					}

				}
			}
		}
	}
	strLog += fmtDeep(nCurDeep) + "}\n"
	return
}

func fmtDeep(nDeep int) (s string) {

	for i := 0; i < nDeep; i++ {
		s += fmt.Sprintf("... ")
	}
	return
}

func writeLogDriectly(strLog string) {
	if logger != nil {
		logger.Print(strLog)
	}
}
