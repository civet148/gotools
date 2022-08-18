package tcp

import (
	"github.com/civet148/gotools/tcp/endian"
	"net"

	log "github.com/civet148/gotools/log"
)

const (
	//PackMaxSize 定义单个通信包最大字节数
	PackMaxSize = uint32(6553500)
	//定义包头flag
	PackHdrFlag = uint16(0xAABB)
)

//DataHdr 数据包头定义
type DataHdr struct {
	Flag    uint16
	ChanID  uint16
	DataLen uint32
}

//获取头部数据对齐后的字节数
func getHdrSize() int {

	var hdr DataHdr
	data, _ := EncodeDataHdr(hdr)
	return len(data)
}

//MakeDataHdr 生成一个数据包头结构对象
func MakeDataHdr(ChanID uint16, DataLen uint32) DataHdr {

	hdr := DataHdr{
		Flag:    PackHdrFlag,
		ChanID:  ChanID,
		DataLen: DataLen,
	}
	return hdr
}

// EncodeDataHdr 将数据包头信息序列化为网络数据
// +++++++++++++++++++++++++++++++++++++----------------------------
// | 2B flag | 2B chanid |4B data length |          (n) data         |
// +++++++++++++++++++++++++++++++++++++----------------------------
func EncodeDataHdr(hdr DataHdr) ([]byte, error) {

	return endian.EncodeEndian(false, &hdr)
}

// DecodeDataHdr 将网络数据反序列化为数据包头
// +++++++++++++++++++++++++++++++++++++----------------------------
// | 2B flag | 2B chanid |4B data length |        (n) data           |
// +++++++++++++++++++++++++++++++++++++----------------------------
func DecodeDataHdr(hdrdata []byte) (hdr DataHdr, err error) {

	err = endian.DecodeEndian(false, hdrdata, &hdr)
	return hdr, err
}

// ReadDataHdr 读取包头数据
func ReadDataHdr(conn net.Conn) (data []byte, err error) {

	hdrsize := getHdrSize()
	data, err = ReadData(conn, hdrsize)

	return data, err
}

// ReadData 从连接读取数据
func ReadData(conn net.Conn, datalen int) ([]byte, error) {

	slice := 0
	left := datalen
	buffer := make([]byte, datalen)
	//log.Debug("Reading user data")

	for left > 0 {
		if n, err := conn.Read(buffer[slice:datalen]); err == nil {
			//log.Debug("Read %d bytes from %s", n, conn.RemoteAddr())
			slice += n
			left -= n
		} else {
			log.Error("[%s] read error [%s]", conn.RemoteAddr(), err)
			return nil, err
		}
	}
	return buffer, nil
}

// Send : TCP 发送数据接口
// Params :	data 用户数据缓冲区
// 			ChanID 数据通道标识(用户自定义用途)
func WriteData(conn net.Conn, data []byte, chanid uint16) (bool, error) {

	datalen := len(data)
	left := datalen
	hdr := MakeDataHdr(chanid, uint32(datalen))
	buf, _ := EncodeDataHdr(hdr) //序列化包头数据
	n, err := conn.Write(buf)    //发送包头数据
	if err != nil {
		return false, err
	}
	n = 0
	slice := 0
	for left > 0 {
		n, err = conn.Write(data[slice:]) //发送用户数据
		if err != nil {
			return false, err
		}
		left -= n
		slice += n
	}

	return true, nil
}

