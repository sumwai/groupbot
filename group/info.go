package group

import (
	"groupbot/api"
	"strconv"
)

type GroupInfo struct {
	GroupCreateTime int    `json:"group_create_time"`
	GroupId         int    `json:"group_id"`
	GroupLevel      int    `json:"group_level"`
	GroupMemo       string `json:"group_memo"`
	GroupName       string `json:"group_name"`
	MaxMemberCount  int    `json:"max_member_count"`
	MemberCount     int    `json:"member_count"`
}

// GetGroupInfo 获取群信息
//
// groupId: 指定群号码
//
// noCache: 是否不使用缓存
//
func GetGroupInfo(groupId int64, noCache bool) {
	var params = make(map[string]string)
	var noCacheString string = "false"
	if noCache {
		noCacheString = "true"
	}
	params["group_id"] = strconv.Itoa(int(groupId))
	params["no_cache"] = noCacheString
	api.Call("get_group_info", params)
}
