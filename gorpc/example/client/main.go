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
	SERVICE_NAME    = "echo"
	END_POINTS_ETCD = "http://127.0.0.1:2379"
)

func main() {
	//c := NewClientWithEtcd()
	c := NewClientWithMDNS()
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
	endPoints := strings.Split(END_POINTS_ETCD, ",")
	return gorpc.NewClient(gorpc.EndpointType_ETCD, endPoints...)
}

func NewClientWithMDNS() (c client.Client) {
	return gorpc.NewClient(gorpc.EndpointType_MDNS)
}
