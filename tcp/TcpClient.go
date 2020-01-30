package tcp

import (
	"errors"
	"net"
	"time"

	log "github.com/civet148/gotools/log"
)

// TcpClient TCP客户端结构体定义
type TcpClient struct {
	conn *net.TCPConn
}

// Connect : 连接函数，目标服务器地址类型"192.168.1.105:9981"
//args... 第1个参数是bool类型，是否将地址解析为TCP6格式(默认false)
func (tc *TcpClient) Connect(addr string, args ...interface{}) error {
	var tcp6 bool
	var tcpver = string("tcp4")

	if len(args) == 1 {
		tcp6 = args[0].(bool)
		log.Warn("Connecting to a TCP6 address [%s]", addr)
	}
	if tcp6 {
		tcpver = "tcp6"
	}
	tcpAddr, err := net.ResolveTCPAddr(tcpver, addr)
	if err != nil {
		log.Error("Resolve TCP address [%s] failed, error [%s]", addr, err)

		return err
	}
	tc.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Error("DialTCP to [%s] failed", addr)
		return err
	}

	return nil
}

// Close : 关闭连接
func (tc *TcpClient) Close() {

	if tc.conn != nil {
		tc.conn.Close()
	}
}

// Recv : 接收数据
//return Data, Type, DataLen and err
func (tc *TcpClient) Recv() (Data []byte, ChanID uint16, DataLen int, err error) {

	var hdr DataHdr
	hdrbuf, err := ReadDataHdr(tc.conn)
	if err != nil {
		log.Error("%s", err)
		return nil, 0, 0, err
	}
	//log.Debug("Got header ok")
	//解码包头数据
	hdr, err = DecodeDataHdr(hdrbuf)
	if err != nil || hdr.Flag != PackHdrFlag {
		tc.Send([]byte("Socket pack header decode error"), ChanID)
		time.Sleep(50 * time.Millisecond)
		return nil, 0, -1, errors.New("Decode header error")
	}
	//log.Debug("User data length [%d]", hdr.DataLen)
	if hdr.DataLen > PackMaxSize {
		log.Error("Data length out of range [%d]", PackMaxSize)
		return nil, 0, -1, errors.New("Data length out of range")
	}

	//分配用户数据接收缓冲区
	buffer, err := ReadData(tc.conn, int(hdr.DataLen))
	if err != nil {
		log.Error("[%s] read error [%s]", tc.conn.RemoteAddr(), err)
		return nil, 0, 0, err
	}
	return buffer, hdr.ChanID, int(hdr.DataLen), nil
}

// Send : TCP 发送数据接口
// Params :	data 用户数据缓冲区
// 			Type 数据包标志(用户自定义用途)
func (tc *TcpClient) Send(data []byte, ChanID uint16) (bool, error) {
	var res bool
	var err error
	if res, err = WriteData(tc.conn, data, ChanID); err != nil {
		log.Error("Send data error [%s]", err)
		return res, err
	}
	return res, nil
}
