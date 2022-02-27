package fastdfs_client_go

import (
	"errors"
)

//DownloadFileByFileId 通过文件id下载文件
// @fileId 文件id，格式：group1/M00/00/01/MeiRdmISSSqASUsRAJA3gI2KXVQ867.mp3
// @saveFileName 指定下载后保存的名称
// @offset 偏移字节数量
// @downloadBytes  指定需要下载的字节大小，超过此字节不会下载
func (c *trackerServerTcpClient) DownloadFileByFileId(fileId string, saveFileName string) error {
	// 下载文件时，首先必须确保被下载的文件不存在
	if _, err := getFileInfoByFileName(saveFileName); err == nil {
		return errors.New(ERROR_FILE_DOWNLOAD_RELA_FILENAME_NOT_EMPTY)
	}
	groupName, remoteFilename, err := splitStorageServerFileId(fileId)
	if err != nil {
		return err
	}
	storageServInfo, err := c.getStorageInfoByTracker(TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE, groupName, remoteFilename)
	if err != nil {
		return err
	}

	down := &storageDownloadHeaderBody{}
	down.groupName = groupName
	down.remoteFilename = remoteFilename
	down.offset = 0        //offset 偏移的字节数，设置0，从头开始下载
	down.downloadBytes = 0 // downloadBytes 需要下载的字节数，设置为0，下载整个文件
	down.saveFileName = saveFileName

	return c.sendCmdToStorageServer(down, storageServInfo)
}
