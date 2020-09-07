package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/wss"
	_ "github.com/civet148/gotools/wss/websock"
	"time"
)

const (
	WEBSOCKET_DATA_PING = "ping"
	WEBSOCKET_DATA_PONG = "pong"
)

func main() {
	client("wss://127.0.0.1:6668/websocket")
}

func client(strUrl string) (err error) {
	c := wss.NewClient()
	if err = c.Connect(strUrl); err != nil {
		log.Errorf(err.Error())
		return
	}

	for {
		var data []byte
		var from string

		if _, err := c.Send([]byte(WEBSOCKET_DATA_PING)); err != nil {
			log.Errorf(err.Error())
			break
		}

		if data, from, err = c.Recv(-1); err != nil {
			log.Errorf(err.Error())
			break
		}
		log.Infof("web socket client received data [%s] length [%v] from [%v]", data, len(data), from)
		time.Sleep(3 * time.Second)
	}
	return
}
