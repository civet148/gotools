module github.com/civet148/gotools

go 1.13

require (
	github.com/Shopify/sarama v1.24.1
	github.com/bsm/sarama-cluster v2.1.15+incompatible
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/garyburd/redigo v1.6.0
	github.com/gin-gonic/gin v1.5.0
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/websocket v1.4.1
	github.com/mattn/go-colorable v0.1.7
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins/registry/consul v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/registry/zookeeper v0.0.0-20200119172437-4fe21aa238fd
	github.com/robertkrimen/otto v0.0.0-20191219234010-c382bd3c16ff
	github.com/sideshow/apns2 v0.20.0
	github.com/spf13/viper v1.6.1
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	github.com/zheng-ji/goSnowFlake v0.0.0-20180906112711-fc763800eec9
	go.etcd.io/etcd v3.3.18+incompatible
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
)

replace github.com/micro/go-micro => github.com/Lofanmi/go-micro v1.16.1-0.20210804063523-68bbf601cfa4
