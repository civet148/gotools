package log

import (
	"encoding/json"
	"fmt"
	"github.com/civet148/gotools/comm"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

/*
  统计: 每个文件的方法被调用的总次数、总调用时间、平均执行时间（微秒）
  statistic: every function called total counts, executed total time, average execute time on micro seconds
*/

const (
	EXPIRE_TIME_MICRO_SECONDS = 24 * 3600 * 1e6
)

var (
	FUNCNAME_ALL = "all"
	FUNCNAME_NIL = ""
)

type statistic struct {
	callers sync.Map
	results sync.Map
	mutex sync.Mutex
}

type caller struct {
	KeyName    string `json:"key_name"`
	FileName   string `json:"file_name"`
	LineNo     int    `json:"line_no"`
	FuncName   string `json:"func_name"`
	EnterTime  int64  `json:"enter_time"`
	LeaveTime  int64  `json:"leave_time"`
	SpendTime  int64  `json:"spend_time"`
	ExpireTime int64  `json:"expire_time"`
	CallOk     bool   `json:"call_ok"`
}

type result struct {
	FileName   string `json:"file_name"`   //code file of function
	LineNo     int    `json:"line_no"`     //line no of function
	FuncName   string `json:"func_name"`   //function name
	CallCount  int64  `json:"call_count"`  //call total times
	ErrorCount int64  `json:"error_count"` //call error times
	ExeTime    int64  `json:"exe_time"`    //micro seconds
	AvgTime    int64  `json:"avg_time"`    //micro seconds
	CreateTime int64  `json:"create_time"` //unix timestamp on seconds
	UpdateTime int64  `json:"update_time"` //unix timestamp on seconds
}

type summary struct {
	TimeUnit string    `json:"time_unit"`
	Results  []*result `json:"statistics"`
}

//create a new statistic object
func newStatistic() *statistic {
	return &statistic{
	}
}

func getUnixSecond() int64 {
	return time.Now().Unix()
}

func getMilliSec() int64 {

	return time.Now().UnixNano() / 1e6 //milliseconds
}

func getMicroSec() int64 {

	return time.Now().UnixNano() / 1e3 //microseconds
}

func getDecimal3(v float64) float64 {

	return comm.Round(v, 3)
}

func getSpendTime(microseconds int64) (h, m, s int, ms float32) {
	nSpend := microseconds / 1e6
	h = int(nSpend / 3600)
	m = int((nSpend % 3600) / 60)
	s = int((nSpend % 3600) % 60)
	ms = float32(microseconds-(nSpend*1e6)) / 1000
	return
}

func getCallerStoreKey(strFile, strFunc string) string {

	return fmt.Sprintf("%v %v %v", getRoutine(), strFile, strFunc)
}

func getResultStoreKey(strFile, strFunc string) string {

	return fmt.Sprintf("%v:%v", strFile, strFunc)
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
		FileName:   strFile,
		LineNo:     nLineNo,
		FuncName:   strFunc,
		EnterTime:  now64,
		LeaveTime:  0,
		SpendTime:  0,
		ExpireTime: now64 + EXPIRE_TIME_MICRO_SECONDS,
		CallOk:     true,
	}
	c.KeyName = getCallerStoreKey(strFile, strFunc)
	s.callers.Store(c.KeyName, &c)

	//Debug("caller store ok")
	strResultKey := getResultStoreKey(strFile, strFunc)
	var r *result
	if _, ok := s.results.Load(strResultKey); !ok {
		r = &result{
			FileName:   strFile,
			LineNo:     nLineNo,
			FuncName:   strFunc,
			CallCount:  0,
			ErrorCount: 0,
			ExeTime:    0,
			AvgTime:    0,
			CreateTime: getUnixSecond(),
			UpdateTime: getUnixSecond(),
		}
		s.results.Store(strResultKey, r)
	}
}

//退出方法(leave function)
func (s *statistic) leave(strFile, strFunc string, nLineNo int) (int64, bool) {

	now64 := getMicroSec()
	strCallerKey := getCallerStoreKey(strFile, strFunc)
	strResultKey := getResultStoreKey(strFile, strFunc)

	if v, ok := s.callers.Load(strCallerKey); ok {

		s.callers.Delete(strCallerKey)
		c := v.(*caller)

		s.mutex.Lock() //lock
		c.LeaveTime = now64
		c.SpendTime = c.LeaveTime - c.EnterTime

		var r *result
		if v2, ok2 := s.results.Load(strResultKey); ok2 {

			r = v2.(*result)
			r.CallCount++
			if c.SpendTime > 0 {
				r.ExeTime += c.SpendTime
				r.AvgTime = r.ExeTime / r.CallCount
			}
			if !c.CallOk {
				r.ErrorCount++
			}
			r.UpdateTime = getUnixSecond()
		}
		s.mutex.Unlock() //unlock
		//Json("result: ", r)
		return c.SpendTime, ok
	}

	return 0, false
}

//统计error次数(incr error counts)
func (s *statistic) error(strFile, strFunc string, nLineNo int) {
	strKey := getCallerStoreKey(strFile, strFunc)
	if v, ok := s.callers.Load(strKey); ok {

		c := v.(*caller)
		c.CallOk = false
	}
}

//统计信息汇总(statistic summary)
func (s *statistic) summary(args...interface{}) string {

	var strFuncName string
	if len(args) == 0 {
		strFuncName = FUNCNAME_NIL
	} else {
		strFuncName = args[0].(string)
	}

	var summ = summary{
		TimeUnit: "micro seconds",
	}

	s.results.Range(
		func(k, v interface{}) bool {

			r := v.(*result)
			if strFuncName == FUNCNAME_ALL || strFuncName == FUNCNAME_NIL || strings.Contains(r.FuncName, strFuncName) {
				summ.Results = append(summ.Results, v.(*result))
			}
			return true
		},
	)

	data, _ := json.MarshalIndent(summ, "", "\t")
	return string(data)
}
