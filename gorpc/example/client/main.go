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
	SERVICE_NAME           = "echo"
	END_POINTS_HTTP_ETCD   = "http://127.0.0.1:2379"
	END_POINTS_HTTP_CONSUL = "http://127.0.0.1:8500"
	END_POINTS_ZOOKEEPER   = "127.0.0.1:2181"
)

func main() {

	c := NewClientWithMDNS()
	//c := NewClientWithEtcd()
	//c := NewClientWithConsul()
	//c := NewClientWithZk()
	service := echopb.NewEchoServerService(SERVICE_NAME, c)
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"X-User-Id": "lory",
		"X-From-Id": "10086",
	})
	for i := 0; i < 10; i++ {

		if pong, err := service.Call(ctx, &echopb.Ping{Text: "Ping"}); err != nil {
			log.Error(err.Error())
		} else {
			log.Infof("server response [%+v]", pong)
		}
		time.Sleep(2 * time.Second)
	}
}

func NewClientWithEtcd() (c client.Client) {
	var g *gorpc.GoRPC
	g = gorpc.NewGoRPC(gorpc.EndpointType_ETCD)
	endPoints := strings.Split(END_POINTS_HTTP_ETCD, ",")
	return g.NewClient(endPoints...)
}

func NewClientWithMDNS() (c client.Client) {
	var g *gorpc.GoRPC
	g = gorpc.NewGoRPC(gorpc.EndpointType_MDNS)
	return g.NewClient()
}

func NewClientWithConsul() (c client.Client) {
	var g *gorpc.GoRPC
	g = gorpc.NewGoRPC(gorpc.EndpointType_CONSUL)
	endPoints := strings.Split(END_POINTS_HTTP_CONSUL, ",")
	return g.NewClient(endPoints...)
}

func NewClientWithZk() (c client.Client) {
	var g *gorpc.GoRPC
	g = gorpc.NewGoRPC(gorpc.EndpointType_ZOOKEEPER)
	endPoints := strings.Split(END_POINTS_ZOOKEEPER, ",")
	return g.NewClient(endPoints...)
}
