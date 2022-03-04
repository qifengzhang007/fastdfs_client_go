package fastdfs_client_go

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

//IntToBytes 整形转换成字节
// 注意： 8字节整数整形转换为 字节
// @n 需要转换的整数
//func IntToBytes(n int64) []byte {
//	bytesBuffer := bytes.NewBuffer([]byte{})
//	_ = binary.Write(bytesBuffer, binary.BigEndian, n)
//	return bytesBuffer.Bytes()
//}

//bytesToInt 字节转换成整形
// @bys 需要转换的字节
func bytesToInt(bys []byte) int64 {
	bytesBuffer := bytes.NewBuffer(bys)
	// 注意：这里转换的结果是 ： 8字节整数
	var x int64
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return x
}

// getBytesByPosition  截取指定长度的字节切片
// @bys 原始字节切片
// @start 开始位置，
// @num 截取的字节数目
func getBytesByPosition(bys []byte, start, num int) []byte {
	var newBytes = bys[start:]
	endPosition := bytes.IndexByte(newBytes, 0x0)
	if endPosition > 0 {
		num = endPosition
	}
	return newBytes[:num]
}

// getFileExtNameStr  获取文件扩展名的文本
// @fileName 文件的完整名
func getFileExtNameStr(fileName string) string {
	index := strings.LastIndexByte(fileName, '.')
	if index != -1 {
		return fileName[index+1:]
	}
	return ""
}

// specialFileExtNameConvBytes  指定的文件扩展名转换为二进制
// @specialFileExtName 指定文件的扩展名，例如：png、tar.gz，最开始位置不包括点(.)
func specialFileExtNameConvBytes(specialFileExtName string) []byte {
	var fileExtName = make([]byte, FILE_EXTNAME_FIX_LEN)
	copy(fileExtName, specialFileExtName[:])
	return fileExtName
}

//groupNameConvBytes    storage server  组名转为二进制
// @groupName  storage server 文件存储组名
func groupNameConvBytes(groupName string) []byte {
	var gName = make([]byte, FDFS_GROUP_NAME_FIX_LEN)
	copy(gName, groupName[:])
	return gName
}

//sendBytesByFilePtr 通过文件指针以及tcp连接将数据发送出去
// 注释事项：这里对于文件指针只是使用，开发者必须在文件打开的地方，记得调用 defer 关闭
// @filePtr 文件指针
// @tcpConn tcp 连接
func sendBytesByFilePtr(filePtr *os.File, tcpConn net.Conn) (int64, error) {
	fInfo, err := filePtr.Stat()
	if err != nil {
		return 0, err
	}
	totalSize := fInfo.Size()
	var oneReadBuf = make([]byte, TCP_READ_BUFFER_SIZE)
	var remainSize int64 = 1
	var actualReadSize = 0
	var alreadySendSize int64 = 0
	_ = tcpConn.SetWriteDeadline(time.Time{})
	for ; remainSize > 0; remainSize = totalSize - alreadySendSize {
		if actualReadSize, err = filePtr.Read(oneReadBuf); err != nil {
			if err == io.EOF {
				return alreadySendSize, err
			} else {
				return alreadySendSize, err
			}
		} else {
			alreadySendSize += int64(actualReadSize)
			//fmt.Printf("tcp发送的数据: %#+v\n", oneReadBuf[:actualReadSize])
			if _, tcpErr := tcpConn.Write(oneReadBuf[:actualReadSize]); tcpErr != nil {
				return alreadySendSize, tcpErr
			}
		}
	}
	return alreadySendSize, nil
}

//writeBufferFromTcpConn 从tcp连接的内核缓冲区把数据写到内存缓冲区
// @conn tcp连接
// @writer 内存缓冲区io
// @totalSize 待接收的二进制总长度
func writeBufferFromTcpConn(conn net.Conn, bufWriter *bufio.Writer, totalSize int64) (err error) {
	// 单次读取的字节
	var oneReadSize int
	// 假设有剩余字节需要读取
	var remainSize int64 = 1
	// 已经接受到的字节
	var alreadyReceivedSize int64 = 0
	// 初始化一个临时存储缓冲区
	buf := make([]byte, TCP_RECEIVE_BUFFER_SIZE)
	_ = conn.SetReadDeadline(time.Time{})
	var i = 1 //记录数据读取的次数
	for ; remainSize > 0; remainSize = totalSize - alreadyReceivedSize {
		i++
		//最大每次读取 4096 个字节
		if remainSize > TCP_RECEIVE_BUFFER_SIZE {
			remainSize = TCP_RECEIVE_BUFFER_SIZE
		}
		//fmt.Printf("从tcp内核缓冲区读取的字节数目：%d, 二进制：%#+v\n", oneReadSize, buf[:remainSize])
		oneReadSize, err = conn.Read(buf[:remainSize])
		if err != nil {
			return err
		}
		// 从 tcp 内核的缓冲区写入开发者定义的接受变量对应的内存缓冲区
		_, err = bufWriter.Write(buf[:oneReadSize])
		if err != nil {
			return err
		}
		alreadyReceivedSize += int64(oneReadSize)
		// 假设每次从 tcp 内核缓冲区读取的内容都是最大值 4096 字节，那么 1000 次大概是 4M 左右
		// 每隔 4M 左右将内存缓冲区数据写入到底层的硬盘, 确保大文下载件时，内存占用始终处于低位
		if i%1000 == 0 {
			//fmt.Printf("%d - 每隔大约 4M 左右的数据，刷新到硬盘，已经接受的字节量：%d\n", i, alreadyReceivedSize)
			if err = bufWriter.Flush(); err != nil {
				return errors.New(ERROR_STORAGE_SERVER_DOWN_FILE_WRITE_FLUSH + err.Error())
			}
		}
	}
	if err = bufWriter.Flush(); err != nil {
		return errors.New(ERROR_STORAGE_SERVER_DOWN_FILE_WRITE_FLUSH + err.Error())
	} else {
		return nil
	}
}

// splitStorageServerFileId 分割 group 和 文件存储路径使用
// @fileId 服务器存储的文件Id （fileId）
func splitStorageServerFileId(fileId string) (string, string, error) {
	pos1 := strings.IndexByte(fileId, '/')
	// 如果文件的id(fileId) 以 / 开头, 就把最开始的 / 删除
	if pos1 == 0 {
		fileId = fileId[1:]
	}
	str := strings.SplitN(fileId, "/", 2)
	if len(str) != 2 {
		return "", "", errors.New(ERROR_STORAGE_SERVER_FILE_NAME_FORMAT2)
	}
	return str[0], str[1], nil
}
