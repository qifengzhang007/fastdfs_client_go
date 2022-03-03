package fastdfs_client_go

import (
	"errors"
	"os"
)

// 文件结构信息（通过文件名上传、或者通过字节上传文件使用）
type fileInfo struct {
	filePtr     *os.File // 文件指针
	fileExtName string   // 文件扩展名
	fileSize    int64    // 文件大小
	buffer      []byte   // 文件字节（文件内容本身）
}

// 通过文件名获取文件信息
// @fileName 文件全路径名称
func getFileInfoByFileName(fileName string) (*fileInfo, error) {
	if fileName != "" {
		file, err := os.OpenFile(fileName, os.O_RDONLY, 0755)
		if err != nil {
			return nil, err
		}
		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}
		if int(stat.Size()) == 0 {
			return nil, errors.New(fileName + ERROR_FILE_SIZE_IS_ZERO)
		}

		return &fileInfo{
			fileSize:    stat.Size(),
			filePtr:     file,
			buffer:      nil,
			fileExtName: getFileExtNameStr(fileName),
		}, nil
	} else {
		return nil, errors.New(ERROR_FILE_FILENAME_IS_EMPTY)
	}
}

// 通过文件字节获取文件信息
func getFileInfoByFileByte(buffer []byte, fileExtName string) (*fileInfo, error) {
	if len(buffer) == 0 {
		return nil, errors.New(ERROR_FILE_SIZE_IS_ZERO)
	}
	if len(fileExtName) == 0 {
		return nil, errors.New(ERROR_FILE_EXT_NAME_IS_EMPTY)
	}
	return &fileInfo{
		filePtr:     nil,
		fileSize:    int64(len(buffer)),
		buffer:      buffer,
		fileExtName: fileExtName,
	}, nil
}

// 关闭文件
func (c *fileInfo) Close() {
	if c.filePtr != nil {
		_ = c.filePtr.Close()
	}
	return
}
