package fastdfs_client_go

import (
	"net"
)

// tcpSendReceive  发送、接受消息必须通过 tcp 连接，定义一个统一的接口
type tcpSendReceive interface {
	Send(tcpConn net.Conn) error
	Receive(tcpConn net.Conn) error
}
