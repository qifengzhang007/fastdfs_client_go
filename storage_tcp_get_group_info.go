package fastdfs_client_go

import (
	"bytes"
	"errors"
	"net"
)

// storageGetGroupInfoHeaderBody 查询一个group的基本信息
type storageGetGroupInfoHeaderBody struct {
	header
	// header  额外参数，发送使用
	groupName string // 发送需要限制为16个字符长度
	// 响应信息 113 个字节长度
	GroupInfo
	//respGroupName               int64 // 17字节字符串
	//resptotalMb                 int64 //磁盘空间总量，单位MB
	//respFreeMb                  int64 // 磁盘剩余空间，单位MB
	//respReservedMb              int64 // 磁盘预留空间，单位MB（since V6.13.1）
	//resptrunkFreeMb             int64 // trunk文件剩余空间，单位MB（合并存储开启时有效）
	//respServerCount             int64 // storage server数量
	//respServerPort              int64 // storage server端口号
	//respReadableServerCount     int64 // 当前可读的storage server数量（since V6.13）
	//respWritableServerCount     int64 // 当前可写的storage server数量（since V6.13）
	//respCurrentWriteServerCount int64 // 当前写入的 storage server顺序号
	//respstorePathCount          int64 // storage server 存储路径数
	//respsubdirCountPerPath      int64 // 存储路径下的子目录数（FastDFS采用两级子目录），如 256
	//respcurrentTrunkFileId      int64 // 当前使用的trunk文件ID（合并存储开启时有效）
}

// Send 发送删除文件命令
// @tcpConn tcp连接
func (s *storageGetGroupInfoHeaderBody) Send(tcpConn net.Conn) error {
	// 设置删除文件时的 header 参数
	//@group_name：16字节字符串，组名
	//@filename：不定长字符串，文件名
	s.header.pkgLen = 16 // 协议规定 group_name 长度为必须是 16 字节
	s.header.cmd = TRACKER_PROTO_CMD_SERVER_LIST_ONE_GROUP
	s.header.status = 0

	if err := s.header.sendHeader(tcpConn); err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	buffer.Write(groupNameConvBytes(s.groupName))
	if _, err := tcpConn.Write(buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

// Receive  接受删除命令发送服务端的响应头
// @tcpConn tcp连接
func (s *storageGetGroupInfoHeaderBody) Receive(tcpConn net.Conn) error {
	if err := s.header.receiveHeader(tcpConn); err != nil {
		if int(s.header.status) != 0 {
			return errors.New(ERROR_GET_GROUP_INFO_FAILED)
		}
	}
	buf := make([]byte, 113)
	if _, err := tcpConn.Read(buf); err != nil {
		return err
	}
	// 从返回的长度逐个解析出具体的返回字段
	s.respGroupName = bytesToInt(getBytesByPosition(buf, 0, 17))
	s.resptotalMb = bytesToInt(getBytesByPosition(buf, 17, 8))
	s.respFreeMb = bytesToInt(getBytesByPosition(buf, 25, 8))
	s.respReservedMb = bytesToInt(getBytesByPosition(buf, 33, 8))
	s.resptrunkFreeMb = bytesToInt(getBytesByPosition(buf, 41, 8))
	s.respServerCount = bytesToInt(getBytesByPosition(buf, 49, 8))
	s.respServerPort = bytesToInt(getBytesByPosition(buf, 57, 8))
	s.respReadableServerCount = bytesToInt(getBytesByPosition(buf, 65, 8))
	s.respWritableServerCount = bytesToInt(getBytesByPosition(buf, 73, 8))
	s.respCurrentWriteServerCount = bytesToInt(getBytesByPosition(buf, 81, 8))
	s.respstorePathCount = bytesToInt(getBytesByPosition(buf, 89, 8))
	s.respsubdirCountPerPath = bytesToInt(getBytesByPosition(buf, 97, 8))
	s.respcurrentTrunkFileId = bytesToInt(getBytesByPosition(buf, 105, 8))
	return nil
}
