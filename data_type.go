package fastdfs_client_go

// TrackerStorageServerConfig  fast dfs  服务端参数配置
type TrackerStorageServerConfig struct {
	TrackerServer []string
	MaxConns      int
}

// RemoteFileInfo 查询远程服务器的文件信息
type RemoteFileInfo struct {
	fileSize        int64
	createTimestamp int64
	crc32           int64
	SourceIpAddr    string
}

// storageServerInfo 服务器信息（需要通过 tracker server获取）
type storageServerInfo struct {
	addrPort         string
	storagePathIndex byte
}
