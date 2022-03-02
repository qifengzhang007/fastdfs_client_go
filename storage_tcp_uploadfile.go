package fastdfs_client_go

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

type storageServerUploadHeaderBody struct {
	header
	// header  额外参数，发送使用
	fileInfo         *fileInfo
	storagePathIndex byte
	// 文件id (fileId) 响应信息
	fileId string
}

// Send 文件上传 body 参数详情
// @tcpConn  tcp连接
func (s *storageServerUploadHeaderBody) Send(tcpConn net.Conn) (err error) {
	// @store_path_index  ， byte 1字节 ，存储目录序号
	// @fileSize int64，8字节，     传输的二进制数据的字节数目
	// @file_ext_name  文件扩展名6个字节
	// @file_content：file_size 字节二进制内容，文件内容
	// 15 的组成 ：@store_path_index (1字节整数) + @fileSize ( 8字节整数 ) + @file_ext_name (6字节字符串)
	s.header.pkgLen = s.fileInfo.fileSize + 15
	s.header.cmd = STORAGE_PROTO_CMD_UPLOAD_FILE
	s.header.status = 0

	if err = s.header.sendHeader(tcpConn); err != nil {
		return err
	}
	// 创建 body 数据发送时的缓冲区
	buffer := new(bytes.Buffer)

	// @store_path_index (1字节整数)
	buffer.WriteByte(s.storagePathIndex)

	//@file_size：8字节整数
	if err = binary.Write(buffer, binary.BigEndian, s.fileInfo.fileSize); err != nil {
		return err
	}
	// 文件扩展名 6 字节,必须要写满6个字节
	buffer.Write(specialFileExtNameConvBytes(s.fileInfo.fileExtName))
	if _, err = tcpConn.Write(buffer.Bytes()); err != nil {
		return err
	}

	//发送文件内容本身的二进制数据
	// 1.如果文件结构信息对应的指针不为空，表示该文件是通过文件名方式打开操作的，那么就根据文件指针读取数据，发送出去
	// 2.如果文件结构信息中的文件指针为空，表示该文件需要通过字节流发送出去
	// 3.最后将文件真正的内容发出去，其实 tcp 底层会对大文件分块多次发送，服务器端会按照收到的报文头读取对应的数量字节才结束.
	if s.fileInfo.filePtr != nil {
		if _, err = sendBytesByFilePtr(s.fileInfo.filePtr, tcpConn); err != nil {
			return errors.New(ERROR_STORAGE_SERVER_FILE_UPLOAD_SEND_BYTES + err.Error())
		}
	} else {
		_, err = tcpConn.Write(s.fileInfo.buffer)
	}

	if err != nil {
		return err
	}
	return nil
}

// Receive 发送文件上传命令之后接受服务器的响应头
// @tcpConn  tcp连接
func (s *storageServerUploadHeaderBody) Receive(tcpConn net.Conn) error {
	// @group_name：16字节字符串，组名
	// @filename：不定长字符串，文件名
	if err := s.header.receiveHeader(tcpConn); err != nil {
		return err
	}
	if s.pkgLen <= 16 {
		return errors.New(ERROR_HEADER_RECEV_LEN_LT16_ERROR)
	}

	buf := make([]byte, s.pkgLen)
	if _, err := tcpConn.Read(buf); err != nil {
		return err
	}
	s.fileId = string(getBytesByPosition(buf, 0, 16)) + "/" + string(getBytesByPosition(buf, 16, int(s.pkgLen)-16))
	return nil
}
