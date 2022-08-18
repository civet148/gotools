package udpsock

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/parser"
	"github.com/civet148/gotools/wss"
	"net"
	"strings"
)

type socket struct {
	ui     *parser.UrlInfo
	conn   *net.UDPConn
	closed bool
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

	var strAddr = s.ui.GetHost()
	var udpAddr *net.UDPAddr
	var network = s.getNetwork()

	if udpAddr, err = net.ResolveUDPAddr(network, strAddr); err != nil {
		log.Errorf("resolve UDP addr [%v] error [%v]", strAddr, err.Error())
		return
	}

	if s.conn, err = net.ListenUDP(network, udpAddr); err != nil {
		log.Errorf("listen UDP addr [%v] error [%v]", strAddr, err.Error())
		return
	}
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

	var udpAddr *net.UDPAddr
	var network = s.getNetwork()

	if len(to) == 0 {
		return 0, fmt.Errorf("UDP send method to parameter required")
	}

	strToAddr := to[0]
	nSep := len(parser.URL_SCHEME_SEP)
	if strings.Contains(strToAddr, parser.URL_SCHEME_SEP) {
		nIndex := strings.Index(strToAddr, parser.URL_SCHEME_SEP)
		strToAddr = strToAddr[nIndex+nSep:]
	}

	if udpAddr, err = net.ResolveUDPAddr(network, strToAddr); err != nil {
		log.Errorf("resolve UDP addr [%v] error [%v]", strToAddr, err.Error())
		return
	}
	return s.conn.WriteToUDP(data, udpAddr)
}

func (s *socket) Recv(length int) (data []byte, from string, err error) {
	var n int
	var udpAddr *net.UDPAddr
	data = s.makeBuffer(wss.PACK_FRAGMENT_MAX)
	if n, udpAddr, err = s.conn.ReadFromUDP(data); err != nil {
		log.Errorf("read from UDP error [%v]", err.Error())
		return
	}
	return data[:n], udpAddr.String(), nil
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

func (s *socket) GetLocalAddr() string {
	return s.conn.LocalAddr().String()
}

func (s *socket) GetRemoteAddr() (addr string) {
	return
}

func (s *socket) GetSocketType() wss.SocketType {
	return wss.SocketType_UDP
}

func (s *socket) getNetwork() string {
	if s.isUDP6() {
		return wss.NETWORK_UDPv6
	}
	return wss.NETWORK_UDPv4
}

func (s *socket) isUDP6() (ok bool) {
	scheme := s.ui.GetScheme()
	if scheme == wss.URL_SCHEME_UDP6 {
		return true
	}
	return
}

func (s *socket) makeBuffer(length int) []byte {
	return make([]byte, length)
}
