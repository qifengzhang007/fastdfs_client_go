package fastdfs_client_go

// ConvAppendFileToRegularFile  获取远程服务器的文件信息
// @remoteAppendFileId 远程服务器的append文件Id, 需要完整的文件Id，示例：group1/M00/00/00/ZzOT4WkXSCqALRYbAAAAGT9wmzs044.txt
func (c *trackerServerTcpClient) ConvAppendFileToRegularFile(remoteAppendFileId string) (newFileId string, err error) {
	groupName, remoteAppendFileId, err := splitStorageServerFileId(remoteAppendFileId)
	if err != nil {
		return "", err
	}

	storageServInfo, err := c.getStorageInfoByTracker(TRACKER_PROTO_CMD_SERVICE_QUERY_UPDATE, groupName, remoteAppendFileId)
	if err != nil {
		return "", err
	}
	resAppendFileInfo := &storageConvAppendFileHeaderBody{}
	resAppendFileInfo.appenderFilename = remoteAppendFileId

	if err = c.sendCmdToStorageServer(resAppendFileInfo, storageServInfo); err != nil {
		return "", err
	}

	return resAppendFileInfo.groupName + "/" + resAppendFileInfo.filename, nil
}
