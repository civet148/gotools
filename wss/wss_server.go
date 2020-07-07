package wss

import (
	"fmt"
	"github.com/civet148/gotools/log"
	"sync"
)

type SocketServer struct {
	url       string                   //web socket url
	sock      Socket                   //server socket
	handler   SocketHandler            //server callback handler
	accepting chan Socket              //client connection accepted
	receiving chan Socket              //client message received
	quiting   chan Socket              //client connection closed
	clients   map[Socket]*SocketClient //socket clients
	locker    *sync.Mutex              //locker mutex
	running   bool                     //service running ?
}

func init() {

}

func NewServer(url string) *SocketServer {

	var s Socket
	s = createSocket(url)

	return &SocketServer{
		url:       url,
		locker:    &sync.Mutex{},
		sock:      s,
		accepting: make(chan Socket, 1000),
		quiting:   make(chan Socket, 1000),
		clients:   make(map[Socket]*SocketClient, 0),
		running:   true,
	}
}

//TCP       => 		tcp://127.0.0.1:6666
//UDP       => 		udp://127.0.0.1:6667
//WebSocket => 		ws://127.0.0.1:6668 wss://127.0.0.1:6668
func (w *SocketServer) Listen(handler SocketHandler) (err error) {
	w.handler = handler
	if err = w.sock.Listen(); err != nil {
		log.Errorf(err.Error())
		return
	}
	log.Debugf("listen protocol [%v] address [%v] ok", w.sock.GetSocketType(), w.sock.GetLocalAddr())
	if w.sock.GetSocketType() != SocketType_UDP {
		go func() {
			log.Debugf("start goroutine for channel event accepting/quiting")
			for {
				if !w.running {
					log.Debugf("accepting/quiting channel service loop break")
					break //service loop break
				}
				select {
				case s := <-w.accepting: //client connection coming...
					w.onAccept(s)
				case s := <-w.quiting: //client connection closed
					w.onClose(s)
					//default: //disable default because of high CPU performance
				}
			}
		}()

		//new go routine for accept new connections
		go func() {
			log.Debugf("start goroutine for accept new connection")
			for {
				if !w.running {
					log.Debugf("accept service loop break")
					break //service loop break
				}
				s := w.sock.Accept()
				if err != nil {
					log.Fatal("accept failed error [%v]", err.Error())
					return
				}
				//socket quiting...
				w.accepting <- s
			}
		}()
	} else {
		w.onAccept(w.sock)
	}
	return
}

func (w *SocketServer) Close(client *SocketClient) (err error) {
	return w.closeSocket(client.sock)
}

func (w *SocketServer) Send(client *SocketClient, data []byte, to ...string) (n int, err error) {
	return w.sendSocket(client.sock, data, to...)
}

func (w *SocketServer) GetClientCount() int {
	return w.getClientCount()
}

func (w *SocketServer) GetClientAll() (clients []*SocketClient) {
	return w.getClientAll()
}

func (w *SocketServer) closeSocket(s Socket) (err error) {
	if s == nil {
		return fmt.Errorf("close socket is nil")
	}
	w.removeClient(s)
	return s.Close()
}

func (w *SocketServer) sendSocket(s Socket, data []byte, to ...string) (n int, err error) {
	if s == nil || len(data) == 0 {
		err = fmt.Errorf("send socket is nil or data length is 0")
		return
	}
	return s.Send(data, to...)
}

func (w *SocketServer) recvSocket(s Socket) (data []byte, from string, err error) {
	if s == nil {
		err = fmt.Errorf("send socket is nil")
		return
	}
	return s.Recv(-1)
}

func (w *SocketServer) onAccept(s Socket) {
	c := w.addClient(s)
	go w.readSocket(s)
	w.handler.OnAccept(c)
}

func (w *SocketServer) onClose(s Socket) {
	w.handler.OnClose(w.removeClient(s))
}

func (w *SocketServer) onReceive(s Socket, data []byte, length int, from string) {
	c := w.getClient(s)
	w.handler.OnReceive(c, data, length, from)
}

func (w *SocketServer) readSocket(s Socket) {
	for {
		if data, from, err := w.recvSocket(s); err != nil {
			w.quiting <- s
			break
		} else if len(data) > 0 {
			w.onReceive(s, data, len(data), from)
		}
	}
}

func (w *SocketServer) lock() {
	w.locker.Lock()
}

func (w *SocketServer) unlock() {
	w.locker.Unlock()
}

func (w *SocketServer) closeClientAll() {
	w.lock()
	defer w.unlock()
	for s, _ := range w.clients {
		w.onClose(s)
		delete(w.clients, s)
	}
}

func (w *SocketServer) addClient(s Socket) (client *SocketClient) {
	client = &SocketClient{
		sock: s,
	}
	w.lock()
	defer w.unlock()
	w.clients[client.sock] = client
	return client
}

func (w *SocketServer) removeClient(s Socket) (client *SocketClient) {
	w.lock()
	defer w.unlock()
	client = w.clients[s]
	delete(w.clients, s)
	return
}

func (w *SocketServer) getClient(s Socket) (client *SocketClient) {
	var ok bool
	w.lock()
	defer w.unlock()
	if client, ok = w.clients[s]; ok {
		return
	}
	return
}

func (w *SocketServer) getClientCount() int {
	w.lock()
	defer w.unlock()
	return len(w.clients)
}

func (w *SocketServer) getClientAll() (clients []*SocketClient) {
	w.lock()
	defer w.unlock()
	for _, v := range w.clients {
		clients = append(clients, v)
	}
	return
}
