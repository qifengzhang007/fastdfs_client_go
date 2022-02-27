package fastdfs_client_go

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

// 本页面代码功能介绍：
// header 头是所有tcp连接通讯必须组合使用的公共结构体
// 如果 body 参数为空，那么可以直接使用 header 结构体以及提供的方法进行网络通讯

// header(包头) 组成结构：
//  pkg_len：8字节整数，body长度，不包含header，只是body的长度
//  cmd：   1字节整数，命令码
//  status：1字节整数，状态码，0表示成功，非0失败（UNIX错误码）
type header struct {
	pkgLen int64 // 占8字节
	cmd    byte  // 占1字节
	status byte  // 占1字节
}

//sendHeader  发送header消息
// @tcpConn tcp连接
func (h *header) sendHeader(tcpConn net.Conn) error {
	buffer := new(bytes.Buffer)
	// c.pkgLen 整数型（4字节或者8字节）写入到二进制缓冲区
	//将整数类型采用网络字节序（Big-Endian），包括4字节整数(int32)和8字节整数(int64)
	if err := binary.Write(buffer, binary.BigEndian, h.pkgLen); err != nil {
		return err
	}
	buffer.WriteByte(h.cmd)
	buffer.WriteByte(h.status)

	if _, err := tcpConn.Write(buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

//receiveHeader  接受 header 头
// @tcpConn tcp连接
func (h *header) receiveHeader(tcpConn net.Conn) error {
	buf := make([]byte, TCP_HEADER_LEN)
	if _, err := tcpConn.Read(buf); err != nil {
		return err
	}
	buffer := bytes.NewBuffer(buf)
	// 读取已接受字节的实际长度，赋值给 pkgLen，pkgLen 按照协议长度必须=10
	if err := binary.Read(buffer, binary.BigEndian, &h.pkgLen); err != nil || h.pkgLen != TCP_HEADER_LEN {
		return err
	} else {
		h.cmd = (buf[8:9])[0]
		h.status = (buf[9:10])[0]
		if h.status != 0 {
			return errors.New(ERROR_HEADER_RECEV_STATUS_NOT_ZERO + ", ErrorMsg: " + string(h.status))
		}
	}
	return nil
}
