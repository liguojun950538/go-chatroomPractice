package processors

import (
	"chatroom/client/menu"
	"chatroom/common"
	"chatroom/common/utils"
	"fmt"
	"net"
	"os"
)

var (
	theMainProcessor *mainProcessor
)

type mainProcessor struct {
	up       *UserProcessor
	sp       *SmsProcessor
	currUser *common.User
	currConn net.Conn
}

func init() {
	theMainProcessor = &mainProcessor{}
}
func GetTheMainProcessor() *mainProcessor {
	return theMainProcessor
}
func (this *mainProcessor) Process(loop bool) {
	viewer := menu.GetTheviewer()
	for loop {
		key := viewer.WelcomeMenu()
		this.processWelcomeMenu(key)
	}
}

func (this *mainProcessor) processWelcomeMenu(key int) {
	viewer := menu.GetTheviewer()
	switch key {
	case 1:
		fmt.Println("登陆聊天室")
		viewer.LoginMenu()
		this.up = &UserProcessor{}
		this.up.Login(viewer.User)

	case 2:
		fmt.Println("注册用户")
		this.up = &UserProcessor{}
		viewer.RegisterMenu()
		this.up.Register(viewer.User)
	case 3:
		fmt.Println("退出系统")
		//loop = false
		os.Exit(0)
	default:
		fmt.Println("你的输入有误，请重新输入")
	}
}

func (this *mainProcessor) processMainMenu(key int) {
	/*
		fmt.Println("\n\n1. 显示在线用户列表")
		fmt.Println("2. 发送信息")
		fmt.Println("3. 信息列表")
		fmt.Println("4. 退出系统")
	*/
	switch key {
	case 1:
		this.showOnlineUserTable()
	case 2:
		key := menu.GetTheviewer().ChatMenu()
		this.processChatMenu(key)
	case 3:
		this.showMessageHistory()
	case 4:
		os.Exit(0)
	default:
		fmt.Println("错误的输入, 请输入(1-4)")
	}
}

func (this *mainProcessor) processChatMenu(key int) {
	/*
		fmt.Println("\n\n1. 发送群聊")
		fmt.Println("2. 点对点聊天")
		fmt.Println("3. 返回")
	*/
	loop := true
	for loop {
		loop = false
		switch key {
		case 1:
			this.sp = &SmsProcessor{}
			fmt.Println("请输入消息内容, 以$结束")
			content := utils.MyScanf('$')
			this.sp.processGroupChat(this.currUser, this.currConn, &content)
		case 2:
			this.sp = &SmsProcessor{}
			fmt.Println("请输入对端ID")
			var toUserID int
			fmt.Scanf("%d", &toUserID)
			fmt.Println("请输入消息内容, 以$结束")
			content := utils.MyScanf('$')
			this.sp.processPointChat(this.currUser, toUserID, this.currConn, &content)
		case 3:
			// doNothing
		default:
			loop = true
			fmt.Println("请输入1-3")
		}
	}
}

func (this *mainProcessor) showMessageHistory() {

}
func (this *mainProcessor) showOnlineUserTable() {
	// 这里直接根据表来显示
	GetUsersTableControler().PrintTable()
}
