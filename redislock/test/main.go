package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/redislock"
	"sync"
	"time"
)

func main() {

	rlock := redislock.NewRedisLock()
	if err := rlock.Open("redis://127.0.0.1:6379"); err != nil {

		log.Error("%v", err.Error())
		return
	}

	var wg = &sync.WaitGroup{}
	wg.Add(2)
	go TestLockA(rlock, 5, 5, wg) //用户A锁5秒，10秒后解锁（可能成功也可能失败）
	go TestLockB(rlock, 5, 5, wg) //用户B锁10秒，5秒后解锁（可能成功也可能失败）

	wg.Wait()
	log.Info("Program is ending...")
}

func TestLockA(rlock *redislock.RedisLock, expireSec, unlockSec int, wg *sync.WaitGroup) {

	var REIDSLOCK_KEY = "REDISLOCK"
	var REDISLOCK_VALUE = "I'am user A"

	defer wg.Done()

	locker, ok := rlock.TryLock(REIDSLOCK_KEY, REDISLOCK_VALUE, expireSec)
	if !ok {
		log.Error("Lock key=[%v] value=[%v] failed", REIDSLOCK_KEY, REDISLOCK_VALUE)
		return
	}
	log.Info("Lock key=[%v] value=[%v] ok", REIDSLOCK_KEY, REDISLOCK_VALUE)

	time.Sleep(time.Duration(unlockSec) * time.Second)

	_ = locker.Unlock()

	log.Info("Unlock key=[%v] value=[%v] ok", REIDSLOCK_KEY, REDISLOCK_VALUE)
}

func TestLockB(rlock *redislock.RedisLock, expireSec, unlockSec int, wg *sync.WaitGroup) {

	var REIDSLOCK_KEY = "REDISLOCK"
	var REDISLOCK_VALUE = "I'am user B"

	defer wg.Done()

	locker, ok := rlock.TryLock(REIDSLOCK_KEY, REDISLOCK_VALUE, expireSec)
	if !ok {
		log.Error("Lock key=[%v] value=[%v] failed", REIDSLOCK_KEY, REDISLOCK_VALUE)
		return
	}
	log.Info("Lock key=[%v] value=[%v] ok", REIDSLOCK_KEY, REDISLOCK_VALUE)

	time.Sleep(time.Duration(unlockSec) * time.Second)

	_ = locker.Unlock()
	log.Info("Unlock key=[%v] value=[%v] ok", REIDSLOCK_KEY, REDISLOCK_VALUE)
}
