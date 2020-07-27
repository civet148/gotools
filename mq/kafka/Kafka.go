package kafka

import (
	"crypto/tls"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/mq"
	"github.com/civet148/gotools/parser"
)

type ReactKafka struct {
	config                                *cluster.Config //Kafka配置
	mode                                  mq.Mode
	url                                   string              //服务器连接URL
	host                                  string              //服务器地址:端口
	closed                                bool                //远程服务器是否已断开连接
	debug                                 bool                //开启或关闭调试信息
	tls                                   bool                //是否启用TLS
	ca, key, cer                          string              //TLS根证书、密钥和客户端证书
	strClientId, strUserName, strPassword string              //客户端ID、用户名和密码
	consumer                              *cluster.Consumer   //Kafka消费者
	producer                              sarama.SyncProducer //Kafka生产者
	ok                                    bool                //第一次连接是否成功
	strBindingKey                         string              //绑定键值
	strQueueName                          string              //队列名
	closing                               chan bool           //主动关闭通知通道
}

func init() {
	mq.Register(mq.Adapter_KafkaMQ, NewMQ)
}

func NewMQ() mq.IReactMQ {

	return &ReactKafka{
		closing: make(chan bool, 1),
	}
}

//MQTT 服务器接收数据回调
func (this *ReactKafka) OnReceive( /*client Kafka.Client, msg Kafka.Message*/ ) {

}

//MQTT 服务器连接成功回调
func (this *ReactKafka) OnConnect( /*Client Kafka.Client*/ ) {

}

//MQTT 服务器断开连接回调
func (this *ReactKafka) OnDisconnect( /*Client Kafka.Client,*/ err error) {

	this.closed = true
}

/*
* @brief 	MQ服务器连接接口定义
* @param 	strUrl 连接服务器URL(格式规范  kafka://127.0.0.1:9092)
* @return 	err 连接失败返回具体错误信息
 */
func (this *ReactKafka) Connect(mode mq.Mode, strURL string) (err error) {
	ui := parser.ParseUrl(strURL)
	this.url = strURL
	this.host = ui.Host
	this.config = cluster.NewConfig()
	this.config.Consumer.Return.Errors = true
	this.config.Group.Return.Notifications = true
	//this.config.Group.Mode = cluster.ConsumerModePartitions
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
		log.Debug("kafka reconnect to MQ broker [%v] ...", this.url)
	}

	if this.closed {
		err = fmt.Errorf("kafka auto reconnecting... please wait for remote server startup")
		log.Info("%v", err.Error())
	}
	return
}

//关闭MQ
func (this *ReactKafka) Close() {
	this.closing <- true
}

/*
* @brief 	判定是否MQ服务器断开连接（异常宕机或重启）
* @param
* @return 	远程服务器连接断开返回true，否则返回false
 */
func (this *ReactKafka) IsClosed() bool {

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
func (this *ReactKafka) Publish(strBindingKey, strQueueName, key string, value string) (err error) {

	if this.producer == nil {
		kc := sarama.NewConfig()
		kc.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
		kc.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
		kc.Producer.Return.Successes = true
		if this.producer, err = sarama.NewSyncProducer([]string{this.host}, kc); err != nil {
			log.Errorf("kafka new producer error [%v]", err.Error())
			return
		} else {
			log.Debugf("kafka new producer topic [%v] group [%v] ok", strBindingKey, strQueueName)
		}
	}
	if _, _, err = this.producer.SendMessage(&sarama.ProducerMessage{
		Topic: strBindingKey,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}); err != nil {
		log.Errorf("kafka producer SendMessage key [%s] error(%v)", key, err)
		return
	}
	//log.Debugf("kafka producer SendMessage key [%s] value [%v]", key, value)
	return
}

/*
* @brief 	消息消费接口定义
* @param 	strBindingKey 	队列绑定Key(topic)
* @param	strQueueName 	队列名称(group)
* @param    handler         消费回调处理对象
* @return   err 成功返回nil，失败返回返回具体错误信息
* @remark   服务器异常或重启时内部会自动重连服务器
 */
func (this *ReactKafka) Consume(strBindingKey, strQueueName string, handler mq.ReactHandler) (err error) {

	if this.consumer, err = cluster.NewConsumer([]string{this.host}, strQueueName, []string{strBindingKey}, this.config); err != nil {
		log.Errorf("kafka new consumer from [%v] error [%v]", this.host, err.Error())
		return
	}
	defer this.consumer.Close()
	log.Debugf("kafka new consumer from [%v] topic [%v] group [%v] ok", this.host, strBindingKey, strQueueName)

	// consume partitions
	for {
		select {
		case err := <-this.consumer.Errors():
			log.Errorf("kafka consumer error(%v)", err)
		case <-this.consumer.Notifications():
			//log.Infof("kafka consumer rebalanced")
		case msg, ok := <-this.consumer.Messages():
			//log.Debugf("kafka consumer from [%v] topic [%v] group [%v] message received ok [%v]", this.host, strBindingKey, strQueueName, ok)
			if ok {
				this.consumer.MarkOffset(msg, "")
				handler.OnConsume(mq.Adapter_KafkaMQ, strBindingKey, strQueueName, string(msg.Key), string(msg.Value))
			}
		//case part, ok := <-this.consumer.Partitions():
		//	log.Debugf("kafka consumer from [%v] topic [%v] group [%v] partitions [%+v] received ok (%v)", this.host, strBindingKey, strQueueName, part, ok)
		//	if !ok {
		//		log.Debugf("kafka partitions read channel is closed")
		//		return
		//	}
		//	// start a separate goroutine to consume messages
		//	go func(pc cluster.PartitionConsumer) {
		//		for msg := range pc.Messages() {
		//			log.Debugf("kafka topic [%s] partition [%d] offset [%d] key [%s] value [%s]", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
		//			this.consumer.MarkOffset(msg, "") // mark message as processed
		//			handler.OnConsume(mq.Adapter_KafkaMQ, strBindingKey, strQueueName, string(msg.Value))
		//		}
		//	}(part)
		case <-this.closing:
			log.Debugf("kafka consumer closing...")
			return
		}
	}

	log.Debugf("kafka new consumer from [%v] topic [%v] group [%v] returning", this.host, strBindingKey, strQueueName)
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

/*
* @brief 	获取当前MQ类型
* @param 	adapter  MQ类型
 */
func (this *ReactKafka) GetAdapter() (adapter mq.Adapter) {

	return mq.Adapter_KafkaMQ
}
