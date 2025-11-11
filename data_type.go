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

type GroupInfo struct {
	respGroupName               int64 // 17字节字符串
	resptotalMb                 int64 //磁盘空间总量，单位MB
	respFreeMb                  int64 // 磁盘剩余空间，单位MB
	respReservedMb              int64 // 磁盘预留空间，单位MB（since V6.13.1）
	resptrunkFreeMb             int64 // trunk文件剩余空间，单位MB（合并存储开启时有效）
	respServerCount             int64 // storage server数量
	respServerPort              int64 // storage server端口号
	respReadableServerCount     int64 // 当前可读的storage server数量（since V6.13）
	respWritableServerCount     int64 // 当前可写的storage server数量（since V6.13）
	respCurrentWriteServerCount int64 // 当前写入的 storage server顺序号
	respstorePathCount          int64 // storage server 存储路径数
	respsubdirCountPerPath      int64 // 存储路径下的子目录数（FastDFS采用两级子目录），如 256
	respcurrentTrunkFileId      int64 // 当前使用的trunk文件ID（合并存储开启时有效）
}
