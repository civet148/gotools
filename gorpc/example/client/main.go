package main

import (
	"fmt"
	"github.com/civet148/gotools/gorpc"
	"github.com/civet148/gotools/gorpc/example/echopb"
	"github.com/civet148/gotools/log"
	"github.com/micro/go-micro/client"
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

	c := NewGoMicroClient(gorpc.EndpointType_MDNS)
	service := echopb.NewEchoServerService(SERVICE_NAME, c)

	for i := 0; i < 10; i++ {
		ctx := gorpc.NewContext(map[string]string{
			"X-User-Id": "lory",
			"X-From-Id": fmt.Sprintf("%d", 10000+i),
		}, 5)
		log.Debugf("send request [%v]", i)
		if pong, err := service.Call(ctx, &echopb.Ping{Text: "Ping"}); err != nil {
			log.Error(err.Error())
		} else {
			log.Infof("server response [%+v]", pong)
		}
		time.Sleep(2 * time.Second)
	}
}

func NewGoMicroClient(typ gorpc.EndpointType) (c client.Client) {
	var g *gorpc.GoRPC
	var endPoints []string
	g = gorpc.NewGoRPC(typ)
	switch typ {
	case gorpc.EndpointType_MDNS:
	case gorpc.EndpointType_ETCD:
		endPoints = strings.Split(END_POINTS_HTTP_ETCD, ",")
	case gorpc.EndpointType_CONSUL:
		endPoints = strings.Split(END_POINTS_HTTP_CONSUL, ",")
	case gorpc.EndpointType_ZOOKEEPER:
		endPoints = strings.Split(END_POINTS_ZOOKEEPER, ",")
	}
	return g.NewClient(endPoints...)
}
