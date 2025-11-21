package test

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/qifengzhang007/fastdfs_client_go"
)

var conf = &fastdfs_client_go.TrackerStorageServerConfig{
	// 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
	// 1.配置 trackerServer 地址，端口默认为：22122
	// 2. trackerServer 服务器会返回storage_server 服务器地址： xx.xx.xx.xx: 23000，
	// 3.因此如果是外网测试，请保证 trackerServer 服务器和 storage_server 服务器的ip、端口都能访问到
	// 4.上线部署以后，请使用内网ip、端口，保证安全
	// 5.IPV4地址：192.168.10.10:22122，IPV6示例： [2402::xxxx:c032:xxxx:xxxx:xxxx:d44e:0]:22122 ，注意： 在云服务器商控制面板开通互联网可访问的IPV6地址才行。
	TrackerServer: []string{"[2402::xxxx:c032:xxxx:xxxx:xxxx:d44e:0]:22122"},
	// tcp 连接池最大允许的连接数（trackerServer 和 storageServer 连接池共用该参数）
	MaxConns: 128,
}

// 设置测试文件的根目录，测试使用
//var curDir = "E:/Project/2020/fastdfs_client_go/"
//var fileName = "1024.txt"

// var curDir = "F:/BaiduNetdiskDownload/MySQL高级/"
var curDir = "E:/tmp/"
var fileName = "20251114-001.txt" // 28M 左右

// 通过文件名上传文件
func TestUploadByFileName(t *testing.T) {
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	if err != nil {
		t.Log("单元测试失败，创建TCP连接出错：" + err.Error())
		return
	}
	defer fdfsClient.Destroy()
	fileId, err := fdfsClient.UploadByFileName(curDir+fileName, 1)
	if err != nil {
		t.Errorf("单元测试失败，上传文件出错：%s", err.Error())
		return
	} else {
		t.Logf("单元测试成功，成功上传文件：%s", fileId)
	}
}

// 通过字节集上传文件
func TestUploadByBytes(t *testing.T) {
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	if err != nil {
		t.Log("单元测试失败，创建TCP连接出错：" + err.Error())
		return
	}
	defer fdfsClient.Destroy()
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(no int) {
			defer wg.Done()
			if fileId, err := fdfsClient.UploadByBuffer([]byte(strconv.Itoa(no+1)+" - 字节集直接上传"), "txt"); err != nil {
				t.Error("通过字节集文件流上传文件失败, ERROR:" + err.Error())
			} else {
				t.Log("通过字节集文件流上传文件成功！文件名：" + fileId)
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
	fileId := "group1/M00/00/00/ZzOT4WkXR6uAODW7AAAAN7EZJzQ965.txt"
	//fileId := "group1/M00/00/01/MeiRdmISDUiAaURaAsRMrFnLJoE317.wav" // 大小 9451392
	if err = fdfsClient.DownloadFileByFileId(fileId, curDir+"下载-demo-003.txt"); err != nil {
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
	fileId := "group1/M00/00/00/ZzOT4WkXSCqAPtfhAAAAGVSr950529.txt"
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
	fileId := "group1/M00/00/00/ZzOT4WkXR6uAODW7AAAAN7EZJzQ965.txt"
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
	serverAppendFileName := "M00/00/00/ZzOT4WkXZg2ETKwCAAAAAG7DoII745.txt" // 删除group名以后的 append文件名
	localFileName := "E:/tmp/20251114-001.txt"                             // 客户端文件名
	if err = fdfsClient.UploadAppendFileByFileName(serverAppendFileName, localFileName); err != nil {
		t.Error("单元测试失败，上传append文件出错：" + err.Error())
	} else {
		t.Logf("单元测试完成，上传append文件完成")
	}
}

// 基于已经存在的append文件名，上传追加文件 - 通过字节集数据追加
func TestUploadAppendFileByBuffer(t *testing.T) {
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	if err != nil {
		t.Error("单元测试失败，创建TCP连接出错：" + err.Error())
		return
	}
	defer fdfsClient.Destroy()

	// 一个文件在服务器的完整路径为：/group1/M00/00/00/cnQ3KGkTXAKAZWCPAbnLqPIYSzQ544.log
	serverAppendFileName := "M00/00/00/ZzOT4WkXZpqAB17hAAAAb8KVDIE745.txt" // 删除group名以后的 append文件名
	buffer := []byte("\r\nappend - 上传字节集 - 追加的内容-222")
	if err = fdfsClient.UploadAppendFileByBuffer(serverAppendFileName, buffer); err != nil {
		t.Error("单元测试失败，上传append-字节集文件出错：" + err.Error())
	} else {
		t.Logf("单元测试完成，上传append-字节集文件完成")
	}
}

// 转换append类型的文件为普通文件
func TestConvAppendFileToRegularFile(t *testing.T) {
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	if err != nil {
		t.Error("单元测试失败，创建TCP连接出错：" + err.Error())
		return
	}
	defer fdfsClient.Destroy()

	// 一个文件在服务器的完整路径为：/group1/M00/00/00/ZzOT4WkXSCqALRYbAAAAGT9wmzs044.log
	serverAppendFileName := "group1/M00/00/00/ZzOT4WkXZg2ETKwCAAAAAG7DoII745.txt" // 删除group名以后的 append文件名
	if newFileId, err := fdfsClient.ConvAppendFileToRegularFile(serverAppendFileName); err != nil {
		t.Error("单元测试失败，转换append类型的文件为普通文件出错：" + err.Error())
	} else {
		t.Logf("单元测试完成，转换append类型的文件为普通文件完成，新文件ID：%s", newFileId)
	}
}

// 生成fastdfs 资源访问token
func TestGetAccessToken(t *testing.T) {
	// 注意：token生成时应该使用不包含 group 的文件ID
	fileId := "M00/00/00/ZzOT4WkW2XGAQOgAAAAAUjEe6oc012.txt"
	//有关token设置的配置文件路径：
	// /etc/fdfs/http.conf ->
	//   > http.anti_steal.secret_key  =  设置加密key ，长度不要超过128字节
	// /etc/fdfs/mod_fastdfs.conf ->
	// > group_name = group1
	// >url_have_group_name = true

	// 闭坑指南 - 服务器时区与客户端时区不一致时，会导致token验证失败
	//# 设置时区为北京时间
	//timedatectl set-timezone Asia/Shanghai
	//echo "Asia/Shanghai" > /etc/timezone
	//
	//# 验证
	//echo "当前时区: $(date +%Z)"
	//当前时区: CST  // 注意： 时区为 CST ，表示东八区时间
	//timedatectl status | grep "Time zone"
	//Time zone: Asia/Shanghai (CST, +0800)

	secretKey := "your_secret_key" // 设置加密key
	var timestamp int64 = time.Now().Unix()
	token, ts, err := fastdfs_client_go.GetAccessToken(fileId, secretKey, timestamp)
	if err != nil {
		t.Error("单元测试失败，生成token出错：" + err.Error())
	} else {
		t.Logf("单元测试完成，生成token：%s, timestamp：%d, err：%v", token, ts, err)
		// 输出完整的访问URL示例
		t.Logf("完整访问URL示例：http://103.xxx.xxx.xxx:端口/group1/%s?token=%s&ts=%d", fileId, token, ts)
	}
}
