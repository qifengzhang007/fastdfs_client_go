# fastdfs_client_go

## 1. 概述

`FastDFS` 是一个开源的轻量级分布式文件系统，采用二进制 `TCP` 通信协议进行高效的数据传输。  
本项目是基于 Go 语言实现的 FastDFS 客户端，重点在于实现 FastDFS 的二进制通信协议，方便 Go 语言开发者快速集成 FastDFS 文件存储功能。

[点击了解 FastDFS 分布式文件存储系统](https://github.com/happyfish100/fastdfs)

---

## 2. FastDFS 二进制通信协议细节

本项目实现了 FastDFS 的完整通信协议，包括 Tracker Server 和 Storage Server 的交互流程。

[点击查看 Go 实现的协议细节](./tcp_protocal_detail.md)

---

## 3. 安装本包

建议使用 `go get` 安装最新版本（请参考 [Gitee Tags](https://gitee.com/daitougege/fastdfs_client_go/tags) 获取最新版本号）：

```bash
go get gitee.com/daitougege/fastdfs_client_go@v1.1.0
```

---

## 4. 已封装的函数列表

| 序号 | 函数名                            | 功能描述                             | 适用版本 |
|------|-----------------------------------|--------------------------------------|----------|
| 1    | UploadByFileName                  | 上传文件（通过文件路径）             | 所有版本 |
| 2    | UploadByBuffer                    | 上传文件（通过字节流）               | 所有版本 |
| 3    | DownloadFileByFileId              | 下载文件                             | 所有版本 |
| 4    | DeleteFile                        | 删除文件                             | 所有版本 |
| 5    | GetRemoteFileInfo                 | 获取远程文件信息                     | 所有版本 |
| 6    | GetGroups                         | 获取所有组信息                       | ≥ 6.13   |
| 7    | GetGroupInfo                      | 获取单个组信息                       | ≥ 6.13   |
| 8    | GetStorageServersByGroup          | 获取组内所有存储节点信息             | ≥ 6.13   |
| 9    | ConvAppendFileToRegularFile       | 将追加文件转换为普通文件             | ≥ 6.13   |

---

## 5. 使用示例

### 5.1 上传文件（指定文件名）

```go
var conf = &fastdfs_client_go.TrackerStorageServerConfig{
    TrackerServer: []string{"192.168.10.10:22122"},
    MaxConns:      128,
}
fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
fileId, err := fdfsClient.UploadByFileName("/path/to/file", 0)
```

### 5.2 上传文件（传递字节流）

```go
fdfsClient, err := fastdfs_client_go.CreateFdfsClient(conf)
fileId, err := fdfsClient.UploadByBuffer([]byte("文件内容"), "txt", 0)
```

### 5.3 上传追加文件

```go
err := fdfsClient.UploadAppendFileByFileName("M00/00/00/xxx.log", "/path/to/local.log")
```

### 5.4 文件下载

```go
err := fdfsClient.DownloadFileByFileId("group1/M00/00/01/xxx.txt", "/path/to/save.txt")
```

### 5.5 获取远程文件信息

```go
remoteFileInfo, err := fdfsClient.GetRemoteFileInfo("group1/M00/00/01/xxx.txt")
```

### 5.6 删除文件

```go
err := fdfsClient.DeleteFile("group1/M00/00/01/xxx.txt")
```

### 5.7 转换追加文件为普通文件

```go
newFileId, err := fdfsClient.ConvAppendFileToRegularFile("group1/M00/00/00/xxx.log")
```

### 5.8 获取组信息

```go
groups, err := fdfsClient.GetGroups()
groupInfo, err := fdfsClient.GetGroupInfo("group1")
storages, err := fdfsClient.GetStorageServersByGroup("group1")
```

[点击查看单元测试代码](./test/fdfscClient_test.go)

---

## 6. 注意事项

- FastDFS 通常部署在内网环境中，不建议直接对外网开放。
- 建议用户上传文件时先保存到临时目录，再通过本客户端上传到 FastDFS。
- 文件访问建议通过 Nginx 等反向代理提供对外服务。

---

## 7. 联系与支持

- QQ群：129885228

---

## 协议

本项目遵循 MIT 协议。