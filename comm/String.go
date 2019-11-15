package comm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//字符串转float
func Str2Float(in string) (out float64) {

	out, _ = strconv.ParseFloat(in, 64)
	return
}

//字符串转int
func Str2Int64(in string) (out int64) {
	out, _ = strconv.ParseInt(in, 10, 64)
	return
}

func Str2Int(in string) (out int) {
	i32 , _ := strconv.ParseInt(in, 10, 32)
	out = int(i32)
	return
}

//数值类型转字符串
func Decimal2Str(in interface{}) (out string) {

	out = fmt.Sprintf("%v", in)
	return
}

//UTF-8字符串查找(按UTF-8字符算一个位置)
func StrIndexUtf8(strIn, strSub string) (RuneIdx int){

	// 子串在字符串的字节位置
	RuneIdx = strings.Index(strIn, strSub)
	if RuneIdx >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := []byte(strIn)[0:RuneIdx]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(string(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		RuneIdx = len(rs)
		fmt.Println("StrIndexUtf8 原始字符串[", strIn,"]子串[", strSub,"]所在位置[", RuneIdx, "]")
	}

	return RuneIdx
}

//UTF-8字符串截取子串后字符串(bInclude true 表示包括子串，false 表示不包括)
func Utf8StrGetTail(strIn, strSub string, bInclude bool) (bFound bool, strOut string){

	var idx int
	strOut = strIn
	idx = strings.Index(strIn, strSub)
	//fmt.Println("Utf8StrGetFront len(strSub) =", len(strSub))
	if idx != -1 {

		bFound = true
		if bInclude {
			strOut = string(strIn[idx:])
			//fmt.Println("Utf8StrGetTail 原始字符串[", strIn,"]子串[", strSub,"]截取字符(包含子串)串结果[", strOut,"]")
		}else{
			strOut = string(strIn[idx+len(strSub):])
			//fmt.Println("Utf8StrGetTail 原始字符串[", strIn,"]子串[", strSub,"]截取字符串结果[", strOut,"]")
		}
	}

	return
}


//UTF-8字符串截取子串前面的字符串(bInclude true 表示包括子串，false 表示不包括)
func Utf8StrGetFront(strIn, strSub string, bInclude bool) (bFound bool, strOut string){

	var idx int
	strOut = strIn
	idx = strings.Index(strIn, strSub)
	//fmt.Println("Utf8StrGetFront len(strSub) =", len(strSub))
	if idx != -1 {

		bFound = true
		if bInclude {
			strOut = string(strIn[:idx+len(strSub)])
			//fmt.Println("Utf8StrGetFront 原始字符串[", strIn,"]子串[", strSub,"]截取字符(包含子串)串结果[", strOut,"]")
		}else{
			strOut = string(strIn[:idx])
			//fmt.Println("Utf8StrGetFront 原始字符串[", strIn,"]子串[", strSub,"]截取字符串结果[", strOut,"]")
		}
	}
	return
}

//通过正则表达式提取符合规则的子串
func RegexpSubStr(strIn, strRegExp string) ([]string) {

	comp := regexp.MustCompile(strRegExp)
	return comp.FindStringSubmatch(strIn)
}