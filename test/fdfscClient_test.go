package test

import (
	"github.com/qifengzhang007/fastdfs_client_go"
	"strconv"
	"sync"
	"testing"
)

var conf = &fastdfs_client_go.TrackerStorageServerConfig{
	// 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
	TrackerServer: []string{"192.168.10.10:22122"},
	MaxConns:      128,
}

// 设置测试文件的根目录，测试使用
//var curDir = "E:/Project/2020/fastdfs_client_go/"
//var fileName = "1024.txt"

var curDir = "E:/音乐资源/"
var fileName = "15 俩俩相忘.mp3" // 9M 左右

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
	//fileId := "group1/M00/00/01/MeiRdmISCRuASW3BABrBPcp7oMo520.jpg"
	fileId := "group1/M00/00/01/MeiRdmISDUiAaURaAsRMrFnLJoE317.wav" // 大小 9451392
	if err = fdfsClient.DownloadFileByFileId(fileId, curDir+"下载测试-橄榄树.wav"); err != nil {
		t.Error("下载文件单元测试出错, ERROR:" + err.Error())
	} else {
		t.Log("下载文件单元测试成功 !")
	}
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
