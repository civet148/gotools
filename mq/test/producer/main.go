package main

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/mq"
	_ "github.com/civet148/gotools/mq/rabbit"
	"time"
)

func main() {
	var err error
	var strTopicRoutingKey = "TOPIC.NEWS.ROUTINGKEY"
	//var strTopicBindingKey = "TOPIC.NEWS.#" //#表示一个或多个单词（以.分割），*表示一个单词
	var strQueueName = "TOPIC.QUEUE"

	strConnUrl := "amqp://127.0.0.1:5672"
	rabbitMQ, _ := mq.NewMQ(mq.Adapter_RabbitMQ)
	if err = rabbitMQ.Connect(mq.Mode_Topic, strConnUrl); err != nil {
		log.Errorf("connect to MQ broker error [%v]", err.Error())
		return
	}

	_ = TopicProducer(rabbitMQ, strTopicRoutingKey, strQueueName)
}

func TopicProducer(ReactMQ mq.IReactMQ, strBindingKey, strQueueName string) (err error) {

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
