### fastdfs_client_for_go

### 1.TCP 通信协议详情  
本篇文档在编写过程中遇到了几个的错误，主要是因为作者在 2009年 和 2019年 分别发布了两份协议参数，但是两份协议参数不尽相同，有些协议的对接反而是旧版本成功了，而按照新版本的协议参数对接却始终失败.   
参考地址(作者2009年发布)：`    http://bbs.chinaunix.net/thread-2001015-1-1.html   `    
参考地址(作者2019年发布)：`    https://mp.weixin.qq.com/s/lpWEv3NCLkfKmtzKJ5lGzQ   `      
本篇文档是综合了以上两份协议参数编写而成，有时候按照新版本协议对接失败，请切换为旧版本参数对接调试即可 .    

### 2.header 和 body 组成

FastDFS 采用二进制 TCP 通信协议。一个数据包由 包头（header）和包体（body）组成。

#### 2.1 包头(header)由10个字节，格式如下：

```code   

    @ pkg_len：8字节整数，body体长度，不包含header, 只是body的长度
    @ cmd：1字节整数，命令码
    @ status：1字节整数，状态码，0表示成功，非0失败（UNIX错误码）
    
```
发送header 头参数核心函数   

```code   
//sendHeader  发送header消息
func (h *header) sendHeader(conn net.Conn) error {

    // 创建一个发送数据前的字节缓冲区
	buffer := new(bytes.Buffer)
	
	// c.pkgLen 整数型（4字节或者8字节）写入到二进制缓冲区
	//将整数类型采用网络字节序（Big-Endian），包括4字节整数(int32)和8字节整数(int64)
	if err := binary.Write(buffer, binary.BigEndian,h.pkgLen); err != nil {
		return err
	}
	// 单字节写入缓冲区
	buffer.WriteByte(h.cmd)
	buffer.WriteByte(h.status)

    // 通过 tcp 连接将缓冲区的全部字节数据发送出去
	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

```
包头参数解释:  
1.整数类型采用网络字节序（Big-Endian）, 包括4字节整数和8字节整数.  
2.字节整数不存在字节序问题 , 在 GO  中直接映射为 byte 类型，C/C++ 中为 char 类型.   
3.固定长度的字符串类型以 ASCII 码 0 结尾 , 对于固定长度字符串, 开发者必须在取出的二进制文件中，找到 `0x0` 然后截取前面的字符串.   
4.变长字符串的长度可以直接拿到或者根据包长度计算出来, 不以 ASCII 0 结尾.

#### 2.2 包体(body)组成
不同的命令，`body` 体发出去的参数不同, 这就相当于调用不同的函数，传递的参数不同是一样的道理，但是我们只要按照官方公布的 `tcp` 协议格式传递参数就可以实现相应的功能.    
下面我们将以具体的命令为例来说明 `body` 参数的使用.  

###  3.命令列表

#### 3.1 文件上传 - STORAGE_PROTO_CMD_UPLOAD_FILE
> 通过指定已经存在的文件名上传到服务器  
- 1.`header` 头命令
```code  
    @pkg_len : int64  ，8字节 ，body 体的总字节数目，不包括 header 任何部分
    @cmd： byte , 1字节 ， 上传命令常量  STORAGE_PROTO_CMD_UPLOAD_FILE 
    @status： byte，1字节 ，该参数主要用于返回判断状态，发送时默认为 0 即可 
```
`header`  头参数组装成二进制核心函数
```code    
    s.header.status = 0  // 主要用于返回做判断，发送时位置为0即可
    s.header.cmd = STORAGE_PROTO_CMD_UPLOAD_FILE   //  具体的协议常量值
    
    // pkgLen 表示 body 体的总字节长度，计算依据请参考后面 body 体的参数组成
    // 15 的组成 ：@store_path_index (1字节整数) + @fileSize ( 8字节整数 ) + @file_ext_name (6字节字符串)
    s.header.pkgLen = s.fileInfo.fileSize + 15
    if err = s.sendHeader(conn); err != nil {
        return err
    }
```

- 2.`body` 体发送参数
```code  
    # body 发送参数构成
    @store_path_index :   byte ，基于 0 的存储路径顺序号 ，服务器存储目录有很多个，0 表示存储目录的序号，具体值由 tracker server 返回
    @file_size ： int64  ，被上传的文件大小(总字节量)
    @file_ext_name  string ，最大6字节，不包括小数点的文件扩展名，例如  jpeg、tar.gz
    @file_content  []byte , 字节切片，不定长 
```
`body`  体参数组装成二进制核心过程
```code  

// 创建 body 数据发送时的缓冲区
	buffer := new(bytes.Buffer)

	// @store_path_index (1字节整数)
	buffer.WriteByte(s.storagePathIndex)

	//@file_size：8字节整数
	if err = binary.Write(buffer, binary.BigEndian, s.fileInfo.fileSize); err != nil {
		return err
	}
	// 文件扩展名 6 字节
	buffer.Write(GetFileExtNameBytes(s.fileInfo.fileName))

    // 这里首先把报文头发给服务器，服务器就会根据报文头，一直接受完毕约定的字节量
	if _, err = conn.Write(buffer.Bytes()); err != nil {
		return err
	}
	
	//发送文件内容本身的二进制数据
	// 1.如果文件结构信息对应的指针不为空，表示该文件是通过文件名方式打开操作的，那么就根据文件指针读取数据，发送出去
	// 2.如果文件结构信息中的文件指针为空，表示该文件需要通过字节流发送出去
	// 3.最后将文件真正的内容发出去，其实 tcp 底层会对大文件分块多次发送，服务器端会按照收到的报文头读取对应的数量字节才结束.
	if s.fileInfo.filePtr != nil {
		_, err = conn.(*net.TCPConn).ReadFrom(s.fileInfo.filePtr)
	} else {
		_, err = conn.Write(s.fileInfo.buffer)
	}
```

- 3.`body` 体接受参数
```code  
    @ group_name ： string，16字节字符串，组名 
    @ filename： string 不定长字符串，文件名
```
首先接受 `header` 头参数, 先行判断服务器响应的状态码、数据长度必须符合协议约定，然后再读取 `body` 体的内容.   
响应的body体协议规定的格式如下：
```code   
#body 响应格式
@ group_name：16字节字符串，组名 
@ filename：不定长字符串，文件名 
```
根据 `body` 体的参数，那么 `header` 的响应就必须满足如下格式
```code   
    # 响应 header 头参数 
    pkgLen int64  // 长度 必须>16
    cmd    byte  // 发送命令的参数，接受时忽略即可
    status byte  // 状态值必须是0，其他值表示有错误

```

#### 3.2 文件下载 - STORAGE_PROTO_CMD_DOWNLOAD_FILE
- 1.`header` 头命令
```code  
    @pkg_len : int64  ，8字节 ，body 体的总字节数目，不包括 header 任何部分
    @cmd： byte , 1字节 ， 上传命令常量  STORAGE_PROTO_CMD_UPLOAD_FILE 
    @status： byte，1字节 ，该参数主要用于返回判断状态，发送时默认为 0 即可 
```
`header` 下载命令头参数赋值
```code   
	// 构建 header 头参数
	s.header.status = 0
	s.header.cmd = STORAGE_PROTO_CMD_DOWNLOAD_FILE
	// 32 = body 下载参数的前3个参数二进制总长度
	s.header.pkgLen = int64(len(s.remoteFilename) + 32)

	if err := s.sendHeader(conn); err != nil {
		return err
	}

```
- 2.`body` 参数
```code   
// Send 下载文件 body 参数
// @file_offset 8字节整数，文件偏移量,从指定的位置开始下载
// @download_bytes：8字节整数，需要下载字节数
// @group name：16字节字符串，组名
// @filename：不定长字符串，文件名
```
`body` 参数赋值
```code   

		// 构建 body 头参数
	buffer := new(bytes.Buffer)
	if err := binary.Write(buffer, binary.BigEndian, s.offset); err != nil {
		return err
	}
	if err := binary.Write(buffer, binary.BigEndian, s.downloadBytes); err != nil {
		return err
	}
	buffer.Write(groupNameConvBytes(s.groupName))
	buffer.WriteString(s.remoteFilename)
	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return err
	}
	
```


#### 3.3  文件删除 - STORAGE_PROTO_CMD_DELETE_FILE
- 1.`header` 头命令格式
```code  
    @pkg_len : int64  ，8字节 ，body 体的总字节数目，不包括 header 任何部分
    @cmd： byte , 1字节 ， 上传命令常量  STORAGE_PROTO_CMD_UPLOAD_FILE 
    @status： byte，1字节 ，该参数主要用于返回判断状态，发送时默认为 0 即可 
```
`header` 删除命令头参数赋值
```code   
	// 设置删除文件时的 header 参数
	s.header.status = 0
	s.header.cmd = STORAGE_PROTO_CMD_DELETE_FILE
	s.header.pkgLen = int64(len(s.remoteFilename) + 16)

	if err := s.sendHeader(conn); err != nil {
		return err
	}

```
- 2.`body` 参数
```code   
//Send 发送删除文件命令
//@group_name：16字节字符串，组名
//@filename：不定长字符串，文件名
```
`body` 参数赋值
```code   
	// 写入body参数 
	buffer := new(bytes.Buffer)
	buffer.Write(groupNameConvBytes(s.groupName))
	buffer.WriteString(s.remoteFilename)
    // 发送body参数
	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return err
	}
	
```

### 命令常量列表

| 常量含义                  | 常量代码                                                    | 相关值 |
|:----------------------|:--------------------------------------------------------|:----|
| tracker 正确响应码         | TRACKER_PROTO_CMD_RESP                                  | 100 |
| storage 正确响应码         | STORAGE_PROTO_CMD_RESP                                  | 100 |
| 激活测试,通常用于检测连接是否有效     | FDFS_PROTO_CMD_ACTIVE_TEST                              | 111 |
| 待补充                   | TRACKER_PROTO_CMD_SERVER_LIST_ONE_GROUP                 | 90  |
| 获取组列表                 | TRACKER_PROTO_CMD_SERVER_LIST_ALL_GROUPS                | 91  |
| 不需要组名获取一个存储节点         | TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE | 101 |
| 获取下载节点QUERY_FETCH_ONE | TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE               | 102 |
| 获取更新节点QUERY_UPDATE    | TRACKER_PROTO_CMD_SERVICE_QUERY_UPDATE                  | 103 |
| 按组获取存储节点              | TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ONE    | 104 |
| 待补充                   | TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ALL               | 105 |
| 待补充                   | TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ALL | 106 |
| 待补充                   | TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ALL    | 10  |
| 文件上传                  | STORAGE_PROTO_CMD_UPLOAD_FILE                           | 11  |
| 删除文件                  | STORAGE_PROTO_CMD_DELETE_FILE                           | 12  |
| 设置文件元数据               | STORAGE_PROTO_CMD_SET_METADATA                          | 13  |
| 文件下载                  | STORAGE_PROTO_CMD_DOWNLOAD_FILE                         | 14  |
| 获取文件元数据               | STORAGE_PROTO_CMD_GET_METADATA                          | 15  |
| 上传附属文件                | STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE                     | 21  |
| 查询文件信息                | STORAGE_PROTO_CMD_QUERY_FILE_INFO                       | 22  |
| 创建支持断点续传的文件           | STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE                  | 23  |
| 断点续传                  | STORAGE_PROTO_CMD_APPEND_FILE                           | 24  |
| 文件修改                  | STORAGE_PROTO_CMD_MODIFY_FILE                           | 34  |
| 清除文件                  | STORAGE_PROTO_CMD_TRUNCATE_FILE                         | 36  |
| 待补充                   | STORAGE_PROTO_CMD_REGENERATE_APPENDER_FILENAME          | 38  |
