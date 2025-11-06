### `FastDFS` 通信协议详解

> 在线地址：https://blog.csdn.net/happy_fish100/article/details/152223097

`FastDFS` 采用二进制 `TCP` 通信协议。一个数据包由包头（`header`）和包体（`body`）组成，包头只有 `10` 个字节，格式如下：

**header 包组成**

| 参数      | 类型     | 说明                             |
| ------- | ------ | ------------------------------ |
| pkg_len | 8 字节整数 | body 长度，不包含 header，只是 body 的长度 |
| cmd     | 1 字节整数 | 命令码                            |
| status  | 1 字节整数 | 状态码，0 表示成功，非 0 失败（UNIX 错误码）    |

数据包中的类型说明：

1. 整数类型采用网络字节序（`Big-Endian`），包括 `4` 字节整数和 `8` 字节整数；
2. `1` 字节整数不存在字节序问题，在 `Java` 中直接映射为 `byte` 类型，`C/C++` 中为 `char` 类型；
3. 固定长度的字符串类型以 `ASCII` 码 `0` 结尾，对于 `Java` 等语言需要调用 `trim` 处理返回的字符串。变长字符串的长度可以直接拿到或者根据包长度计算出来，不以 `ASCII 0` 结尾。

**特别说明**：FDFS v6.15 返回的 IP 地址长度可以在 `tracker.conf` 中设置 `response_ip_addr_size` 这个参数，默认配置为 `auto`，表示如果没有使用 IPv6，则返回的地址长度为 IPv4（15 字节）。可以显式设置为 IPv6，表示返回的是 IPv6 地址长度（45 字节）。

下面将列举 `client` 发送给 `FastDFS server` 的命令码及其 `body`（包体）结构。

### 公共命令码

#### FDFS_PROTO_CMD_ACTIVE_TEST

> 激活测试，通常用于检测连接是否有效。客户端使用连接池的情况下，建立连接后发送一次 `active test` 即可和 `server` 端保持长连接。

- 请求参数 

| 参数  | 类型  | 说明    |
| --- | --- | ----- |
| 无   | 无   | 无请求参数 |

- 响应结果 
  
  ```code
  无
  ```

### 发送给 tracker server 的命令码

#### TRACKER_PROTO_CMD_SERVER_LIST_ONE_GROUP

> 查看一个 `group` 状态 

- 请求参数 

| 参数         | 类型       | 说明  |
| ---------- | -------- | --- |
| group_name | 16 字节字符串 | 组名  |

- 响应结果

```code
@group_name：17字节字符串
@total_mb：8字节整数，磁盘空间总量，单位MB
@free_mb：8字节整数，磁盘剩余空间，单位MB
@reserved_mb：8字节整数，磁盘预留空间，单位MB（since V6.13.1）
@trunk_free_mb：8字节整数，trunk文件剩余空间，单位MB（合并存储开启时有效）
@server_count：8字节整数，storage server数量
@storage_port：8字节整数，storage server端口号
@storage_http_port：8字节整数，storage server上的HTTP端口号（V6.13删除该字段）
@active_count：8字节整数，当前活着的storage server数量（V6.13删除该字段）
@readable_server_count：8字节整数，当前可读的storage server数量（since V6.13）
@writable_server_count：8字节整数，当前可写的storage server数量（since V6.13）
@current_write_server：8字节整数，当前写入的storage server顺序号
@store_path_count：8字节整数，storage server存储路径数
@subdir_count_per_path：8字节整数，存储路径下的子目录数（FastDFS采用两级子目录），如 256
@current_trunk_file_id：8字节整数，当前使用的trunk文件ID（合并存储开启时有效）
```

#### TRACKER_PROTO_CMD_SERVER_LIST_ALL_GROUPS

> 列举所有 `group`

- 请求参数

| 参数  | 类型  | 说明    |
| --- | --- | ----- |
| 无   | 无   | 无请求参数 |

- 返回结果

```code
n 个group实体信息，n >= 0。每个group的数据结构参见上面的TRACKER_PROTO_CMD_SERVER_LIST_ONE_GROUP。
```

#### TRACKER_PROTO_CMD_SERVER_LIST_STORAGE

> 列举一个 group 下的 storage server`

- 请求参数

| 参数         | 类型              | 说明                     |
| ---------- | --------------- | ---------------------- |
| group_name | 16 字节字符串        | 组名                     |
| server_id  | 不定长，最大长度为 15 字节 | storage server id，可选参数 |

- 返回结果

```code
n 个storage server实体信息，n >= 0。每个storage实体结构如下：
@status：1字节整数，storage server状态
@rw_mode：1字节整数，读写模式（since v6.13）
@id：16字节字符串，server ID
@ip_addr：16/46字节字符串，IP地址（V6.11前为16字节，之后为46字节）
@domain_name：128字节字符串，域名（V6.13删除该字段）
@src_storage_id：16字节字符串，同步源storage的server ID
@version：6/8字节字符串，运行的FastDFS版本号，例如6.04（V6.11前为6字节，V6.11开始为8字节）
@join_time: 8字节整数，加入集群时间
@up_time: 8字节整数，fdfs_storaged启动时间
@total_mb: 8字节整数，磁盘空间总量，单位MB
@free_mb: 8字节整数，磁盘剩余空间，单位MB
@reserved_mb：8字节整数，磁盘预留空间，单位MB（since V6.13.1）
@upload_priority: 8字节整数，上传文件优先级
@store_path_count: 8字节整数，存储路径数
@subdir_count_per_path: 8字节整数，存储路径下的子目录数（FastDFS采用两级子目录），如 256
@current_write_path: 8字节整数，当前写入的存储路径（顺序号）
@storage_port: 8字节整数，storage server服务端口号
@storage_http_port: 8字节整数，HTTP服务端口号（V6.13删除该字段）
@alloc_count：4字节整数，已分配的连接buffer数目
@current_count：4字节整数，当前连接数
@max_count：4字节整数，曾经达到过的最大连接数
@total_upload_count：8字节整数，上传文件总数
@success_upload_count：8字节整数，成功上传文件数
@total_append_count：8字节整数，调用append总次数
@success_append_count：8字节整数，成功调用append次数
@total_modify_count：8字节整数，调用modify总次数
@success_modify_count：8字节整数，成功调用modify次数
@total_truncate_count：8字节整数，调用truncate总次数
@success_truncate_count：8字节整数，成功调用truncate次数
@total_set_meta_count：8字节整数，设置文件附加属性（meta data）总次数
@success_set_meta_count：8字节整数，成功设置文件附加属性（meta data）次数
@total_delete_count：8字节整数，删除文件总数
@success_delete_count：8字节整数，成功删除文件数
@total_download_count：8字节整数，下载文件总数
@success_download_count：8字节整数，成功下载文件数
@total_get_meta_count：8字节整数，获取文件附加属性（meta data）总次数
@success_get_meta_count：8字节整数，成功获取文件附加属性（meta data）次数
@total_create_link_count：8字节整数，创建文件符号链接总数
@success_create_link_count：8字节整数，成功创建文件符号链接数
@total_delete_link_count：8字节整数，删除文件符号链接总数
@success_delete_link_count：8字节整数，成功删除文件符号链接数
@total_upload_bytes：8字节整数，上传文件总字节数
@success_upload_bytes：8字节整数，成功上传文件字节数
@total_append_bytes：8字节整数，append总字节数
@success_append_bytes：8字节整数，成功append字节数
@total_modify_bytes：8字节整数，modify总字节数
@success_modify_bytes：8字节整数，成功modify字节数
@total_download_bytes：8字节整数，下载总字节数
@success_download_bytes：8字节整数，成功下载字节数
@total_sync_in_bytes：8字节整数，文件同步流入总字节数
@success_sync_in_bytes：8字节整数，文件同步成功流入字节数
@total_sync_out_bytes：8字节整数，文件同步流出总字节数
@success_sync_out_bytes：8字节整数，文件同步成功流出字节数
@total_file_open_count：8字节整数，文件打开总次数
@success_file_open_count：8字节整数，文件成功打开次数
@total_file_read_count：8字节整数，文件读总次数
@success_file_read_count：8字节整数，文件成功读次数
@total_file_write_count：8字节整数，文件写总次数
@success_file_write_count：8字节整数，文件成功写次数
@last_source_update：8字节整数，最近一次源头更新时间
@last_sync_update：8字节整数，最近一次同步更新时间
@last_synced_timestamp：8字节整数，最近一次被同步到的时间戳
@last_heart_beat_time：8字节整数，最近一次心跳时间
```

#### TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE

> 获取一个 storage server 用来存储文件（不指定 group name）

- 请求参数

| 参数  | 类型  | 说明    |
| --- | --- | ----- |
| 无   | 无   | 无请求参数 |

- 返回结果

```code
@group_name：16字节字符串，组名
@ip_addr：15/45字节字符串，storage server IP地址（V6.11前为15字节，V6.11开始为45字节）
@port：8字节整数，storage server端口号
@store_path_index：1字节整数，基于0的存储路径顺序号
```

#### TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ONE

> 获取一个 storage server 用来存储文件（指定组名）

- 请求参数

| 参数         | 类型       | 说明  |
| ---------- | -------- | --- |
| group_name | 16 字节字符串 | 组名  |

- 返回结果

```code
@ip_addr：15/45字节字符串，storage server IP地址（V6.11前为15字节，V6.11开始为45字节）
@port：8字节整数，storage server端口号
@store_path_index：1字节整数，基于0的存储路径顺序号
```

#### TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ALL

> 获取 storage server 列表用来存储文件（不指定组名）

- 请求参数

| 参数  | 类型  | 说明    |
| --- | --- | ----- |
| 无   | 无   | 无请求参数 |

- 返回结果

```code
@group_name：16字节字符串，组名
@ip_addr：15/45字节字符串，storage server IP地址（V6.11前为15字节，V6.11开始为45字节）
@port：8字节整数，storage server端口号
@store_path_index：1字节整数，基于0的存储路径顺序号
```

#### TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ALL

> 获取 storage server 列表用来存储文件（指定组名）

- 请求参数

| 参数         | 类型       | 说明  |
| ---------- | -------- | --- |
| group_name | 16 字节字符串 | 组名  |

- 返回结果

```code
@group_name：16字节字符串，组名
@ip_addr：15/45字节字符串，storage server IP地址（V6.11前为15字节，V6.11开始为45字节）
@port：8字节整数，storage server端口号
@store_path_index：1字节整数，基于0的存储路径顺序号
```

#### TRACKER_PROTO_CMD_SERVICE_QUERY_UPDATE

> 获取 storage server 列表用来修改文件或文件附加信息

- 请求参数

| 参数         | 类型       | 说明  |
| ---------- | -------- | --- |
| group_name | 16 字节字符串 | 组名  |
| filename   | 不定长字符串   | 文件名 |

- 返回结果

```code
@group_name：16字节字符串，组名
@ip_addr：15/45字节字符串，storage server IP地址（V6.11前为15字节，V6.11开始为45字节）
@port：8字节整数，storage server端口号
```

#### TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ALL

> 获取 storage server 列表用来下载文件

- 请求参数

| 参数         | 类型       | 说明  |
| ---------- | -------- | --- |
| group_name | 16 字节字符串 | 组名  |
| filename   | 不定长字符串   | 文件名 |

- 返回结果

```code
@group_name：16字节字符串，组名
@ip_addr：15/45字节字符串，第一个storage server IP地址（V6.11前为15字节，V6.11开始为45字节）
@port：8字节整数，storage server端口号
@ip_addr：15/45字节字符串，其他storage server IP地址（V6.11前为15字节，V6.11开始为45字节）
```

### 发送给 `storage server` 的命令码

#### STORAGE_PROTO_CMD_UPLOAD_FILE / STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE

> 上传普通文件 / 上传 `appender` 类型文件

- 请求参数

| 参数               | 类型                     | 说明                                                            |
| ---------------- | ---------------------- | ------------------------------------------------------------- |
| store_path_index | 1 字节整数                 | 基于 0 的存储路径顺序号                                                 |
| meta_data_length | 8 字节整数                 | meta data（文件附加属性）长度，可以为 0                                     |
| file_size        | 8 字节整数                 | 文件大小                                                          |
| file_ext_name    | 6 字节字符串                | 不包括小数点的文件扩展名，例如 jpeg、tar.gz                                   |
| meta_data        | meta_data_length 字节字符串 | 文件附加属性，每个属性用字符 `\x01` 分隔，名称 `key` 和取值 `value` 之间用字符 `\x02` 分隔 |
| file_content     | file_size 字节二进制内容      | 文件内容                                                          |

- 返回结果

```code
@group_name：16字节字符串，组名
@filename：不定长字符串，文件名
```

#### STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE

> 上传 `slave` 文件

- 请求参数

| 参数                     | 类型                           | 说明                                                              |
| ---------------------- | ---------------------------- | --------------------------------------------------------------- |
| master_filename_length | `8` 字节整数                     | 主文件名长度                                                          |
| meta_data_length       | `8` 字节整数                     | `meta data`（文件附加属性）长度，可以为 `0`                                   |
| file_size              | `8` 字节整数                     | 文件大小                                                            |
| filename_prefix        | `16` 字节字符串                   | 从文件前缀名                                                          |
| file_ext_name          | `6` 字节字符串                    | 不包括小数点的文件扩展名，例如 `jpeg`、`tar.gz`                                 |
| master_filename        | master_filename_length 字节字符串 | 主文件名                                                            |
| meta_data              | meta_data_length 字节字符串       | 文件附加属性，每个属性记录用字符 `\x01` 分隔，名称 `key` 和取值 `value` 之间用字符 `\x02` 分隔 |
| file_content           | file_size 字节二进制内容            | 文件内容                                                            |

- 返回结果

```code
@group_name：16字节字符串，组名
@filename：不定长字符串，文件名
```

#### STORAGE_PROTO_CMD_DELETE_FILE

> 删除文件

- 请求参数

| 参数           | 类型         | 说明  |
| ------------ | ---------- | --- |
| `group_name` | `16` 字节字符串 | 组名  |
| `filename`   | 不定长字符串     | 文件名 |

- 返回结果

```code
无
```

#### STORAGE_PROTO_CMD_SET_METADATA

> 设置 `meta data`（文件附加属性）

- 请求参数

| 参数               | 类型                     | 说明                                                              |
| ---------------- | ---------------------- | --------------------------------------------------------------- |
| filename_length  | 8 字节整数                 | 文件名长度                                                           |
| meta_data_length | 8 字节整数                 | meta data（文件附加属性）长度，可以为 `0`                                     |
| op_flag          | 1 字节字符                 | 操作标记：'O' - 覆盖方式，'M'` - merge 方式                                 |
| group_name       | 16 字节字符串               | 组名                                                              |
| filename         | filename_length 字节的字符串 | 文件名                                                             |
| meta_data        | meta_data_length 字节字符串 | 文件附加属性，每个属性记录用字符 `\x01` 分隔，名称 `key` 和取值 `value` 之间用字符 `\x02` 分隔 |

- 返回结果

```code
无
```

#### STORAGE_PROTO_CMD_DOWNLOAD_FILE

> 下载文件

- 请求参数

| 参数             | 类型       | 说明    |
| -------------- | -------- | ----- |
| file_offset    | 8 字节整数   | 文件偏移量 |
| download_bytes | 8 字节整数   | 下载字节数 |
| group_name     | 16 字节字符串 | 组名    |
| filename       | 不定长字符串   | 文件名   |

- 返回结果

```code
@file_content：不定长二进制内容，文件内容
```

#### STORAGE_PROTO_CMD_GET_METADATA

> 获取 `meta data`（文件附加属性）

- 请求参数

| 参数         | 类型       | 说明  |
| ---------- | -------- | --- |
| group_name | 16 字节字符串 | 组名  |
| filename   | 不定长字符串   | 文件名 |

- 返回结果

```code
@meta_data：不定长字符串，文件附加属性，每个属性记录用字符 \x01分隔，名称key和取值value之间用字符 \x02分隔
```

#### STORAGE_PROTO_CMD_QUERY_FILE_INFO

>  获取文件信息

- 请求参数

| 参数         | 类型       | 说明  |
| ---------- | -------- | --- |
| group_name | 16 字节字符串 | 组名  |
| filename   | 不定长字符串   | 文件名 |

- 返回结果

```code
@file_size：8字节整数，文件大小
@create_timestamp：8字节整数，文件创建时间（Unix时间戳）
@crc32：8字节整数，文件内容CRC32校验码
@source_ip_addr：16字节字符串，源storage server IP地址
```

#### STORAGE_PROTO_CMD_APPEND_FILE

> 追加文件内容

- 请求参数

| 参数                       | 类型                              | 说明             |
| ------------------------ | ------------------------------- | -------------- |
| appender_filename_length | 8 字节整数                          | appender 文件名长度 |
| file_size                | 8 字节整数                          | 文件大小           |
| appender_filename        | appender_filename_length 字节数字符串 | appender 文件名   |
| file_content             | file_size 字节数的二进制内容             | 追加的文件内容        |

- 返回结果

```code
无
```

#### STORAGE_PROTO_CMD_MODIFY_FILE

> 修改文件内容

- 请求参数

| 参数                       | 类型                              | 说明             |
| ------------------------ | ------------------------------- | -------------- |
| appender_filename_length | 8 字节整数                          | appender 文件名长度 |
| file_offset              | 8 字节整数                          | 文件偏移量          |
| file_size                | 8 字节整数                          | 文件大小           |
| appender_filename        | appender_filename_length 字节数字符串 | appender 文件名   |
| file_content             | file_size 字节数的二进制内容             | 新（目标）文件内容      |

- 返回结果

```code
无
```

#### STORAGE_PROTO_CMD_TRUNCATE_FILE

> 改变文件大小

- 请求参数

| 参数                       | 类型                              | 说明              |
| ------------------------ | ------------------------------- | --------------- |
| appender_filename_length | 8 字节整数                          | appender 文件名长度  |
| truncated_file_size      | 8 字节整数                          | truncate 后的文件大小 |
| appender_filename        | appender_filename_length 字节数字符串 | appender 文件名    |

- 返回结果

```code
无
```

#### STORAGE_PROTO_CMD_REGENERATE_APPENDER_FILENAME

> `appender` 类型文件改名为普通文件

- 请求参数

| 参数                | 类型     | 说明             |
| ----------------- | ------ | -------------- |
| appender_filename | 不定长字符串 | `appender` 文件名 |

- 返回结果

```code
@group_name：16字节字符串，组名
@filename：不定长字符串，重新生成的文件名（普通类型）
```

### 命令码常量列表

| 常量名                                                     | 值   | 说明                        |
| ------------------------------------------------------- | --- | ------------------------- |
| TRACKER_PROTO_CMD_RESP                                  | 100 | tracker 响应码               |
| STORAGE_PROTO_CMD_RESP                                  | 100 | storage 响应码               |
| FDFS_PROTO_CMD_ACTIVE_TEST                              | 111 | 激活测试                      |
| TRACKER_PROTO_CMD_SERVER_LIST_ONE_GROUP                 | 90  | 查看一个group状态               |
| TRACKER_PROTO_CMD_SERVER_LIST_ALL_GROUPS                | 91  | 列举所有group                 |
| TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE | 101 | 获取storage server（不指定组名）   |
| TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE               | 102 | 获取一个storage server用来下载文件  |
| TRACKER_PROTO_CMD_SERVICE_QUERY_UPDATE                  | 103 | 获取storage server列表用来修改文件  |
| TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ONE    | 104 | 获取storage server（指定组名）    |
| TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ALL               | 105 | 获取storage server列表用来下载文件  |
| TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ALL | 106 | 获取storage server列表（不指定组名） |
| TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ALL    | 107 | 获取storage server列表（指定组名）  |
| STORAGE_PROTO_CMD_UPLOAD_FILE                           | 11  | 上传普通文件                    |
| STORAGE_PROTO_CMD_DELETE_FILE                           | 12  | 删除文件                      |
| STORAGE_PROTO_CMD_SET_METADATA                          | 13  | 设置meta data               |
| STORAGE_PROTO_CMD_DOWNLOAD_FILE                         | 14  | 下载文件                      |
| STORAGE_PROTO_CMD_GET_METADATA                          | 15  | 获取meta data               |
| STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE                     | 21  | 上传slave文件                 |
| STORAGE_PROTO_CMD_QUERY_FILE_INFO                       | 22  | 获取文件信息                    |
| STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE                  | 23  | 上传appender类型文件            |
| STORAGE_PROTO_CMD_APPEND_FILE                           | 24  | 追加文件内容                    |
| STORAGE_PROTO_CMD_MODIFY_FILE                           | 34  | 修改文件内容                    |
| STORAGE_PROTO_CMD_TRUNCATE_FILE                         | 36  | 改变文件大小                    |
| STORAGE_PROTO_CMD_REGENERATE_APPENDER_FILENAME          | 38  | appender类型文件改名为普通文件       |