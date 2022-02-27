package fastdfs_client_go

// TrackerStorageServerConfig  fast dfs  服务端参数配置
type TrackerStorageServerConfig struct {
	TrackerServer []string
	MaxConns      int
}

// storageServerInfo 服务器信息（需要通过 tracker server获取）
type storageServerInfo struct {
	addrPort         string
	storagePathIndex byte
}
