package main

import (
	"github.com/civet148/gotools/comm/timer"
	"github.com/civet148/gotools/log"
	"time"
)

var TimerId_Test = 1

type MyTest struct {
}

//定时任务执行接口
func (this *MyTest) OnTimer(id int, param interface{}) {

	log.Debug("timer id %v param %v", id, param)
}

func main() {

	var test = &MyTest{}

	log.Info("set timer %v", TimerId_Test)
	//测试定时器ID=1，时间间隔1秒，参数为int(10086)
	timer.SetTimer(test, TimerId_Test, 1000, 3, 10086)

	time.Sleep(10 * time.Second)

	timer.KillTimer(test, TimerId_Test)
	log.Info("kill timer %v", TimerId_Test)
	time.Sleep(60 * time.Second)
}
