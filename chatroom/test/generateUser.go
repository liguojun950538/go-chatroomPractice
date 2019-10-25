// 序列化一个user, 用来存放到redis中去

package main

import (
	"chatroom/common"
	"encoding/json"
	"io/ioutil"
)

func main() {
	user := common.User{
		UserID:   100,
		UserName: "jmx",
		Passwd:   "123",
	}
	userByte, _ := json.Marshal(user)
	ioutil.WriteFile("user100", userByte, 0666)
}
