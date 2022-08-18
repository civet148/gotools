package comm

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//DateTimeUnix 生成当前时间Unix时间戳
func DateTimeUnix() int64 {
	return time.Now().Unix()
}

//DateTimeStr 生成当前时间字符串
func DateTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//DateTimeUnix2Str Unix时间戳转时间字符串(格式："2017-01-02 03:04:05")
//Timestamp 时间戳可以是int64或string类型
//strFmt 明确指出时间字符串格式，"/"表示 "2017/01/06 03:04:05" 默认不传参数为"2017-01-02 03:04:05"格式
func DateTimeUnix2Str(Timestamp interface{}, timeFmt... interface{}) string {

	var nTimestamp int64
	switch Timestamp.(type) {
	case string:
		nTimestamp, _ = strconv.ParseInt(Timestamp.(string), 10, 64)
	case int64:
		nTimestamp = Timestamp.(int64)
	}
	if nTimestamp == 0{
		return ""
	}
	t := time.Unix(nTimestamp, 0)

	if len(timeFmt) > 0 {
		if timeFmt[0].(string) == "/" {
			return t.Format("2006/01/02 15:04:05")
		}
	}
	return t.Format("2006-01-02 15:04:05")
}

//DateTimeStr2Unix 时间字符串转时间戳(时间字符串格式："2017-01-02 03:04:05")
//strFmt 明确指出时间字符串格式，"/"表示 "2017/01/06 03:04:05" 默认不传参数为"2017-01-02 03:04:05"格式
func DateTimeStr2Unix(strDateTime string, timeFmt... interface{}) (unixTime int64) {

	var t time.Time
	var bNormal = true


	if len(timeFmt) > 0 {

		if timeFmt[0].(string) == "/" {
			bNormal = false
		}
	} else {
		if strings.Contains(strDateTime, "/") {
			bNormal = false
		}
	}

	if len(strDateTime) != 19 { //时间格式不正确(自动补全)

		nIndex := strings.Index(strDateTime, " ") //找到空格
		if nIndex == -1 {//无效时间
			fmt.Println("error: DateTimeStr2Unix invalid datetime format")
			return 0
		}

		sdt := strings.Split(strDateTime, " ")
		if len(sdt)  == 1 {
			fmt.Println("error: DateTimeStr2Unix invalid datetime format")
			return 0
		}

		ymd := sdt[0] //年月日
		hms := sdt[1] //时分秒

		var s1, s2 []string
		if bNormal {//补全格式："2017-01-02 03:04:05"
			s1 = strings.Split(ymd, "-")

		}else{//补全格式："2017/01/06 03:04:05"
			s1 = strings.Split(ymd, "/")
		}

		s2 = strings.Split(hms, ":")

		if len(s1) != 3 || len(s2) != 3 {
			fmt.Println("error: DateTimeStr2Unix invalid datetime format, not match 'YYYY-MM-DD hh:mm:ss' or 'YYYY/MM/DD hh:mm:ss'")
			return 0
		}
		year := s1[0]; month := s1[1]; day := s1[2]
		hour := s2[0]; min := s2[1]; sec := s2[2]
		if len(year) != 4 {
			fmt.Println("error: DateTimeStr2Unix invalid year format, not match YYYY")
			return 0
		}
		if len(month) == 1 {
			month = "0"+month
		}
		if len(day) == 1 {
			day = "0"+day
		}
		if len(hour) == 1 {
			hour = "0"+hour
		}
		if len(min) == 1 {
			min = "0"+min
		}
		if len(sec) == 1 {
			sec = "0"+sec
		}

		if bNormal {//补全格式："2017-01-02 03:04:05"
			strDateTime = fmt.Sprintf("%v-%v-%v %v:%v:%v", year, month, day, hour, min, sec)

		}else{//补全格式："2017/01/06 03:04:05"
			strDateTime = fmt.Sprintf("%v/%v/%v %v:%v:%v", year, month, day, hour, min, sec)
		}
	}

	if strDateTime != "" {

		//UTC时区要转为中国
		loc, _ := time.LoadLocation("Local")

		if bNormal {
			t, _ = time.ParseInLocation("2006-01-02 15:04:05", strDateTime, loc)
		} else {
			t, _ = time.ParseInLocation("2006/01/02 15:04:05", strDateTime, loc)
		}

		unixTime = t.Unix()
	}

	return
}

func ExcelTime2Unix(v interface{}) (nUnixTime int64){
	/*
	  日期在EXCEL中是以数值存放的。他的基准点就是1900.1.1 00：00：00.按秒进行存放。
	  5.23，正好是1900-1-5 5:45:36，你可以算，小数点前代表天。小数点后带表时间。即1代表1天，1天的0.23正好是5小时45分36秒
	  Excel中43300.4559435185 实际时间为 2018-07-19 10:56:34
	*/

	strDateTime := ExcelTime2DateTime(v)

	return DateTimeStr2Unix(strDateTime)
}

//Excel存储的数值时间转为正常的时间字符串
//支持参数类型: string、float32、float64、int、int32、int64
func ExcelTime2DateTime(v interface{}) (strDateTime string){
	/*
	  日期在EXCEL中是以数值存放的。他的基准点就是1900.1.1 00：00：00.按秒进行存放。
	  5.23，正好是1900-1-5 5:45:36，你可以算，小数点前代表天。小数点后带表时间。即1代表1天，1天的0.23正好是5小时45分36秒
	  Excel中43300.4559435185 实际时间为 2018-07-19 10:56:34(因1900年有366天)
	*/
	var timeFloat float64
	var nBaseTime = DateTimeStr2Unix("1900-01-01 00:00:00")
	switch v.(type){

	case string:
		{
			timeFloat, _ = strconv.ParseFloat(v.(string),64)
		}
	case float32:
		{
			timeFloat = float64(v.(float32))
		}
	case float64:
		{
			timeFloat = v.(float64)
		}
	case int:
		{
			timeFloat = float64(v.(int))
		}
	case int32:
		{
			timeFloat = float64(v.(int32))
		}
	case int64:
		{
			timeFloat = float64(v.(int64))
		}
	default:
		{
			return "" //不支持的类型返回空字符串
		}
	}

	//log.Debug("timeFloat=[%v]", timeFloat)
	nIntervalTime := int64(timeFloat*24*3600) - 2*24*3600 //因1900年有366天，需减去2天的时间
	//log.Debug("nIntervalTime=[%v]", nIntervalTime)
	nDateTime := nBaseTime+nIntervalTime
	//log.Debug("nDateTime=[%v]", nDateTime)
	strDateTime = DateTimeUnix2Str(nDateTime)
	return
}


//将时间戳转为time.Time类型(Timestamp参数支持string和int64类型)
func DateTimeUnix2Time(Timestamp interface{}) (t time.Time) {

	var nTimestamp int64
	switch Timestamp.(type) {
	case string:
		nTimestamp, _ = strconv.ParseInt(Timestamp.(string), 10, 64)
	case int64:
		nTimestamp = Timestamp.(int64)
	}
	t = time.Unix(nTimestamp, 0)
	return
}

//将时间转日期(支持输入'2810-10-28 23:32:00'或 int32/64类型的时间戳)
//format 可输入字符串"/"或"-"，不输入默认使用"-"格式化日期
func DateTime2Date(DateTime interface{}, formats ...interface{}) (strDate string) {

	var nTimestamp int64
	var strDateTime string
	var strFormat = "-"

	switch DateTime.(type) {
	case string:
		strDateTime = DateTime.(string)
	case int:
		{
			nTimestamp = int64(DateTime.(int))
		}
	case int32:
		{
			nTimestamp = int64(DateTime.(int32))
		}
	case int64:
		{
			nTimestamp = DateTime.(int64)
		}
	}
	//fmt.Println("DateTime2Date nTimestamp=", nTimestamp)
	if nTimestamp == 0 {
		nTimestamp = DateTimeStr2Unix(strDateTime)
	}

	if nTimestamp == 0 {
		nTimestamp = DateTimeUnix()//如果没有输入字符串日期或正确的时间戳则返回当天日期
	}

	if len(formats) > 0 {
		switch formats[0].(type) {
		case string:
			strFormat = formats[0].(string)
		default:
			strFormat = "-"
		}
	}

	t := time.Unix(nTimestamp, 0)
	year, month, day := t.Date()
	if strFormat == "/" {
		strDate = fmt.Sprintf("%04d/%02d/%02d", year, month, day)
	}else if strFormat == "-" {
		strDate = fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	}

	return
}
