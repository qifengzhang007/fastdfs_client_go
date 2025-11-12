package fastdfs_client_go

import "errors"

// GetGroupInfo  指定 group 名称获取 group 信息
// @groupName 组名，例如： group1
func (c *trackerServerTcpClient) GetGroupInfo(groupName string) (resGroupInfo *GroupInfo, err error) {
	if len(groupName) < 1 {
		return resGroupInfo, errors.New(ERROR_GET_GROUP_INFO_FAILED_GROUP_NOT_EMPTY)
	}
	resGroupInfo, err = c.getGroupInfoByTracker(TRACKER_PROTO_CMD_SERVER_LIST_ONE_GROUP, groupName)
	if err != nil {
		return resGroupInfo, err
	}
	return resGroupInfo, nil
}

// getGroupInfoByTracker  主要通过 tracker server 获取 group 组信息
// @groupName 组名 (body 参数)
func (c *trackerServerTcpClient) getGroupInfoByTracker(cmd byte, groupName string) (*GroupInfo, error) {
	trackerSendParams := &trackerTcpConn{}
	// 将命令参数设置在 header 头部分
	trackerSendParams.header.pkgLen = 16
	trackerSendParams.header.cmd = cmd
	trackerSendParams.header.status = 0
	trackerSendParams.groupName = groupName
	trackerSendParams.remoteFilename = ""

	if err := c.sendHeaderByTrackerServer(trackerSendParams); err != nil {
		return nil, err
	}
	return &GroupInfo{
		groupName:               trackerSendParams.groupInfo.groupName,
		totalMb:                 trackerSendParams.groupInfo.totalMb,
		freeMb:                  trackerSendParams.groupInfo.freeMb,
		reservedMb:              trackerSendParams.groupInfo.reservedMb,
		trunkFreeMb:             trackerSendParams.groupInfo.trunkFreeMb,
		serverCount:             trackerSendParams.groupInfo.serverCount,
		serverPort:              trackerSendParams.groupInfo.serverPort,
		readableServerCount:     trackerSendParams.groupInfo.readableServerCount,
		writableServerCount:     trackerSendParams.groupInfo.writableServerCount,
		currentWriteServerCount: trackerSendParams.groupInfo.currentWriteServerCount,
		storePathCount:          trackerSendParams.groupInfo.storePathCount,
		subdirCountPerPath:      trackerSendParams.groupInfo.subdirCountPerPath,
		currentTrunkFileId:      trackerSendParams.groupInfo.currentTrunkFileId,
	}, nil
}
