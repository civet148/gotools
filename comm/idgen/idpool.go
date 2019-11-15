package idgen

import (
	"errors"
	"fmt"
	gsf "github.com/zheng-ji/goSnowFlake"
)

var g_iw *gsf.IdWorker = nil

func init() {
	// Params: Given the workerId, 0 < workerId < 1024
	iw, err := gsf.NewIdWorker(1)
	if err != nil {
		fmt.Println(fmt.Sprintf("snowflake new id worker err: %s", err.Error()))
		return
	}
	g_iw = iw
}

// 账号系统 workerid = 1
// 订单系统 workerid = 2
// 运单系统 workerid = 3
// 车辆id   workerid = 4
// 结算系统  workerid = 6
func Init(workerId int64) {
	// Params: Given the workerId, 0 < workerId < 1024
	iw, err := gsf.NewIdWorker(workerId)
	if err != nil {
		fmt.Println(fmt.Sprintf("snowflake new id worker err: %s", err.Error()))
		return
	}
	g_iw = iw
}

func CreateId() (int64, error) {
	if g_iw == nil {
		return 0, errors.New("snowflake id worker init fail.")
	}

	if id, err := g_iw.NextId(); err != nil {
		return 0, err
	} else {
		return id, nil
	}
}
