package log

import (
	"encoding/json"
	"fmt"
	"github.com/mattn/go-colorable"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var colorStdout = colorable.NewColorableStdout()

var LevelName = []string{"[DEBUG]", "[INFO]", "[WARN]", "[ERROR]", "[FATAL]"}

const (
	LEVEL_DEBUG = 0
	LEVEL_INFO  = 1
	LEVEL_WARN  = 2
	LEVEL_ERROR = 3
	LEVEL_FATAL = 4
	LEVEL_PANIC = 5
)

type LogContent struct {
	FilePath   string `json:"file_path"`
	LogLevel   string `json:"log_level"`
	FileSize   int    `json:"file_size"`
	MaxBackups int    `json:"max_backups"`
	Console    bool   `json:"console"`
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
	locker   sync.RWMutex
}

type Option struct {
	LogLevel     int    //文件日志输出级别
	FileSize     int    //文件日志分割大小(MB)
	MaxBackups   int    //文件最大分割数
	CloseConsole bool   //开启/关闭终端屏幕输出
	filePath     string //文件日志路径

}

//全局变量
var (
	logFile *os.File    //日志文件对象
	logger  *log.Logger //日志输出对象
	logUrl  LogUrl      //URL解析对象
	option  Option      //日志参数选项
)

/**  打开日志
* 1. 通过参数直接指定日志文件、输出级别(DEBUG,INFO,WARN,ERROR, FATAL)和属性
*
*	1.1. 直接输入文件名
*	Open("test.log")
*
*	1.2. 设置文件日志输出级别和分块大小(单位：MB)以及备份文件数
*  	Open("file:///var/log/test.log?log_level=INFO&file_size=50&max_backups=10")
*
* 2. 通过指定json配置文件设置日志级别、日志文件及属性
*
*   2.1. 指定json配置文件
*   Open("json:///etc/test.json")
*
   JSON范例：
   {
      "file_path":"/tmp/test.log",
      "log_level":"INFO",
      "file_size":"1024",
      "max_backups": 10,
      "console": true,
   }
*/

//var colorStdout = colorable.NewColorableStdout()

func init() {
	option.FileSize = 1024 //MB
	option.MaxBackups = 31
	option.LogLevel = LEVEL_DEBUG
	go cleanBackupLog()
}

func Open(strUrl string, opts ...Option) bool {

	if strUrl == "" {

		Error("Open url is nil")
		return false
	}

	err := logUrl.parseUrl(strUrl)
	if err != nil {
		Error("%s", err)
		return false
	}

	if logUrl.Scheme == "json" { //以 'json://' 开头的URL

		err = logUrl.readFromJson() //从JSON配置文件读取
		if err != nil {
			Error("%s", err)
			return false
		}

	} else if logUrl.Scheme == "file" || logUrl.Scheme == "" { //以 'file://' 开头的URL或者没有协议名

		return logUrl.createFile() //创建文件
	} else {
		Error("Unknown scheme [%s]", logUrl.Scheme)
	}

	if len(opts) > 0 {
		option = opts[0]
	}
	return true
}

//关闭日志
func Close() {
	if logFile != nil {

		err := logFile.Close()
		if err != nil {
			Error("%s", err)
			return
		}
		logFile = nil
	}
}

//设置日志文件分割大小（MB)
func SetFileSize(size int) {
	option.FileSize = size
}

//设置日志级别(字符串型: debug/info/warn/error/fatal 数值型: 0=DEBUG 1=INFO 2=WARN 3=ERROR 4=FATAL)
func SetLevel(level interface{}) {

	var nLevel int
	switch level.(type) {
	case string:
		strLevel := strings.ToLower(level.(string))
		switch strLevel {
		case "debug":
			nLevel = 0
		case "info":
			nLevel = 1
		case "warn", "warning":
			nLevel = 2
		case "error":
			nLevel = 3
		case "fatal":
			nLevel = 4
		}
	case int8, int16, int, int32, int64, uint8, uint16, uint, uint32, uint64:
		nLevel, _ = strconv.Atoi(fmt.Sprintf("%v", level))
	default:
		panic("not support yet")
	}
	option.LogLevel = nLevel
}

//设置关闭/开启屏幕输出
func CloseConsole(ok bool) {
	option.CloseConsole = ok
}

//设置最大备份数量
func SetMaxBackup(nMaxBackups int) {
	option.MaxBackups = nMaxBackups
}

//定期清理日志，仅保留MaxBackups个数的日志
func cleanBackupLog() {
	for {

		time.Sleep(1 * time.Hour)
	}
}

//解析Url
func (lu *LogUrl) parseUrl(strUrl string) (err error) {

	var querys url.Values
	var u *url.URL
	u, err = url.Parse(strUrl)
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

	option.filePath = lu.Host + lu.Path

	//Info("scheme [%s] host [%s] path [%s] querys [%s]", lu.Scheme, lu.Host, lu.Path, querys)
	for k, v := range querys {
		//Info("key = [%s] v = [%s]", k, v[0])
		switch k {
		case "file_size":
			option.FileSize, _ = strconv.Atoi(v[0])
		case "log_level":
			option.LogLevel = getLevel(v[0])
		case "max_backups":
			option.MaxBackups, _ = strconv.Atoi(v[0])
		case "console":
			Console, _ := strconv.ParseBool(v[0])
			option.CloseConsole = !Console
		}
	}

	return
}

//创建日志文件
func (lu *LogUrl) createFile() bool {
	var err error

	//判断日志文件后缀名合法性
	//if strings.Index(option.filePath, ".") == -1 {
	//	panic("log file path illegal, must contain dot suffix [日志文件必须带.后缀名]")
	//}
	lu.locker.Lock()
	defer lu.locker.Unlock()
	logFile, err = os.OpenFile(option.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		Error("Open log file ", option.filePath, " failed ", err)
		return false
	}

	logger = log.New(logFile, "", log.Lmicroseconds|log.LstdFlags)

	return true
}

//从json文件加载配置
func (lu *LogUrl) readFromJson() (err error) {

	strJsonFile := option.filePath

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

			option.filePath = logjson.LogCon.FilePath
			option.LogLevel = getLevel(logjson.LogCon.LogLevel)
			option.FileSize = logjson.LogCon.FileSize
			option.CloseConsole = !logjson.LogCon.Console
			if option.CloseConsole {
				fmt.Println(logjson)
				fmt.Println("Console output closed by ", strJsonFile)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	return
}

//截取函数名称
func getFuncName(pc uintptr) (name string) {

	n := runtime.FuncForPC(pc).Name()
	ns := strings.Split(n, ".")
	name = ns[len(ns)-1]
	return
}

//通过级别名称获取索引
func getLevel(name string) (idx int) {

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

	return
}

func getCaller(skip int) (strFile, strFunc string, nLineNo int) {
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		strFile = path.Base(file)
		nLineNo = line
		strFunc = getFuncName(pc)
	}
	return
}

func getStack(skip, n int) string {
	var strStack string
	for i := 0; i < n; i++ {
		pc, file, line, ok := runtime.Caller(skip + i)
		if ok {
			strFile := path.Base(file)
			nLineNo := line
			strFunc := getFuncName(pc)
			strStack += fmt.Sprintf("|- %s:%d %s() \n", strFile, nLineNo, strFunc)
		}
	}
	return strStack
}

//内部格式化输出函数
func output(level int, fmtstr string, args ...interface{}) (strFile, strFunc string, nLineNo int) {
	var inf, code string
	var colorTimeName string

	strTimeFmt := fmt.Sprintf("%v", time.Now().Format("2006-01-02 15:04:05.000000"))
	strRoutine := fmt.Sprintf("{%v}", getRoutineId())
	strPID := fmt.Sprintf("PID:%d", os.Getpid())
	Name := LevelName[level]
	switch level {
	case LEVEL_DEBUG:
		colorTimeName = fmt.Sprintf("\033[34m%v %s %s", strTimeFmt, strPID, Name)
	case LEVEL_INFO:
		colorTimeName = fmt.Sprintf("\033[32m%v %s %s", strTimeFmt, strPID, Name)
	case LEVEL_WARN:
		colorTimeName = fmt.Sprintf("\033[33m%v %s %s", strTimeFmt, strPID, Name)
	case LEVEL_ERROR:
		colorTimeName = fmt.Sprintf("\033[31m%v %s %s", strTimeFmt, strPID, Name)
	case LEVEL_FATAL:
		colorTimeName = fmt.Sprintf("\033[35m%v %s %s", strTimeFmt, strPID, Name)
	default:
		colorTimeName = fmt.Sprintf("\033[34m%v %s %s", strTimeFmt, strPID, Name)
	}

	if fmtstr != "" {
		inf = fmt.Sprintf(fmtstr, args...)
	} else {
		inf = fmt.Sprint(args...)
	}

	strFile, strFunc, nLineNo = getCaller(3)
	code = "<" + strFile + ":" + strconv.Itoa(nLineNo) + " " + strFunc + "()" + ">"
	if level < option.LogLevel {
		return
	}

	var output string

	switch runtime.GOOS {
	//case "windows": //Windows终端不再支持颜色显示
	//output = strTimeFmt + " " + Name + " " + strRoutine + " " + code + " " + inf
	default: //Unix类终端支持颜色显示
		output = "\033[1m" + colorTimeName + " " + strRoutine + " " + code + "\033[0m " + inf
	}

	if level >= LEVEL_ERROR {
		output += fmt.Sprintf("\n" + "call stack \n" + getStack(3, 10))
	}
	//打印到终端屏幕
	if !option.CloseConsole {
		_, _ = fmt.Fprintln(colorStdout /*os.Stdout*/, output)
	}

	//输出到文件（如果Open函数传入了正确的文件路径）
	if logger != nil {
		fi, e := os.Stat(option.filePath)
		if e == nil {
			fs := fi.Size()
			if fs > int64(option.FileSize*1024*1024) {

				logFile.Close()
				datetime := time.Now().Format("20060102_150405")
				var newpath string
				newpath = fmt.Sprintf("%v.%v", option.filePath, datetime) //日志文件有后缀(日志备份文件名格式不能随意改动)
				e = os.Rename(option.filePath, newpath)                   //将文件备份
				if e != nil {
					Error("%s", e)
					return
				} else {
					logUrl.createFile() //重新创建文件
				}
			}
		} else {
			logUrl.createFile() //重新创建文件
		}

		logger.Println(Name + " " + strRoutine + " " + code + " " + inf)
	}
	return
}

func fmtString(args ...interface{}) (strOut string) {
	if len(args) > 0 {
		switch args[0].(type) {
		case string:
			if strings.Contains(args[0].(string), "%") {
				strOut = fmt.Sprintf(args[0].(string), args[1:]...)
			} else {
				strOut = fmt.Sprint(args...)
			}
		default:
			strOut = fmt.Sprint(args...)
		}
	}
	return
}

//输出调试级别信息
func Debug(args ...interface{}) {
	output(LEVEL_DEBUG, fmtString(args...))
}

//输出运行级别信息
func Info(args ...interface{}) {
	output(LEVEL_INFO, fmtString(args...))
}

//输出警告级别信息
func Warn(args ...interface{}) {
	output(LEVEL_WARN, fmtString(args...))
}

//输出警告级别信息
func Warning(args ...interface{}) {
	output(LEVEL_WARN, fmtString(args...))
}

//输出错误级别信息
func Error(args ...interface{}) {
	stic.error(output(LEVEL_ERROR, fmtString(args...)))
}

//输出危险级别信息
func Fatal(args ...interface{}) {
	stic.error(output(LEVEL_FATAL, fmtString(args...)))
}

//panic
func Panic(args ...interface{}) {
	panic(fmt.Sprintf(fmtString(args...)))
}

//输出调试级别信息
func Debugf(fmtstr string, args ...interface{}) {
	output(LEVEL_DEBUG, fmtstr, args...)
}

//输出运行级别信息
func Infof(fmtstr string, args ...interface{}) {
	output(LEVEL_INFO, fmtstr, args...)
}

//输出警告级别信息
func Warnf(fmtstr string, args ...interface{}) {
	output(LEVEL_WARN, fmtstr, args...)
}

//输出警告级别信息
func Warningf(fmtstr string, args ...interface{}) {
	output(LEVEL_WARN, fmtstr, args...)
}

//输出错误级别信息
func Errorf(fmtstr string, args ...interface{}) {
	stic.error(output(LEVEL_ERROR, fmtstr, args...))
}

//输出危险级别信息
func Fatalf(fmtstr string, args ...interface{}) {
	stic.error(output(LEVEL_FATAL, fmtstr, args...))
}

//输出到空设备
func Null(fmtstr string, args ...interface{}) {
}

//进入方法（统计）
func Enter(args ...interface{}) {
	output(LEVEL_INFO, "enter ", args...)
	stic.enter(getCaller(2))
}

//离开方法（统计）
//返回执行时间：h 时 m 分 s 秒 ms 毫秒 （必须先调用Enter方法才能正确统计执行时间）
func Leave() (h, m, s int, ms float32) {

	if nSpendTime, ok := stic.leave(getCaller(2)); ok {
		h, m, s, ms = getSpendTime(nSpendTime)
		output(LEVEL_INFO, "leave (%vh %vm %vs %.3fms)", h, m, s, ms)
	}
	return
}

//打印结构体JSON
func Json(args ...interface{}) {

	var strOutput string

	for _, v := range args {

		data, _ := json.MarshalIndent(v, "", "\t")
		strOutput += "\n...................................................\n" + string(data)
	}

	output(LEVEL_DEBUG, strOutput+"\n...................................................\n")
}

func JsonDebugString(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "\t")
	return string(data)
}

//args: a string of function name or nil for all
func Report(args ...interface{}) string {
	return stic.report(args...)
}

//打印结构体
func Struct(args ...interface{}) {

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
		case reflect.String:
			strLog += fmt.Sprintf("%v (string) = \"%+v\" \n", typ.Name(), val.Interface())
		default:
			strLog += fmt.Sprintf("%v (%v) = <%+v> \n", typ.Name(), typ.Kind(), val.Interface())
		}

		output(LEVEL_DEBUG, strLog)
	}
}

//将字段值存到其他类型的变量中
func fmtStruct(deep int, typ reflect.Type, val reflect.Value, args ...interface{}) (strLog string) {

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

	if bPointer { //this variant is a struct pointer
		//strLog = fmt.Sprintf("%v%v (*%v) {\n", fmtDeep(deep) , typ.Kind().String(), typ.String())
		strLog = fmt.Sprintf("%v%v (*%v) {\n", fmtDeep(deep), strParentName, typ.String())
	} else {
		strLog = fmt.Sprintf("%v%v (%v) {\n", fmtDeep(deep), strParentName, typ.String())
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
					strLog += fmtDeep(deep) + fmt.Sprintf("%v (%v) = <nil> \n", typField.Name, typField.Type)
				} else if !valField.CanInterface() { //非导出字段
					strLog += fmtDeep(deep) + fmt.Sprintf("%v (%v) = <%+v> \n", typField.Name, typField.Type, valField)
				} else {

					switch typField.Type.Kind() {
					case reflect.String:
						strLog += fmtDeep(deep) + fmt.Sprintf("%v (%v) = \"%+v\" \n", typField.Name, typField.Type, valField.Interface())
					default:
						strLog += fmtDeep(deep) + fmt.Sprintf("%v (%v) = <%+v> \n", typField.Name, typField.Type, valField.Interface())
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
