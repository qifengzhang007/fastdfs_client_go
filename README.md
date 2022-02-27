# fastdfs_client_go

### 1.概述

- `FastDFS` 采用二进制 `TCP` 通信协议.
- 开发本包的重点与核心主要是实现二进制通讯协议.

### 2.FastDFS二进制通讯协议细节

[点击查看GO 实现过程](./tcp_protocal_detail.md)

### 3.安装本包

```code 
// 请在本仓库的 gitTag 中查看最新版本,永远建议大家使用最新版本.  
// 查看地址：
go  get  github.com/qifengzhang007/fastdfs_client_go@v1.0.0

```

### 4. 已封装的函数列表

- 1.0.0 版本我们提供了核心功能，基本上可以解决绝大部分的需求，同时提供了非常详细的二进制协议对接细节、go 示例代码, 其他开发者可以仿照我们的项目结构自己扩展不常用功能.
- 关于其他未实现的不常用功能，如果您需要，可以提 issue , 我们会在下个版本更新进去.
- 在使用中出现的其他问题，都可以提 issue ，我们会在第一时间处理.

#### 4.1 文件上传(指定文件名)

```code  
    // 设置 trackerServer 配置参数
    var conf = &fastdfs_client_go.TrackerStorageServerConfig{
	// 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
        TrackerServer: []string{"192.168.10.10:22122"},
        MaxConns:      128,
    }
    # 文件上传核心函数
    fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
    fileId, err := fdfsClient.UploadByFileName(curDir + fileName)

```

#### 4.2 文件上传(传递二进制)

```code  
    // 设置 trackerServer 配置参数
    var conf = &fastdfs_client_go.TrackerStorageServerConfig{
	// 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
        TrackerServer: []string{"192.168.10.10:22122"},
        MaxConns:      128,
    }
    # 文件上传核心函数
    fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
    // 直接传递二进制上传文件，适合文件比较小的场景使用
    fileId, err := fdfsClient.UploadByBuffer([]byte(strconv.Itoa(no+1)+" - 二进制直接上传"),

```

#### 4.3 文件下载

```code  

    // 设置 trackerServer 配置参数
    var conf = &fastdfs_client_go.TrackerStorageServerConfig{
	// 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
        TrackerServer: []string{"192.168.10.10:22122"},
        MaxConns:      128,
    }
    // 指定需要被下载的文件id （fileId）
    fileId := "group1/M00/00/01/MeiRdmISDUiAaURaAsRMrFnLJoE317.wav" // 大小 46419116，约 46M 左右
    // 创建 fdfs 客户端
    fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
    
    // 指定需要下载的文件id（fileId），最终的保存路径，开始下载
	fdfsClient.DownloadFileByFileId(fileId, "E:/音乐文件夹/下载测试-俩俩相忘.wav")

```

#### 4.4 文件删除

```code   
    // 设置 trackerServer 配置参数
    var conf = &fastdfs_client_go.TrackerStorageServerConfig{
	// 替换为自己的 storagerServer ip 和端口即可，保证在开发阶段外网可访问
        TrackerServer: []string{"192.168.10.10:22122"},
        MaxConns:      128,
    }
	fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
	// 指定需要删除的文件Id
	fileId := "group1/M00/00/01/MeiRdmISSbuAZwwSAAAAD_Q4O2U879.txt"
	// 指定删除命令
	 err = fdfsClient.DeleteFile(fileId);
	 
```

####  以上命令的使用示例，[点击查看单元测试详情](./test/fdfscClient_test.go)  


#### 5.最后一些说明  
- 5.1 `fastdfs` 分布式文件系统应该部署在内网环境,  整个系统原则上是不对互联网直接开放访问权限的（除了开发调试之外）.    
- 5.2 基于以上原因，开发者可以将用户上传的文件，首先保存在临时目录，然后调用本客户端将临时目录的文件上传到 `fastdfs` 文件系统, 获取可访问的文件id，最终返回给用户访问地址(建议通过nginx代理访问资源)  

