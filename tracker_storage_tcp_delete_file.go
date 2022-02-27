package fastdfs_client_go

//DeleteFile 删除存储在  storage server 服务器的文件
// @fileId storage server 服务器存储的完整文件名，例如：group1/M00/00/01/MeiRdmIQ7N-AOJmJAAAAD1y-ivQ067.png
func (c *trackerServerTcpClient) DeleteFile(fileId string) error {
	groupName, remoteFilename, err := splitStorageServerFileId(fileId)
	if err != nil {
		return err
	}

	storageServInfo, err := c.getStorageInfoByTracker(TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE, groupName, remoteFilename)
	if err != nil {
		return err
	}

	del := &storageDeleteHeaderBody{}
	del.groupName = groupName
	del.remoteFilename = remoteFilename

	return c.sendCmdToStorageServer(del, storageServInfo)
}
