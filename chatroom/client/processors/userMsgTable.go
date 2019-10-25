package processors

import (
	"chatroom/common"
	"fmt"
)

var (
	onlineUsersTable *userMsgTable
)

type userMsgTable struct {
	usersTable map[int]*common.User
}

func init() {
	onlineUsersTable = &userMsgTable{
		usersTable: make(map[int]*common.User, 1024),
	}
}

func GetUsersTableControler() *userMsgTable {
	return onlineUsersTable
}

// 设计增删改查
// 这里先不考虑安全问题, 即重复的ID登陆该怎么办
// 增即是改
func (this *userMsgTable) AddOnlineUsers(user *common.User) {
	this.usersTable[user.UserID] = user
}

// 删
func (this *userMsgTable) DeleteOnlineUsers(userid int) {
	delete(this.usersTable, userid)
}

// 查
func (this *userMsgTable) GetUserProcessorByID(userid int) (user *common.User, err error) {
	user, ok := this.usersTable[userid]
	if !ok {
		err = fmt.Errorf("用户%d不存在", userid)
	}
	return
}

func (this *userMsgTable) GetTable() map[int]*common.User {
	return this.usersTable
}

func (this *userMsgTable) PrintTable() {
	for idx, v := range this.usersTable {
		fmt.Printf("%d: ------ id = %d, name = %v\n", idx, v.UserID, v.UserName)
	}
	fmt.Println()
}
