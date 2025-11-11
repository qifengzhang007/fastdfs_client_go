package fastdfs_client_go

import "errors"

// GetGroupInfo  获取远程服务器的文件信息
// @remoteFileId 远程服务器的文件Id
func (c *trackerServerTcpClient) GetGroupInfo(groupName string) (resGroupInfo GroupInfo, err error) {

	if len(groupName) < 1 {
		return resGroupInfo, errors.New("groupName 不能为空")
	}

	storageServInfo, err := c.getStorageInfoByTracker(TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ONE, groupName, "")
	if err != nil {
		return resGroupInfo, err
	}
	queryGroupInfo := &storageGetGroupInfoHeaderBody{}
	// 请求时需要的业务参数赋值
	queryGroupInfo.groupName = groupName

	if err = c.sendCmdToStorageServer(queryGroupInfo, storageServInfo); err != nil {
		return resGroupInfo, err
	}
	resGroupInfo.respGroupName = queryGroupInfo.respGroupName
	resGroupInfo.resptotalMb = queryGroupInfo.resptotalMb
	resGroupInfo.respFreeMb = queryGroupInfo.respFreeMb
	resGroupInfo.respReservedMb = queryGroupInfo.respReservedMb
	resGroupInfo.resptrunkFreeMb = queryGroupInfo.resptrunkFreeMb
	resGroupInfo.respServerCount = queryGroupInfo.respServerCount
	resGroupInfo.respServerPort = queryGroupInfo.respServerPort
	resGroupInfo.respReadableServerCount = queryGroupInfo.respReadableServerCount
	resGroupInfo.respWritableServerCount = queryGroupInfo.respWritableServerCount
	resGroupInfo.respCurrentWriteServerCount = queryGroupInfo.respCurrentWriteServerCount
	resGroupInfo.respstorePathCount = queryGroupInfo.respstorePathCount
	resGroupInfo.respsubdirCountPerPath = queryGroupInfo.respsubdirCountPerPath
	resGroupInfo.respcurrentTrunkFileId = queryGroupInfo.respcurrentTrunkFileId

	return resGroupInfo, nil
}
