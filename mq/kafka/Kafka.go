package kafka

import (
	"crypto/tls"
	"fmt"
log "github.com/civet148/gotools/log"
"github.com/civet148/gotools/mq"
KafkaCluster "github.com/bsm/sarama-cluster"
)

var PARAM_TLS = "tls"
var PARAM_TLS_KEY = "key"
var PARAM_TLS_CER = "cer"
var PARAM_TLS_CA = "ca"
var PARAM_CLIENT_ID = "client-id"
var TLS_IS_ENABLE = "true"

var g_mapClientID = make(map[string]bool, 1) //全局客户端ID信息记录(对于冲突的客户端ID会无法创建连接)


type KafkaClient struct {

}

type ReactKafka struct {
	mode                       mq.Mode
	url                        string         //服务器连接URL
	host                       string         //服务器地址:端口
	closed                     bool           //远程服务器是否已断开连接
	debug                      bool           //开启或关闭调试信息
	tls                        bool           //是否启用TLS
	ca, key, cer               string         //TLS根证书、密钥和客户端证书
	clientid, username, passwd string         //客户端ID、用户名和密码
	conn                       KafkaCluster.Client    //Kafka客户端连接
	ok                         bool           //第一次连接是否成功
	cb                         mq.FnConsumeCb //消息回调方法
	bindingkey, queuename      string         //订阅关键字
}

func init() {
	mq.Register(mq.Adapter_KafkaMQ, NewMQ)
}

func NewMQ() (mq.IReactMQ) {

	return &ReactKafka{}
}

//MQTT 服务器接收数据回调
func (this *ReactKafka) OnReceive(/*client Kafka.Client, msg Kafka.Message*/) {


}

//MQTT 服务器连接成功回调
func (this *ReactKafka) OnConnect(/*Client Kafka.Client*/) {

}

//MQTT 服务器断开连接回调
func (this *ReactKafka) OnDisconnect(/*Client Kafka.Client,*/ err error) {

	this.closed = true
	//log.Error("MQTT client %p connection closed by remote server", Client)
}

/*
* @brief 	MQ服务器连接接口定义
* @param 	strUrl 连接服务器URL( 格式规范  kafka://username:password@host:port[/config?tls=[true|false]&&ca=ca.crt&key=client.key&cer=client.crt] )
* @return 	err 连接失败返回具体错误信息
 */
func (this *ReactKafka) Connect(mode mq.Mode, strURL string) (err error) {

	return
}

/*
* @brief 	MQ服务器重新连接接口定义
* @param
* @return 	err 连接失败返回具体错误信息
* @remark   当Publish返回错误且IsClosed()方法亦返回true时调用此方法重连MQ服务器
*           如果已使用Consume订阅过，内部会监听连接断开事件并自动重连。
 */
func (this *ReactKafka) Reconnect() (err error) {

	if this.debug {
		log.Debug("Reconnect to MQ broker [%v] ...", this.url)
	}

	if this.closed {
		err = fmt.Errorf("Auto reconnecting... please wait for remote server startup.")
		log.Info("%v", err.Error())
	}
	return
}

/*
* @brief 	判定是否MQ服务器断开连接（异常宕机或重启）
* @param
* @return 	远程服务器连接断开返回true，否则返回false
 */
func (this *ReactKafka) IsClosed() (bool) {

	return this.closed
}

/*
* @brief 	消息发布接口定义(仅支持字符串类型消息)
* @param 	strRoutingKey 	路由Key
* @param	strData 		消息数据
* @return   err 发布失败返回具体错误信息
 */
func (this *ReactKafka) Publish(strRoutingKey, strData string) (err error) {

	return
}

/*
* @brief 	消息消费接口定义
* @param 	strBidingKey 	队列绑定Key
* @param	strQueueName 	队列名称(MQTT忽略此参数)
* @return   err 成功返回nil，失败返回返回具体错误信息
* @remark   服务器异常或重启时内部会自动重连服务器
 */
func (this *ReactKafka) Consume(strBidingKey, strQueueName string, cb mq.FnConsumeCb) (err error) {


	return
}

/*
* @brief 	开启或关闭调式模式
* @param 	enable 	true开启/false关闭
 */
func (this *ReactKafka) Debug(enable bool) {

	this.debug = enable
}

func (this *ReactKafka) NewTLSConfig(ca, key, crt string) *tls.Config {

	return nil
}
