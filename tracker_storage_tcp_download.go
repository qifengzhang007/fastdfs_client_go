package fastdfs_client_go

import (
	"errors"
	"fmt"
	"time"
)

//DownloadFileByFileId 通过文件id下载文件
// @fileId 文件id，格式：group1/M00/00/01/MeiRdmISSSqASUsRAJA3gI2KXVQ867.mp3
// @saveFileName 指定下载后保存的名称
// @offset 偏移字节数量
// @downloadBytes  指定需要下载的字节大小，超过此字节不会下载
func (c *trackerServerTcpClient) DownloadFileByFileId(fileId, saveFileName string) error {
	// 下载文件时，首先必须确保被下载的文件不存在
	if _, err := getFileInfoByFileName(saveFileName); err == nil {
		return errors.New(ERROR_FILE_DOWNLOAD_RELA_FILENAME_NOT_EMPTY)
	}
	groupName, remoteFilename, err := splitStorageServerFileId(fileId)
	if err != nil {
		return err
	}

	var tmpFileInfo *fileInfo
	//启动下载中断后、重试、续传功能，
	// 因为本套系统的应用场景主要还是内网，因此只考虑内网短暂的抖动导致异常情况去解决问题，设置最大尝试次数为 5，每次间隔5秒
	for i := 1; i <= 5; i++ {
		// 首先查询被下载文件已经缓存在硬盘的数据量（二进制大小）
		if i >= 2 && err != nil {
			if tmpFileInfo, err = getFileInfoByFileName(saveFileName); err == nil {
				offset := tmpFileInfo.fileSize
				fmt.Printf("下载出错了，开始重新偏移下载，偏移量：%d\n", offset)
				tmpFileInfo.Close()
				if err = c.startDownFiles(groupName, remoteFilename, saveFileName, offset); err != nil {
					time.Sleep(time.Second * 5)
				}
			}
		} else if i == 1 {
			fmt.Println("首次开始下载任务")
			err = c.startDownFiles(groupName, remoteFilename, saveFileName, 0)
		} else {
			fmt.Println("解除下载中出现的错误...断点续传（下载）完成！")
			break
		}
	}
	return err
}

// 下载文件主要逻辑
func (c *trackerServerTcpClient) startDownFiles(groupName, remoteFilename, saveFileName string, offset int64) error {
	storageServInfo, err := c.getStorageInfoByTracker(TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE, groupName, remoteFilename)
	if err != nil {
		return err
	}

	down := &storageDownloadHeaderBody{}
	down.groupName = groupName
	down.remoteFilename = remoteFilename
	down.offset = offset   //offset 偏移的字节数，设置0，从头开始下载
	down.downloadBytes = 0 // downloadBytes 需要下载的字节数，设置为0，下载整个文件
	down.saveFileName = saveFileName

	return c.sendCmdToStorageServer(down, storageServInfo)
}
