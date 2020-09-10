package main

import (
	"context"
	"github.com/civet148/gotools/gorpc"
	"github.com/civet148/gotools/gorpc/example/echopb"
	"github.com/civet148/gotools/log"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"strings"
	"time"
)

const (
	SERVICE_NAME              = "echo"
	END_POINTS_HTTP_ETCD      = "127.0.0.1:2379"
	END_POINTS_HTTP_CONSUL    = "127.0.0.1:8500"
	END_POINTS_HTTP_ZOOKEEPER = "127.0.0.1:2181"
)

func main() {
	c := NewClientWithEtcd()
	//c := NewClientWithConsul()
	//c := NewClientWithZk()
	//c := NewClientWithMDNS()

	service := echopb.NewEchoServerService(SERVICE_NAME, c)
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"X-User-Id": "lory",
		"X-From-Id": "10086",
	})
	for i := 0; i < 10; i++ {

		if pong, err := service.Call(ctx, &echopb.Ping{Text: "Ping"}); err != nil {
			log.Error(err.Error())
			return
		} else {
			log.Infof("server response [%+v]", pong)
		}
		time.Sleep(2 * time.Second)
	}
}

func NewClientWithEtcd() (c client.Client) {
	endPoints := strings.Split(END_POINTS_HTTP_ETCD, ",")
	return gorpc.NewClient(gorpc.EndpointType_ETCD, endPoints...)
}

func NewClientWithMDNS() (c client.Client) {
	return gorpc.NewClient(gorpc.EndpointType_MDNS)
}

func NewClientWithConsul() (c client.Client) {
	endPoints := strings.Split(END_POINTS_HTTP_CONSUL, ",")
	return gorpc.NewClient(gorpc.EndpointType_CONSUL, endPoints...)
}

func NewClientWithZk() (c client.Client) {
	endPoints := strings.Split(END_POINTS_HTTP_ZOOKEEPER, ",")
	return gorpc.NewClient(gorpc.EndpointType_ZOOKEEPER, endPoints...)
}
