package timer

import (
	"fmt"
	"sync"
	"time"
)

var (
	RepeatForever = -1 //重复执行
	RepeatDone    = 0  //执行完毕
	RepeatOnce    = 1  //执行一次
)

const (
	TIKER_INTERVAL_MS = 500
)

type ITimer interface {
	//定时任务执行接口
	OnTimer(id int, param interface{})
}

//定时任务内部对象
type timerInner struct {
	this    interface{} //需要执行定时任务的对象指针
	id      int         //定时任务ID
	elapse  int         //定时任务执行时间间隔(单位: 毫秒)
	param   interface{} //参数
	repeat  int         //重复次数
	closing chan bool   //关闭消息
}

var gMapTimer = sync.Map{}

//初始化
func init() {

}

func getTimerKey(this interface{}, id int) string {
	return fmt.Sprintf("%p->%v", this, id)
}

//设置定时任务
//this 		实现OnTimer接口的对象指针
//id   		定时任务ID
//elapse 	执行间隔时间(最小单位：毫秒)
//repeat 	重复次数(-1表示重复执行，大于0则表示执行具体次数)
//param 	定时任务附带参数(尽量不要传递对象指针)
func SetTimer(this interface{}, id int, elapse int, repeat int, param interface{}) bool {

	if repeat <= 0 && repeat != RepeatForever { //仅-1允许

		return false
	}

	if elapse < TIKER_INTERVAL_MS {
		elapse = TIKER_INTERVAL_MS
	}

	strKey := getTimerKey(this, id)
	inner := &timerInner{
		this:    this,
		id:      id,
		elapse:  elapse,
		repeat:  repeat,
		param:   param,
		closing: make(chan bool, 1),
	}
	gMapTimer.Store(strKey, inner)
	go startTimer(inner)
	return true
}

//停止定时任务
func KillTimer(this interface{}, id int) {
	var strKey = getTimerKey(this, id)
	if v, ok := gMapTimer.Load(strKey); ok {
		inner := v.(*timerInner)
		inner.closing <- true
		gMapTimer.Delete(strKey)
	}
}

func startTimer(inner *timerInner) {

	ticker := time.NewTicker(time.Duration(inner.elapse) * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-inner.closing:
			return
		case <-ticker.C:
			{
				cb := inner.this.(ITimer)
				cb.OnTimer(inner.id, inner.param) //调用OnTimer方法执行定时任务
				if inner.repeat != RepeatForever && inner.repeat != RepeatDone {
					inner.repeat-- //计数器减一
				}
				if inner.repeat == 0 {
					KillTimer(inner.this, inner.id)
				}
			}
		}
	}
}
