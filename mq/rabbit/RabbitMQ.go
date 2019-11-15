package rabbit

import (
	"fmt"
	log "github.com/civet148/gotools/log"
	"github.com/civet148/gotools/comm"
	"github.com/civet148/gotools/mq"
	"github.com/streadway/amqp"
	"strings"
	"sync"
	"time"
)

type RabbitMQ struct {
	mode     mq.Mode
	ex       string           //交换器名称
	tag      string           //消费者标签
	url      string           //服务器连接URL
	conn     *amqp.Connection //AMQP连接会话
	producer *amqp.Channel    //生产者信道
	lock   sync.Mutex //断线重连锁
	closed bool       //远程服务器是否已断开连接
	debug  bool //开启或关闭调试信息
}

func init() {
	mq.Register(mq.Adapter_RabbitMQ, NewMQ)
}

func NewMQ() (mq.IReactMQ) {

	return &RabbitMQ{}
}

//判定是否服务器重启断开客户端连接
func (this *RabbitMQ) IsClosed() bool {

	return this.closed
}

func (this *RabbitMQ) Connect(mode mq.Mode, strURL string) (err error) {

	this.mode = mode
	this.url = strURL
	this.ex = this.exchangeName(mode)
	this.tag = this.consumerTag(mode)

	if err = this.Reconnect(); err != nil {
		if this.debug {
			log.Error("Connect [%v] failed, error [%v]", this.url, err.Error())
		}
		return
	}

	return
}

//当调用Publish返回错误并且IsClosed()返回true时，可调用此方法发起重连
func (this *RabbitMQ) Reconnect() (err error) {

	this.lock.Lock()
	defer this.lock.Unlock()

	if this.closed { //断线重连之前关闭旧连接和通道
		if this.conn != nil {
			this.conn.Close()
		}

		if this.producer != nil {
			this.producer.Close()
		}
	}
	if this.debug {
		log.Debug("Try to connect MQ server...")
	}
	if this.conn, err = amqp.Dial(this.url); err != nil {

		return
	}

	if this.debug {
		log.Debug("Try to get channel ...")
	}
	if this.producer, err = this.getChannel(); err != nil {
		if this.debug {
			log.Error("Get producer channel error [%v]", err.Error())
		}
		return
	}

	if this.closed {
		this.closed = false //连接成功重置标记状态
	}
	return
}

func (this *RabbitMQ) Publish(strRoutingKey, strData string) (err error) {

	if this.closed {

		err = fmt.Errorf("Connection still invalid...")
		return
	}
	if this.debug {
		log.Info("[PRODUCER] mode [%v] key [%v] data [%v]", this.mode.String(), strRoutingKey, strData)
	}

	if err = this.producer.Publish (
		this.ex, 	// publish to an exchange
		strRoutingKey,            // routing to 0 or more queues
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(strData),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		if  strings.Contains(err.Error(), "is not open") ||
			strings.Contains(err.Error(),"504") {//服务器断开连接
			this.closed = true //标记连接断开
		}
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	//log.Debug("Publish direct to exchange [%v] routing key [%v] data [%v] ok", this.ex, strRoutingKey, strData)
	return
}


func (this *RabbitMQ) Consume(strBidingKey, strQueueName string, cb mq.FnConsumeCb) (err error) {

CONSUME_BEGIN:
	if this.debug {
		log.Info("[CONSUMER] mode [%v] key [%v] queue [%v]", this.mode.String(), strBidingKey, strQueueName)
	}
	var channel *amqp.Channel
	if channel, err = this.getChannel(); err != nil {
		if this.debug {
			log.Error("Get consumer channel error [%v]", err.Error())
		}
		return
	}

	queue, errQueue := channel.QueueDeclare (
		strQueueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if errQueue != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	if this.debug {
		log.Debug("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
			queue.Name, queue.Messages, queue.Consumers, strBidingKey)
	}

	if err = channel.QueueBind(
		queue.Name, // name of the queue
		strBidingKey,        // bindingKey
		this.ex,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	if this.debug {
		log.Debug("Queue bound to Exchange, starting Consume with binding key [%v] queue [%v]", strBidingKey, strQueueName)
	}
	deliveries, errConsume := channel.Consume(
		queue.Name, // name
		this.tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if errConsume != nil {
		return fmt.Errorf("Queue Consume: %s", errConsume)
	}

	if this.debug {
		log.Debug("Consumer bind queue ok, ready to receive message")
	}

	for d := range deliveries {
		if this.debug {
			log.Debug("got %dB delivery: [%v] %q", len(d.Body), d.DeliveryTag, d.Body)
		}
		d.Ack(false)
		cb(string(d.Body))//收到数据后通过调用者传入的方法完成数据回调通知
	}

	for {
		this.closed = true //标记连接断开
		channel.Close() //释放信道
		if err = this.Reconnect(); err == nil {//重连成功，重新执行消费消息代码
			if this.debug {
				log.Info("Reconnect rabbitmq server ok, continue receive message...")
			}
			goto CONSUME_BEGIN
		}
		time.Sleep(5*time.Second) //5秒重连一次
	}

	if this.debug {
		log.Warn("Consumer mode [%v] channel closed by remote server", this.mode.String())
	}
	return
}

/*
* @brief 	开启或关闭调式模式
* @param 	enable 	true开启/false关闭
*/
func (this *RabbitMQ) Debug(enable bool) {

	this.debug = enable
}

func (this *RabbitMQ) getChannel() (channel *amqp.Channel, err error) {

	if this.conn == nil {
		return nil, fmt.Errorf("MQ connection is nil")
	}

	channel, err = this.conn.Channel()
	if err != nil {
		if this.debug {
			log.Error("create channel error [%v]", err.Error())
		}
		return
	}

	if err = channel.ExchangeDeclare (
		this.ex,            // exchange name
		this.mode.String(), // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // noWait
		nil,                // arguments
	); err != nil {
		if this.debug {
			log.Error("ExchangeDeclare return [%v]", err.Error())
		}
		return
	}
	return
}

func (this *RabbitMQ) exchangeName(mode mq.Mode) (exchg string) {

	switch mode {
	case mq.Mode_Direct:
		exchg = mq.EXCHANGE_NAME_DIRECT
	case mq.Mode_Topic:
		exchg = mq.EXCHANGE_NAME_TOPIC
	case mq.Mode_Fanout:
		exchg = mq.EXCHANGE_NAME_FANOUT
	}
	return
}

func (this *RabbitMQ) consumerTag(mode mq.Mode) (tag string) {

	tag = fmt.Sprintf("consumer.tag.%v#%v", mode.String(), comm.GenRandStrMD5(16))
	return
}
