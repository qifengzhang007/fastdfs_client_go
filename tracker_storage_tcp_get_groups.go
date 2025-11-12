package fastdfs_client_go

// GetGroups  获取所有group信息
func (c *trackerServerTcpClient) GetGroups() (resGroupsInfo []GroupInfo, err error) {

	trackerSendParams := &trackerTcpConn{}
	// 将命令参数设置在 header 头部分
	trackerSendParams.header.pkgLen = 10
	trackerSendParams.header.cmd = TRACKER_PROTO_CMD_SERVER_LIST_ALL_GROUPS
	trackerSendParams.header.status = 0
	trackerSendParams.groupName = ""
	trackerSendParams.remoteFilename = ""

	if err = c.sendHeaderByTrackerServer(trackerSendParams); err != nil {
		return nil, err
	}

	return trackerSendParams.groups, nil
}
