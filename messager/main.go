package messager

import (
	"fmt"
	"groupbot/api"
	"strconv"
)

func main() {
	fmt.Println("This is a message from messager")
}

// SendPrivate 发送私聊消息到指定用户
//
// user_id: 用户账号
//
// group_id: 如果为群临时消息，则此处填写群号，否则为0
//
// message: 消息内容
//
func SendPrivate(userId int64, groupId int64, message string) {
	var params = make(map[string]string)
	params["user_id"] = strconv.Itoa(int(userId))
	params["group_id"] = strconv.Itoa(int(groupId))
	params["message"] = message
	api.Call("send_private_msg", params)
}

// SendGroup 发送群消息到指定群
//
// user_id: 用户账号
//
// group_id: 如果为群临时消息，则此处填写群号，否则为0
//
// message: 消息内容
//
func SendGroup(groupId int64, message string) {
	var params = make(map[string]string)
	params["group_id"] = strconv.Itoa(int(groupId))
	params["message"] = message
	api.Call("send_group_msg", params)
}
