package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/wss"
	_ "github.com/civet148/gotools/wss/tcpsock" //required (register tcp socket instance)
)

type TcpServer struct {
	service *wss.SocketServer
}

func main() {
	var c = make(chan bool, 1)
	var url = "tcp://127.0.0.1:6666"
	server(url)
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

}

func (s *TcpServer) OnAccept(c *wss.SocketClient) {
	log.Infof("connection accepted [%v]", c.GetRemoteAddr())
}

func (s *TcpServer) OnReceive(c *wss.SocketClient, data []byte, length int, from string) {
	log.Infof("data received data [%s] length [%v] from [%v]", data, length, from)
}

func (s *TcpServer) OnClose(c *wss.SocketClient) {
	log.Infof("connection closed [%v]", c.GetRemoteAddr())
}
