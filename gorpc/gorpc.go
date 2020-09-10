package gorpc

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/micro/go-micro/registry/mdns"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/service/grpc"
	"time"
)

const (
	DISCOVERY_DEFAULT_INTERVAL = 3
	DISCOVERY_DEFAULT_TTL      = 10
)

type EndpointType int

const (
	EndpointType_MDNS      EndpointType = 0 // multicast DNS
	EndpointType_ETCD      EndpointType = 1 // etcd
	EndpointType_CONSUL    EndpointType = 2 // consul
	EndpointType_ZOOKEEPER EndpointType = 3 // zookeeper
)

func (t EndpointType) String() string {
	switch t {
	case EndpointType_MDNS:
		return "EndpointType_MDNS"
	case EndpointType_ETCD:
		return "EndpointType_ETCD"
	case EndpointType_CONSUL:
		return "EndpointType_CONSUL"
	case EndpointType_ZOOKEEPER:
		return "EndpointType_ZOOKEEPER"
	}
	return "EndpointType_Unknown"
}

type Discovery struct {
	ServiceName string   // register service name [required]
	RpcAddr     string   // register server RPC address [required]
	Interval    int      // register interval default 3 seconds [optional]
	TTL         int      // register TTL default 10 seconds [optional]
	Endpoints   []string // register endpoints of etcd/consul/zookeeper eg. ["192.168.0.10:2379","192.168.0.11:2379"]
}

func NewClient(endpointType EndpointType, endPoints ...string) (c client.Client) { // returns go-micro client object

	var options []micro.Option

	log.Debugf("endpoint type [%v] end points [%+v]", endpointType, endPoints)

	reg := newRegistry(endpointType, endPoints...)
	if reg != nil {
		options = append(options, micro.Registry(reg))
	}
	service := grpc.NewService(options...)
	return service.Client()
}

func NewServer(endpointType EndpointType, discovery *Discovery) (s server.Server) { // returns go-micro server object
	log.Debugf("endpoint type [%v] discovery [%+v]", endpointType, discovery)
	if len(discovery.Endpoints) == 0 && endpointType != EndpointType_MDNS {
		panic("discovery end points is nil and not EndpointType_MDNS")
	}
	if discovery.ServiceName == "" {
		panic("discover service name is nil")
	}
	if discovery.Interval == 0 {
		discovery.Interval = DISCOVERY_DEFAULT_INTERVAL
	}
	if discovery.TTL == 0 {
		discovery.TTL = DISCOVERY_DEFAULT_TTL
	}

	reg := newRegistry(endpointType, discovery.Endpoints...)

	var options []micro.Option
	if reg == nil {
		panic(fmt.Errorf("[%+v] discovery [%+v] -> registry is nil", endpointType, discovery))
	}
	options = append(options, micro.Registry(reg))
	options = append(options, micro.RegisterInterval(time.Duration(discovery.Interval)*time.Second))
	options = append(options, micro.RegisterTTL(time.Duration(discovery.TTL)*time.Second))
	options = append(options, micro.Name(discovery.ServiceName))
	options = append(options, micro.Address(discovery.RpcAddr))
	service := grpc.NewService(options...)
	return service.Server()
}

func newRegistry(endpointType EndpointType, endPoints ...string) (r registry.Registry) {

	regAddrs := registry.Addrs(endPoints...)
	switch endpointType {
	case EndpointType_MDNS:
		r = mdns.NewRegistry()
	case EndpointType_ETCD:
		r = etcd.NewRegistry(regAddrs)
	case EndpointType_CONSUL:
	case EndpointType_ZOOKEEPER:
	}
	log.Debugf("[%+v] end points [%+v] -> registry [%+v]", endpointType, endPoints, r)
	return
}
