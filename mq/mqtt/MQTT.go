package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/civet148/gotools/comm"
	log "github.com/civet148/gotools/log"
	"github.com/civet148/gotools/mq"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"net/url"
	"time"
)

var PARAM_TLS = "tls"
var PARAM_TLS_KEY = "key"
var PARAM_TLS_CER = "cer"
var PARAM_TLS_CA = "ca"
var PARAM_CLIENT_ID = "client-id"
var TLS_IS_ENABLE = "true"

var g_mapClientID = make(map[string]bool, 1) //全局客户端ID信息记录(对于冲突的客户端ID会无法创建连接)

type ReactMQTT struct {
	mode                       mq.Mode
	url                        string          //服务器连接URL
	host                       string          //服务器地址:端口
	closed                     bool            //远程服务器是否已断开连接
	debug                      bool            //开启或关闭调试信息
	tls                        bool            //是否启用TLS
	ca, key, cer               string          //TLS根证书、密钥和客户端证书
	clientid, username, passwd string          //客户端ID、用户名和密码
	conn                       MQTT.Client     //MQTT客户端连接
	ok                         bool            //第一次连接是否成功
	handler                    mq.ReactHandler //消息回调方法
	strBindingKey              string          //绑定键值
	strQueueName               string          //队列名
	closing                    chan bool       //主动关闭通知通道
}

func init() {
	mq.Register(mq.Adapter_MQTT, NewMQ)
}

func NewMQ() mq.IReactMQ {

	return &ReactMQTT{
		closing: make(chan bool, 1),
	}
}

//MQTT 服务器接收数据回调
func (this *ReactMQTT) OnReceive(client MQTT.Client, msg MQTT.Message) {

	if this.debug {
		log.Debug("TOPIC: [%s] MSG: [%s]", msg.Topic(), msg.Payload())
	}
	if this.handler != nil {
		this.handler.OnConsume(mq.Adapter_ETCD, this.strBindingKey, this.strQueueName, mq.DEFAULT_DATA_KEY, string(msg.Payload()))
	}
}

//MQTT 服务器连接成功回调
func (this *ReactMQTT) OnConnect(Client MQTT.Client) {

	if this.ok && this.closed { //断线重连成功后，重新订阅

		if this.strBindingKey != "" { //断线之前已订阅过
			this.Consume(this.strBindingKey, this.strQueueName, this.handler)
		}
	}
	this.ok = true
	this.closed = false

	if this.debug {
		log.Debug("Client %p this.conn=[%p] connect to MQ ok", Client, this.conn)
	}
}

//MQTT 服务器断开连接回调
func (this *ReactMQTT) OnDisconnect(Client MQTT.Client, err error) {

	this.closed = true
	log.Error("MQTT client %p connection closed by remote server", Client)
}

/*
* @brief 	MQ服务器连接接口定义
* @param 	strUrl 连接服务器URL( 格式规范  mqtt://username:password@host:port[/config?tls=[true|false]&&ca=ca.crt&key=client.key&cer=client.crt&client-id=MyClientID] )
* @return 	err 连接失败返回具体错误信息
 */
func (this *ReactMQTT) Connect(mode mq.Mode, strURL string) (err error) {

	this.url = strURL
	this.mode = mode

	var u *url.URL
	if u, err = url.Parse(strURL); err != nil {
		log.Error("URL format error (%v)", err.Error())
		return
	}
	if this.debug {
		log.Debug("URL [%#v]", u)
	}

	if u.User != nil {
		this.username = u.User.Username()
		this.passwd, _ = u.User.Password()
	}
	m, _ := url.ParseQuery(u.RawQuery)

	if this.debug {
		for k, v := range m {
			log.Debug("param key=[%v] value=%v", k, v)
		}
	}

	if v, ok := m[PARAM_TLS]; ok {

		if v[0] == TLS_IS_ENABLE { //tls参数设置为true，启用TLS连接
			this.tls = true
		}
	}

	if v, ok := m[PARAM_CLIENT_ID]; ok { //URL中查找客户端ID
		this.clientid = v[0]
	}
	if v, ok := m[PARAM_TLS_KEY]; ok { //URL中查找TLS密钥KEY文件名
		this.key = v[0]
	}
	if v, ok := m[PARAM_TLS_CER]; ok { //URL中查找TLS证书文件名
		this.cer = v[0]
	}
	if v, ok := m[PARAM_TLS_CA]; ok { //URL中查找CA证书文件名
		this.ca = v[0]
	}

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	var strConnUrl string
	if !this.tls {
		strConnUrl = fmt.Sprintf("tcp://%v", u.Host)
	} else {
		strConnUrl = fmt.Sprintf("ssl://%v", u.Host)
	}

	if !this.ok || this.closed { //尚未连接或连接已断开

		opts := MQTT.NewClientOptions()
		if this.clientid != "" {
			if _, ok := g_mapClientID[this.clientid]; ok {

				err = fmt.Errorf("Client ID [%v] is in use, can't connect MQ server with duplicate client id", this.clientid)
				log.Error("%v", err.Error())
				return
			}
		} else {
			this.clientid = comm.GenRandStrMD5(32) //随机客户端ID
		}
		opts.SetClientID(this.clientid) //设置客户端ID，同一个客户端ID可能会导致其他连接被踢掉线
		opts.AddBroker(strConnUrl)      //添加MQ服务器连接URL
		opts.SetDefaultPublishHandler(this.OnReceive)
		opts.SetAutoReconnect(true)
		opts.SetMaxReconnectInterval(3 * time.Second)
		opts.SetKeepAlive(5 * time.Second)
		opts.SetWriteTimeout(15 * time.Second)
		opts.SetConnectionLostHandler(this.OnDisconnect)
		opts.SetOnConnectHandler(this.OnConnect)

		if this.username != "" {
			opts.Username = this.username
		}
		if this.passwd != "" {
			opts.Password = this.passwd
		}

		if this.tls {
			opts.SetTLSConfig(this.NewTLSConfig(this.ca, this.key, this.cer))
		}

		if this.debug {
			log.Debug("Ready to connect MQTT %v...", strConnUrl)
		}

		//create and start a client using the above ClientOptions
		c := MQTT.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			log.Error("Connect error (%v)", token.Error())
			return token.Error()
		}
		this.conn = c
		g_mapClientID[this.clientid] = true //保存客户端ID
	}

	if this.debug {
		log.Debug("%#v", *this)
	}

	if this.debug {
		log.Info("MQTT [%v] connect ok", this.url)
	}
	return
}

/*
* @brief 	MQ服务器重新连接接口定义
* @param
* @return 	err 连接失败返回具体错误信息
* @remark   当Publish返回错误且IsClosed()方法亦返回true时调用此方法重连MQ服务器
*           如果已使用Consume订阅过，内部会监听连接断开事件并自动重连。
 */
func (this *ReactMQTT) Reconnect() (err error) {

	if this.debug {
		log.Debug("Reconnect to MQ broker [%v] ...", this.url)
	}

	if this.closed {
		err = fmt.Errorf("Auto reconnecting... please wait for remote server startup.")
		log.Info("%v", err.Error())
	}
	return
}

//关闭MQ
func (this *ReactMQTT) Close() {
	this.closing <- true
}

/*
* @brief 	判定是否MQ服务器断开连接（异常宕机或重启）
* @param
* @return 	远程服务器连接断开返回true，否则返回false
 */
func (this *ReactMQTT) IsClosed() bool {

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
func (this *ReactMQTT) Publish(strBindingKey, strQueueName, key string, value string) (err error) {

	token := this.conn.Publish(strBindingKey, 1, false, value)

	if token.Wait(); token.Error() != nil {
		err = token.Error()
		log.Error("error (%s)", token.Error())
		return
	}
	if this.debug {
		log.Info("publish data to [%v] routing key [%v] ok", this.mode.String(), strBindingKey)
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
func (this *ReactMQTT) Consume(strBindingKey, strQueueName string, handler mq.ReactHandler) (err error) {

	if this.debug {
		log.Debug("Subscribe start with bindingkey [%v] queue [%v] ", strBindingKey, strQueueName)
	}

	if this.strBindingKey != "" {
		this.conn.Unsubscribe(this.strBindingKey) //断线重连重置订阅状态
	}

	this.strBindingKey = strBindingKey
	this.strQueueName = strQueueName
	this.handler = handler
	if token := this.conn.Subscribe(strBindingKey, 2, nil); token.Wait() && token.Error() != nil {

		err = token.Error()
		log.Error("Subscribe error (%v)", err.Error())
		return
	}
	return
}

/*
* @brief 	开启或关闭调式模式
* @param 	enable 	true开启/false关闭
 */
func (this *ReactMQTT) Debug(enable bool) {

	this.debug = enable
}

func (this *ReactMQTT) NewTLSConfig(ca, key, crt string) *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(ca)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	if this.debug {
		log.Debug("read pemCerts success")
	}

	// Import client certificate/key pair
	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		panic(err)
	}

	if this.debug {
		log.Debug("read LoadX509KeyPair from key and crt success")
	}

	// Just to print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}
	if this.debug {
		log.Debug("read cert.Leaf(ParseCertificate) from cert.Certificate success")
	}
	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}
}

/*
* @brief 	获取当前MQ类型
* @param 	adapter  MQ类型
 */
func (this *ReactMQTT) GetAdapter() (adapter mq.Adapter) {

	return mq.Adapter_MQTT
}
