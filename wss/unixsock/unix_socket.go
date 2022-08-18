package unixsock

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/parser"
	"github.com/civet148/gotools/wss"
	"net"
	"os"
	"strings"
)

type socket struct {
	ui       *parser.UrlInfo
	conn     net.Conn
	listener *net.UnixListener
	closed   bool
}

func init() {
	_ = wss.Register(wss.SocketType_UNIX, NewSocket)
}

func NewSocket(ui *parser.UrlInfo) wss.Socket {

	return &socket{
		ui: ui,
	}
}

func (s *socket) Listen() (err error) {
	var network = s.getNetwork()
	addr := s.getUnixSockFile()
	if err = os.Remove(addr); err != nil {
		log.Errorf("remove file error [%v]", err.Error())
		return
	}
	var unixAddr *net.UnixAddr
	unixAddr, err = net.ResolveUnixAddr(network, s.ui.GetPath())
	if err != nil {
		err = fmt.Errorf("Cannot resolve unix addr: " + err.Error())
		log.Errorf(err.Error())
		return
	}
	log.Debugf("trying listen [%v] protocol [%v]", addr, s.ui.GetScheme())
	if s.listener, err = net.ListenUnix("unix", unixAddr); err != nil {
		log.Errorf("listen tcp address [%s] failed", addr)
		return
	}
	return
}

func (s *socket) Accept() wss.Socket {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil
	}
	return &socket{
		conn: conn,
		ui:   s.ui,
	}
}

func (s *socket) Connect() (err error) {
	var network = s.getNetwork()
	addr := s.getUnixSockFile()
	var unixAddr *net.UnixAddr
	unixAddr, err = net.ResolveUnixAddr(network, addr)
	if err != nil {
		log.Errorf("resolve tcp address [%s] failed, error [%s]", addr, err)
		return err
	}

	s.conn, err = net.DialUnix(network, nil, unixAddr)
	if err != nil {
		log.Errorf("dial tcp to [%s] failed", addr)
		return err
	}
	return
}

func (s *socket) Send(data []byte, to ...string) (n int, err error) {
	return s.conn.Write(data)
}

// length <= 0, default PACK_FRAGMENT_MAX=1500 bytes
func (s *socket) Recv(length int) (data []byte, from string, err error) {

	var once bool
	var recv, left int
	if length <= 0 {
		once = true
		length = wss.PACK_FRAGMENT_MAX
	}
	left = length
	data = s.makeBuffer(length)
	var n int
	if once {
		if n, err = s.conn.Read(data); err != nil {
			log.Errorf("read data error [%v]", err.Error())
			return
		}
		recv = n
	} else {

		for left > 0 {
			if n, err = s.conn.Read(data[recv:]); err != nil {
				log.Errorf("read data error [%v]", err.Error())
				return
			}
			left -= n
			recv += n
		}
	}

	if recv < length {
		data = data[:recv]
	}
	from = s.GetLocalAddr()
	return
}

func (s *socket) Close() (err error) {
	if s.closed {
		err = fmt.Errorf("socket already closed")
		return
	}
	if s.conn == nil {
		err = fmt.Errorf("socket is nil")
		log.Error(err.Error())
		return
	}
	s.closed = true
	return s.conn.Close()
}

func (s *socket) GetLocalAddr() (strAddr string) {
	return s.getUnixSockFile()
}

func (s *socket) GetRemoteAddr() (strAddr string) {
	return s.getUnixSockFile()
}

func (s *socket) GetSocketType() wss.SocketType {
	return wss.SocketType_UNIX
}

func (s *socket) getUnixSockFile() (strSockFile string) {

	if s.ui == nil {
		return
	}
	strSockFile = s.ui.GetPath()
	if !strings.HasSuffix(strSockFile, "sock") {
		panic("unix socket must .sock as file suffix")
	}
	return
}

func (s *socket) getNetwork() string {
	return wss.NETWORK_UNIX
}

func (s *socket) makeBuffer(length int) []byte {
	return make([]byte, length)
}
