package main

import (
	"fmt"
	log "github.com/civet148/gotools/log"
	"github.com/civet148/gotools/mq"
	_ "github.com/civet148/gotools/mq/etcd"
	_ "github.com/civet148/gotools/mq/kafka"
	_ "github.com/civet148/gotools/mq/mqtt"
	_ "github.com/civet148/gotools/mq/rabbit"
	_ "github.com/civet148/gotools/mq/redis"
	"time"
)

type MyConsumer struct {
}

func main() {

	var err error
	var strConnUrl string

	log.Info("Program is running on...")

	mq.SetLogLevel("debug") //set debug mode

	var mode mq.Adapter
	mode = mq.Adapter_KafkaMQ //mq.Adapter_ETCD/mq.Adapter_RabbitMQ/mq.Adapter_RedisMQ/mq.Adapter_KafkaMQ/mq.Adapter_MQTT

	switch mode {
	case mq.Adapter_RabbitMQ:
		{
			// RabbitMQ无认证信息 amqp://192.168.1.15:5672
			// RabbitMQ带认证信息 amqp://test:123456@192.168.1.15:5672
			strConnUrl = "amqp://127.0.0.1:5672"
		}
	case mq.Adapter_RedisMQ:
		{
			// Redis无认证信息 redis://192.168.1.15:6379
			// Redis带认证信息 redis://123456@192.168.1.15:6379
			strConnUrl = "redis://127.0.0.1:6379"
		}
	case mq.Adapter_MQTT:
		{
			// 格式规范 mqtt://username:password@host:port[/config?tls=[true|false]&&ca=ca.crt&key=client.key&cer=client.crt&client-id=MyClientID]
			// SSL加密连接URL范例: "mqtt://192.168.1.15:8883/config?tls=true&ca=ca.crt&key=client.key&cer=client.crt&client-id=MyNameIsLory"
			// 非加密连接URL规范
			strConnUrl = "mqtt://127.0.0.1:1883/config?client-id="
		}
	case mq.Adapter_ETCD:
		{
			strConnUrl = "etcd://127.0.0.1:2379"
		}
	case mq.Adapter_KafkaMQ:
		{
			strConnUrl = "kafka://127.0.0.1:9092"
		}
	}

	/*
	 警告：MQTT协议下direct、topic、fanout模式最终行为一致，无需创建不同对象
	*/
	//if err = TestDirect(mode, strConnUrl); err != nil {
	//	log.Error("%v", err.Error())
	//}
	if err = TestTopic(mode, strConnUrl); err != nil {
		log.Error("%v", err.Error())
	}
	//if err = TestFanout(mode, strConnUrl); err != nil {
	//	log.Error("%v", err.Error())
	//}

	time.Sleep(1000 * time.Hour)
}

func (c *MyConsumer) OnConsume(adapter mq.Adapter, strBindingKey, strQueueName, strKey, strValue string) {
	log.Infof("[%+v] binding key [%v] queue name [%v] message key [%+v] value [%+v]", adapter, strBindingKey, strQueueName, strKey, strValue)
}

func TestDirect(mode mq.Adapter, strUrl string) (err error) {

	DirectMQ, _ := mq.NewMQ(mode)
	DirectMQ.Debug(false)
	if err = DirectMQ.Connect(mq.Mode_Direct, strUrl); err != nil {

		log.Error("%v", err.Error())
		return
	}
	log.Info("[DIRECT] Connect to [%v] MQ [%v] broker ok...", mode, strUrl)

	var strDirectKey, strQueueName string
	strQueueName = "DIRECT/QUEUE"
	if mode == mq.Adapter_MQTT {
		//MQTT模式测试
		strDirectKey = "DIRECT/ROUTINGKEY"
	} else {
		//其他MQ模式测试: 先使用消费者对象绑定队列，否则可能会导致生产者发布的消息被丢弃
		strDirectKey = "DIRECT.ROUTINGKEY"
	}

	go ConsumeDirect(DirectMQ, strDirectKey, strQueueName)
	go PublishDirect(DirectMQ, strDirectKey, strQueueName)

	return
}

func TestTopic(mode mq.Adapter, strUrl string) (err error) {
	TopicMQ, _ := mq.NewMQ(mode)

	TopicMQ.Debug(false)
	if err = TopicMQ.Connect(mq.Mode_Topic, strUrl); err != nil {

		log.Error("%v", err.Error())
		return
	}

	log.Info("[TOPIC] Connect to [%v] MQ [%v] broker ok...", mode, strUrl)

	var strTopicBindingKey, strTopicRoutingKey, strQueueName string

	if mode == mq.Adapter_MQTT {
		//MQTT模式测试
		strTopicRoutingKey = "TOPIC/NEWS/ROUTINGKEY"
		strTopicBindingKey = "TOPIC/NEWS/#" //#通配多级
		strQueueName = "TOPIC/QUEUE"

	} else if mode == mq.Adapter_KafkaMQ {
		//kafka模式测试(kafka不支持通配)
		strTopicRoutingKey = "TOPIC-NEWS-ROUTINGKEY" //topic
		strTopicBindingKey = strTopicRoutingKey      //topic
		strQueueName = "TOPIC.QUEUE"                 //group
	} else {
		//其他MQ模式测试
		strTopicRoutingKey = "TOPIC.NEWS.ROUTINGKEY"
		strTopicBindingKey = "TOPIC.NEWS.#" //#表示一个或多个单词（以.分割），*表示一个单词
		strQueueName = "TOPIC.QUEUE"
	}

	go ConsumeTopic(TopicMQ, strTopicBindingKey, strQueueName)
	go PublishTopic(TopicMQ, strTopicRoutingKey, strQueueName)
	return
}

func TestFanout(mode mq.Adapter, strUrl string) (err error) {

	FanoutMQ, _ := mq.NewMQ(mode)
	FanoutMQ.Debug(false)
	if err = FanoutMQ.Connect(mq.Mode_Fanout, strUrl); err != nil {

		log.Error("%v", err.Error())
		return
	}
	log.Info("[FANOUT] Connect to [%v] MQ [%v] broker ok...", mode, strUrl)
	var strFanoutBindingKey, strFanoutRoutingKey, strFanoutQueue string
	if mode == mq.Adapter_MQTT {
		//MQTT模式测试
		strFanoutRoutingKey = "FANOUT/NEWS/ROUTINGKEY"
		strFanoutBindingKey = "FANOUT/NEWS/ROUTINGKEY"
		strFanoutQueue = "FANOUT/QUEUE"
	} else {
		//其他MQ模式测试
		strFanoutRoutingKey = "FANOUT.NEWS.ROUTINGKEY"
		strFanoutBindingKey = "FANOUT.NEWS.ROUTINGKEY"
		strFanoutQueue = "FANOUT.QUEUE"
	}
	go ConsumeFanout(FanoutMQ, strFanoutBindingKey, strFanoutQueue)
	go PublishFanout(FanoutMQ, strFanoutRoutingKey, strFanoutQueue)

	return
}

func PublishDirect(ReactMQ mq.IReactMQ, strBindingKey, strQueueName string) (err error) {
	var nMsgIndex int
	var strData string = "This is direct data"
	time.Sleep(5 * time.Second)

	var strKey = mq.DEFAULT_DATA_KEY

	prod := ReactMQ

	for {

		strMsg := fmt.Sprintf("%v[%v]", strData, nMsgIndex)
		if err = prod.Publish(strBindingKey, strQueueName, strKey, strMsg); err != nil {
			log.Error("Publish [direct] data to broker error(%v)", err.Error())

			goto CONTINUE
		}
		log.Info("Publish [direct] data [%v] to broker ok", strMsg)
		nMsgIndex++

	CONTINUE:
		if prod.IsClosed() {
			if err = prod.Reconnect(); err != nil { //重新做一次连接尝试
				log.Error("Reconnect to MQ server error [%v]", err.Error())
			} else {
				log.Info("Reconnect to MQ server ok")
			}
		}

		time.Sleep(3 * time.Second)
	}
}

func ConsumeDirect(ReactMQ mq.IReactMQ, strBindingKey, strQueueName string) {

	c1 := ReactMQ
	errConsume := c1.Consume(strBindingKey, strQueueName, &MyConsumer{})
	if errConsume != nil {
		log.Error("%v", errConsume.Error())
		return
	}
}

func PublishTopic(ReactMQ mq.IReactMQ, strBindingKey, strQueueName string) (err error) {

	var bConnDown bool //MQ服务器异常
	var nMsgIndex int

	var strKey = mq.DEFAULT_DATA_KEY

	prod := ReactMQ
	for {

		if bConnDown { //服务器异常宕机或重启
			//ReactMQ.Connect(mq.Mode_Topic, )
		}
		var strData string = "This is topic data"
		strMsg := fmt.Sprintf("%v[%v]", strData, nMsgIndex)
		if err = prod.Publish(strBindingKey, strQueueName, strKey, strMsg); err != nil {
			log.Error("%v", err.Error())
			bConnDown = true
			goto CONTINUE
		}
		log.Info("Publish [topic] data [%v] to broker ok", strMsg)
		nMsgIndex++

	CONTINUE:
		if prod.IsClosed() {
			if err = prod.Reconnect(); err != nil { //重新做一次连接尝试
				log.Error("Reconnect to MQ server error [%v]", err.Error())
			} else {
				log.Info("Reconnect to MQ server ok")
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func ConsumeTopic(ReactMQ mq.IReactMQ, strBindingKey, strQueueName string) {
	c1 := ReactMQ
	errConsume := c1.Consume(strBindingKey, strQueueName, &MyConsumer{})
	if errConsume != nil {
		log.Error("%v", errConsume.Error())
		return
	}
}

func PublishFanout(ReactMQ mq.IReactMQ, strBindingKey, strQueueName string) (err error) {
	var nMsgIndex int
	var strData string = "This is fanout data"

	prod := ReactMQ

	var strKey = mq.DEFAULT_DATA_KEY
	if ReactMQ.GetAdapter() == mq.Adapter_KafkaMQ {
		strKey = "kafka-data-key"
	}

	for {
		strMsg := fmt.Sprintf("%v[%v]", strData, nMsgIndex)
		if err = prod.Publish(strBindingKey, strQueueName, strKey, strMsg); err != nil {
			log.Error("%v", err.Error())
			goto CONTINUE
		}
		log.Info("Publish [fanout] data [%v] to broker ok", strMsg)
		nMsgIndex++

	CONTINUE:
		if prod.IsClosed() {
			if err = prod.Reconnect(); err != nil { //重新做一次连接尝试
				log.Error("Reconnect to MQ server error [%v]", err.Error())
			} else {
				log.Info("Reconnect to MQ server ok")
			}
		}
		time.Sleep(9 * time.Second)

	}
}

func ConsumeFanout(ReactMQ mq.IReactMQ, strBindingKey, strQueueName string) {

	c1 := ReactMQ

	go func() {
		//fanout模式消费队列1
		errConsume := c1.Consume(strBindingKey, strQueueName, &MyConsumer{})
		if errConsume != nil {
			log.Error("%v", errConsume.Error())
			return
		}
	}()
}
