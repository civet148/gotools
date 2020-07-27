package redis

import (
	"fmt"
	log "github.com/civet148/gotools/log"
	"github.com/civet148/gotools/mq"
	"github.com/garyburd/redigo/redis"
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
	closed       bool        //远程服务器是否已断开连接
	debug        bool        //开启或关闭调试信息
	closing      chan bool   //主动关闭通知通道
}

func init() {
	mq.Register(mq.Adapter_RedisMQ, NewMQ)
}

func NewMQ() mq.IReactMQ {

	return &RedisMQ{
		closing: make(chan bool, 1),
	}
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

		err = fmt.Errorf("url scheme illegal, prefix must be redis://")
		return
	}
	strConnUrl = strings.Replace(strConnUrl, strRedisScheme, "", -1)
	if strings.Contains(strConnUrl, "@") {

		this.passwd = strings.Split(strConnUrl, "@")[0]
		this.host = strRedisScheme + strings.Split(strConnUrl, "@")[1]
	} else {
		this.host = strRedisScheme + strConnUrl
	}

	if this.closed { //服务器断开连接，尝试PING服务器

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
			err = fmt.Errorf("create redis pool failed")
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

func (this *RedisMQ) Close() {
	this.closing <- true
}

func (this *RedisMQ) IsClosed() bool {

	return this.closed
}

/*
* @brief 	消息发布接口定义(仅支持字符串类型消息)
* @param 	strBindingKey 	队列绑定Key(topic)
* @param	strQueueName 	队列名称(group)
* @param	key 		消息KEY(仅kafka必填，其他MQ类型默认填PRODUCER_KEY_NULL)
* @param	value 		消息数据
* @return   err 发布失败返回具体错误信息
 */
func (this *RedisMQ) Publish(strBindingKey, strQueueName, key string, value string) (err error) {

	//构造一个队列名
	strBindingKey = this.getQueueName(strBindingKey)
	switch this.mode {
	case mq.Mode_Direct:
		return this.publishDirect(strBindingKey, key, value)
	case mq.Mode_Topic:
		return this.publishTopic(strBindingKey, key, value)
	case mq.Mode_Fanout:
		return this.publishFanout(strBindingKey, key, value)
	default:
		return fmt.Errorf("Unknown mode [%v]", this.mode)
	}
	return
}

/*
* @brief 	消息消费接口定义
* @param 	strBindingKey 	队列绑定Key
* @param	strQueueName 	队列名称
* @param    handler         消费回调处理对象
* @return   err 成功返回nil，失败返回返回具体错误信息
* @remark   服务器异常或重启时内部会自动重连服务器
 */
func (this *RedisMQ) Consume(strBindingKey, strQueueName string, handler mq.ReactHandler) (err error) {
	//构造一个队列名
	strBindingKey = this.getQueueName(strBindingKey)

	switch this.mode {
	case mq.Mode_Direct:
		return this.consumeDirect(strBindingKey, strQueueName, handler)
	case mq.Mode_Topic:
		return this.consumeTopic(strBindingKey, strQueueName, handler)
	case mq.Mode_Fanout:
		return this.consumeFanout(strBindingKey, strQueueName, handler)
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

func (this *RedisMQ) publishDirect(strBindingKey, key string, value string) (err error) {

	strData := string(value)

	if this.debug {
		log.Info("[PRODUCER] mode [%v] key [%v] data [%v]", this.mode.String(), strBindingKey, strData)
	}

	if this.IsClosed() {
		err = fmt.Errorf("MQ is closed")
		log.Error(err.Error())
		return
	}
	c := this.pool.Get()
	defer c.Close()
	if _, err = redis.Int(c.Do(REDIS_CMD_LPUSH, strBindingKey, strData)); err != nil {

		if this.debug {
			log.Error("exec command [%v] error [%v]", REDIS_CMD_LPUSH, err.Error())
		}

		return
	}

	return
}

func (this *RedisMQ) consumeDirect(strBindingKey, strQueueName string, handler mq.ReactHandler) (err error) {

	var strData string

RETRY_CONSUME:
	c := this.pool.Get()
	for {
		select {
		case <-this.closing:
			{
				log.Debugf("MQ is closing...")
				return
			}
		default:
			{
				var datas [][]byte
				if datas, err = redis.ByteSlices(c.Do(REDIS_CMD_BRPOP, strBindingKey, 0)); err != nil {

					this.closed = this.checkDisconnectedByError(err)
					if this.closed {
						_ = c.Close()
						time.Sleep(1 * time.Second)
						goto RETRY_CONSUME
					}
				}

				if len(datas) > 0 {
					strData = string(datas[len(datas)-1])
					handler.OnConsume(mq.Adapter_RedisMQ, strBindingKey, strQueueName, mq.DEFAULT_DATA_KEY, strData)
				}

				if this.debug {
					log.Info("[CONSUMER] mode [%v] key [%v] data [%v]", this.mode.String(), strBindingKey, strData)
				}
			}
		}

	}

	return
}

func (this *RedisMQ) publishTopic(strBindingKey, key string, value string) (err error) {

	strData := string(value)

	if this.debug {
		log.Info("[PRODUCER] mode [%v] key [%v] data [%v]", this.mode.String(), strBindingKey, strData)
	}
	c := this.pool.Get()
	defer c.Close()
	if _, err = redis.Int(c.Do(REDIS_CMD_PUBLISH, strBindingKey, strData)); err != nil {
		if this.debug {
			log.Error("Exec command [%v] error [%v]", REDIS_CMD_PUBLISH, err.Error())
		}
		return
	}

	return
}

func (this *RedisMQ) consumeTopic(strBindingKey, strQueueName string, handler mq.ReactHandler) (err error) {

	strBindingKey = strings.Replace(strBindingKey, "#", "*", -1) //将#符号替换成*(redis不支持#符号匹配模式)
	var strData string

RETRY_CONSUME:
	c := this.pool.Get()
	for {
		select {
		case <-this.closing:
			{
				log.Debugf("MQ is closing...")
				return
			}
		default:
			{
				var datas []string
				if datas, err = redis.Strings(c.Do(REDIS_CMD_PSUBSCRIBE, strBindingKey)); err != nil {

					this.closed = this.checkDisconnectedByError(err)
					if this.closed {
						c.Close()
						time.Sleep(3 * time.Second)
						goto RETRY_CONSUME
					}
					time.Sleep(50 * time.Millisecond)
				} else {
					if this.debug {
						log.Debug("reply %+v", datas)
					}
					if len(datas) > 0 {
						strData = string(datas[len(datas)-1])
						if this.debug {
							log.Info("[CONSUMER] mode [%v] key [%v] data [%v]", this.mode.String(), strBindingKey, strData)
						}
						handler.OnConsume(mq.Adapter_RedisMQ, strBindingKey, strQueueName, mq.DEFAULT_DATA_KEY, strData)
					}
				}
			}
		}
	}

	return
}

func (this *RedisMQ) publishFanout(strBindingKey, key string, value string) (err error) {

	if this.debug {
		log.Info("[PRODUCER] mode [%v] key [%v] data [%v]", this.mode.String(), strBindingKey, value)
	}
	c := this.pool.Get()
	defer c.Close()
	if _, err = redis.Int(c.Do(REDIS_CMD_PUBLISH, strBindingKey, value)); err != nil {
		if this.debug {
			log.Error("Exec command [%v] error [%v]", REDIS_CMD_PUBLISH, err.Error())
		}
		return
	}

	return
}

func (this *RedisMQ) consumeFanout(strBindingKey, strQueueName string, handler mq.ReactHandler) (err error) {

	var strData string
RETRY_CONSUME:
	c := this.pool.Get()
	for {
		select {
		case <-this.closing:
			{
				log.Debugf("MQ is closing...")
				return
			}
		default:
			{
				var datas []string
				if datas, err = redis.Strings(c.Do(REDIS_CMD_SUBSCRIBE, strBindingKey)); err != nil {

					this.closed = this.checkDisconnectedByError(err)
					if this.closed {
						c.Close()
						time.Sleep(3 * time.Second)
						goto RETRY_CONSUME
					}
					time.Sleep(50 * time.Millisecond)
				} else {

					if this.debug {
						log.Debug("reply %+v", datas)
					}
					if len(datas) > 0 {
						strData = string(datas[len(datas)-1])
						if this.debug {
							log.Info("[CONSUMER] mode [%v] key [%v] data [%v]", this.mode.String(), strBindingKey, strData)
						}
						handler.OnConsume(mq.Adapter_RedisMQ, strBindingKey, strQueueName, mq.DEFAULT_DATA_KEY, strData)
					}
				}
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
	case mq.Mode_Direct:
		name = fmt.Sprintf("%v.%v", mq.EXCHANGE_NAME_DIRECT, strKeyName)
	case mq.Mode_Topic:
		name = fmt.Sprintf("%v.%v", mq.EXCHANGE_NAME_TOPIC, strKeyName)
	case mq.Mode_Fanout:
		name = fmt.Sprintf("%v.%v", mq.EXCHANGE_NAME_FANOUT, strKeyName)
	}
	return
}

func (this *RedisMQ) checkDisconnectedByError(err error) bool {

	if strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "connection") { //标记服务器断开连接
		if this.debug {
			log.Error("MQ server is disconnected...")
		}
		return true
	}

	return false
}

/*
* @brief 	获取当前MQ类型
* @param 	adapter  MQ类型
 */
func (this *RedisMQ) GetAdapter() (adapter mq.Adapter) {

	return mq.Adapter_RedisMQ
}
