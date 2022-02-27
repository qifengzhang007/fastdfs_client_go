package fastdfs_client_go

import (
	"errors"
	"strconv"
	"sync"
)

func CreateFdfsClient(trackerServerOptions *TrackerStorageServerConfig) (*trackerServerTcpClient, error) {

	tcpClient := &trackerServerTcpClient{
		trackerServerConfig: trackerServerOptions,
		trackerPools:        make(map[string]*tcpConnPool),
		storagePoolLock:     &sync.Mutex{},
		storagePools:        make(map[string]*tcpConnPool),
	}
	for _, addr := range trackerServerOptions.TrackerServer {
		trackerServerPool, err := initTcpConnPool(addr, trackerServerOptions.MaxConns)
		if err != nil {
			return nil, err
		}
		tcpClient.trackerPools[addr] = trackerServerPool
	}

	return tcpClient, nil
}

// trackerServerTcpClient 创建一个go语言连接 fastdfs 服务的 tcp 客户端
// 一个客户端可以同时连接到 tracker server 和  storage server
type trackerServerTcpClient struct {
	trackerServerConfig *TrackerStorageServerConfig
	trackerPools        map[string]*tcpConnPool
	storagePools        map[string]*tcpConnPool
	storagePoolLock     *sync.Mutex
}

// getTrackerConn 从连接池获取一个 trackerServer 的 tcp 连接
// @ 参数 ：无
// 返回参数解释：
// tcpConnPool 连接池地址
// tcpConnBaseInfo 从连接池中获取的tcp连接
// error 可能的错误
func (c *trackerServerTcpClient) getTrackerConn() (*tcpConnPool, *tcpConnBaseInfo, error) {
	// 连接池地址
	var trackerPool *tcpConnPool
	// 从连接池获取的tcp连接
	var trackerConn *tcpConnBaseInfo
	var err error
	var getOne bool
	for _, trackerPool = range c.trackerPools {
		trackerConn, err = trackerPool.get()
		if err == nil {
			getOne = true
			break
		}
	}
	if getOne {
		// 返回连接池地址、连接池地址获取的tcp连接对象、错误
		return trackerPool, trackerConn, nil
	}
	if err == nil {
		return nil, nil, errors.New(ERROR_CONN_POOL_NO_ACTIVE_CONN)
	}
	return nil, nil, err
}

// Destroy  整个客户端销毁时，关闭连接池中的所有tcp连接（包括 trackerServer 和 storageServer）
func (c *trackerServerTcpClient) Destroy() {
	for _, pool := range c.trackerPools {
		pool.Destroy()
	}
	for _, pool := range c.storagePools {
		pool.Destroy()
	}
}

// getStorageInfoByTracker  主要通过 tracker server 获取 storage server 服务的ip、端口等信息，然后通过 storage server 传输文件
// @ body 参数 ： 不需要
func (c *trackerServerTcpClient) getStorageInfoByTracker(cmd byte, groupName string, remoteFilename string) (*storageServerInfo, error) {
	trackerSendParmas := &trackerTcpConn{}

	// 将命令参数设置在 header 头部分
	trackerSendParmas.header.pkgLen = 0
	trackerSendParmas.header.cmd = cmd
	trackerSendParmas.header.status = 0
	trackerSendParmas.groupName = groupName
	trackerSendParmas.remoteFilename = remoteFilename

	if err := c.sendHeaderByTrackerServer(trackerSendParmas); err != nil {
		return nil, err
	}
	return &storageServerInfo{
		addrPort:         trackerSendParmas.storageInfo.ipAddr + ":" + strconv.FormatInt(trackerSendParmas.storageInfo.port, 10),
		storagePathIndex: trackerSendParmas.storageInfo.storePathIndex,
	}, nil
}

// sendHeaderByTrackerServer  通过trackerServer 的header 头参数发送特定命令获取 storageServer 服务器
// @trackerTcpConn trackerServer 的 tcp连接
func (c *trackerServerTcpClient) sendHeaderByTrackerServer(trackerTcpConn tcpSendReceive) error {
	trackerTcpPoolPtr, trackerTcp, err := c.getTrackerConn()
	if err != nil {
		return err
	}
	defer func() {
		trackerTcpPoolPtr.put(trackerTcp)
	}()
	if err = trackerTcpConn.Send(trackerTcp); err != nil {
		return err
	}
	if err = trackerTcpConn.Receive(trackerTcp); err != nil {
		return err
	}
	return nil
}

// getStorageConn 通过 trackerServer 获取的参数，创建 StorageServer 的tcp连接
// @storageServInfo   trackerServer 获取的 storageServer 参数
// 返回参数解释：
// storageTcpConnPool 连接池地址
// tcpConnBaseInfo 从连接池中获取的tcp连接
// err 可能的错误
func (c *trackerServerTcpClient) getStorageConn(storageServInfo *storageServerInfo) (storageTcpConnPool *tcpConnPool, tcpConnBaseInfo *tcpConnBaseInfo, err error) {
	c.storagePoolLock.Lock()
	defer c.storagePoolLock.Unlock()
	var isOk bool
	storageTcpConnPool, isOk = c.storagePools[storageServInfo.addrPort]
	if isOk {
		tcpConnBaseInfo, err = storageTcpConnPool.get()
		if err == nil {
			c.storagePools[storageServInfo.addrPort] = storageTcpConnPool
		}
		return
	}
	storageTcpConnPool, err = initTcpConnPool(storageServInfo.addrPort, c.trackerServerConfig.MaxConns)
	if err == nil {
		tcpConnBaseInfo, err = storageTcpConnPool.get()
		if err == nil {
			c.storagePools[storageServInfo.addrPort] = storageTcpConnPool
		}
		return
	}
	return nil, nil, err
}

// sendCmdToStorageServer  给 storageServer 发送具体的业务命令
// @headerBody  实现了 tcpSendReceive 接口的 header 和 body 参数组装的结构体
// @storageInfo  storageServer 的服务器信息，用于创建到  storageServer 的tcp连接
func (c *trackerServerTcpClient) sendCmdToStorageServer(headerBody tcpSendReceive, storageInfo *storageServerInfo) error {
	storageTcpPool, storageTcpConn, err := c.getStorageConn(storageInfo)
	if err != nil {
		return err
	}
	defer func() {
		storageTcpPool.put(storageTcpConn)
	}()

	if err = headerBody.Send(storageTcpConn); err != nil {
		return err
	}
	if err = headerBody.Receive(storageTcpConn); err != nil {
		return err
	}

	return nil
}
