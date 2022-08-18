package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/mq"
	_ "github.com/civet148/gotools/mq/rabbit"
)

type ConsumerHandler struct {
}

func main() {
	var err error
	//var strTopicRoutingKey = "TOPIC.NEWS.ROUTINGKEY"
	var strTopicBindingKey = "TOPIC.NEWS.#" //#表示一个或多个单词（以.分割），*表示一个单词
	var strQueueName = "TOPIC.QUEUE"

	strConnUrl := "amqp://127.0.0.1:5672"
	rabbitMQ, _ := mq.NewMQ(mq.Adapter_RabbitMQ)
	if err = rabbitMQ.Connect(mq.Mode_Topic, strConnUrl); err != nil {
		log.Errorf("connect to MQ broker error [%v]", err.Error())
		return
	}

	_ = TopicConsumer(rabbitMQ, strTopicBindingKey, strQueueName)
}

func (c *ConsumerHandler) OnConsume(adapter mq.Adapter, strBindingKey, strQueueName, strKey, strValue string) {
	log.Infof("[%+v] binding key [%v] queue name [%v] message key [%+v] value [%+v]", adapter, strBindingKey, strQueueName, strKey, strValue)
}

func TopicConsumer(rabbitMQ mq.IReactMQ, strBindingKey, strQueueName string) (err error) {
	err = rabbitMQ.Consume(strBindingKey, strQueueName, &ConsumerHandler{})
	if err != nil {
		log.Error("%v", err.Error())
		return
	}
	return
}
