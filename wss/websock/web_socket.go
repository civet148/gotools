package websock

import (
	"crypto/tls"
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/gotools/parser"
	"github.com/civet148/gotools/wss"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type socket struct {
	ui        *parser.UrlInfo
	conn      *websocket.Conn
	accepting chan *websocket.Conn
	closed    bool
}

func init() {
	_ = wss.Register(wss.SocketType_WEB, NewSocket)
}

func NewSocket(ui *parser.UrlInfo) wss.Socket {
	return &socket{
		ui:        ui,
		accepting: make(chan *websocket.Conn, 1000),
	}
}

func (s *socket) Listen() (err error) {
	engine := gin.Default()
	if s.ui.GetPath() == "" {
		s.ui.Path = "/"
	}

	engine.GET(s.ui.Path, s.webSocketRegister)
	strCertFile := s.ui.Queries[wss.WSS_TLS_CERT]
	strKeyFile := s.ui.Queries[wss.WSS_TLS_KEY]

	go func() {
		if s.ui.Scheme == wss.URL_SCHEME_WSS {
			log.Debugf("listen GET [%s://%s%s] -> cert [%s] key [%s]", s.ui.Scheme, s.ui.Host, s.ui.Path, strCertFile, strKeyFile)
			err = engine.RunTLS(s.ui.Host, strCertFile, strKeyFile)
		} else {
			log.Debugf("listen GET [%s://%s%s]", s.ui.Scheme, s.ui.Host, s.ui.Path)
			err = engine.Run(s.ui.Host)
		}

		if err != nil {
			s.closed = true
			log.Errorf("listen websocket closing with error [%v]", err.Error())
			return
		}
	}()

	return
}

func (s *socket) Accept() wss.Socket {

	var c *websocket.Conn
	select {
	case c = <-s.accepting:
		{
			log.Debugf("accepted client [%v]", c.RemoteAddr().String())
			c.SetCloseHandler(s.webSocketCloseHandler)
			c.SetPingHandler(s.websocketPingHandler)
			c.SetPongHandler(s.websocketPongHandler)
			return &socket{
				conn: c,
				ui:   s.ui,
			}
		}
	}
	return nil
}

func (s *socket) Connect() (err error) {
	url := fmt.Sprintf("%v://%v%v", s.ui.Scheme, s.ui.Host, s.ui.Path)
	log.Debugf("connect to url [%v]", url)
	dialer := &websocket.Dialer{}
	if s.ui.Scheme == wss.URL_SCHEME_WSS {
		dialer.TLSClientConfig = &tls.Config{RootCAs: nil, InsecureSkipVerify: true}
	}
	if s.conn, _, err = dialer.Dial(url, nil); err != nil {
		log.Errorf(err.Error())
		return
	}
	return
}

func (s *socket) Send(data []byte, to ...string) (n int, err error) {
	if s.conn == nil {
		err = fmt.Errorf("web socket connection is nil")
		return
	}

	if err = s.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return
	}
	n = len(data)
	log.Debugf("data [%v] length [%v]", string(data), n)
	return
}

func (s *socket) Recv(length int) (data []byte, from string, err error) {
	if s.conn == nil {
		err = fmt.Errorf("web socket connection is nil")
		return
	}

	var msgType int
	if msgType, data, err = s.conn.ReadMessage(); err != nil {
		log.Errorf(err.Error())
		return
	}
	s.debugMessageType(msgType)
	from = s.conn.RemoteAddr().String()
	return
}

func (s *socket) Close() (err error) {

	if s.closed {
		err = fmt.Errorf("socket already closed")
		return
	}
	s.closed = true
	if s.conn == nil {
		err = fmt.Errorf("socket is nil")
		log.Error(err.Error())
		return
	}
	s.closed = true
	return s.conn.Close()
}

func (s *socket) GetLocalAddr() (addr string) {
	if s.conn == nil {
		return s.ui.Host //web socket server connection is nil
	}
	addr = s.conn.LocalAddr().String()
	return
}

func (s *socket) GetRemoteAddr() (addr string) {
	if s.conn == nil {
		return //web socket client connection can't be nil
	}
	addr = s.conn.RemoteAddr().String()
	return
}

func (s *socket) GetSocketType() wss.SocketType {
	return wss.SocketType_WEB
}

func (s *socket) debugMessageType(msgType int) {

	switch msgType {
	case websocket.TextMessage:
		log.Debugf("message type [TextMessage]")
	case websocket.BinaryMessage:
		log.Debugf("message type [BinaryMessage]")
	case websocket.CloseMessage:
		log.Debugf("message type [CloseMessage]")
	case websocket.PingMessage:
		log.Debugf("message type [PingMessage]")
	case websocket.PongMessage:
		log.Debugf("message type [PongMessage]")
	}
}

func (s *socket) webSocketRegister(ctx *gin.Context) {
	var err error
	log.Debugf("request ctx [%v]", ctx)
	upGrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}
	var c *websocket.Conn
	if c, err = upGrader.Upgrade(ctx.Writer, ctx.Request, nil); err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Debugf("client [%v] registered", c.RemoteAddr().String())
	s.accepting <- c
}

func (s *socket) webSocketCloseHandler(code int, text string) (err error) {
	log.Debugf("close code [%v] text [%v]", code, text)
	return
}

func (s *socket) websocketPingHandler(appData string) (err error) {
	log.Debugf("ping app data [%v]", appData)
	return
}

func (s *socket) websocketPongHandler(appData string) (err error) {
	log.Debugf("pong app data [%v]", appData)
	return
}
