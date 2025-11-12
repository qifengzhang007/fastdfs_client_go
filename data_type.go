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
	groupName               string // 17字节字符串
	totalMb                 int64  //磁盘空间总量，单位MB
	freeMb                  int64  // 磁盘剩余空间，单位MB
	reservedMb              int64  // 磁盘预留空间，单位MB（since V6.13.1）
	trunkFreeMb             int64  // trunk文件剩余空间，单位MB（合并存储开启时有效）
	serverCount             int64  // storage server数量
	serverPort              int64  // storage server端口号
	readableServerCount     int64  // 当前可读的storage server数量（since V6.13）
	writableServerCount     int64  // 当前可写的storage server数量（since V6.13）
	currentWriteServerCount int64  // 当前写入的 storage server顺序号
	storePathCount          int64  // storage server 存储路径数
	subdirCountPerPath      int64  // 存储路径下的子目录数（FastDFS采用两级子目录），如 256
	currentTrunkFileId      int64  // 当前使用的trunk文件ID（合并存储开启时有效）
}

type StorageServer struct {
	Status                 byte   // 1字节整数，storage server状态
	RwMode                 byte   // 1字节整数，读写模式（since v6.13）
	Id                     string // 16字节字符串，server ID
	IpAddr                 string // 46字节字符串，IP地址（V6.11前为16字节，之后为46字节）
	SrcStorageId           string // 16字节字符串，同步源storage的server ID
	Version                string // 8字节字符串，运行的FastDFS版本号，例如6.04（V6.11前为6字节，V6.11开始为8字节）
	JoinTime               int64  // 8字节整数，加入集群时间
	UpTime                 int64  // 8字节整数，fdfs_storaged启动时间
	TotalMb                int64  // 8字节整数，磁盘空间总量，单位MB
	FreeMb                 int64  // 8字节整数，磁盘剩余空间，单位MB
	ReservedMb             int64  // 8字节整数，磁盘预留空间，单位MB（since V6.13.1）
	UploadPriority         int64  // 8字节整数，上传文件优先级
	StorePathCount         int64  // 8字节整数，存储路径数
	SubdirCountPerPath     int64  // 8字节整数，存储路径下的子目录数（FastDFS采用两级子目录），如 256
	CurrentWritePath       int64  // 8字节整数，当前写入的存储路径（顺序号）
	StoragePort            int64  // 8字节整数，storage server服务端口号
	AllocCount             int32  // 4字节整数，已分配的连接buffer数目
	CurrentCount           int32  // 4字节整数，当前连接数
	MaxCount               int32  // 4字节整数，曾经达到过的最大连接数
	TotalUploadCount       int64  // 8字节整数，上传文件总数
	SuccessUploadCount     int64  // 8字节整数，成功上传文件数
	TotalAppendCount       int64  // 8字节整数，调用append总次数
	SuccessAppendCount     int64  // 8字节整数，成功调用append次数
	TotalModifyCount       int64  // 8字节整数，调用modify总次数
	SuccessModifyCount     int64  // 8字节整数，成功调用modify次数
	TotalTruncateCount     int64  // 8字节整数，调用truncate总次数
	SuccessTruncateCount   int64  // 8字节整数，成功调用truncate次数
	TotalSetMetaCount      int64  // 8字节整数，设置文件附加属性（meta data）总次数
	SuccessSetMetaCount    int64  // 8字节整数，成功设置文件附加属性（meta data）次数
	TotalDeleteCount       int64  // 8字节整数，删除文件总数
	SuccessDeleteCount     int64  // 8字节整数，成功删除文件数
	TotalDownloadCount     int64  // 8字节整数，下载文件总数
	SuccessDownloadCount   int64  // 8字节整数，成功下载文件数
	TotalGetMetaCount      int64  // 8字节整数，获取文件附加属性（meta data）总次数
	SuccessGetMetaCount    int64  // 8字节整数，成功获取文件附加属性（meta data）次数
	TotalCreateLinkCount   int64  // 8字节整数，创建文件符号链接总数
	SuccessCreateLinkCount int64  // 8字节整数，成功创建文件符号链接数
	TotalDeleteLinkCount   int64  // 8字节整数，删除文件符号链接总数
	SuccessDeleteLinkCount int64  // 8字节整数，成功删除文件符号链接数
	TotalUploadBytes       int64  // 8字节整数，上传文件总字节数
	SuccessUploadBytes     int64  // 8字节整数，成功上传文件字节数
	TotalAppendBytes       int64  // 8字节整数，append总字节数
	SuccessAppendBytes     int64  // 8字节整数，成功append字节数
	TotalModifyBytes       int64  // 8字节整数，modify总字节数
	SuccessModifyBytes     int64  // 8字节整数，成功modify字节数
	TotalDownloadBytes     int64  // 8字节整数，下载总字节数
	SuccessDownloadBytes   int64  // 8字节整数，成功下载字节数
	TotalSyncInBytes       int64  // 8字节整数，文件同步流入总字节数
	SuccessSyncInBytes     int64  // 8字节整数，文件同步成功流入字节数
	TotalSyncOutBytes      int64  // 8字节整数，文件同步流出总字节数
	SuccessSyncOutBytes    int64  // 8字节整数，文件同步成功流出字节数
	TotalFileOpenCount     int64  // 8字节整数，文件打开总次数
	SuccessFileOpenCount   int64  // 8字节整数，文件成功打开次数
	TotalFileReadCount     int64  // 8字节整数，文件读总次数
	SuccessFileReadCount   int64  // 8字节整数，文件成功读次数
	TotalFileWriteCount    int64  // 8字节整数，文件写总次数
	SuccessFileWriteCount  int64  // 8字节整数，文件成功写次数
	LastSourceUpdate       int64  // 8字节整数，最近一次源头更新时间
	LastSyncUpdate         int64  // 8字节整数，最近一次同步更新时间
	LastSyncedTimestamp    int64  // 8字节整数，最近一次被同步到的时间戳
	LastHeartBeatTime      int64  // 8字节整数，最近一次心跳时间
}
