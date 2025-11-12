package fastdfs_client_go

import (
	"errors"
	"os"
)

// UploadByFileName  创建fdfs客户端, 上传文件
// @fileName 指定已经存在的文件名
// @fileType 文件类型，默认值为 0,表示普通文件， 1 = append 追加文件
func (c *trackerServerTcpClient) UploadByFileName(fileName string, fileType ...int) (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return "", errors.New(ERROR_FILE_NOT_FOUND + " : " + fileName)
	}
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
	uploadServ.uploadFileType = STORAGE_PROTO_CMD_UPLOAD_FILE // 默认上传的是普通文件
	if len(fileType) > 0 {
		if fileType[0] == 1 {
			uploadServ.uploadFileType = STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE // 如果设置了文件类型参数为 1 ，则上传的是追加文件
		}
	}

	if err = c.sendCmdToStorageServer(uploadServ, storageServInfo); err != nil {
		return "", err
	}
	return uploadServ.fileId, nil
}

// UploadByBuffer  创建fdfs客户端, 上传文件
// @buffer 二进制数据,适合小文件一次性发送
// @fileExtName 指定文件在服务器端保存时的文件名
// @fileType 文件类型，默认值为 0,表示普通文件， 1 = append 追加文件
func (c *trackerServerTcpClient) UploadByBuffer(buffer []byte, fileExtName string, fileType ...int) (string, error) {
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
	uploadServ.uploadFileType = STORAGE_PROTO_CMD_UPLOAD_FILE // 默认上传的是普通文件
	if len(fileType) > 0 {
		if fileType[0] == 1 {
			uploadServ.uploadFileType = STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE // 如果设置了文件类型参数为 1 ，则上传的是追加文件
		}
	}
	if err = c.sendCmdToStorageServer(uploadServ, storageServInfo); err != nil {
		return "", err
	}
	return uploadServ.fileId, nil
}
