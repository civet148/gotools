package main

import (
	"github.com/civet148/gotools/tcp"
	"time"

	log "github.com/civet148/gotools/log"
)

var LocalAddr = "127.0.0.1:13343"

func main() {
	var tcpClient tcp.TcpClient
	if err := tcpClient.Connect(LocalAddr); err != nil {
		log.Error("%s", err)
		return
	}
	for {
		tcpClient.Send([]byte("PING"), 1)
		Data, ChanID, nDataLen, err := tcpClient.Recv()
		if err != nil {
			log.Error("%s", err)
			break
		}
		log.Info("ChanID [%v] nDataLen [%v] Data [%s]", ChanID, nDataLen, Data)
		time.Sleep(5 * time.Second)
	}
}
