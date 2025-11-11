package fastdfs_client_go

import "errors"

// GetGroupInfo  获取远程服务器的文件信息
// @remoteFileId 远程服务器的文件Id
func (c *trackerServerTcpClient) GetGroupInfo(groupName string) (resGroupInfo *GroupInfo, err error) {
	if len(groupName) < 1 {
		return resGroupInfo, errors.New(ERROR_GET_GROUP_INFO_FAILED_GROUP_NOT_EMPTY)
	}
	resGroupInfo, err = c.getGroupInfoByTracker(TRACKER_PROTO_CMD_SERVER_LIST_ONE_GROUP, groupName)
	if err != nil {
		return resGroupInfo, err
	}
	return resGroupInfo, nil
}
