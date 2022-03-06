package fastdfs_client_go

import "time"

// fastdfs通信协议参考地址（作者2019年发布）： https://mp.weixin.qq.com/s/lpWEv3NCLkfKmtzKJ5lGzQ
// fastdfs通信协议参考地址（作者2009年发布）：http://bbs.chinaunix.net/thread-2001015-1-1.html
// FastDFS采用二进制TCP通信协议。一个数据包由 包头（header）和包体（body）组成。

// 程序配置常量
const (
	// 包头固定长度10个字节
	TCP_HEADER_LEN = 10
	//tracker 响应码
	TRACKER_PROTO_CMD_RESP = 100
	//获取一个storage server用来存储文件（指定组名
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE = 101
	// 获取一个 storage server 用来下载文件
	TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE = 102
	//  上传文件
	STORAGE_PROTO_CMD_UPLOAD_FILE = 11
	//  删除文件
	STORAGE_PROTO_CMD_DELETE_FILE = 12
	//  下载文件
	STORAGE_PROTO_CMD_DOWNLOAD_FILE = 14
	// 获取文件信息
	STORAGE_PROTO_CMD_QUERY_FILE_INFO = 22

	//  激活测试，通常用于检测连接是否有效
	// 客户端使用连接池的情况下，建立连接后发送一次active test即可和server端保持长连接。
	FDFS_PROTO_CMD_ACTIVE_TEST = 111

	//  groupName 长度常量
	FDFS_GROUP_NAME_FIX_LEN = 16

	// TCP连接池最小连接数
	TCP_CONNS_MIN_NUM = 3
	// tcp 连接超时时间
	TCP_CONN_TIMEOUT = time.Second * 10
	// tcp 连接最大空闲时间(秒)
	TCP_CONN_IDLE_TIMEOUT float64 = 15
	// tcp  心跳的秒数
	HEART_BEAT_SECOND = time.Second * 10

	// 文件的扩展名长度常量
	FILE_EXTNAME_FIX_LEN = 6

	// TCP 从内核读取数据到内存中转缓冲区的大小
	TCP_RECEIVE_BUFFER_SIZE = 4096

	// TCP 从内核读取数据到内存中转缓冲区的大小
	TCP_READ_BUFFER_SIZE = 512

	// 使用进程的状态常量，模拟 tcp 连接状态

	// 执行状态（从连接池取后的状态）
	TCP_STATUS_RUNNING = 'R'
	// 可中断的休眠状态（使用后归还到连接池的状态）
	TCP_STATUS_INTERRUPTIBLE = 'S'
	// 不可中断的睡眠状态（刚创建后的状态）
	TCP_STATUS_UNINTERRUPTIBLE = 'D'
)

// 错误常量
const (
	ERROR_CONN_POOL_OVER_MAX                    = "连接池超过已设置的最大连接数 , 当前最大连接数："
	ERROR_GET_TCP_CONN_FAIL                     = "从连接池获取一个 tcp 连接失败 , 出错明细："
	ERROR_TCP_SERVER_RESPONSE_NOT_ZERO          = "tcp 服务器返回的status状态值不为0"
	ERROR_FILE_FILENAME_IS_EMPTY                = "待上传的文件名不允许为空"
	ERROR_FILE_SIZE_IS_ZERO                     = "待上传的文件大小不允许为0字节"
	ERROR_FILE_DOWNLOAD_RELA_FILENAME_NOT_EMPTY = "下载文件时, 接受数据对应的文件名不能已经存在，否则可能会影响已经存在的文件数据"
	ERROR_FILE_EXT_NAME_IS_EMPTY                = "通过文件流(二进制)上传文件时，必须手动指定文件扩展名"
	ERROR_CONN_POOL_NO_ACTIVE_CONN              = "tcp 连接池中没有有效对象"
	ERROR_HEADER_RECEV_STATUS_NOT_ZERO          = "收到的消息头（receive header）中 status 值不为0"
	ERROR_HEADER_RECEV_ERROR                    = "收到的消息头（receive header）中有错误"
	ERROR_HEADER_RECEV_LEN_LT16_ERROR           = "收到的消息头（receive header）长度必须 > 16"
	ERROR_STORAGE_SERVER_FILE_NAME_FORMAT2      = "storage server 文件名格式不正确,  文件Id(fileId) 中必须至少存在一个斜杠( /)，不能不能在开头位置"
	ERROR_STORAGE_SERVER_DOWN_HEADER            = "storage server 下载时获取服务器响应头出错: "
	ERROR_STORAGE_SERVER_DOWN_IS_EMPTY          = "storage server 被下载的文件在服务器端不存在(或文件内容为空)"
	ERROR_STORAGE_SERVER_DOWN_FILENAME_EMPTY    = "storage server 下载文件时，必须指定文件名才能保存"
	ERROR_STORAGE_SERVER_DOWN_RECEIVE           = "storage server 下载时接受文件出错: "
	ERROR_STORAGE_SERVER_DOWN_FILE_RECEIVE      = "storage server 下载的文件在读取数据过程中出错："
	ERROR_STORAGE_SERVER_DOWN_FILE_WRITE_FLUSH  = "storage server 下载的文件写出到硬盘出错："
	ERROR_STORAGE_SERVER_FILE_UPLOAD_SEND_BYTES = "storage server 上传文件在通过tcp连接发送二进制文件时出错："
	ERROR_STORAGE_SERVER_GET_FILEINFO           = "storage server 查询文件信息时出错："
	ERROR_STORAGE_SERVER_GET_FILEINFO_BODY_LEN  = "storage server 查询文件信息时获取响应 body 长度不符合长度为 40 字节的标准"
	ERROR_TCP_CONN_ASSERT_FAIL                  = "从连接池中获取的 tcp 连接断言为结构体 tcpConnBaseInfo 失败"
)
