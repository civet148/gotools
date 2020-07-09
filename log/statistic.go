package log

import (
	"encoding/json"
	"fmt"
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
	mutex   sync.RWMutex
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
	FileName     string `json:"file_name"`      //code file of function [代码文件名]
	LineNo       int    `json:"line_no"`        //line no of function [行号]
	FuncName     string `json:"func_name"`      //function name [方法名称]
	CallCount    int64  `json:"call_count"`     //call total times [调用次数]
	ErrorCount   int64  `json:"error_count"`    //call error times [错误次数]
	TotalTime    int64  `json:"total_time"`     //micro seconds [执行总时间]
	TotalTimeStr string `json:"total_time_str"` //time string format [执行总时间-日期字符串格式]
	AvgTime      int64  `json:"avg_time"`       //micro seconds [平均执行时间]
	AvgTimeStr   string `json:"avg_time_str"`   //time string format [平均执行时间-日期字符串格式]
	MaxTime      int64  `json:"max_time"`       //max time elapse once [单次最大执行时间]
	MaxTimeStr   string `json:"max_time_str"`   //max time elapse once [单次最大执行时间-日期字符串格式]
	CreateTime   int64  `json:"create_time"`    //unix timestamp on seconds [计时开始时间]
	UpdateTime   int64  `json:"update_time"`    //unix timestamp on seconds [计时更新时间]
}

type summary struct {
	TimeUnit string    `json:"time_unit"`
	Results  []*result `json:"statistics"`
}

var stic *statistic //数据统计对象

func init() {

	stic = newStatistic()
	go checkExpire(stic)
}

//create a new statistic object
func newStatistic() *statistic {
	return &statistic{}
}

func getUnixSecond() int64 {
	return time.Now().Unix()
}

func getDatetime() string {

	return time.Now().Format("2006-01-02 15:04:05")
}

func getMilliSec() int64 {

	return time.Now().UnixNano() / 1e6 //milliseconds
}

func getMicroSec() int64 {

	return time.Now().UnixNano() / 1e3 //microseconds
}

func getSpendTime(microseconds int64) (h, m, s int, ms float32) {

	if microseconds > 0 {
		nSpend := microseconds / 1e6

		if nSpend > 0 {
			h = int(nSpend / 3600)
			m = int((nSpend % 3600) / 60)
			s = int((nSpend % 3600) % 60)
			rem := microseconds - (nSpend * 1e6)
			if rem > 0 {
				ms = float32(rem) / 1000
			}
		} else {
			ms = float32(microseconds) / 1000
		}
	}

	return
}

func getCallerStoreKey(strFile, strFunc string) string {
	return fmt.Sprintf("%v %v %v", getRoutineId(), strFile, strFunc)
}

func makeCallerStoreKey(strRoutineId string, strFile, strFunc string) string {
	return fmt.Sprintf("%v %v %v", strRoutineId, strFile, strFunc)
}

func getResultStoreKey(strFile, strFunc string) string {

	return fmt.Sprintf("%v:%v", strFile, strFunc)
}

func getRoutineId() (strRoutine string) {

	strStack := string(debug.Stack())
	nIdx := strings.IndexAny(strStack, ":\r\n")
	if nIdx > 0 {
		strRoutine = strStack[:nIdx]
		strings.TrimSpace(strRoutine)
		nIdx = strings.Index(strRoutine, "[")
		if nIdx > 0 {
			strRoutine = strRoutine[:nIdx-1]
		}
		return
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
			TotalTime:  0,
			AvgTime:    0,
			CreateTime: getUnixSecond(),
			UpdateTime: getUnixSecond(),
		}
		s.results.Store(strResultKey, r)
	}
	return
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
				r.TotalTime += c.SpendTime
				r.AvgTime = r.TotalTime / r.CallCount
				if c.SpendTime > r.MaxTime {
					r.MaxTime = c.SpendTime //单次调用最大耗时
				}
			}
			if !c.CallOk {
				r.ErrorCount++
			}
			h, m, s, ms := getSpendTime(r.TotalTime)
			r.TotalTimeStr = fmt.Sprintf("%vh %vm %vs %.3fms", h, m, s, ms)
			h, m, s, ms = getSpendTime(r.AvgTime)
			r.AvgTimeStr = fmt.Sprintf("%vh %vm %vs %.3fms", h, m, s, ms)
			h, m, s, ms = getSpendTime(r.MaxTime)
			r.MaxTimeStr = fmt.Sprintf("%vh %vm %vs %.3fms", h, m, s, ms)
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
func (s *statistic) report(args ...interface{}) string {

	var strFuncName string
	if len(args) == 0 {
		strFuncName = FUNCNAME_NIL
	} else {
		strFuncName = args[0].(string)
	}

	var summ = summary{
		TimeUnit: "micro seconds",
	}
	s.mutex.RLock()
	s.results.Range(
		func(k, v interface{}) bool {

			r := v.(*result)
			if strFuncName == FUNCNAME_ALL || strFuncName == FUNCNAME_NIL || strings.Contains(r.FuncName, strFuncName) {
				summ.Results = append(summ.Results, r)
			}
			return true
		},
	)
	data, _ := json.MarshalIndent(summ, "", "\t")
	s.mutex.RUnlock()
	return string(data)
}

func checkExpire(stic *statistic) {

	for {
		stic.callers.Range(
			func(k, v interface{}) bool {

				now64 := getMicroSec()

				c := v.(*caller)
				if now64 > c.ExpireTime {
					Warn("caller key [%v] expired at [%v]", k.(string), getDatetime())
					stic.callers.Delete(k)
				}
				return true
			},
		)
		time.Sleep(time.Hour)
	}
}
