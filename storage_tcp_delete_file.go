package fastdfs_client_go

import (
	"bytes"
	"net"
)

// storageDeleteHeaderBody 文件删除
type storageDeleteHeaderBody struct {
	header
	// header  额外参数，发送使用
	groupName      string
	remoteFilename string
}

//Send 发送删除文件命令
// @tcpConn tcp连接
func (s *storageDeleteHeaderBody) Send(tcpConn net.Conn) error {
	// 设置删除文件时的 header 参数
	//@group_name：16字节字符串，组名
	//@filename：不定长字符串，文件名
	s.header.pkgLen = int64(len(s.remoteFilename) + 16)
	s.header.cmd = STORAGE_PROTO_CMD_DELETE_FILE
	s.header.status = 0

	if err := s.header.sendHeader(tcpConn); err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	buffer.Write(groupNameConvBytes(s.groupName))
	buffer.WriteString(s.remoteFilename)

	if _, err := tcpConn.Write(buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

// Receive  接受删除命令发送服务端的响应头
// @tcpConn tcp连接
func (s *storageDeleteHeaderBody) Receive(tcpConn net.Conn) error {
	return s.header.receiveHeader(tcpConn)
}
