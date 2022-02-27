package fastdfs_client_go

// UploadByFileName  创建fdfs客户端, 上传文件
// @fileName 指定已经存在的文件名
func (c *trackerServerTcpClient) UploadByFileName(fileName string) (string, error) {
	file, err := getFileInfoByFileName(fileName)
	defer file.Close()
	if err != nil {
		return "", err
	}

	storageServInfo, err := c.getStorageInfoByTracker(TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE, "", "")
	if err != nil {
		return "", err
	}

	uploadServ := &storageServerUploadHeaderBody{}
	uploadServ.fileInfo = file
	uploadServ.storagePathIndex = storageServInfo.storagePathIndex

	if err = c.sendCmdToStorageServer(uploadServ, storageServInfo); err != nil {
		return "", err
	}
	return uploadServ.fileId, nil
}

// UploadByBuffer  创建fdfs客户端, 上传文件
// @buffer 二进制数据,适合小文件一次性发送
// @fileExtName 指定文件在服务器端保存时的文件名
func (c *trackerServerTcpClient) UploadByBuffer(buffer []byte, fileExtName string) (string, error) {
	tmpFileInfo, err := getFileInfoByFileByte(buffer, fileExtName)
	defer tmpFileInfo.Close()
	if err != nil {
		return "", err
	}
	storageServInfo, err := c.getStorageInfoByTracker(TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE, "", "")
	if err != nil {
		return "", err
	}

	uploadServ := &storageServerUploadHeaderBody{}
	uploadServ.fileInfo = tmpFileInfo
	uploadServ.storagePathIndex = storageServInfo.storagePathIndex

	if err = c.sendCmdToStorageServer(uploadServ, storageServInfo); err != nil {
		return "", err
	}
	return uploadServ.fileId, nil
}
