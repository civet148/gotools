package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/wss"
	_ "github.com/civet148/gotools/wss/unixsock" //required (register tcp socket instance)
	"time"
)

const (
	UNIX_DATA_PING = "ping"
	UNIX_DATA_PONG = "pong"
)

func main() {
	var url = "unix:///tmp/unix.sock"
	client(url)
}

func client(strUrl string) {
	c := wss.NewClient()
	if err := c.Connect(strUrl); err != nil {
		log.Errorf(err.Error())
		return
	}

	for {
		if _, err := c.Send([]byte(UNIX_DATA_PING)); err != nil {
			log.Errorf(err.Error())
			break
		}

		if data, from, err := c.Recv(len(UNIX_DATA_PONG)); err != nil {
			log.Error(err.Error())
			break
		} else {
			log.Infof("unix client received data [%s] length [%v] from [%v]", string(data), len(data), from)
		}

		time.Sleep(3 * time.Second)
	}
}
