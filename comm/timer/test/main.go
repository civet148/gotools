package main

import (
	"game/Common/log"
	"game/Common/timer"
	"time"
)

var TimerId_Test timer.TimerId = 1

type MyTest struct {

}

//定时任务执行接口
func (this *MyTest) OnTimer(id timer.TimerId, param interface{}) {

	log.Debug("timer id %v param %v", id, param)
}

func main() {

	var test = &MyTest{}

	log.Info("set timer %v", TimerId_Test)
	timer.SetTimer(test, TimerId_Test, 1000, timer.RepeatForever, 10086)

	time.Sleep(10*time.Second)

	timer.KillTimer(test, TimerId_Test)
	log.Info("kill timer %v", TimerId_Test)
	time.Sleep(60*time.Second)
}
