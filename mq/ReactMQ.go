package mq

import (
	"fmt"
	"github.com/civet148/gotools/log"
)

const (
	DEFAULT_DATA_KEY = "MQ-DATA-KEY"
)

type Mode int8

const (
	Mode_Direct Mode = 1 //根据routing-key完全匹配消息队列
	Mode_Topic  Mode = 2 //根据topic规则匹配消息队列
	Mode_Fanout Mode = 3 //广播模式
)

type Adapter int8

const (
	Adapter_RabbitMQ Adapter = 1
	Adapter_RedisMQ  Adapter = 2
	Adapter_KafkaMQ  Adapter = 3
	Adapter_RocketMQ Adapter = 4
	Adapter_MQTT     Adapter = 5
	Adapter_ETCD     Adapter = 6
)

var UNSET = "<UNSET>"

//var MAX_POOL_SIZE = 512 //连接池最大连接数
//var DEFAULT_POOL_SIZE = int(5) //连接池默认连接数
var EXCHANGE_NAME_DIRECT = "REACT_MQ_DIRECT_EXCHANGE" //Driect模式默认交换器名称定义
var EXCHANGE_NAME_TOPIC = "REACT_MQ_TOPIC_EXCHANGE"   //Topic模式默认交换器名称定义
var EXCHANGE_NAME_FANOUT = "REACT_MQ_FANOUT_EXCHANGE" //Fanout模式默认交换器名称定义

func (m Mode) String() string {
	switch m {
	case Mode_Direct:
		return "direct"
	case Mode_Topic:
		return "topic"
	case Mode_Fanout:
		return "fanout"
	}
	return UNSET
}

func (a Adapter) String() string {
	switch a {
	case Adapter_RabbitMQ:
		return "RabbitMQ"
	case Adapter_RedisMQ:
		return "RedisMQ"
	case Adapter_KafkaMQ:
		return "KafkaMQ"
	case Adapter_RocketMQ:
		return "RocketMQ"
	case Adapter_MQTT:
		return "MQTT"
	case Adapter_ETCD:
		return "ETCD"
	}
	return UNSET
}

type ReactHandler interface {
	OnConsume(adapter Adapter, strBindingKey, strQueueName, strKey string, strValue string)
}

type IReactMQ interface {
	/*
	* @brief 	MQ服务器连接接口定义
	* @param 	strUrl 连接服务器URL(格式规范 [amqp|redis|rocket|kafka|mqtt]://user:password@host:port)
	* @return 	err 连接失败返回具体错误信息
	 */
	Connect(mode Mode, strURL string) (err error)

	/*
	* @brief 	MQ服务器重新连接接口定义
	* @param
	* @return 	err 连接失败返回具体错误信息
	* @remark   当Publish返回错误且IsClosed()方法亦返回true时调用此方法重连MQ服务器
	 */
	Reconnect() (err error)

	/*
	* @brief 	主动关闭
	 */
	Close()

	/*
	* @brief 	判定是否MQ服务器断开连接（异常宕机或重启）
	* @param
	* @return 	远程服务器连接断开返回true，否则返回false
	 */
	IsClosed() bool

	/*
	* @brief 	消息发布接口定义(仅支持字符串类型消息)
	* @param 	strBindingKey 	路由Key
	* @param	strQueueName 	队列名称(redis/mqtt协议非必填)
	* @param	key 		消息标识(kafka必填，其他MQ填DEFAULT_DATA_KEY)
	* @param	value 		消息数据
	* @return   err 发布失败返回具体错误信息
	 */
	Publish(strBindingKey, strQueueName, key string, value string) (err error)

	/*
	* @brief 	消息消费接口定义
	* @param 	strBindingKey 	队列绑定Key
	* @param	strQueueName 	队列名称(redis/mqtt协议非必填)
	* @param    consumer        ReactHandler方法实现对象
	* @return   err 成功返回nil，失败返回返回具体错误信息
	* @remark   服务器异常或重启时内部会自动重连服务器
	 */
	Consume(strBindingKey, strQueueName string, consumer ReactHandler) (err error)

	/*
	* @brief 	开启或关闭调式模式
	* @param 	enable 	true开启/false关闭
	 */
	Debug(enable bool)

	/*
	* @brief 	获取当前MQ类型
	* @param 	adapter  MQ类型
	 */
	GetAdapter() (adapter Adapter)
}

type Instance func() IReactMQ

var AdapterMap = make(map[Adapter]Instance)

func init() {
	log.SetLevel("info")
}

//strLevel  ->   "debug"/"info"/"warn"/"error"  default "info"
func SetLogLevel(strLevel string) {
	log.SetLevel(strLevel)
}

//strLogPath  log file path
func SetLogPath(strLogPath string) {
	log.Open(strLogPath)
}

//注册对象创建方法
func Register(adapter Adapter, ins Instance) {

	if _, ok := AdapterMap[adapter]; !ok {

		AdapterMap[adapter] = ins
		return
	}
	log.Errorf("Adapter [%v] instance already exists", adapter)
	return
}

//适配器: 创建MQ对象
func NewMQ(adapter Adapter) (IReactMQ, error) {

	ins, ok := AdapterMap[adapter]
	if !ok {
		return nil, fmt.Errorf("Adapter [%v] instance not exists", adapter)
	}

	return ins(), nil
}
