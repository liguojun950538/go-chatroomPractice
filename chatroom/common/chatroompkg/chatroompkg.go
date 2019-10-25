package chatroompkg

import (
	"chatroom/common"
)

const (
	KLoginMsgType = iota
	KLoginResMsgType
	KRegisterMsgType
	KRegisterResMsgType
	KOnlineNotifyMsgType
	KGroupChatMsgType
	KExitMsgType
	KPointChatMsgType
	KPointChatResMsgType
	KPointChatResToOtherMsgType
)

const (
	KLoginOK      = 200
	KUserNotExits = 500
)

const (
	KStatusUserOnline = iota
	KStatusUserOffline
	KStatusUserBusy
)

// 所有的消息
type Pkg struct {
	Type int    `json:"type"` // 消息的类型
	Data string `json:"data"` // 消息的数据
}

// 根据协议确定
type LoginMsg struct {
	UserId int    `json:"userid"`
	Passwd string `json:"passwd"`
}

type LoginResMsg struct {
	Code        int    `json:"code"`
	Error       string `json:"error"`
}

type RegisterMsg struct {
	common.User
}

type RegisterResMsg struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type OnlineNotifyMsg struct {
	UserID     int `json:"userid"`
	UserStatus int `json:"userstatus"`
}

type GroupChatMsg struct {
	common.UserCanBeSeen
	Content string `json:"content"`
}
type PointChatMsg struct {
	common.UserCanBeSeen
	ToUserID int    `json:"touserid"`
	Content  string `json:"content"`
}
type PointChatResMsg struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
type PointChatToOtherMsg GroupChatMsg

type ExitMsg struct {
}
