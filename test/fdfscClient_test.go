package test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/qifengzhang007/fastdfs_client_go"
)

var conf = &fastdfs_client_go.TrackerStorageServerConfig{
	// 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
	// 1.配置 trackerServer 地址，端口默认为：22122
	// 2. trackerServer 服务器会返回storage_server 服务器地址： xx.xx.xx.xx: 23000，
	// 3.因此如果是外网测试，请保证 trackerServer 服务器和 storage_server 服务器的ip、端口都能访问到
	// 4.上线部署以后，请使用内网ip、端口，保证安全
	TrackerServer: []string{"114.116.55.40:22122"},
	// tcp 连接池最大允许的连接数（trackerServer 和 storageServer 连接池共用该参数）
	MaxConns: 128,
}

// 设置测试文件的根目录，测试使用
//var curDir = "E:/Project/2020/fastdfs_client_go/"
//var fileName = "1024.txt"

// var curDir = "F:/BaiduNetdiskDownload/MySQL高级/"
var curDir = "E:/tmp/"
var fileName = "test-001.mp4" // 28M 左右

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
	fileId := "group1/M00/00/00/cnQ3KGkTXAKAZWCPAbnLqPIYSzQ544.mp4"
	//fileId := "group1/M00/00/01/MeiRdmISDUiAaURaAsRMrFnLJoE317.wav" // 大小 9451392
	if err = fdfsClient.DownloadFileByFileId(fileId, curDir+"下载-001.mp4"); err != nil {
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
	fileId := "group1/M00/00/00/cnQ3KGkTVyCAHxovAbnLqPIYSzQ745.mp4"
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
	// 通过指定 文件id(fileId) 查询远程文件信息
	fileId := "group1/M00/00/00/cnQ3KGkTXAKAZWCPAbnLqPIYSzQ544.mp4"
	if remoteFileInfo, err := fdfsClient.GetRemoteFileInfo(fileId); err != nil {
		t.Error("单元测试失败，查询远程文件信息出错：" + err.Error())
	} else {
		t.Logf("远程文件查询结果：%#+v\n", remoteFileInfo)
	}
}

// 获取服务器上所有的 groups 组列表信息
func TestQueryGetGroups(t *testing.T) {
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	if err != nil {
		t.Error("单元测试失败，创建TCP连接出错：" + err.Error())
		return
	}
	defer fdfsClient.Destroy()
	// 通过指定 组名(groupName) 查询组信息
	if groups, err := fdfsClient.GetGroups(); err != nil {
		t.Error("单元测试失败，查询组列表信息出错：" + err.Error())
	} else {
		t.Logf("单元测试完成，groups 组列表信息：%#+v\n", groups)
	}

}

// 查询指定的组(group)-查询组基本信息
func TestQueryGroupInfo(t *testing.T) {
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	if err != nil {
		t.Error("单元测试失败，创建TCP连接出错：" + err.Error())
		return
	}
	defer fdfsClient.Destroy()
	// 通过指定一个组名(groupName) 查询组信息
	groupName := "group1"
	if groupInfo, err := fdfsClient.GetGroupInfo(groupName); err != nil {
		t.Error("单元测试失败，查询组信息出错：" + err.Error())
	} else {
		t.Logf("单元测试完成，group基本信息：%#+v\n", groupInfo)
	}
}

// 查询指定的组(group)-查询组下的所有storage server信息
func TestQueryStorageServersFromGroup(t *testing.T) {
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	if err != nil {
		t.Error("单元测试失败，创建TCP连接出错：" + err.Error())
		return
	}
	defer fdfsClient.Destroy()
	// 通过指定 组名(groupName) 查询组下的所有storage server信息
	groupName := "group1"
	if storageServers, err := fdfsClient.GetStorageServersByGroup(groupName); err != nil {
		t.Error("单元测试失败，查询组下的所有storage server信息出错：" + err.Error())
	} else {
		t.Logf("单元测试完成，group下的所有storage server信息：%#+v\n", storageServers)
	}
}

// 基于已经存在的append文件名，上传追加文件
func TestUploadAppendFileByFileName(t *testing.T) {
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	if err != nil {
		t.Error("单元测试失败，创建TCP连接出错：" + err.Error())
		return
	}
	defer fdfsClient.Destroy()

	// 一个文件在服务器的完整路径为：/group1/M00/00/00/cnQ3KGkTXAKAZWCPAbnLqPIYSzQ544.log
	serverAppendFileName := "M00/00/00/cnQ3KGkTXAKAZWCPAbnLqPIYSzQ544.log" // 删除group名以后的 append文件名
	localFileName := "F:/tmp/123.log"                                      // 客户端文件名
	if err = fdfsClient.UploadAppendFileByFileName(serverAppendFileName, localFileName); err != nil {
		t.Error("单元测试失败，上传append文件出错：" + err.Error())
	} else {
		t.Logf("单元测试完成，上传append文件完成")
	}
}
