package test

import (
	"github.com/qifengzhang007/fastdfs_client_go"
	"strconv"
	"sync"
	"testing"
	"time"
)

var conf = &fastdfs_client_go.TrackerStorageServerConfig{
	// 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
	TrackerServer: []string{"49.232.145.118:22122"},
	// tcp 连接池最大允许的连接数（trackerServer 和 storageServer 连接池共用该参数）
	MaxConns: 128,
}

// 设置测试文件的根目录，测试使用
//var curDir = "E:/Project/2020/fastdfs_client_go/"
//var fileName = "1024.txt"

var curDir = "F:/BaiduNetdiskDownload/MySQL高级/"
var fileName = "mysql-8.0.18.tar.gz" // 9M 左右

// 通过文件名上传文件
func TestUploadByFileName(t *testing.T) {
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	if err != nil {
		t.Log("单元测试失败，创建TCP连接出错：" + err.Error())
		return
	}
	defer fdfsClient.Destroy()
	fileId, err := fdfsClient.UploadByFileName(curDir + fileName)
	if err != nil {
		t.Errorf("单元测试失败，上传文件出错：%s", err.Error())
		return
	} else {
		t.Logf("单元测试成功，成功上传文件：%s", fileId)
	}
}

// 通过二进制上传文件
func TestUploadByBytes(t *testing.T) {
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	if err != nil {
		t.Log("单元测试失败，创建TCP连接出错：" + err.Error())
		return
	}
	defer fdfsClient.Destroy()
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(no int) {
			defer wg.Done()
			if fileId, err := fdfsClient.UploadByBuffer([]byte(strconv.Itoa(no+1)+" - 二进制直接上传"), "txt"); err != nil {
				t.Error("通过二进制文件流上传文件失败, ERROR:" + err.Error())
			} else {
				t.Log("通过二进制文件流上传文件成功！文件名：" + fileId)
			}
		}(i)
	}
	wg.Wait()
}

// 下载文件测试
func TestDownLoadFile(t *testing.T) {
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	if err != nil {
		t.Log("单元测试失败，创建TCP连接出错：" + err.Error())
		return
	}
	defer fdfsClient.Destroy()
	// 通过指定 文件id 下载文件
	fileId := "group1/M00/00/01/MeiRdmIbNaKATp-GAJA3gI2KXVQ428.mp3"
	//fileId := "group1/M00/00/01/MeiRdmISDUiAaURaAsRMrFnLJoE317.wav" // 大小 9451392
	if err = fdfsClient.DownloadFileByFileId(fileId, curDir+"音乐001.mp3"); err != nil {
		t.Error("下载文件单元测试出错, ERROR:" + err.Error())
	} else {
		t.Log("下载文件单元测试成功 !")
	}

	time.Sleep(time.Minute * 5)
}

// 删除文件
func TestDeleteFile(t *testing.T) {
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	if err != nil {
		t.Error("单元测试失败，创建TCP连接出错：" + err.Error())
		return
	}
	defer fdfsClient.Destroy()
	// 通过指定 文件id(fileId) 删除文件
	fileId := "group1/M00/00/01/MeiRdmISSbuAZwwSAAAAD_Q4O2U879.txt"
	if err = fdfsClient.DeleteFile(fileId); err != nil {
		t.Error("单元测试失败，删除文件出错：" + err.Error())
	} else {
		t.Log("删除文件 - 单元测试成功!")
	}
}

// 查询远程文件信息
func TestQueryRemoteFileInfo(t *testing.T) {
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	if err != nil {
		t.Error("单元测试失败，创建TCP连接出错：" + err.Error())
		return
	}
	defer fdfsClient.Destroy()
	// 通过指定 文件id(fileId) 删除文件
	fileId := "group1/M00/00/00/MeiRdmIOFHuAWrOVAAdps6fU_X8506.png"
	if remoteFileInfo, err := fdfsClient.GetRemoteFileInfo(fileId); err != nil {
		t.Error("单元测试失败，删除文件出错：" + err.Error())
	} else {
		t.Logf("远程文件查询结果：%#+v\n", remoteFileInfo)
	}
}
