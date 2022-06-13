package group

import (
	"encoding/json"
	"fmt"
	"groupbot/api"
)

type GetGroupListResponse struct {
	Data    []GroupInfo `json:"data"`
	Retcode int         `json:"retcode"`
	Status  string      `json:"status"`
}

type GroupList []GroupInfo

// GetGroupList 获取群列表
//
// noCache: 是否不使用缓存
//
func GetGroupList(noCache bool) (groupList GroupList) {
	var params = make(map[string]string)
	params["no_cache"] = "false"
	if noCache {
		params["no_cache"] = "true"
	}
	result := api.Call("get_group_list", params)

	if result == "" {
		fmt.Println("failed to get group list response, it return a empty string")
		return nil
	}

	var groupListResponse *GetGroupListResponse

	err := json.Unmarshal([]byte(result), &groupListResponse)
	if err != nil {
		fmt.Println("failed to unmarshal json string")
		return nil
	}

	groupList = groupListResponse.Data
	return

}
