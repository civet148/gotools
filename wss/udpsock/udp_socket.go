package udpsock

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/parser"
	"github.com/civet148/gotools/wss"
)

type socket struct {
	ui *parser.UrlInfo
}

func init() {
	_ = wss.Register(wss.SocketType_UDP, NewSocket)
}

func NewSocket(ui *parser.UrlInfo) wss.Socket {

	return &socket{
		ui: ui,
	}
}

func (s *socket) Listen() (err error) {
	return
}

func (s *socket) Accept() wss.Socket {
	log.Warnf("accept method not for UDP socket")
	return nil
}

func (s *socket) Connect() (err error) {
	return fmt.Errorf("only for TCP/WEB socket")
}

func (s *socket) Send(data []byte, to ...string) (n int, err error) {
	return
}

func (s *socket) Recv(length int) (data []byte, from string, err error) {
	return
}

func (s *socket) Close() (err error) {

	return
}

func (s *socket) GetLocalAddr() string {
	return "not yet"
}

func (s *socket) GetRemoteAddr() string {
	return "not yet"
}
