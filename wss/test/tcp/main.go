package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/wss"
	_ "github.com/civet148/gotools/wss/tcpsock" //required (register tcp socket instance)
	"time"
)

const (
	TCP_DATA_PING = "ping"
	TCP_DATA_PONG = "pong"
)

type TcpServer struct {
	service *wss.SocketServer
}

func main() {
	var c = make(chan bool, 1)
	var url = "tcp://127.0.0.1:6666"
	server(url)
	time.Sleep(1 * time.Second)
	client(url)
	<-c //block main go routine
}

func server(strUrl string) {

	var tcpSvr TcpServer
	tcpSvr.service = wss.NewServer(strUrl)
	if err := tcpSvr.service.Listen(&tcpSvr); err != nil {
		log.Errorf(err.Error())
		return
	}
}

func client(strUrl string) {
	c := wss.NewClient()
	if err := c.Connect(strUrl); err != nil {
		log.Errorf(err.Error())
		return
	}
	go func() {
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
	}()
}

func (s *TcpServer) OnAccept(c *wss.SocketClient) {
	log.Infof("connection accepted [%v]", c.GetRemoteAddr())
}

func (s *TcpServer) OnReceive(c *wss.SocketClient, data []byte, length int, from string) {
	log.Infof("tcp server received data [%s] length [%v] from [%v]", data, length, from)
	if string(data) == TCP_DATA_PING {
		if _, err := c.Send([]byte(TCP_DATA_PONG)); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *TcpServer) OnClose(c *wss.SocketClient) {
	log.Infof("connection [%v] closed", c.GetRemoteAddr())
}
