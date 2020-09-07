package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/wss"
	_ "github.com/civet148/gotools/wss/udpsock" //required (register socket instance)
	"time"
)

const (
	UDP_CLIENT_URL  = "udp://127.0.0.1:6664"
	UDP_SERVER_ADDR = "udp://127.0.0.1:6665"
)

func main() {
	client(UDP_CLIENT_URL)
}

const (
	UDP_DATA_PING = "ping"
	UDP_DATA_PONG = "pong"
)

func client(strUrl string) (err error) {
	c := wss.NewClient()
	if err = c.Listen(strUrl); err != nil {
		log.Errorf(err.Error())
		return
	}
	for {
		log.Debugf("local address [%v] remote address [%v]", c.GetLocalAddr(), c.GetRemoteAddr())
		if _, err := c.Send([]byte(UDP_DATA_PING), UDP_SERVER_ADDR); err != nil {
			log.Errorf(err.Error())
			break
		}

		if data, from, err := c.Recv(-1); err != nil {
			log.Error(err.Error())
			break
		} else {
			log.Infof("udp client received data [%s] length [%v] from [%v]", string(data), len(data), from)
		}
		time.Sleep(3 * time.Second)
	}
	return
}
