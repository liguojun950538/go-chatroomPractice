package processors

import "fmt"

var (
	onlineUsersTable *userMsgTable
)

type userMsgTable struct {
	usersTable map[int]*UserProcessor
}

func init() {
	onlineUsersTable = &userMsgTable{
		usersTable: make(map[int]*UserProcessor, 1024),
	}
}

func GetUsersTableControler() *userMsgTable {
	return onlineUsersTable
}

// 设计增删改查
// 这里先不考虑安全问题, 即重复的ID登陆该怎么办
// 增即是改
func (this *userMsgTable) AddOnlineUsers(up *UserProcessor) {
	this.usersTable[up.UserID] = up
}

// 删
func (this *userMsgTable) DeleteOnlineUsers(userid int) {
	delete(this.usersTable, userid)
}

// 查
func (this *userMsgTable) GetUserProcessorByID(userid int) (retUP *UserProcessor, err error) {
	retUP, ok := this.usersTable[userid]
	if !ok {
		err = fmt.Errorf("用户%d不存在", userid)
	}
	return
}

func (this *userMsgTable) GetTable() map[int]*UserProcessor {
	return this.usersTable
}
