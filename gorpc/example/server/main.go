package main

import (
	"context"
	"github.com/civet148/gotools/gorpc"
	"github.com/civet148/gotools/gorpc/example/echopb"
	"github.com/civet148/gotools/log"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"strings"
)

const (
	SERVICE_NAME           = "echo"
	END_POINTS_HTTP_ETCD   = "http://127.0.0.1:2379"
	END_POINTS_HTTP_CONSUL = "http://127.0.0.1:8500"
	END_POINTS_ZOOKEEPER   = "127.0.0.1:2181"
	RPC_ADDR               = "127.0.0.1:8899" //RPC service listen address
)

type EchoServerImpl struct {
}

func main() {
	ch := make(chan bool, 1)
	srv := NewServerWithMDNS()
	//srv := NewServerWithEtcd()
	//srv := NewServerWithConsul()
	//srv := NewServerWithZK()

	if err := echopb.RegisterEchoServerHandler(srv, new(EchoServerImpl)); err != nil {
		log.Error(err.Error())
		return
	}
	//go-micro v1.16 call srv.Run() v1.18 call srv.Start()
	if err := srv.Start(); err != nil {
		log.Error(err.Error())
		return
	}

	<-ch //block infinite
}

func NewServerWithEtcd() (s server.Server) {
	var g *gorpc.GoRPC
	g = gorpc.NewGoRPC(gorpc.EndpointType_ETCD)
	return g.NewServer(&gorpc.Discovery{
		ServiceName: SERVICE_NAME,
		RpcAddr:     RPC_ADDR,
		Interval:    3,
		TTL:         10,
		Endpoints:   strings.Split(END_POINTS_HTTP_ETCD, ","),
	})
}

func NewServerWithMDNS() (s server.Server) {
	var g *gorpc.GoRPC
	g = gorpc.NewGoRPC(gorpc.EndpointType_MDNS)
	return g.NewServer(&gorpc.Discovery{
		ServiceName: SERVICE_NAME,
		RpcAddr:     RPC_ADDR,
		Interval:    3,
		TTL:         10,
		Endpoints:   []string{},
	})
}

func NewServerWithConsul() (s server.Server) {
	var g *gorpc.GoRPC
	g = gorpc.NewGoRPC(gorpc.EndpointType_CONSUL)
	return g.NewServer(&gorpc.Discovery{
		ServiceName: SERVICE_NAME,
		RpcAddr:     RPC_ADDR,
		Interval:    3,
		TTL:         10,
		Endpoints:   strings.Split(END_POINTS_HTTP_CONSUL, ","),
	})
}

func NewServerWithZK() (s server.Server) {
	var g *gorpc.GoRPC
	g = gorpc.NewGoRPC(gorpc.EndpointType_ZOOKEEPER)
	return g.NewServer(&gorpc.Discovery{
		ServiceName: SERVICE_NAME,
		RpcAddr:     RPC_ADDR,
		Interval:    3,
		TTL:         10,
		Endpoints:   strings.Split(END_POINTS_ZOOKEEPER, ","),
	})
}

func (s *EchoServerImpl) Call(ctx context.Context, ping *echopb.Ping, pong *echopb.Pong) (err error) {
	md, _ := metadata.FromContext(ctx)
	log.Infof("md [%+v] req [%+v]", md, ping)
	pong.Text = "Pong"
	return
}
