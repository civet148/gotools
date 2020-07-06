package wss

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/parser"
	"strings"
)

const (
	URL_SCHEME_TCP  = "tcp"
	URL_SCHEME_TCP4 = "tcp4"
	URL_SCHEME_TCP6 = "tcp6"
	URL_SCHEME_UDP  = "udp"
	URL_SCHEME_UDP4 = "udp4"
	URL_SCHEME_UDP6 = "udp6"
	URL_SCHEME_WS   = "ws"
	URL_SCHEME_WSS  = "wss"
)

const (
	PACK_FRAGMENT_MAX = 1500
)

type SocketHandler interface {
	OnAccept(c *SocketClient)
	OnReceive(c *SocketClient, data []byte, length int, from string)
	OnClose(c *SocketClient)
}

type Socket interface {
	Listen() (err error)                                   // bind and listen on address and port
	Accept() Socket                                        // accept connection...
	Connect() (err error)                                  // for tcp/web socket
	Send(data []byte, to ...string) (n int, err error)     // send to...
	Recv(length int) (data []byte, from string, err error) // receive from... if length > 0, will receive the bytes specified.
	Close() (err error)                                    // close socket
	GetLocalAddr() string                                  // get socket local address
	GetRemoteAddr() string                                 // get socket remote address
	GetSocketType() SocketType                             //get socket type
}

type SocketType int

const (
	SocketType_TCP SocketType = 1
	SocketType_WEB SocketType = 2
	SocketType_UDP SocketType = 3
)

func (s SocketType) GoString() string {
	return s.String()
}

func (s SocketType) String() string {
	switch s {
	case SocketType_TCP:
		return "SocketType_TCP"
	case SocketType_WEB:
		return "SocketType_WEB"
	case SocketType_UDP:
		return "SocketType_UDP"
	}
	return "SocketType_Unknown"
}

type SocketInstance func(ui *parser.UrlInfo) Socket

var instances = make(map[SocketType]SocketInstance)

func Register(sockType SocketType, inst SocketInstance) (err error) {
	if _, ok := instances[sockType]; !ok {

		instances[sockType] = inst
		return
	}
	err = fmt.Errorf("socket type [%v] instance already exists", sockType)
	log.Errorf("%v", err.Error())
	return
}

func newSocket(sockType SocketType, ui *parser.UrlInfo) (s Socket) {
	if inst, ok := instances[sockType]; !ok {
		log.Errorf("socket type [%v] instance not register", sockType)
		return nil
	} else {
		s = inst(ui)
	}
	return
}

func SetLogDebug(enable bool) {
	if enable {
		log.SetLevel("debug")
	} else {
		log.SetLevel("warn")
	}
}

func SetLogFile(strPath string) {
	if strPath != "" {
		log.Open(strPath)
	}
}

func createSocket(url string) (s Socket) {

	url = strings.ToLower(url)
	ui := parser.ParseUrl(url)
	switch ui.Scheme {
	case URL_SCHEME_TCP, URL_SCHEME_TCP4, URL_SCHEME_TCP6:
		s = newSocket(SocketType_TCP, ui)
	case URL_SCHEME_WS, URL_SCHEME_WSS:
		s = newSocket(SocketType_WEB, ui)
	case URL_SCHEME_UDP, URL_SCHEME_UDP4, URL_SCHEME_UDP6:
		s = newSocket(SocketType_UDP, ui)
	default:
		{
			log.Errorf("unknown scheme [%v] in url [%v]", ui.Scheme, url)
			panic("unknown scheme, url illegal")
		}
	}
	return
}
