package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/wss"
	_ "github.com/civet148/gotools/wss/tcpsock" //required (register socket instance)
	"time"
)

const (
	TCP_DATA_PING = "ping"
	TCP_DATA_PONG = "pong"
)

func main() {
	var url = "tcp://127.0.0.1:6666"
	client(url)
}

func client(strUrl string) {
	c := wss.NewClient()
	if err := c.Connect(strUrl); err != nil {
		log.Errorf(err.Error())
		return
	}

	for {
		if _, err := c.Send([]byte(TCP_DATA_PING)); err != nil {
			log.Errorf(err.Error())
			break
		}

		if data, from, err := c.Recv(len(TCP_DATA_PONG)); err != nil {
			log.Error(err.Error())
			break
		} else {
			log.Infof("tcp client received data [%s] length [%v] from [%v]", string(data), len(data), from)
		}

		time.Sleep(3 * time.Second)
	}
}
