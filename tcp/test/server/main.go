package main

import (
	"net"

	log "github.com/civet148/gotools/log"
	tcp "github.com/civet148/gotools/tcp"
)

//全局变量

func main() {
	var addr = string(":13343")

	tcpsvr := tcp.CreateTcpServer(OnAccpet, OnClose, OnRead)

	//启动服务
	if !tcpsvr.Listen(addr) {

		log.Info("Listen [%s] failed", addr)
	} else {

		log.Info("Listen [%s] ok", addr)
	}

	tcpsvr.Loop()
	log.Fatal("Program exitting...")
}

//收到连接通知回调
func OnAccpet(ts *tcp.TcpServer, conn net.Conn) {
	log.Info("[%s] on accept", conn.RemoteAddr())
}

//消息接收完毕回调
func OnRead(ts *tcp.TcpServer, conn net.Conn, data []byte, ChanID uint16) {

	log.Info("[%s] ChanID [%d] nDataLen [%d] Data [%s] ", conn.RemoteAddr(), ChanID, len(data), data)
	ts.Send(conn, []byte("PONG"), ChanID)
}

//连接断开回调
func OnClose(ts *tcp.TcpServer, conn net.Conn) {
	log.Info("[%s] on close", conn.RemoteAddr())
}
