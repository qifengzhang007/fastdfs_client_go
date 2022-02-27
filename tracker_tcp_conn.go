package fastdfs_client_go

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

// fastdfs 设计原理：
// 1.客户端首先连接到 trackerServer，通过发送指定的命令（cmd常量值）获取 storageServer 的信息
// 2.客户端根据第一步获取的 storageServer =继续建立长连接
// 3.最后，所有的文件上传、下载都是通过第二步中获取的 storageServer 连接继续发送具体命令实现的

// 本页面的代码主要功能：
// 1.与 trackerServer 实现网络通讯

//  trackerTcpConn trackerServer的tcp连接信息
type trackerTcpConn struct {
	// send 参数 ↓
	header
	groupName      string
	remoteFilename string
	// receive 接受返回结果 ↓
	storageInfo storageInfo
}

//  trackerStorageInfo 通过 trackerServer 获取的 storageServer 信息
type storageInfo struct {
	ipAddr         string
	port           int64
	groupName      string
	storePathIndex byte
}

// Send 向对方发送数据
// @tcpConn tcp 连接
func (t *trackerTcpConn) Send(tcpConn net.Conn) error {
	if t.groupName == "" {
		if err := t.header.sendHeader(tcpConn); err != nil {
			return err
		}
	} else {
		t.header.pkgLen = int64(FDFS_GROUP_NAME_FIX_LEN + len(t.remoteFilename))
		buffer := new(bytes.Buffer)
		if err := binary.Write(buffer, binary.BigEndian, t.pkgLen); err != nil {
			return err
		}
		buffer.WriteByte(t.header.cmd)
		buffer.WriteByte(t.header.status)
		buffer.Write(groupNameConvBytes(t.groupName))
		buffer.WriteString(t.remoteFilename)
		if _, err := tcpConn.Write(buffer.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

// Receive 接受对方返回结果
func (t *trackerTcpConn) Receive(tcpConn net.Conn) error {
	if err := t.receiveHeader(tcpConn); err != nil {
		return errors.New(ERROR_HEADER_RECEV_ERROR + err.Error())
	}
	buf := make([]byte, t.pkgLen)
	if _, err := tcpConn.Read(buf); err != nil {
		return err
	}
	// 通信协议详情地址： https://mp.weixin.qq.com/s/lpWEv3NCLkfKmtzKJ5lGzQ
	// 响应body：
	//@group_name：16字节字符串，组名
	//@ip_addr：15字节字符串， storage server IP地址
	//@port：8字节整数，storage server端口号
	//@store_path_index：1字节整数，基于0的存储路径顺序号
	t.groupName = string(getBytesByPosition(buf, 0, 15))
	t.storageInfo.ipAddr = string(getBytesByPosition(buf, 16, 15))
	t.storageInfo.port = bytesToInt(getBytesByPosition(buf, 31, 8))
	switch t.header.cmd {
	// 后面的参数需要根据具体的命令去设置
	case STORAGE_PROTO_CMD_UPLOAD_FILE:
		t.storageInfo.storePathIndex = getBytesByPosition(buf, 39, 1)[0]
	default:
		//
	}
	return nil
}
