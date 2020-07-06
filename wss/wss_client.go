package wss

import "fmt"

type SocketClient struct {
	sock   Socket
	addr   string
	closed bool
}

func init() {

}

func NewClient() *SocketClient {
	return &SocketClient{}
}

//IPv4      => 		tcp://127.0.0.1:6666 [tcp4://127.0.0.1:6666]
//WebSocket => 		ws://127.0.0.1:6668 [wss://127.0.0.1:6668]
func (w *SocketClient) Connect(url string) (err error) {
	var s Socket
	if s = createSocket(url); s == nil {
		return fmt.Errorf("create socket by url [%v] failed", url)
	}
	w.sock = s
	return w.sock.Connect()
}

func (w *SocketClient) Send(data []byte, to ...string) (n int, err error) {
	return w.send(w.sock, data, to...)
}

func (w *SocketClient) Recv() (data []byte, from string, err error) {
	return w.recv(w.sock)
}

func (w *SocketClient) GetLocalAddr() (addr string) {
	return w.sock.GetLocalAddr()
}

func (w *SocketClient) GetRemoteAddr() (addr string) {
	return w.sock.GetRemoteAddr()
}

func (w *SocketClient) Closed() bool {
	return w.closed
}

func (w *SocketClient) send(s Socket, data []byte, to ...string) (n int, err error) {
	return s.Send(data, to...)
}

func (w *SocketClient) recv(s Socket) (data []byte, from string, err error) {
	return s.Recv(-1)
}
