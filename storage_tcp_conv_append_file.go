package fastdfs_client_go

import (
	"bytes"
	"errors"
	"net"
)

type storageConvAppendFileHeaderBody struct {
	header
	// header  额外参数，发送使用
	appenderFilename string //不定长字符串
	// 响应信息
	groupName string //16字节字符串，组名
	filename  string //不定长字符串
}

// Send 获取文件信息 body 参数
// @tcpConn  tcp连接
func (s *storageConvAppendFileHeaderBody) Send(tcpConn net.Conn) error {
	// 构建 header 头参数
	// @appender_filename：不定长字符串，appender文件名
	s.header.pkgLen = int64(len(s.appenderFilename))
	s.header.cmd = STORAGE_PROTO_CMD_REGENERATE_APPENDER_FILENAME
	s.header.status = 0

	if err := s.header.sendHeader(tcpConn); err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	buffer.WriteString(s.appenderFilename) //  不定长文本写入方式
	if _, err := tcpConn.Write(buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

// Receive 通过tcp连接接受数据
// @tcpConn  tcp连接
// 查询文件信息，服务端响body，总计 40 字节长度（固定值）
// @group name：16字节字符串，组名
// @filename：不定长字符串，重新生成的文件名（普通类型）
func (s *storageConvAppendFileHeaderBody) Receive(tcpConn net.Conn) error {
	if err := s.header.receiveHeader(tcpConn); err != nil {
		return errors.New(ERROR_STORAGE_SERVER_REGENERATE_APPENDER_FILENAME + err.Error())
	}
	//if int(s.header.status) != 0 {
	//	return errors.New(ERROR_STORAGE_SERVER_REGENERATE_APPENDER_FILENAME + "响应状态status:" + strconv.Itoa(int(s.header.status)))
	//}
	buf := make([]byte, s.header.pkgLen)
	if _, err := tcpConn.Read(buf); err != nil {
		return err
	}
	s.groupName = string(getBytesByPosition(buf, 0, 16))
	s.filename = string(getBytesByPosition(buf, 16, int(s.header.pkgLen-16)))
	return nil
}
