package log

import (
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

/*
  统计每个文件的方法被调用的总次数、总调用时间、平均执行时间（毫秒）
*/

const (
	EXPIRE_TIME_MICRO_SECONDS = 24*3600*1000000
)

type statistic struct {
	callers sync.Map
}

type caller struct {
	KeyName    string `json:"key_name"`
	FileName   string `json:"file_name"`
	LineNo     int    `json:"line_no"`
	FuncName   string `json:"func_name"`
	EnterTime  int64  `json:"enter_time"`
	LeaveTime  int64  `json:"leave_time"`
	SpendTime  int    `json:"spend_time"`
	ExpireTime int64  `json:"expire_time"`
}

//create a new statistic object
func newStatistic() *statistic {
	return &statistic{
	}
}

func getMilliSec() int64 {

	return time.Now().UnixNano() / 1e6 //milliseconds
}

func getMicroSec() int64 {

	return time.Now().UnixNano() / 1e3 //micro seconds
}

func getCallerStoreKey(strFile, strFunc string) string {

	return fmt.Sprintf("%v %v %v", getRoutine(), strFile, strFunc)
}

func getRoutine() string {

	strStack := string(debug.Stack())
	nIdx := strings.IndexAny(strStack, "\r\n")
	if nIdx > 0 {
		return strStack[:nIdx]
	}
	return "<unknown routine>"
}

//进入方法(enter function)
func (s *statistic) enter(strFile, strFunc string, nLineNo int) {

	now64 := getMicroSec()

	c := caller{
		FileName:  strFile,
		LineNo:    nLineNo,
		FuncName:  strFunc,
		EnterTime: now64,
		LeaveTime: 0,
		SpendTime: 0,
		ExpireTime: now64+EXPIRE_TIME_MICRO_SECONDS,
	}
	c.KeyName = getCallerStoreKey(strFile, strFunc)
	s.callers.Store(c.KeyName, &c)
	//Debug("caller store ok")
}

//退出方法(leave function)
func (s *statistic) leave(strFile, strFunc string, nLineNo int) (int, bool) {

	now64 := getMicroSec()
	strKey := getCallerStoreKey(strFile, strFunc)
	if v, ok := s.callers.Load(strKey); ok {

		s.callers.Delete(strKey)
		c := v.(*caller)
		c.LeaveTime = now64
		c.SpendTime = int(c.LeaveTime - c.EnterTime)
		Json("leave caller", c)
		return c.SpendTime, ok
	} else {
		//没有找到调用log.Enter记录，需注意log.Enter和log.Leave必须同时在一个方法中使用
		Warn("[zh_CN] 请在%v中调用log.Leave之前调用log.Enter [en_US] please call log.Enter before log.Leave in %v function.", strFunc, strFunc)
		//panic(fmt.Sprintf("[en_US] not called log.Enter before log.Leave in [%v] function ? Key [%v]", strFunc, strKey)) //not called Enter before call Leave ?
	}

	return -1, false
}

//统计信息汇总(statistic summary)
func (s *statistic) summary(args ...interface{}) (strSummary string) {

	return
}
