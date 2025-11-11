package fastdfs_client_go

import "errors"

// GetGroupInfo  获取远程服务器的文件信息
// @remoteFileId 远程服务器的文件Id
func (c *trackerServerTcpClient) GetGroupInfo(groupName string) (ResGroupInfo GroupInfo, err error) {

	if len(groupName) < 1 {
		return ResGroupInfo, errors.New("groupName 不能为空")
	}

	storageServInfo, err := c.getStorageInfoByTracker(TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ONE, groupName, "")
	if err != nil {
		return ResGroupInfo, err
	}
	queryGroupInfo := &storageGetGroupInfoHeaderBody{}
	// 请求时需要的业务参数赋值
	queryGroupInfo.groupName = groupName

	if err = c.sendCmdToStorageServer(queryGroupInfo, storageServInfo); err != nil {
		return ResGroupInfo, err
	}
	ResGroupInfo.respGroupName = queryGroupInfo.respGroupName
	ResGroupInfo.resptotalMb = queryGroupInfo.resptotalMb
	ResGroupInfo.respFreeMb = queryGroupInfo.respFreeMb
	ResGroupInfo.respReservedMb = queryGroupInfo.respReservedMb
	ResGroupInfo.resptrunkFreeMb = queryGroupInfo.resptrunkFreeMb
	ResGroupInfo.respServerCount = queryGroupInfo.respServerCount
	ResGroupInfo.respServerPort = queryGroupInfo.respServerPort
	ResGroupInfo.respReadableServerCount = queryGroupInfo.respReadableServerCount
	ResGroupInfo.respWritableServerCount = queryGroupInfo.respWritableServerCount
	ResGroupInfo.respCurrentWriteServerCount = queryGroupInfo.respCurrentWriteServerCount
	ResGroupInfo.respstorePathCount = queryGroupInfo.respstorePathCount
	ResGroupInfo.respsubdirCountPerPath = queryGroupInfo.respsubdirCountPerPath
	ResGroupInfo.respcurrentTrunkFileId = queryGroupInfo.respcurrentTrunkFileId

	return ResGroupInfo, nil
}
