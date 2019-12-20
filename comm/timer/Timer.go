package timer

import (
	"fmt"
	"sync"
	"time"
)

var (
	RepeatForever = -1 //重复执行
	RepeatDone    = 0 //执行完毕
	RepeatOnce    = 1  //执行一次
)

type ITimer interface {
	//定时任务执行接口
	OnTimer(id int, param interface{})
}

//定时任务内部对象
type timerInner struct {
	this   interface{} //需要执行定时任务的对象指针
	id     int     //定时任务ID
	elapse int         //定时任务执行时间间隔(单位: 毫秒)
	param  interface{} //参数
	repeat int         //重复次数
	msec   int64       //下次执行时间(毫秒)
}

var gMapTimer = sync.Map{}

//初始化
func init() {

	go startMonitor() //启动定时任务监视协程
}

func startMonitor() {

	ticker := time.NewTicker(1000*time.Microsecond) //1毫秒秒触发一次

	for {
		select {
		case <-ticker.C: //遍历一次需要执行的定时任务
			gMapTimer.Range(
				func(key, value interface{}) bool {

					nCurMs := getCurTime()

					inner := value.(*timerInner)

					if nCurMs >= inner.msec {

						//log.Debug("key [%v] value [%+v] cur ms[%v] duration [%v]", key, value, nCurMs, nCurMs-inner.msec)
						cb := inner.this.(ITimer)
						cb.OnTimer(inner.id, inner.param) //调用OnTimer方法执行定时任务
						if inner.repeat != RepeatForever && inner.repeat != RepeatDone {
							inner.repeat-- //计数器减一
						}

						if inner.repeat == RepeatDone {//计数为0，删除定时任务对象
							gMapTimer.Delete(key)
						} else {
							inner.msec = getNextTime(inner.elapse) //设置下一次执行时间
						}
					}

					return true
				},
			)
		}
	}

}

//获取当前时间(毫秒)
func getCurTime() (msec int64) {

	now := time.Now()
	return now.UnixNano()/1e6
}

//获取下次执行时间(毫秒)
func getNextTime(elapse int) (msec int64) {

	now := time.Now()
	msec = now.UnixNano()/1e6 + int64(elapse)
	return
}

func getTimerKey(this interface{}, id int) string {
	return  fmt.Sprintf("%p->%v", this, id)
}

//设置定时任务
//this 		实现OnTimer接口的对象指针
//id   		定时任务ID
//elapse 	执行间隔时间(最小单位：毫秒)
//repeat 	重复次数(-1表示重复执行，大于0则表示执行具体次数)
//param 	定时任务附带参数(尽量不要传递对象指针)
func SetTimer(this interface{}, id int, elapse int, repeat int, param interface{}) (bool) {

	if repeat <= 0 && repeat != RepeatForever {//仅-1允许

		return false
	}

	strKey := getTimerKey(this, id)
	inner := timerInner{
		this:   this,
		id:     id,
		elapse: elapse,
		repeat: repeat,
		param:  param,
		msec: getNextTime(elapse), //首次执行时间（毫秒）
	}

	gMapTimer.Store(strKey, &inner)
	return  true
}

//停止定时任务
func KillTimer(this interface{}, id int) {
	gMapTimer.Delete(getTimerKey(this, id))
}
