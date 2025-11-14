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

// trackerTcpConn trackerServer的tcp连接信息
type trackerTcpConn struct {
	// send 参数 ↓
	header
	groupName      string
	remoteFilename string
	// receive 接受返回结果 ↓
	storageInfo    storageInfo
	groupInfo      GroupInfo
	groups         []GroupInfo
	StorageServers []StorageServer
}

// trackerStorageInfo 通过 trackerServer 获取的 storageServer 信息
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
		if t.remoteFilename != "" {
			buffer.WriteString(t.remoteFilename)
		}
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
	// fastdfs v6.11 以后的版本支持到IP v6，所以IP地址的长度为15字节 + 新偏移30字节
	// 现在直接按照新的协议标准来解析
	//fmt.Println("header-cmd:", t.header.cmd)
	//fmt.Println("header-respCmd:", t.header.respCmd)
	//fmt.Println("header-len:", t.pkgLen)
	var ipV6Offset = 30
	switch t.header.cmd {
	// 后面的参数需要根据具体的命令去设置
	case TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE:
		//@group_name：16字节字符串，组名
		//@ip_addr：15/45字节字符串， storage server IP地址（V6.11前为15字节，V6.11开始为45字节）
		//@port：8字节整数，storage server端口号
		if t.pkgLen <= 39 {
			ipV6Offset = 0
		}
		t.groupName = string(getBytesByPosition(buf, 0, 15))
		t.storageInfo.ipAddr = string(getBytesByPosition(buf, 16, 15+ipV6Offset))
		t.storageInfo.port = bytesToInt(getBytesByPosition(buf, 31+ipV6Offset, 8))
	case TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE:
		// 响应body：
		//@group_name：16字节字符串，组名
		//@ip_addr：15字节字符串， storage server IP地址, IPv6在服务端启用时，IP地址的偏移量+30字节
		//@port：8字节整数，storage server端口号
		//@store_path_index：1字节整数，基于0的存储路径顺序号

		if t.pkgLen <= 40 {
			ipV6Offset = 0
		}
		t.groupName = string(getBytesByPosition(buf, 0, 15))
		t.storageInfo.ipAddr = string(getBytesByPosition(buf, 16, 15+ipV6Offset))
		t.storageInfo.port = bytesToInt(getBytesByPosition(buf, 31+ipV6Offset, 8))
		t.storageInfo.storePathIndex = getBytesByPosition(buf, 39+ipV6Offset, 1)[0]
	case TRACKER_PROTO_CMD_SERVER_LIST_ONE_GROUP:
		// 113 字节表示单个组信息的长度
		if t.header.pkgLen == 113 {
			t.groupInfo.groupName = string(getBytesByPosition(buf, 0, 17))
			t.groupInfo.totalMb = bytesToInt(getBytesByPosition(buf, 17, 8))
			t.groupInfo.freeMb = bytesToInt(getBytesByPosition(buf, 25, 8))
			t.groupInfo.reservedMb = bytesToInt(getBytesByPosition(buf, 33, 8))
			t.groupInfo.trunkFreeMb = bytesToInt(getBytesByPosition(buf, 41, 8))
			t.groupInfo.serverCount = bytesToInt(getBytesByPosition(buf, 49, 8))
			t.groupInfo.serverPort = bytesToInt(getBytesByPosition(buf, 57, 8))
			t.groupInfo.readableServerCount = bytesToInt(getBytesByPosition(buf, 65, 8))
			t.groupInfo.writableServerCount = bytesToInt(getBytesByPosition(buf, 73, 8))
			t.groupInfo.currentWriteServerCount = bytesToInt(getBytesByPosition(buf, 81, 8))
			t.groupInfo.storePathCount = bytesToInt(getBytesByPosition(buf, 89, 8))
			t.groupInfo.subdirCountPerPath = bytesToInt(getBytesByPosition(buf, 97, 8))
			t.groupInfo.currentTrunkFileId = bytesToInt(getBytesByPosition(buf, 105, 8))
		}
	case TRACKER_PROTO_CMD_SERVER_LIST_ALL_GROUPS:
		// 113 字节表示单个组信息的长度
		var groupBodyLen = 113
		groupNum := int(t.header.pkgLen / int64(groupBodyLen))
		// 调整groups切片长度，确保有足够空间存储所有组信息
		t.groups = make([]GroupInfo, groupNum)
		for i := 0; i < groupNum; i++ {
			t.groups[i].groupName = string(getBytesByPosition(buf, 0+(i*groupBodyLen), 17))
			t.groups[i].totalMb = bytesToInt(getBytesByPosition(buf, 17+(i*groupBodyLen), 8))
			t.groups[i].freeMb = bytesToInt(getBytesByPosition(buf, 25+(i*groupBodyLen), 8))
			t.groups[i].reservedMb = bytesToInt(getBytesByPosition(buf, 33+(i*groupBodyLen), 8))
			t.groups[i].trunkFreeMb = bytesToInt(getBytesByPosition(buf, 41+(i*groupBodyLen), 8))
			t.groups[i].serverCount = bytesToInt(getBytesByPosition(buf, 49+(i*groupBodyLen), 8))
			t.groups[i].serverPort = bytesToInt(getBytesByPosition(buf, 57+(i*groupBodyLen), 8))
			t.groups[i].readableServerCount = bytesToInt(getBytesByPosition(buf, 65+(i*groupBodyLen), 8))
			t.groups[i].writableServerCount = bytesToInt(getBytesByPosition(buf, 73+(i*groupBodyLen), 8))
			t.groups[i].currentWriteServerCount = bytesToInt(getBytesByPosition(buf, 81+(i*groupBodyLen), 8))
			t.groups[i].storePathCount = bytesToInt(getBytesByPosition(buf, 89+(i*groupBodyLen), 8))
			t.groups[i].subdirCountPerPath = bytesToInt(getBytesByPosition(buf, 97+(i*groupBodyLen), 8))
			t.groups[i].currentTrunkFileId = bytesToInt(getBytesByPosition(buf, 105+(i*groupBodyLen), 8))
		}
	case TRACKER_PROTO_CMD_SERVER_LIST_STORAGE:
		// 516 表示服务器返回的body体单个 storage server 信息的长度
		var storagerServerBodyLen = 516 // 最新版本返回实际长度：487，与文档不一致，等待后续修正代码
		storageServerNum := int(t.header.pkgLen / int64(storagerServerBodyLen))
		// 调整 StorageServers 切片长度，确保有足够空间存储所有组信息
		t.StorageServers = make([]StorageServer, storageServerNum)
		for i := 0; i < storageServerNum; i++ {
			t.StorageServers[i].Status = getBytesByPosition(buf, 0+(i*storagerServerBodyLen), 1)[0]
			t.StorageServers[i].RwMode = getBytesByPosition(buf, 1+(i*storagerServerBodyLen), 1)[0]
			t.StorageServers[i].Id = string(getBytesByPosition(buf, 2+(i*storagerServerBodyLen), 16))
			t.StorageServers[i].IpAddr = string(getBytesByPosition(buf, 18+(i*storagerServerBodyLen), 46))
			t.StorageServers[i].SrcStorageId = string(getBytesByPosition(buf, 64+(i*storagerServerBodyLen), 16))
			t.StorageServers[i].Version = string(getBytesByPosition(buf, 80+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].JoinTime = bytesToInt(getBytesByPosition(buf, 88+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].UpTime = bytesToInt(getBytesByPosition(buf, 96+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalMb = bytesToInt(getBytesByPosition(buf, 104+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].FreeMb = bytesToInt(getBytesByPosition(buf, 112+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].ReservedMb = bytesToInt(getBytesByPosition(buf, 120+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].UploadPriority = bytesToInt(getBytesByPosition(buf, 128+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].StorePathCount = bytesToInt(getBytesByPosition(buf, 136+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SubdirCountPerPath = bytesToInt(getBytesByPosition(buf, 144+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].CurrentWritePath = bytesToInt(getBytesByPosition(buf, 152+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].StoragePort = bytesToInt(getBytesByPosition(buf, 160+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].AllocCount = int32(bytesToInt(getBytesByPosition(buf, 168+(i*storagerServerBodyLen), 4)))
			t.StorageServers[i].CurrentCount = int32(bytesToInt(getBytesByPosition(buf, 172+(i*storagerServerBodyLen), 4)))
			t.StorageServers[i].MaxCount = int32(bytesToInt(getBytesByPosition(buf, 176+(i*storagerServerBodyLen), 4)))
			t.StorageServers[i].TotalUploadCount = bytesToInt(getBytesByPosition(buf, 180+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessUploadCount = bytesToInt(getBytesByPosition(buf, 188+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalAppendCount = bytesToInt(getBytesByPosition(buf, 196+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessAppendCount = bytesToInt(getBytesByPosition(buf, 204+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalModifyCount = bytesToInt(getBytesByPosition(buf, 212+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessModifyCount = bytesToInt(getBytesByPosition(buf, 220+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalTruncateCount = bytesToInt(getBytesByPosition(buf, 228+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessTruncateCount = bytesToInt(getBytesByPosition(buf, 236+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalSetMetaCount = bytesToInt(getBytesByPosition(buf, 244+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessSetMetaCount = bytesToInt(getBytesByPosition(buf, 252+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalDeleteCount = bytesToInt(getBytesByPosition(buf, 260+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessDeleteCount = bytesToInt(getBytesByPosition(buf, 268+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalDownloadCount = bytesToInt(getBytesByPosition(buf, 276+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessDownloadCount = bytesToInt(getBytesByPosition(buf, 284+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalGetMetaCount = bytesToInt(getBytesByPosition(buf, 292+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessGetMetaCount = bytesToInt(getBytesByPosition(buf, 300+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalCreateLinkCount = bytesToInt(getBytesByPosition(buf, 308+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessCreateLinkCount = bytesToInt(getBytesByPosition(buf, 316+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalDeleteLinkCount = bytesToInt(getBytesByPosition(buf, 324+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessDeleteLinkCount = bytesToInt(getBytesByPosition(buf, 332+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalUploadBytes = bytesToInt(getBytesByPosition(buf, 340+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessUploadBytes = bytesToInt(getBytesByPosition(buf, 348+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalAppendBytes = bytesToInt(getBytesByPosition(buf, 356+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessAppendBytes = bytesToInt(getBytesByPosition(buf, 364+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalModifyBytes = bytesToInt(getBytesByPosition(buf, 372+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessModifyBytes = bytesToInt(getBytesByPosition(buf, 380+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalDownloadBytes = bytesToInt(getBytesByPosition(buf, 388+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessDownloadBytes = bytesToInt(getBytesByPosition(buf, 396+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalSyncInBytes = bytesToInt(getBytesByPosition(buf, 404+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessSyncInBytes = bytesToInt(getBytesByPosition(buf, 412+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalSyncOutBytes = bytesToInt(getBytesByPosition(buf, 420+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessSyncOutBytes = bytesToInt(getBytesByPosition(buf, 428+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalFileOpenCount = bytesToInt(getBytesByPosition(buf, 436+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessFileOpenCount = bytesToInt(getBytesByPosition(buf, 444+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalFileReadCount = bytesToInt(getBytesByPosition(buf, 452+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessFileReadCount = bytesToInt(getBytesByPosition(buf, 460+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].TotalFileWriteCount = bytesToInt(getBytesByPosition(buf, 468+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].SuccessFileWriteCount = bytesToInt(getBytesByPosition(buf, 476+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].LastSourceUpdate = bytesToInt(getBytesByPosition(buf, 484+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].LastSyncUpdate = bytesToInt(getBytesByPosition(buf, 492+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].LastSyncedTimestamp = bytesToInt(getBytesByPosition(buf, 500+(i*storagerServerBodyLen), 8))
			t.StorageServers[i].LastHeartBeatTime = bytesToInt(getBytesByPosition(buf, 508+(i*storagerServerBodyLen), 8))
		}
	default:
		//
	}
	return nil
}
