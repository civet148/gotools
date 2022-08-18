package tcp

import (
	"net"

	"time"

	log "github.com/civet148/gotools/log"
)

type AcceptCbFunc func(ts *TcpServer, conn net.Conn)
type CloseCbFunc func(ts *TcpServer, conn net.Conn)
type ReadCbFunc func(ts *TcpServer, conn net.Conn, data []byte, channel uint16)

// TcpServer TCP服务对象
type TcpServer struct {
	listener  net.Listener  //TCP监听对象
	accepting chan net.Conn //客户端连接通知
	quiting   chan net.Conn //客户端退出通知
	running   bool          //服务器是否继续运行
	OnAccept  AcceptCbFunc  //收到连接通知回调
	OnRead    ReadCbFunc    //消息接收完毕回调
	OnClose   CloseCbFunc   //连接断开回调
}

// CreateTcpServer 创建TCP服务
// return ts TCP服务对象
func CreateTcpServer(cbAccept AcceptCbFunc, cbClose CloseCbFunc, cbRead ReadCbFunc) (ts *TcpServer) {

	ts = &TcpServer{
		accepting: make(chan net.Conn),
		quiting:   make(chan net.Conn),
		running:   true,
		OnAccept:  cbAccept,
		OnRead:    cbRead,
		OnClose:   cbClose,
	}

	return ts
}

// Loop 阻塞线程死循环
func (ts *TcpServer) Loop() {

	for {
		if !ts.running {
			break
		}
		time.Sleep(1000 * time.Second)
	}
}

// Listen 监听并接收客户端连接
// addr参数  1、非加密模式服务监听URL "tcp://192.168.1.105:8899"监听指定IP端口或"tcp://:8899"监听本机所有IP端口
//           2、加密模式服务监听URL "tls://192.168.1.105:8899" 监听指定IP端口或"tls://:8899"监听本机所有IP端口
// return conn接收的连接，err 错误信息
func (ts *TcpServer) Listen(addr string) bool {
	var err error
	ts.listener, err = net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Listen TCP address failed, %s", addr)
		return false
	}

	//启动goroutine用于监听客户端连接相关channel事件(accepting/quiting/incoming)
	go func() {
		log.Info("Start goroutine for channel event accepting/quiting")
		for {
			if !ts.running {
				break //服务停止
			}
			select {
			case conn := <-ts.accepting: //客户端连接成功事件(此事件用于内部会话管理)
				ts.Join(conn)
			case conn := <-ts.quiting: //客户端连接断开事件(此事件用于内部会话管理)
				ts.Close(conn)
 			//default: //开启default处理时必须sleep几毫秒，否则会造成CPU使用率过高

			}
		}
	}()

	//启动goroutine开始接收客户端连接
	go func() {
		log.Info("Start goroutine for accept new connection")
		for {
			if !ts.running {
				break //服务停止
			}
			conn, err := ts.listener.Accept()
			if err != nil {
				log.Fatal("Accept failed, %s", err)
				return
			}
			//将连接客户端的conn通知到channel，这里不做额外的处理
			ts.accepting <- conn
		}
	}()

	return true
}

//Stop 强制停止TCP服务
func (ts *TcpServer) Stop() {
	ts.running = false
}

//KickConn 强制断开客户端连接
func (ts *TcpServer) KickConn(conn net.Conn) {
	ts.quiting <- conn
}

//ReadConn 从客户端读取数据
func (ts *TcpServer) ReadConn(conn net.Conn) {

	var err error
	var hdr DataHdr
	var hdrbuf []byte
	//log.Debug("Reading data from connection")
	for {
		//接收包头数据
		hdrbuf, err = ReadDataHdr(conn)
		if err != nil {
			log.Error("%s", err)
			ts.quiting <- conn
			return
		}
		//log.Debug("Got header ok")

		//解码包头数据
		hdr, err = DecodeDataHdr(hdrbuf)
		if err != nil || hdr.Flag != PackHdrFlag {
			log.Error("Decode failed or header flag invalid flag [0x%X] err [%s] ", hdr.Flag, err)
			ts.Send(conn, []byte("Socket pack header decode error! "), 0)
			time.Sleep(500 * time.Millisecond)
			ts.quiting <- conn
			return
		}

		if hdr.DataLen > PackMaxSize {
			log.Error("Data length out of range [%d]", PackMaxSize)
			ts.quiting <- conn
			return
		}

		//分配用户数据接收缓冲区
		buffer, err := ReadData(conn, int(hdr.DataLen))
		if err != nil {
			ts.quiting <- conn
			log.Error("[%s] read error [%s]", conn.RemoteAddr(), err)
			return
		}

		ts.OnRead(ts, conn, buffer[:hdr.DataLen], hdr.ChanID)
	}
}

//Join 新连接加入管理map
func (ts *TcpServer) Join(conn net.Conn) {

	//通知客户端连接成功事件
	ts.OnAccept(ts, conn)

	//启动一个goroutine读取客户端socket数据
	go ts.ReadConn(conn)

}

//Close 客户端连接关闭
func (ts *TcpServer) Close(conn net.Conn) {

	ts.OnClose(ts, conn)
	conn.Close()
}

//Send 向指定客户端连接发送数据
func (ts *TcpServer) Send(conn net.Conn, data []byte, ChanID uint16) (bool, error) {
	var res bool
	var err error
	if res, err = WriteData(conn, data, ChanID); err != nil {
		log.Error("Send data error [%s]", err)
		ts.quiting <- conn
		return res, err
	}
	return res, nil
}
