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

type statistic struct {
	callers sync.Map
}

type caller struct {
	KeyName   string `json:"key_name"`
	FileName  string `json:"file_name"`
	LineNo    int    `json:"line_no"`
	FuncName  string `json:"func_name"`
	EnterTime int64  `json:"enter_time"`
	LeaveTime int64  `json:"leave_time"`
	SpendTime int    `json:"spend_time"`
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
	}

	return -1, false
}

//统计信息汇总(statistic summary)
func (s *statistic) summary(args ...interface{}) (strSummary string) {

	return
}
