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
	SERVICE_NAME    = "echo"
	END_POINTS_ETCD = "http://127.0.0.1:2379" //more end points "http://127.0.0.1:2379,http://127.0.0.1:3379"
	RPC_ADDR        = "127.0.0.1:8899"
)

type EchoServerImpl struct {
}

func main() {
	ch := make(chan bool, 1)
	//srv := NewServerWithEtcd()
	srv := NewServerWithMDNS()

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
	return gorpc.NewServer(gorpc.EndpointType_ETCD, &gorpc.Discovery{
		ServiceName: SERVICE_NAME,
		RpcAddr:     RPC_ADDR,
		Interval:    3,
		TTL:         10,
		Endpoints:   strings.Split(END_POINTS_ETCD, ","),
	})
}

func NewServerWithMDNS() (s server.Server) {
	return gorpc.NewServer(gorpc.EndpointType_MDNS, &gorpc.Discovery{
		ServiceName: SERVICE_NAME,
		RpcAddr:     RPC_ADDR,
		Interval:    3,
		TTL:         10,
		Endpoints:   []string{},
	})
}

func (s *EchoServerImpl) Call(ctx context.Context, ping *echopb.Ping, pong *echopb.Pong) (err error) {
	md, _ := metadata.FromContext(ctx)
	log.Infof("md [%+v] req [%+v]", md, ping)
	pong.Text = "Pong"
	return
}
