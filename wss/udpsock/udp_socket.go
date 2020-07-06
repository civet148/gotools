package udpsock

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/parser"
	"github.com/civet148/gotools/wss"
	"net"
	"strings"
)

const (
	UDPv4 = "udp4"
	UDPv6 = "udp6"
)

type socket struct {
	ui   *parser.UrlInfo
	conn *net.UDPConn
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
	var networkVer = s.getVer()

	if udpAddr, err = net.ResolveUDPAddr(networkVer, strAddr); err != nil {
		log.Errorf("resolve UDP addr [%v] error [%v]", strAddr, err.Error())
		return
	}

	if s.conn, err = net.ListenUDP(networkVer, udpAddr); err != nil {
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
	var networkVer = s.getVer()

	if len(to) == 0 {
		return 0, fmt.Errorf("UDP send method to parameter required")
	}

	strToAddr := to[0]
	strToAddr = strings.Replace(strToAddr, "://", "", -1)
	if udpAddr, err = net.ResolveUDPAddr(networkVer, strToAddr); err != nil {
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
	if s.conn == nil {
		return fmt.Errorf("socket is nil")
	}
	return s.conn.Close()
}

func (s *socket) GetLocalAddr() string {
	return s.conn.LocalAddr().String()
}

func (s *socket) GetRemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

func (s *socket) GetSocketType() wss.SocketType {
	return wss.SocketType_UDP
}

func (s *socket) getVer() (ver string) {
	if s.isUDP6() {
		ver = UDPv6
	}
	return UDPv4
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
