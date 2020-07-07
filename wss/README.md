# a socket wrapper for TCP/UDP/WEB socket

# 1. tcp socket example

## 1.1 tcp socket client 

```go
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

```
## 1.2 tcp socket server 

```go
package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/wss"
	_ "github.com/civet148/gotools/wss/tcpsock" //required (register tcp socket instance)
)

const (
	TCP_DATA_PING = "ping"
	TCP_DATA_PONG = "pong"
)

type Server struct {
	service *wss.SocketServer
}

func main() {

	var url = "tcp://0.0.0.0:6666"
	server(url)

	var c = make(chan bool, 1)
	<-c //block main go routine
}

func server(strUrl string) {

	var server Server
	server.service = wss.NewServer(strUrl)
	if err := server.service.Listen(&server); err != nil {
		log.Errorf(err.Error())
		return
	}
}

func (s *Server) OnAccept(c *wss.SocketClient) {
	log.Infof("connection accepted [%v]", c.GetRemoteAddr())
}

func (s *Server) OnReceive(c *wss.SocketClient, data []byte, length int, from string) {
	log.Infof("tcp server received data [%s] length [%v] from [%v]", data, length, from)
	if string(data) == TCP_DATA_PING {
		if _, err := c.Send([]byte(TCP_DATA_PONG)); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *Server) OnClose(c *wss.SocketClient) {
	log.Infof("connection [%v] closed", c.GetRemoteAddr())
}

```

# 2. udp socket example

## 2.1 udp socket client 

```go
package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/wss"
	_ "github.com/civet148/gotools/wss/udpsock" //required (register udp socket instance)
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

```
## 2.2 udp socket server 

```go
package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/wss"
	_ "github.com/civet148/gotools/wss/udpsock" //required (register udp socket instance)
)

const (
	UDP_DATA_PING  = "ping"
	UDP_DATA_PONG  = "pong"
	UDP_SERVER_URL = "udp://0.0.0.0:6665"
)

type Server struct {
	service *wss.SocketServer
}

func main() {

	server(UDP_SERVER_URL)

	var c = make(chan bool, 1)
	<-c //block main go routine
}

func server(strUrl string) {

	var server Server
	server.service = wss.NewServer(strUrl)
	if err := server.service.Listen(&server); err != nil {
		log.Errorf(err.Error())
		return
	}
}

func (s *Server) OnAccept(c *wss.SocketClient) {
	//log.Infof("connection accepted [%v]", c.GetRemoteAddr())
}

func (s *Server) OnReceive(c *wss.SocketClient, data []byte, length int, from string) {
	log.Infof("udp server received data [%s] length [%v] from [%v] ", data, length, from)
	if string(data) == UDP_DATA_PING {
		if _, err := c.Send([]byte(UDP_DATA_PONG), from); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *Server) OnClose(c *wss.SocketClient) {
	//log.Infof("connection [%v] closed", c.GetRemoteAddr())
}

```

# 3. web socket example

## 3.1 web socket client 

```go
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
	client("ws://127.0.0.1:6668/websocket")
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

```
## 3.2 web socket server 

```go
package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/wss"
	_ "github.com/civet148/gotools/wss/websock"
)

const (
	WEBSOCKET_DATA_PING = "ping"
	WEBSOCKET_DATA_PONG = "pong"
)

type Server struct {
	service *wss.SocketServer
}

func main() {
	server("ws://0.0.0.0:6668/websocket")
	var c = make(chan bool, 1)
	<-c //block main go routine
}

func server(strUrl string) {
	var server Server
	server.service = wss.NewServer(strUrl)
	server.service.Listen(&server)
}

func (s *Server) OnAccept(c *wss.SocketClient) {
	log.Infof("connection accepted [%v]", c.GetRemoteAddr())
}

func (s *Server) OnReceive(c *wss.SocketClient, data []byte, length int, from string) {
	log.Infof("web socket server received data [%s] length [%v] from [%v]", data, length, from)
	if string(data) == WEBSOCKET_DATA_PING {
		if _, err := c.Send([]byte(WEBSOCKET_DATA_PONG)); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *Server) OnClose(c *wss.SocketClient) {
	log.Infof("connection [%v] closed", c.GetRemoteAddr())
}

```

# 4. unix socket example

## 4.1 unix socket client 

```go
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

```
## 4.2 unix socket server 

```go
package main

import (
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/wss"
	_ "github.com/civet148/gotools/wss/unixsock" //required (register tcp socket instance)
)

const (
	UNIX_DATA_PING = "ping"
	UNIX_DATA_PONG = "pong"
)

type Server struct {
	service *wss.SocketServer
}

func main() {

	var url = "unix:///tmp/unix.sock"
	server(url)

	var c = make(chan bool, 1)
	<-c //block main go routine
}

func server(strUrl string) {

	var server Server
	server.service = wss.NewServer(strUrl)
	if err := server.service.Listen(&server); err != nil {
		log.Errorf(err.Error())
		return
	}
}

func (s *Server) OnAccept(c *wss.SocketClient) {
	log.Infof("connection accepted [%v]", c.GetRemoteAddr())
}

func (s *Server) OnReceive(c *wss.SocketClient, data []byte, length int, from string) {
	log.Infof("unix server received data [%s] length [%v] from [%v]", data, length, from)
	if string(data) == UNIX_DATA_PING {
		if _, err := c.Send([]byte(UNIX_DATA_PONG)); err != nil {
			log.Errorf(err.Error())
		}
	}
}

func (s *Server) OnClose(c *wss.SocketClient) {
	log.Infof("connection [%v] closed", c.GetRemoteAddr())
}

```

