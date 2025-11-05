package fastdfs_client_go

// TrackerStorageServerConfig  fast dfs  服务端参数配置
type TrackerStorageServerConfig struct {
	TrackerServer []string
	MaxConns      int
}

// RemoteFileInfo 查询远程服务器的文件信息
type RemoteFileInfo struct {
	fileSize        int64  //文件大小，单位 byte
	createTimestamp int64  //文件创建时间（Unix时间戳）
	crc32           int64  // 文件内容CRC32校验码
	SourceIpAddr    string //16字节字符串，源storage server IP地址
}

// storageServerInfo 服务器信息（需要通过 tracker server获取）
type storageServerInfo struct {
	addrPort         string
	storagePathIndex byte
}
