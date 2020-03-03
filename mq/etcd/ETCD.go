package etcd

import (
	"context"
	"fmt"
	"github.com/civet148/gotools/comm"
	log "github.com/civet148/gotools/log"
	"github.com/civet148/gotools/mq"
	"go.etcd.io/etcd/clientv3"
	"strings"
	"sync"
	"time"
)

const (
	ETCD_SCHEMA_PRIFIX = "etcd://"
)

type EtcdMQ struct {
	mode   mq.Mode
	ex     string //交换器名称
	tag    string //消费者标签
	url    string //服务器连接URL
	config clientv3.Config
	conn   *clientv3.Client //连接会话
	lock   sync.Mutex       //断线重连锁
	closed bool             //远程服务器是否已断开连接
	debug  bool             //开启或关闭调试信息
}

func init() {
	mq.Register(mq.Adapter_ETCD, NewMQ)
}

func NewMQ() mq.IReactMQ {

	return &EtcdMQ{}
}

//判定是否服务器重启断开客户端连接
func (this *EtcdMQ) IsClosed() bool {

	return this.closed
}

//strUrl格式 单机"etcd://127.0.0.1:2379" 集群 "etcd://127.0.0.1:2379,127.0.0.1:2479,..."
func (this *EtcdMQ) Connect(mode mq.Mode, strURL string) (err error) {

	if strings.Contains(strURL, ETCD_SCHEMA_PRIFIX) {

		strURL = func(args ...string) string {
			if len(args) == 2 {
				return args[1]
			}
			return strURL
		}(strings.Split(strURL, ETCD_SCHEMA_PRIFIX)...)
	}

	s := strings.Split(strURL, ",")
	if len(s) < 1 {
		log.Error("Connect strUrl Split failed,strURL = %s", strURL)
		panic(strURL)
	}
	this.config = clientv3.Config{
		Endpoints:   s,
		DialTimeout: 5 * time.Second,
	}
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
func (this *EtcdMQ) Reconnect() (err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.closed { //断线重连之前关闭旧连接和通道
		if this.conn != nil {
			this.conn.Close()
		}
	}
	if this.debug {
		log.Debug("Try to connect MQ server...")
	}
	if this.conn, err = clientv3.New(this.config); err != nil {
		if this.debug {
			log.Error("Connect [%v] failed, error [%v]", this.url, err.Error())
			return
		}
	}
	if this.closed {
		this.closed = false //连接成功重置标记状态
	}
	return
}

func (this *EtcdMQ) Publish(strRoutingKey, strData string) (err error) {
	if this.closed {
		err = fmt.Errorf("Connection still invalid...")
		return
	}
	if this.debug {
		log.Info("[PRODUCER] mode [%v] key [%v] data [%v]", this.mode.String(), strRoutingKey, strData)
	}
	_, err = this.conn.Put(context.Background(), strRoutingKey, strData)
	if err != nil {
		if strings.Contains(err.Error(), "is not open") ||
			strings.Contains(err.Error(), "504") { //服务器断开连接
			this.closed = true //标记连接断开
		}
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	//log.Debug("Publish direct to exchange [%v] routing key [%v] data [%v] ok", this.ex, strRoutingKey, strData)
	return
}

func (this *EtcdMQ) Consume(strBindingKey, strQueueName string, cb mq.FnConsumeCb) (err error) {

	switch this.mode {
	case mq.Mode_Direct:
		return this.consumeDirect(strBindingKey, strQueueName, cb)
	case mq.Mode_Topic:
		return this.consumeTopic(strBindingKey, strQueueName, cb)
	case mq.Mode_Fanout:
		return this.consumeFanout(strBindingKey, strQueueName, cb)
	default:
		return fmt.Errorf("Unknown mode [%v]", this.mode)
	}
	return
}

func (this *EtcdMQ) getQueueName(strKeyName string) (name string) {

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

func (this *EtcdMQ) consumeDirect(strBindingKey, strQueueName string, cb mq.FnConsumeCb) (err error) {
	var strData string
	wc := this.conn.Watch(context.TODO(), strBindingKey)
	/*getResp,err := this.conn.KV.Get(context.TODO(),strBindingKey,clientv3.WithPrevKV())
	if err!= nil{
		log.Error("Get 失败：%s", err.Error())
		panic(err)
	}
	if getResp.Count <= 0{
		_, err := this.conn.KV.Put(context.TODO(), strBindingKey, "")
		if err != nil {
			log.Error("put 失败：%s", err.Error())
			panic(err)
		}
	}*/
RETRY_CONSUME:
	if this.closed {
		if err = this.Reconnect(); err != nil {
			if this.debug {
				log.Error("Connect [%v] failed, error [%v]", this.url, err.Error())
			}
			return
		}
	}
	for {
		for v := range wc {
			for _, e := range v.Events {
				strData = string(e.Kv.Value)
				cb(strData)
			}
		}
		if this.closed {
			time.Sleep(1 * time.Second)
			goto RETRY_CONSUME
		}
		if this.debug {
			log.Info("[CONSUMER] mode [%v] key [%v] data [%v]", this.mode.String(), strBindingKey, strData)
		}
	}
	return
}

func (this *EtcdMQ) consumeTopic(strBindingKey, strQueueName string, cb mq.FnConsumeCb) (err error) {
	var strData string
	go func() {
		wc := this.conn.Watch(context.Background(), strBindingKey) //, clientv3.WithPrefix()
	RETRY_CONSUME:
		if this.closed {
			if err = this.Reconnect(); err != nil {
				if this.debug {
					log.Error("Connect [%v] failed, error [%v]", this.url, err.Error())
				}
				return
			}
		}
		for {
			for v := range wc {
				for _, e := range v.Events {
					strData = string(e.Kv.Value)
					cb(strData)
				}
			}
			if this.closed {
				time.Sleep(3 * time.Second)
				goto RETRY_CONSUME
			}
			if this.debug {
				log.Info("[CONSUMER] mode [%v] key [%v] data [%v]", this.mode.String(), strBindingKey, strData)
			}
		}
	}()
	return
}

func (this *EtcdMQ) consumeFanout(strBindingKey, strQueueName string, cb mq.FnConsumeCb) (err error) {

	var strData string
	wc := this.conn.Watch(context.Background(), strBindingKey, clientv3.WithPrefix(), clientv3.WithPrevKV())
RETRY_CONSUME:
	if this.closed {
		if err = this.Reconnect(); err != nil {
			if this.debug {
				log.Error("Connect [%v] failed, error [%v]", this.url, err.Error())
			}
			return
		}
	}
	for {
		for v := range wc {
			for _, e := range v.Events {
				strData = string(e.Kv.Value)
				cb(strData)
			}
		}
		if this.closed {
			time.Sleep(3 * time.Second)
			goto RETRY_CONSUME
		}
		if this.debug {
			log.Info("[CONSUMER] mode [%v] key [%v] data [%v]", this.mode.String(), strBindingKey, strData)
		}
	}
	return
}

/*
* @brief 	开启或关闭调式模式
* @param 	enable 	true开启/false关闭
 */
func (this *EtcdMQ) Debug(enable bool) {

	this.debug = enable
}

func (this *EtcdMQ) exchangeName(mode mq.Mode) (exchg string) {

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

func (this *EtcdMQ) consumerTag(mode mq.Mode) (tag string) {

	tag = fmt.Sprintf("consumer.tag.%v#%v", mode.String(), comm.GenRandStrMD5(16))
	return
}
