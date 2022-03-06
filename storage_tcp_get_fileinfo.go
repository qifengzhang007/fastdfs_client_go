package fastdfs_client_go

import (
	"bytes"
	"errors"
	"net"
)

type storageGetFileInfoHeaderBody struct {
	header
	// header  额外参数，发送使用
	groupName      string
	remoteFilename string
	// 响应信息
	fileSize        int64
	createTimestamp int64
	crc32           int64
	SourceIpAddr    string
}

// Send 获取文件信息 body 参数
// @tcpConn  tcp连接
func (s *storageGetFileInfoHeaderBody) Send(tcpConn net.Conn) error {
	// 构建 header 头参数
	// 16 = body 下载参数的前3个参数二进制总长度
	// @group name：16字节字符串，组名
	// @remoteFilename：不定长字符串，文件名
	s.header.pkgLen = int64(len(s.remoteFilename) + 16)
	s.header.cmd = STORAGE_PROTO_CMD_QUERY_FILE_INFO
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

// Receive 通过tcp连接接受数据
// @tcpConn  tcp连接
// 查询文件信息，服务端响body，总计 40 字节长度（固定值）
// @file_size  文件大小，8字节
// @create_timestamp  创建时间戳，8字节
// @crc32  文件内容CRC32校验码，8字节
// @source_ip_addr  16字节字符串，源storage server IP地址
func (s *storageGetFileInfoHeaderBody) Receive(tcpConn net.Conn) error {
	if err := s.header.receiveHeader(tcpConn); err != nil {
		return errors.New(ERROR_STORAGE_SERVER_GET_FILEINFO + err.Error())
	}
	if s.header.pkgLen != 40 {
		return errors.New(ERROR_STORAGE_SERVER_GET_FILEINFO_BODY_LEN)
	}
	buf := make([]byte, 40)
	if _, err := tcpConn.Read(buf); err != nil {
		return err
	}
	s.fileSize = bytesToInt(getBytesByPosition(buf, 0, 8))
	s.createTimestamp = bytesToInt(getBytesByPosition(buf, 8, 8))
	s.crc32 = bytesToInt(getBytesByPosition(buf, 16, 8))
	s.SourceIpAddr = string(getBytesByPosition(buf, 24, 16))
	return nil
}
