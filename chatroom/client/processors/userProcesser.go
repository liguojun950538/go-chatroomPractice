package processors

import (
	"chatroom/client/menu"
	"chatroom/common"
	"chatroom/common/chatroompkg"
	"chatroom/common/utils"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type UserProcessor struct {
}

func (this *UserProcessor) Register(user common.User) {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil { // 如果连接服务器失败, 直接返回 TODO: 返回连接服务器失败错误
		log.Println("net.Dial err=", err)
		return
	}
	defer conn.Close()

	var msg chatroompkg.Pkg
	msg.Type = chatroompkg.KRegisterMsgType

	var registerMsg chatroompkg.RegisterMsg
	registerMsg.User = user

	byteSlice, err := json.Marshal(registerMsg)
	if err != nil {
		log.Println("json.Marshal registerMsg err=", err)
		return
	}
	msg.Data = string(byteSlice)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("json.Marshal Msg err=", err)
		return
	}
	tf := &utils.Transfer{
		Conn: conn,
	}
	err = tf.SendMsg(msgBytes)
	if err != nil {
		log.Println("RegisPkg send failed, err=", err)
		return
	}

	this.processRegisterRecv(conn)
}

func (this *UserProcessor) processRegisterRecv(conn net.Conn) {
	tf := &utils.Transfer{
		Conn: conn,
	}
	recvPkg, err := tf.RecvPkg()
	if err != nil {
		log.Println("registerRecvPkg recv failed, err=", err)
		return
	}
	// assert(recvPkg.Type == chatroompkg.KLoginResMsgType)

	var registerResMsg chatroompkg.RegisterResMsg
	err = json.Unmarshal([]byte(recvPkg.Data), &registerResMsg)
	if err != nil {
		log.Println("registerRecvPkg unmarshal failed, err=", err)
	}
	switch registerResMsg.Code {
	case 200:
		log.Println("register success")
	case 600:
		log.Println("register failure, err=", registerResMsg.Error)
	default:
		log.Println("register failure, err=", registerResMsg.Error)
	}
	return
}

func (this *UserProcessor) Login(user common.User) {
	// 需要设置协议
	// 1. 连接到服务器
	// TODO: 改为从配置文件中读取地址
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil { // 如果连接服务器失败, 直接返回 TODO: 返回连接服务器失败错误
		log.Println("net.Dial err=", err)
		return
	}
	defer conn.Close()

	err = this.processLoginMsgSend(&user, conn)
	if err != nil {
		log.Println("发送和处理登陆请求包failed\nerr=", err)
	}

	err = this.processLoginMsgRecv(&user, conn)
	if err != nil {
		log.Println("接收和处理登陆返回响应包failed\nerr=", err)
	}
}
func (this *UserProcessor) processLoginMsgSend(user *common.User, conn net.Conn) (err error) {
	// 2. 准备发送消息给服务器
	// 序列化消息, 以前都是用用的那种长度协议, 这里的利用json来序列化协议没怎么解除过
	// 发送的消息体就是msg, 然后在前面插一个len
	// 如何序列化?
	// 2.1 我们现在发送的一个消息实际上是这样的一种类型 msg.Type + msg.Data
	// msg.Data = LoginMsg, 即来个嵌套的结构体
	// 2.2 序列化我们的loginMsg, 然后将其转为string后, 赋值给msg.Data
	var msg chatroompkg.Pkg
	msg.Type = chatroompkg.KLoginMsgType

	var loginMsg chatroompkg.LoginMsg
	loginMsg.UserId = user.UserID
	loginMsg.Passwd = user.Passwd
	byteSlice, err := json.Marshal(loginMsg)
	if err != nil {
		log.Println("json.Marshal loginMsg err=", err)
		return
	}
	msg.Data = string(byteSlice)
	// 2.3 序列化Msg
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("json.Marshal Msg err=", err)
		return
	}
	// 2.4 发送消息的长度, 这个如果在c++中非常好处理, 用bigEndian转换网络字节序, 这里同理
	// conn.Write(len(data)); error
	tf := &utils.Transfer{
		Conn: conn,
	}
	err = tf.SendMsg(msgBytes)
	if err != nil {
		log.Println("loginPkg send failed, err=", err)
		return
	}
	return nil
}

func (this *UserProcessor) processLoginMsgRecv(currUser *common.User, conn net.Conn) (err error) {
	tf := &utils.Transfer{
		Conn: conn,
	}
	recvPkg, err := tf.RecvPkg()
	if err != nil {
		log.Println("loginRecvPkg recv failed, err=", err)
		return
	}
	// assert(recvPkg.Type == chatroompkg.KLoginResMsgType)

	loginRecvMsg := &chatroompkg.LoginResMsg{}
	err = json.Unmarshal([]byte(recvPkg.Data), loginRecvMsg)
	if err != nil {
		log.Println("loginRecvPkg unmarshal failed, err=", err)
	}

	// 如果成功应该进入下一级, 否则应该返回上一级
	switch loginRecvMsg.Code {
	case 200:
		log.Println("login success")
	case 500:
		log.Println("login failure, err=", loginRecvMsg.Error)
	default:
		log.Println("login failure, err=", loginRecvMsg.Error)
	}
	if loginRecvMsg.Code != 200 {
		return
	}

	// 登陆成功后, 有很多要初始化的东西
	this.initUser(currUser, conn)
	this.initGlobalInfo(loginRecvMsg)

	go backRecv(conn)
	viewer := menu.GetTheviewer()
	for true {
		key := viewer.MainMenu()
		if key == 4 {
			os.Exit(0)
		}
		GetTheMainProcessor().processMainMenu(key)
	}

	return
}

func (this *UserProcessor) initUser(user *common.User, conn net.Conn) {
	mainP := GetTheMainProcessor()
	mainP.currConn = conn
	mainP.currUser = user
}
func (this *UserProcessor) initGlobalInfo(loginRecvMsg *chatroompkg.LoginResMsg) {
	// 1. 初始化自己的User信息, 有必要么?

	// 2. 初始化在线用户列表
	// 这里要显示用户在线列表
	fmt.Println("\n\n\t\t当前在线用户")
	for i, v := range loginRecvMsg.UserIDSlice {
		fmt.Printf("%d: id --- %d\n", i, v)

		// 完成tableControler的初始化
		onlineUserTableControler := GetUsersTableControler()
		currUser := &common.User{
			UserCanBeSeen: common.UserCanBeSeen{
				UserID:     v,
				UserStatus: chatroompkg.KStatusUserOffline,
			},
		}
		onlineUserTableControler.AddOnlineUsers(currUser)
	}
}
