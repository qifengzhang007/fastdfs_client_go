package fastdfs_client_go

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"os"
)

type storageDownloadHeaderBody struct {
	header
	// header  额外参数，发送使用
	groupName      string
	remoteFilename string
	offset         int64
	downloadBytes  int64
	//保存的文件名
	saveFileName string
}

// Send 下载文件 body 参数
// @tcpConn  tcp连接
func (s *storageDownloadHeaderBody) Send(tcpConn net.Conn) error {
	// 构建 header 头参数
	// 32 = body 下载参数的前3个参数二进制总长度
	// @file_offset 8字节整数，文件偏移量,从指定的位置开始下载
	// @download_bytes：8字节整数，需要下载字节数
	// @group name：16字节字符串，组名
	// @filename：不定长字符串，文件名
	s.header.pkgLen = int64(len(s.remoteFilename) + 32)
	s.header.cmd = STORAGE_PROTO_CMD_DOWNLOAD_FILE
	s.header.status = 0

	if err := s.header.sendHeader(tcpConn); err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err := binary.Write(buffer, binary.BigEndian, s.offset); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, s.downloadBytes); err != nil {
		return err
	}
	buffer.Write(groupNameConvBytes(s.groupName))
	buffer.WriteString(s.remoteFilename)
	if _, err := tcpConn.Write(buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

// Receive 通过tcp连接接受数据
// @tcpConn  tcp连接
func (s *storageDownloadHeaderBody) Receive(tcpConn net.Conn) error {
	if err := s.header.receiveHeader(tcpConn); err != nil {
		return errors.New(ERROR_STORAGE_SERVER_DOWN_HEADER + err.Error())
	}
	if s.header.pkgLen == 0 {
		return errors.New(ERROR_STORAGE_SERVER_DOWN_IS_EMPTY)
	}
	if s.saveFileName != "" {
		if err := s.receiveToFile(tcpConn); err != nil {
			return errors.New(ERROR_STORAGE_SERVER_DOWN_RECEIVE + err.Error())
		}
	} else {
		return errors.New(ERROR_STORAGE_SERVER_DOWN_FILENAME_EMPTY)
	}
	return nil
}

//receiveToFile 接受下载的数据流到指定的本地文件(设置文件名接受数据流)，注意：指定的文件名必须不能事先存在，避免对已经存在文件产生影响
// @tcpConn  tcp连接
func (s *storageDownloadHeaderBody) receiveToFile(tcpConn net.Conn) error {
	file, err := os.OpenFile(s.saveFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_EXCL, 0755)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		return err
	}
	writerBuf := bufio.NewWriter(file)
	if err = writeBufferFromTcpConn(tcpConn, writerBuf, s.pkgLen); err != nil {
		return errors.New(ERROR_STORAGE_SERVER_DOWN_FILE_RECEIVE + err.Error())
	}
	return nil
}
