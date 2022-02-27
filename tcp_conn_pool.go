package fastdfs_client_go

import (
	"container/list"
	"errors"
	"net"
	"strconv"
	"sync"
	"time"
)

// tcpConnPool  连接池
type tcpConnPool struct {
	conns    *list.List
	addrPort string
	maxConns int
	count    int
	lock     *sync.Mutex
	isQuit   chan bool
}

// tcpConnBaseInfo 连接应该具有的基本特征
type tcpConnBaseInfo struct {
	// 参见常量 TCP_STATUS 开头的相关值
	status byte
	// 最后一次归还到连接池的时间点
	putTime time.Time
	// 一个 tcp 连接
	net.Conn
}

// 初始化一个tcp连接池，最少数量为 3
func initTcpConnPool(addrPort string, maxConns int) (*tcpConnPool, error) {
	if maxConns < TCP_CONNS_MIN_NUM {
		maxConns = TCP_CONNS_MIN_NUM
	}
	connPool := &tcpConnPool{
		conns:    list.New(),
		addrPort: addrPort,
		maxConns: maxConns,
		lock:     &sync.Mutex{},
		isQuit:   make(chan bool),
	}
	connPool.lock.Lock()
	defer connPool.lock.Unlock()
	for i := 0; i < TCP_CONNS_MIN_NUM; i++ {
		if err := connPool.CreateTcpConn(); err != nil {
			return nil, err
		}
	}
	// 开启一个周期性心跳包
	go func() {
		ticker := time.NewTicker(HEART_BEAT_SECOND)
		for {
			select {
			case isQuit := <-connPool.isQuit:
				if isQuit {
					ticker.Stop()
					return
				}
			case <-ticker.C:
				_ = connPool.checkTcpConnPool()
			}
		}
	}()
	return connPool, nil
}

// Destroy tcp 连接关闭
// 1.首先给管道发送退出消息，结束心跳任务
// 2.关闭一个tcp连接对应的连接池中的全部连接
func (t *tcpConnPool) Destroy() {
	t.isQuit <- true
	for tcpConn := t.conns.Front(); tcpConn != nil; tcpConn = tcpConn.Next() {
		conn := tcpConn.Value.(*tcpConnBaseInfo)
		_ = conn.Close()
		t.conns.Remove(tcpConn)
		t.count--
	}
}

// checkTcpConnPool 检测连接池中所有的 tcp 连接是否有效，把失效的tcp连接删除
func (t *tcpConnPool) checkTcpConnPool() (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	var isOk bool
	var conn *tcpConnBaseInfo
	for tcpConn := t.conns.Front(); tcpConn != nil; tcpConn = tcpConn.Next() {
		if conn, isOk = tcpConn.Value.(*tcpConnBaseInfo); isOk {
			// 1.首先检查休眠状态超时的 tcp 连接，直接从连接池删除
			if conn.status == TCP_STATUS_INTERRUPTIBLE && time.Now().Sub(conn.putTime).Seconds() > TCP_CONN_IDLE_TIMEOUT {
				_ = conn.Close()
				t.conns.Remove(tcpConn)
				t.count--
			}
			// 2.其次，检查没有超时，但是不可用的连接，从连接池删除
			if isOk, err = t.CheckSpecialTcpConnIsActive(conn); !isOk {
				t.conns.Remove(tcpConn)
				t.count--
			}
		}
	}
	return err
}

// CheckSpecialTcpConnIsActive 检查特定的  tcp 连接是否有效
// @tcpConn  tcp连接
func (t *tcpConnPool) CheckSpecialTcpConnIsActive(tcpConn *tcpConnBaseInfo) (bool, error) {
	tmpHeader := &header{
		pkgLen: 0,
		cmd:    FDFS_PROTO_CMD_ACTIVE_TEST,
		status: 0,
	}
	err1 := tmpHeader.sendHeader(tcpConn)
	err2 := tmpHeader.receiveHeader(tcpConn)

	if err1 == nil && err2 == nil && tmpHeader.status == 0 {
		return true, nil
	} else {
		if tmpHeader.status != 0 {
			err3 := ERROR_TCP_SERVER_RESPONSE_NOT_ZERO
			return false, errors.New(err1.Error() + err2.Error() + err3)
		}
		return false, errors.New(err1.Error() + err2.Error())
	}
}

// CreateTcpConn 初始化一个tcp连接，主要用于连接到 fastdfs 的tracker server、storage server
func (t *tcpConnPool) CreateTcpConn() error {
	conn, err := net.DialTimeout("tcp", t.addrPort, TCP_CONN_TIMEOUT)
	if err != nil {
		return err
	}
	var tcpBaseInfo = &tcpConnBaseInfo{
		status: TCP_STATUS_UNINTERRUPTIBLE,
		Conn:   conn,
	}
	t.conns.PushBack(tcpBaseInfo)
	t.count++
	return nil
}

// get 获取一个 tcp 连接
func (t *tcpConnPool) get() (tcpConn *tcpConnBaseInfo, err error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	var isOk bool
	var okTcp *tcpConnBaseInfo
	for {
		conn := t.conns.Front()
		if conn == nil {
			if t.count > t.maxConns {
				err = errors.New(ERROR_CONN_POOL_OVER_MAX + strconv.Itoa(t.maxConns))
				return nil, err
			}
			if err = t.CreateTcpConn(); err == nil {
				continue
			} else {
				break
			}
		}
		// 获取一个 tcp 连接
		if okTcp, isOk = conn.Value.(*tcpConnBaseInfo); isOk {
			if isOk, err = t.CheckSpecialTcpConnIsActive(okTcp); isOk {
				t.conns.Remove(conn)
				// 取出后的 tcp 状态设置为 运行态
				okTcp.status = TCP_STATUS_RUNNING
				return okTcp, nil
			} else {
				t.conns.Remove(conn)
				continue
			}
		} else {
			return nil, errors.New(ERROR_TCP_CONN_ASSERT_FAIL)
		}
	}
	return nil, errors.New(ERROR_GET_TCP_CONN_FAIL + err.Error())
}

// put 将使用完毕的tcp连接放回连接池
// @tcpConn  tcp 连接
func (t *tcpConnPool) put(tcpConn *tcpConnBaseInfo) {
	t.lock.Lock()
	defer t.lock.Unlock()

	tcpConn.status = TCP_STATUS_INTERRUPTIBLE
	tcpConn.putTime = time.Now()
	t.conns.PushBack(tcpConn)
}
