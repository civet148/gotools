package redis

import (
	"fmt"
	"github.com/civet148/gotools/mq"
	"github.com/garyburd/redigo/redis"
	log "github.com/civet148/gotools/log"
	"strings"
	"time"
)

var REDIS_CMD_AUTH = "AUTH"
var REDIS_CMD_SET = "SET"
var REDIS_CMD_GET = "GET"
var REDIS_CMD_DEL = "DEL"
var REDIS_CMD_EVAL = "EVAL"
var REDIS_CMD_EXPIRE = "EXPIRE"
var REDIS_CMD_PING = "PING"
var REDIS_CMD_LPUSH = "LPUSH"
var REDIS_CMD_BRPOP = "BRPOP"
var REDIS_CMD_PUBLISH = "PUBLISH"
var REDIS_CMD_SUBSCRIBE = "SUBSCRIBE"
var REDIS_CMD_PSUBSCRIBE = "PSUBSCRIBE" //模式订阅(使用*作为通配符)
var REDIS_RESP_OK = "OK"

var REDIS_SET_EX = "EX"
var REDIS_SET_NX = "NX"
var REDIS_SET_WITH_EXPIRE_TIME = "PX"

type RedisMQ struct {
	mode         mq.Mode
	url          string      //服务器连接URL
	pool         *redis.Pool //Redis连接池
	host, passwd string      //host Redis服务器地址:端口 passwd Redis认证密码
	closed bool           //远程服务器是否已断开连接
	debug  bool //开启或关闭调试信息
}

func init() {
	mq.Register(mq.Adapter_RedisMQ, NewMQ)
}

func NewMQ() (mq.IReactMQ) {

	return &RedisMQ{}
}


//Redis连接URL 有密码 "redis://123456@192.168.1.10:6379"
//            没有密码则是 "redis://192.168.1.10:6379"
func (this *RedisMQ) Connect(mode mq.Mode, strURL string) (err error) {
	this.mode = mode
	this.url = strURL

	if err = this.Reconnect(); err != nil {
		log.Error("Connect to [%v] error [%v]", strURL, err.Error())
		return
	}
	return
}

func (this *RedisMQ) Reconnect() (err error) {

	var strRedisScheme = "redis://"
	strConnUrl := strings.TrimSpace(this.url)

	if !strings.Contains(strConnUrl, strRedisScheme) {

		err = fmt.Errorf("Url scheme illegal, prefix must be redis://")
		return
	}
	strConnUrl = strings.Replace(strConnUrl, strRedisScheme, "", -1)
	if strings.Contains(strConnUrl, "@") {

		this.passwd = strings.Split(strConnUrl, "@")[0]
		this.host = strRedisScheme + strings.Split(strConnUrl, "@")[1]
	} else {
		this.host = strRedisScheme + strConnUrl
	}

	if this.closed {//服务器断开连接，尝试PING服务器

		c := this.pool.Get()
		defer c.Close()
		_, err = c.Do(REDIS_CMD_PING)
		if err != nil {

			return fmt.Errorf("ping redis error: %v", err.Error())
		}
		this.closed = false
	} else {

		this.pool = this.newPool()
		if this.pool == nil {
			err = fmt.Errorf("Create redis pool failed.")
			return
		}

		c := this.pool.Get()
		defer c.Close()
		_, err = c.Do(REDIS_CMD_PING)
		if err != nil {

			return fmt.Errorf("ping redis error: %v", err.Error())
		}
	}
	return
}

func (this *RedisMQ) IsClosed() (bool) {

	return this.closed
}

func (this *RedisMQ) Publish(strRoutingKey, strData string) (err error) {

	//构造一个队列名
	strRoutingKey = this.getQueueName(strRoutingKey)
	switch this.mode {
	case mq.Mode_Direct: return this.publishDirect(strRoutingKey, strData)
	case mq.Mode_Topic: return this.publishTopic(strRoutingKey, strData)
	case mq.Mode_Fanout: return this.publishFanout(strRoutingKey, strData)
	default:
		return fmt.Errorf("Unknown mode [%v]", this.mode)
	}
	return
}

func (this *RedisMQ) Consume(strBindingKey, strQueueName string, cb mq.FnConsumeCb) (err error) {
	//构造一个队列名
	strBindingKey = this.getQueueName(strBindingKey)

	switch this.mode {
	case mq.Mode_Direct: return this.consumeDirect(strBindingKey, strQueueName, cb)
	case mq.Mode_Topic: return this.consumeTopic(strBindingKey, strQueueName, cb)
	case mq.Mode_Fanout: return this.consumeFanout(strBindingKey, strQueueName, cb)
	default:
		return fmt.Errorf("Unknown mode [%v]", this.mode)
	}
	return
}

/*
* @brief 	开启或关闭调式模式
* @param 	enable 	true开启/false关闭
*/
func (this *RedisMQ) Debug(enable bool) {

	this.debug = enable
}

func (this *RedisMQ) publishDirect(strRoutingKey, strData string) (err error) {

	if this.debug {
		log.Info("[PRODUCER] mode [%v] key [%v] data [%v]", this.mode.String(), strRoutingKey, strData)
	}

	c := this.pool.Get()
	defer c.Close()
	if _, err = redis.Int(c.Do(REDIS_CMD_LPUSH, strRoutingKey, strData)); err != nil {

		if this.debug {
			log.Error("Exec command [%v] error [%v]", REDIS_CMD_LPUSH, err.Error())
		}

		return
	}

	return
}

func (this *RedisMQ) consumeDirect(strBindingKey, strQueueName string, cb mq.FnConsumeCb) (err error) {

	var strData string

RETRY_CONSUME:
	c := this.pool.Get()
	for {

		var datas [][]byte
		if datas, err = redis.ByteSlices(c.Do(REDIS_CMD_BRPOP, strBindingKey, 0)); err != nil {

			this.closed = this.checkDisconnectedByError(err)

			if this.debug {
				//log.Error("Exec command [%v] error [%v]", REDIS_CMD_BRPOP, err.Error())
			}

			if this.closed {
				c.Close()
				time.Sleep(1*time.Second)
				goto RETRY_CONSUME
			}
		}

		if len(datas) > 0 {
			strData = string(datas[len(datas)-1])
			cb(strData)
		}

		if this.debug {
			log.Info("[CONSUMER] mode [%v] key [%v] data [%v]", this.mode.String(), strBindingKey, strData)
		}
	}

	return
}

func (this *RedisMQ) publishTopic(strRoutingKey, strData string) (err error) {
	if this.debug {
		log.Info("[PRODUCER] mode [%v] key [%v] data [%v]", this.mode.String(), strRoutingKey, strData)
	}
	c := this.pool.Get()
	defer c.Close()
	if _, err = redis.Int(c.Do(REDIS_CMD_PUBLISH, strRoutingKey, strData)); err != nil {
		if this.debug {
			log.Error("Exec command [%v] error [%v]", REDIS_CMD_PUBLISH, err.Error())
		}
		return
	}

	return
}

func (this *RedisMQ) consumeTopic(strBindingKey, strQueueName string, cb mq.FnConsumeCb) (err error) {

	strBindingKey = strings.Replace(strBindingKey, "#", "*", -1) //将#符号替换成*(redis不支持#符号匹配模式)
	var strData string

RETRY_CONSUME:
	c := this.pool.Get()
	for {

		var datas []string
		if datas, err = redis.Strings(c.Do(REDIS_CMD_PSUBSCRIBE, strBindingKey)); err != nil {

			this.closed = this.checkDisconnectedByError(err)
			if this.debug {
				//log.Error("Exec command [%v] error [%v]", REDIS_CMD_PSUBSCRIBE, err.Error())
			}
			if this.closed {
				c.Close()
				time.Sleep(3*time.Second)
				goto RETRY_CONSUME
			}
			time.Sleep(50*time.Millisecond)
		} else {
			if this.debug {
				log.Debug("reply %+v", datas)
			}
			if len(datas) > 0 {
				strData  = string(datas[len(datas)-1])
				if this.debug {
					log.Info("[CONSUMER] mode [%v] key [%v] data [%v]", this.mode.String(), strBindingKey, strData)
				}
				cb(strData)
			}
		}
	}

	return
}

func (this *RedisMQ) publishFanout(strRoutingKey, strData string) (err error) {
	if this.debug {
		log.Info("[PRODUCER] mode [%v] key [%v] data [%v]", this.mode.String(), strRoutingKey, strData)
	}
	c := this.pool.Get()
	defer c.Close()
	if _, err = redis.Int(c.Do(REDIS_CMD_PUBLISH, strRoutingKey, strData)); err != nil {
		if this.debug {
			log.Error("Exec command [%v] error [%v]", REDIS_CMD_PUBLISH, err.Error())
		}
		return
	}

	return
}

func (this *RedisMQ) consumeFanout(strBindingKey, strQueueName string, cb mq.FnConsumeCb) (err error) {

	var strData string
RETRY_CONSUME:
	c := this.pool.Get()
	for {

		var datas []string
		if datas, err = redis.Strings(c.Do(REDIS_CMD_SUBSCRIBE, strBindingKey)); err != nil {

			this.closed = this.checkDisconnectedByError(err)
			if this.debug {
				//log.Error("Exec command [%v] error [%v]", REDIS_CMD_SUBSCRIBE, err.Error())
			}
			if this.closed {
				c.Close()
				time.Sleep(3*time.Second)
				goto RETRY_CONSUME
			}
			time.Sleep(50*time.Millisecond)
		} else {

			if this.debug {
				log.Debug("reply %+v", datas)
			}
			if len(datas) > 0 {
				strData  = string(datas[len(datas)-1])
				if this.debug {
					log.Info("[CONSUMER] mode [%v] key [%v] data [%v]", this.mode.String(), strBindingKey, strData)
				}
				cb(strData)
			}
		}
	}

	return
}

// NewRedisPool 返回redis连接池
func (this *RedisMQ) newPool() *redis.Pool {
	return &redis.Pool{

		MaxIdle:     1,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(this.host)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}

			if this.passwd != "" {
				//验证redis密码
				if _, authErr := c.Do(REDIS_CMD_AUTH, this.passwd); authErr != nil {

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

func (this *RedisMQ) getQueueName(strKeyName string) (name string) {

	switch this.mode {
	case mq.Mode_Direct: name = fmt.Sprintf("%v.%v", mq.EXCHANGE_NAME_DIRECT, strKeyName)
	case mq.Mode_Topic:  name = fmt.Sprintf("%v.%v", mq.EXCHANGE_NAME_TOPIC, strKeyName)
	case mq.Mode_Fanout: name = fmt.Sprintf("%v.%v", mq.EXCHANGE_NAME_FANOUT, strKeyName)
	}
	return
}

func (this *RedisMQ) checkDisconnectedByError(err error) (bool) {

	if strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "connection"){//标记服务器断开连接
		if this.debug {
			log.Error("MQ server is disconnected...")
		}
		return true
	}

	return false
}