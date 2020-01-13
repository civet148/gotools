package main

import (
	"github.com/civet148/gotools/redislock"
	log "github.com/civet148/gotools/log"
	"time"
)

func main() {

	rlock := redislock.NewRedisLock()
	if err := rlock.Open("redis://192.168.1.15:6379"); err != nil {

		log.Error("%v", err.Error())
		return
	}

	go TestLockA(rlock, 5, 1)
	time.Sleep(10*time.Second)
	go TestLockB(rlock, 15, 30)

	time.Sleep(1*time.Minute)
	log.Info("Program is ending...")
}


func TestLockA(rlock *redislock.RedisLock, expire, sleep int) {

	var REIDSLOCK_KEY = "REDISLOCK"
	var REDISLOCK_VALUE = "I'am user A"

	locker , errLock := rlock.TryLock(REIDSLOCK_KEY, REDISLOCK_VALUE, expire)
	if errLock != nil {
		log.Error("Lock key=[%v] value=[%v]  failed error (%v)", REIDSLOCK_KEY, REDISLOCK_VALUE, errLock.Error())
		return
	}
	log.Info("Lock key=[%v] value=[%v] ok", REIDSLOCK_KEY, REDISLOCK_VALUE)

	time.Sleep( time.Duration(sleep) *time.Second)

	errLock = locker.Unlock()
	if errLock != nil {
		log.Error("Unlock key=[%v] value=[%v]  failed error (%v)", REIDSLOCK_KEY, REDISLOCK_VALUE, errLock.Error())
		return
	}

	log.Info("Unlock key=[%v] value=[%v] ok", REIDSLOCK_KEY, REDISLOCK_VALUE)
}

func TestLockB(rlock *redislock.RedisLock, expire, sleep int) {

	var REIDSLOCK_KEY = "REDISLOCK"
	var REDISLOCK_VALUE = "I'am user B"

	locker , errLock := rlock.TryLock(REIDSLOCK_KEY, REDISLOCK_VALUE, expire)
	if errLock != nil {
		log.Error("Lock key=[%v] value=[%v]  failed error (%v)", REIDSLOCK_KEY, REDISLOCK_VALUE, errLock.Error())
		return
	}
	log.Info("Lock key=[%v] value=[%v] ok", REIDSLOCK_KEY, REDISLOCK_VALUE)

	time.Sleep( time.Duration(sleep) *time.Second)

	errLock = locker.Unlock()
	if errLock != nil {
		log.Error("Unlock key=[%v] value=[%v]  failed error (%v)", REIDSLOCK_KEY, REDISLOCK_VALUE, errLock.Error())
		return
	}

	log.Info("Unlock key=[%v] value=[%v] ok", REIDSLOCK_KEY, REDISLOCK_VALUE)
}