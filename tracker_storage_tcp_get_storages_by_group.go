package fastdfs_client_go

// GetStorageServersByGroup  获取group下的所有storage server信息
func (c *trackerServerTcpClient) GetStorageServersByGroup(groupName string) (resStorageServer []StorageServer, err error) {

	trackerSendParams := &trackerTcpConn{}
	// 将命令参数设置在 header 头部分
	trackerSendParams.header.cmd = TRACKER_PROTO_CMD_SERVER_LIST_STORAGE
	trackerSendParams.header.status = 0
	trackerSendParams.groupName = groupName
	trackerSendParams.remoteFilename = ""
	// 初始化 StorageServers 切片，防止在 Receive 方法中访问未初始化的切片元素
	trackerSendParams.StorageServers = make([]StorageServer, 0)

	if err = c.sendHeaderByTrackerServer(trackerSendParams); err != nil {
		return nil, err
	}

	return trackerSendParams.StorageServers, nil
}
