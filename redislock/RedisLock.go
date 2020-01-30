package redislock

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strings"
	"time"
)

/*
*  基于Redis实现分布式锁
*  1、互斥性。在任意时刻，只有一个客户端能持有锁。
*  2、不会发生死锁。即使有一个客户端在持有锁的期间崩溃而没有主动解锁，也能保证后续其他客户端能加锁
*  3、具有容错性，只要大部分的Redis节点正常运行，客户端就可以加锁和解锁
*  4、加锁和解锁必须是同一个客户端，客户端自己不能把别人加的锁给解了
*/

var REDIS_CMD_AUTH = "AUTH"
var REDIS_CMD_SET = "SET"
var REDIS_CMD_GET = "GET"
var REDIS_CMD_DEL = "DEL"
var REDIS_CMD_EVAL = "EVAL"
var REDIS_CMD_EXPIRE = "EXPIRE"
var REDIS_CMD_PING = "PING"

var REDIS_LOCK_SUCC = "OK"
var REDIS_UNLOCK_SUCC = int64(1)

var REDIS_SET_EX = "EX" //设置过期时间（单位：秒）
var REDIS_SET_NX = "NX" //设置key-value当指定key不存在时
var REDIS_SET_WITH_EXPIRE_TIME = "PX" //设置过期时间（单位：毫秒）


type RedisLock struct {

	pool *redis.Pool //Redis连接池
	host, passwd string //Host 服务器地址:端口 Passwd 认证密码
	stat  bool //连接状态（true已连接，false未连接）
}

type Locker struct {
	pool *redis.Pool
	key string
	value string
}

func NewRedisLock() *RedisLock {

	return &RedisLock{}
}


/*
* @breif 	打开Redis连接，初始化连接池
* @param 	strURL  Redis连接URL 有用户名密码 "redis://123456@192.168.1.10:6379"
*                   没有用户名密码则是 "redis://192.168.1.10:6379"
* @return 	err		成功返回nil，失败返回错误类型
*/
func (rds *RedisLock) Open(strURL string) (err error) {

	var strRedisScheme  = "redis://"
	strConnUrl := strings.TrimSpace(strURL)

	if !strings.Contains(strConnUrl, strRedisScheme) {

		err = fmt.Errorf("Url scheme illegal, prefix must be redis://")
		return
	}
	strConnUrl = strings.Replace(strConnUrl, strRedisScheme, "", -1)
	if strings.Contains(strConnUrl, "@") {

		rds.passwd = strings.Split(strConnUrl, "@")[0]
		rds.host = strRedisScheme+strings.Split(strConnUrl, "@")[1]
	} else {
		rds.host = strRedisScheme+strConnUrl
	}

	rds.pool = rds.newPool()
	if rds.pool == nil {
		err = fmt.Errorf("Create redis pool failed.")
		return
	}

	c := rds.pool.Get()
	defer c.Close()
	_, err = c.Do(REDIS_CMD_PING)
	if err != nil {

		return fmt.Errorf("ping redis error: %v", err.Error())
	}

	rds.stat = true
	return
}


/*
* @breif 	尝试加锁
* @param 	strKey  		Redis唯一key
* @param 	strValue 		Redis值，用于识别哪个用户加锁
* @param 	nExpireTime 	超时时间（秒）
* @return 	locker			加锁成功后，返回用于解锁的对象
            err				成功返回nil，失败返回error
*/
func (rds *RedisLock) TryLock(strKey, strValue string, nExpireTime int) (locker *Locker, err error) {

	if !rds.stat {
		err = fmt.Errorf("Reids connect status error, please call 'Open' function to connect reids server")
		return nil, err
	}
	c := rds.pool.Get()
	defer c.Close()

	var reply string

	reply, err = redis.String(c.Do(REDIS_CMD_SET, strKey, strValue, REDIS_SET_EX, nExpireTime, REDIS_SET_NX))
	if err != nil || reply != REDIS_LOCK_SUCC {

		err = fmt.Errorf("Try lock failed, replay [%v] error [%v]", reply, err)
		return nil, err
	}

	locker = &Locker{
		pool: rds.pool,
		key: strKey,
		value: strValue,
	}
	return
}

/*
* @breif 	解锁
* @return 	err	 	成功返回nil，失败返回error(失败的情况仅限于redis服务器宕机了或执行脚本出错)
*/
func (locker *Locker) Unlock() (err error) {

	c := locker.pool.Get()
	defer c.Close()

	//让Redis服务器执行Lua脚本如果key-value匹配则删除，防止误删别的用户加的锁
	//场景：A用户加锁成功后在解锁前已超时且B用户加锁成功的情况下，如果只以key作为删除的依据，A用户可能会误删B用户的锁
	strScript := "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	sc := redis.NewScript(1, strScript)
	//除了redis执行脚本返回错误，无需判断返回值，一律视为解锁成功
	if _, err = redis.Int64(sc.Do(c, locker.key, locker.value)); err != nil {

		return
	}
	return
}

// NewRedisPool 返回redis连接池
func (rds *RedisLock) newPool() *redis.Pool {
	return &redis.Pool{

		MaxIdle:     1,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(rds.host)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}

			if rds.passwd != "" {
				//验证redis密码
				if _, authErr := c.Do(REDIS_CMD_AUTH, rds.passwd); authErr != nil {

					return nil, fmt.Errorf("redis auth password error: %s", authErr)
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do(REDIS_CMD_PING)
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
}