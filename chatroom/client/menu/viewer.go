package menu

import (
	"chatroom/common"
	"fmt"
)

var (
	theViewer *viewer
)

type viewer struct {
	common.User
}

func init() {
	theViewer = &viewer{}
}

func GetTheviewer() *viewer {
	return theViewer
}

func (this *viewer) WelcomeMenu() (key int) {
	fmt.Println("\n\n----------------欢迎登陆多人聊天系统------------")
	fmt.Println("\t\t\t 1 登陆聊天室")
	fmt.Println("\t\t\t 2 注册用户")
	fmt.Println("\t\t\t 3 退出系统")
	fmt.Println("\t\t\t 请选择(1-3):")

	fmt.Scanf("%d\n", &key)
	return
}

func (this *viewer) LoginMenu() {
	fmt.Println("请输入用户的id")
	fmt.Scanf("%d\n", &this.User.UserID)
	fmt.Println("请输入用户的密码")
	fmt.Scanf("%s\n", &this.User.Passwd)
}

func (this *viewer) ChatMenu() (key int) {
	fmt.Println("\n\n1. 发送群聊")
	fmt.Println("2. 点对点聊天")
	fmt.Println("3. 返回")

	fmt.Scanf("%d\n", &key)
	return
}

func (this *viewer) MainMenu() (key int) {
	fmt.Println("\n\n1. 显示在线用户列表")
	fmt.Println("2. 发送信息")
	fmt.Println("3. 信息列表")
	fmt.Println("4. 退出系统")

	fmt.Scanf("%d\n", &key)
	return
}

func (this *viewer) RegisterMenu() {
	fmt.Println("请输入用户的id")
	fmt.Scanf("%d\n", &this.User.UserID)
	fmt.Println("请输入用户的密码")
	fmt.Scanf("%s\n", &this.User.Passwd)
	fmt.Println("请输入用户的昵称")
	fmt.Scanf("%s\n", &this.User.UserName)
}

func (this *viewer) ShowMessage(id int, name string, message string) {
	fmt.Printf("\n\nrecv msg from id %v, name %v \n", id, name)
	fmt.Println("message: ---> ", message)
}
