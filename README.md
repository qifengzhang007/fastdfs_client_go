## fastdfs_client_go

### 1.概述

- `FastDFS` 采用二进制 `TCP` 通信协议.  
- 开发本包的重点与核心主要是实现二进制通讯协议.  
- [点击了解 fastdfs 分布式文件存储系统](https://github.com/happyfish100/fastdfs) .

### 2.FastDFS二进制通讯协议细节

[点击查看GO 实现过程](./tcp_protocal_detail.md)

### 3.安装本包

```code 
// 请在本仓库的 gitTag 中查看最新版本,永远建议大家使用最新版本.  
// 查看地址：https://github.com/qifengzhang007/fastdfs_client_go/tags
go  get  github.com/qifengzhang007/fastdfs_client_go@v1.0.7

```

### 4. 已封装的函数列表

| 序号 | 函数 | 功能 | 适用版本说明|
|-----|-----|-----|----|
| 1 | UploadByFileName | 上传文件，生成普通文件ID或者append文件ID，append类型文件后续支持追加内容 | 所有版本|
| 2 | UploadByBuffer | 传递字节集（[]byte） 上传文件，生成普通文件ID或者append文件ID | 所有版本|
| 3 | DownloadFileByFileId | 根据文件ID下载文件 | 所有版本|
| 5 | GetRemoteFileInfo | 根据文件ID获取文件信息 |     所有版本|
| 4 | DeleteFile | 根据文件ID删除文件 | 所有版本|
| 6 | GetGroups | 获取所有组(groups)信息 | ≥ 6.13|
| 7 | GetGroupInfo | 获取一个特定组(group)信息 |  ≥ 6.13|
| 8 | GetStorageServersByGroup | 获取组(group)下的所有存储节点(storage server)信息 |  ≥ 6.13|




#### 4.1 上传文件(指定文件名)

```code  
    // 设置 trackerServer 配置参数
    var conf = &fastdfs_client_go.TrackerStorageServerConfig{
	    // 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
        TrackerServer: []string{"192.168.10.10:22122"},
        // tcp 连接池最大允许的连接数（trackerServer 和 storageServer 连接池共用该参数）
        MaxConns:      128,
    }
    # 文件上传核心函数
    fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
    // curDir + fileName  = 被上传文件的路径，例如：/home/dmmo/test-001.mp4
    //UploadByFileName 该函数第二个参数为0 表示上传普通文件 ； 为1表示append文件类型(后续支持追加内容)
    fileId, err := fdfsClient.UploadByFileName(curDir + fileName,0)  
    // 返回的 fileId 格式：
    // group1/M00/00/00/cnQ3KGkK8BOAUqfCANze3qa2RIk584.mp4

```

#### 4.2 上传文件(传递字节集 []byte)

```code  
    // 设置 trackerServer 配置参数
    var conf = &fastdfs_client_go.TrackerStorageServerConfig{
	    // 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
        TrackerServer: []string{"192.168.10.10:22122"},
        // tcp 连接池最大允许的连接数（trackerServer 和 storageServer 连接池共用该参数）
        MaxConns:      128,
    }
    # 文件上传核心函数
    fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
    // 直接传递二进制上传文件，适合文件比较小的场景使用
    // UploadByBuffer 该函数第二个参数为0 表示上传普通文件 ； 为1表示append文件类型(后续支持追加内容)
    fileId, err := fdfsClient.UploadByBuffer([]byte("测试文本数据转为二进制直接上传"),

```


#### 4.3 上传append文件 - 基于已有append文件二次(多次)上传追加文件
 -  调用本函数时，需要传递一个在服务端已存在的 append 类型的文件名, 例如：/group1/M00/00/00/cnQ3KGkTXAKAZWCPAbnLqPIYSzQ544.log, 该append文件可以通过 4.1、4.2 文件上传函数上传得到.  

```code
    // 设置 trackerServer 配置参数
    var conf = &fastdfs_client_go.TrackerStorageServerConfig{
    // 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
    TrackerServer: []string{"192.168.10.10:22122"},
    // tcp 连接池最大允许的连接数（trackerServer 和 storageServer 连接池共用该参数）
    MaxConns:      128,
    }
    fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
    
	// 一个文件在服务器的完整路径为：/group1/M00/00/00/cnQ3KGkTXAKAZWCPAbnLqPIYSzQ544.log
    serverAppendFileName := "M00/00/00/cnQ3KGkTXAKAZWCPAbnLqPIYSzQ544.log" // 删除 group 名以后的文件名
    localFileName := "F:/tmp/123.log"   // 客户端文件名(等待上传的文件名)
    err = fdfsClient.UploadAppendFileByFileName(serverAppendFileName, localFileName)
    
```

#### 4.4 文件下载

```code  

    // 设置 trackerServer 配置参数
    var conf = &fastdfs_client_go.TrackerStorageServerConfig{
	    // 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
        TrackerServer: []string{"192.168.10.10:22122"},
        // tcp 连接池最大允许的连接数（trackerServer 和 storageServer 连接池共用该参数）
        MaxConns:      128,
    }
    // 指定需要被下载的文件id （fileId）
    fileId := "group1/M00/00/01/MeiRdmISDUiAaURaAsRMrFnLJoE317.wav" // 大小 46419116，约 46M 左右
    // 创建 fdfs 客户端
    fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
    
    // 指定需要下载的文件id（fileId），最终的保存路径，开始下载
	fdfsClient.DownloadFileByFileId(fileId, "E:/音乐文件夹/下载测试-俩俩相忘.wav")

```

#### 4.5 获取远程文件信息

```code   
    // 设置 trackerServer 配置参数
    var conf = &fastdfs_client_go.TrackerStorageServerConfig{
	    // 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
        TrackerServer: []string{"192.168.10.10:22122"},
        // tcp 连接池最大允许的连接数（trackerServer 和 storageServer 连接池共用该参数）
        MaxConns:      128,
    }
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	// 指定需要查询的远程文件Id
	fileId := "group1/M00/00/01/MeiRdmISSbuAZwwSAAAAD_Q4O2U879.txt"
	// 查询远程文件信息，返回一个包含文件信息的结构体
	 remoteFileInfo, err := fdfsClient.GetRemoteFileInfo(fileId)
	 
```

#### 4.6 文件删除

```code   
    // 设置 trackerServer 配置参数
    var conf = &fastdfs_client_go.TrackerStorageServerConfig{
	    // 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
        TrackerServer: []string{"192.168.10.10:22122"},
        // tcp 连接池最大允许的连接数（trackerServer 和 storageServer 连接池共用该参数）
        MaxConns:      128,
    }
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	// 指定需要删除的文件Id
	fileId := "group1/M00/00/01/MeiRdmISSbuAZwwSAAAAD_Q4O2U879.txt"
	// 指定删除命令
	 err = fdfsClient.DeleteFile(fileId);
	 
	 //注意：删除文件后，多次删除仍然会提示成功，不会返回错误，如果您需要严格获取删除结果,您可以先调用远程查询命令确保文件存在，然后再执行删除命令。
	 
```


#### 5.1 查询 groups 信息
- 服务器存储文件首先是按照 group 组织的，每个 group 下可以有多个 storage server 节点.  
```code
    // 设置 trackerServer 配置参数
    var conf = &fastdfs_client_go.TrackerStorageServerConfig{
    // 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
    TrackerServer: []string{"192.168.10.10:22122"},
    // tcp 连接池最大允许的连接数（trackerServer 和 storageServer 连接池共用该参数）
    MaxConns:      128,
    }
    fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
    
    // 返回结果为一个包含所有组信息的结构体切片(数组)
    groups, err := fdfsClient.GetGroups()



```


#### 5.2 查询 group 信息

```code
// 设置 trackerServer 配置参数
var conf = &fastdfs_client_go.TrackerStorageServerConfig{
// 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
TrackerServer: []string{"192.168.10.10:22122"},
// tcp 连接池最大允许的连接数（trackerServer 和 storageServer 连接池共用该参数）
MaxConns:      128,
}
fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)

// 通过指定 组名(groupName) 查询组信息
groupName := "group1"
groupInfo, err := fdfsClient.GetGroupInfo(groupName);
//  groupInfo 表示返回的组基本信息结构体

```


#### 5.3 查询 group 下的所有 storage server 信息

```code
    // 设置 trackerServer 配置参数
    var conf = &fastdfs_client_go.TrackerStorageServerConfig{
    // 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
    TrackerServer: []string{"192.168.10.10:22122"},
    // tcp 连接池最大允许的连接数（trackerServer 和 storageServer 连接池共用该参数）
    MaxConns:      128,
    }
    fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
    
    // 通过指定 组名(groupName) 查询组下的所有 storage server 信息，
    // 返回结果为一个包含所有 storage server 信息的结构体切片(数组)
    groupName := "group1"
    storageServers, err := fdfsClient.GetStorageServersByGroup(groupName)


```

####  以上命令的使用示例，[点击查看单元测试详情](./test/fdfscClient_test.go)  


#### 5.最后一些说明  
- 5.1 `fastdfs` 分布式文件系统应该部署在内网环境,  整个系统原则上是不对互联网直接开放访问权限的（除了开发调试之外）.    
- 5.2 基于以上原因，开发者可以将用户上传的文件，首先暂停在临时目录，然后调用本客户端将临时目录的文件上传到 `fastdfs` 文件系统, 获取可访问的`文件id(fileId)(格式：group1/M00/00/01/MeiRdmISDUiAaURaAsRMrFnLJoE317.wav)`，最终返回给用户访问地址(建议通过nginx代理访问资源)  

#### 使用中遇到问题需要讨论
- QQ群：129885228 