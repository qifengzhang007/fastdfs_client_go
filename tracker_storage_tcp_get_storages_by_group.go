package fastdfs_client_go

// GetStorageServersByGroup  获取group下的所有storage server信息
func (c *trackerServerTcpClient) GetStorageServersByGroup(groupName string) (resStorageServer []StorageServer, err error) {

	trackerSendParams := &trackerTcpConn{}
	// 将命令参数设置在 header 头部分
	trackerSendParams.header.cmd = TRACKER_PROTO_CMD_SERVER_LIST_STORAGE
	trackerSendParams.header.status = 0
	trackerSendParams.groupName = groupName
	trackerSendParams.remoteFilename = ""

	if err = c.sendHeaderByTrackerServer(trackerSendParams); err != nil {
		return nil, err
	}

	return trackerSendParams.StorageServers, nil
}
