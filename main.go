package main

import (
	"fmt"
	"groupbot/api"
	"groupbot/group"
	"groupbot/messager"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const AdminId = 2811187643

var groupList group.GroupList
var groupLink map[int]string

type GroupInfo struct {
	GroupName           string
	GroupId             int
	GroupLink           string
	GroupMemberCount    int
	GroupMaxMemberCount int
	LoadAverage         float64
}

type GroupList map[int]*GroupInfo

type Load struct {
	Load  float64
	Group GroupInfo
}

var message string

func main() {
	//This is the group bot main function
	log.Println("Group Bot Start")
	for {
		min, max, total := load()
		message = fmt.Sprintf(`[画质魔盒群负载情况]
总负载: %.2f%%
成员情况: %d/%d
 
[最高负载]
%s[%d]
成员情况: %d/%d
负载情况: %.2f%%
 
[最低负载]
%s[%d]
成员情况: %d/%d
负载情况: %.2f%%`,
			total.Load, total.Group.GroupMemberCount, total.Group.GroupMaxMemberCount,
			max.Group.GroupName, max.Group.GroupId, max.Group.GroupMemberCount, max.Group.GroupMaxMemberCount, max.Load,
			min.Group.GroupName, min.Group.GroupId, min.Group.GroupMemberCount, min.Group.GroupMaxMemberCount, min.Load)
		messager.SendPrivate(AdminId, 0, message)
		result := api.ChangeGroupLink(min.Group.GroupLink)
		message = fmt.Sprintf("[修改群链接][%s]\n群名称: %s\n群号码: %d\n%s", result, min.Group.GroupName, min.Group.GroupId, min.Group.GroupLink)
		messager.SendPrivate(AdminId, 0, message)
		time.Sleep(time.Minute * 10)
	}
}

func load() (min, max, total Load) {
	maxLoad := Load{}
	minLoad := Load{}
	var totalMemCount float64 = 0
	var totalMaxMemCount float64 = 0
	var groupList = getGroupList()
	log.Println("get group load info ... ")
	// 获取所有群信息，并计算群负载情况
	for _, info := range groupList {
		var memCount = float64(info.GroupMemberCount * 1.00)
		var maxMemCount = float64(info.GroupMaxMemberCount * 1.00)
		totalMemCount += memCount
		totalMaxMemCount += maxMemCount
		var groupLoad = ceil(memCount/maxMemCount*100, 2)
		groupList[info.GroupId].LoadAverage = groupLoad

		if maxLoad.Load < groupLoad {
			maxLoad.Load = groupLoad
			maxLoad.Group = *info
		}

		if minLoad.Load > groupLoad || minLoad.Load == 0 {
			minLoad.Load = groupLoad
			minLoad.Group = *info
		}
	}
	totalLoad := Load{
		Load: ceil(totalMemCount/totalMaxMemCount*100, 2),
		Group: GroupInfo{
			GroupMemberCount:    int(totalMemCount),
			GroupMaxMemberCount: int(totalMaxMemCount),
		},
	}
	log.Printf("[Max] group: %s(%d), mem: %d/%d, load: %.2f\n", maxLoad.Group.GroupName, maxLoad.Group.GroupMemberCount, maxLoad.Group.GroupMaxMemberCount, maxLoad.Group.GroupId, maxLoad.Load)
	log.Printf("[Min] group: %s(%d), mem: %d/%d, load: %.2f\n", minLoad.Group.GroupName, minLoad.Group.GroupMemberCount, minLoad.Group.GroupMaxMemberCount, minLoad.Group.GroupId, minLoad.Load)
	log.Printf("[Total] member: %d/%d, load: %.2f\n", totalLoad.Group.GroupMemberCount, totalLoad.Group.GroupMaxMemberCount, totalLoad.Load)
	return minLoad, maxLoad, totalLoad
}

func ceil(i float64, n int) (ret float64) {
	var l = math.Pow(10, float64(n))

	i *= l

	ret = math.Floor(i)
	// 五入
	if i*10-math.Floor(i)*10 >= 5 {
		ret++
	}

	ret /= l
	return
}

func getGroupList() GroupList {
	log.Print("load group list ....")
	// 1. 获取群列表
	groupList = group.GetGroupList(true)
	log.Printf("load %d groups\n", len(groupList))
	log.Print("load group links ....")
	// 1.2 加载群链接列表
	groupLink = loadLinks("groups.ini")
	log.Printf("load %d links\n", len(groupLink))
	// 2. 获取群信息
	groups := GroupList{}
	for _, groupInfo := range groupList {
		// 2.1 设置群信息【群名称，群号码，群链接】
		groups[groupInfo.GroupId] = &GroupInfo{
			GroupName:           groupInfo.GroupName,
			GroupId:             groupInfo.GroupId,
			GroupLink:           groupLink[groupInfo.GroupId],
			GroupMemberCount:    groupInfo.MemberCount,
			GroupMaxMemberCount: groupInfo.MaxMemberCount,
		}
	}
	return groups
}

func loadLinks(filename string) (result map[int]string) {
	result = make(map[int]string)
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}
	contents := string(bytes)
	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		index := strings.Index(line, "=")
		if index <= 0 {
			continue
		}
		id, _ := strconv.Atoi(line[:index])
		link := line[index+1:]
		result[id] = link
	}
	return
}

//
//func putGroupName(groups group.GroupList) {
//	var groupInfoMap map[int]group.GroupInfo = make(map[int]group.GroupInfo)
//	for _, info := range groups {
//		groupInfoMap[info.GroupId] = info
//	}
//
//	bytes, err := os.ReadFile("groups.ini")
//	if err != nil {
//		return
//	}
//
//	var lines []string = strings.Split(string(bytes), "\n")
//
//	var content string = ""
//
//	for _, line := range lines {
//		var index = strings.Index(line, "=")
//		if index > 0 {
//			id, _ := strconv.Atoi(line[:index])
//			content += fmt.Sprintf("# %s \n", groupInfoMap[id].GroupName)
//		}
//		content += line + "\n"
//	}
//
//	err = os.WriteFile("groups.ini", []byte(content), 0644)
//	if err != nil {
//		return
//	}
//}
