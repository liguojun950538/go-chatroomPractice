package processors

import (
	"chatroom/common"
	"chatroom/common/chatroompkg"
	"chatroom/common/utils"
	"chatroom/server/model"
	"encoding/json"
	"log"
	"net"
)

type UserProcessor struct {
	*common.UserCanBeSeen
	Conn   net.Conn
	Buffer []byte
}

func (this *UserProcessor) ProcessRegisterMsg() (err error) {
	var msg chatroompkg.RegisterMsg
	err = json.Unmarshal(this.Buffer, &msg)
	if err != nil {
		log.Println("json.Unmarshal err=", err)
		return
	}
	// 做返回包
	var resMsg chatroompkg.Pkg
	resMsg.Type = chatroompkg.KRegisterResMsgType
	var registerResMsg chatroompkg.RegisterResMsg

	// 向数据库验证
	err = model.TheUserDao.Register(&msg.User)
	if err != nil {
		switch err {
		case model.ErrUserExist:
			registerResMsg.Code = 600
			registerResMsg.Error = "用户已经存在"
		default:
			registerResMsg.Code = 999
			registerResMsg.Error = "未知错误"
		}
	} else {
		registerResMsg.Code = 200
	}

	// 序列化返回包
	data, err := json.Marshal(registerResMsg)
	if err != nil {
		log.Println("json.Marshal err=", err)
		return
	}
	resMsg.Data = string(data)
	// 序列化msg
	data, err = json.Marshal(resMsg)
	if err != nil {
		log.Println("json.Marshal err=", err)
		return
	}
	// 发回给客户端
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	tf.SendMsg(data)

	return nil
}

func (this *UserProcessor) ProcessLoginMsg() (err error) {
	// 1. 反序列化登陆包
	var msg chatroompkg.LoginMsg
	err = json.Unmarshal(this.Buffer, &msg)
	if err != nil {
		log.Println("json.Unmarshal err=", err)
		return
	}
	// 做返回包
	var resMsg chatroompkg.Pkg
	resMsg.Type = chatroompkg.KLoginResMsgType
	var loginResMsg chatroompkg.LoginResMsg
	/*
		if msg.UserId != "100" || msg.Passwd != "123" {
			// 不合法
			loginResMsg.Code = 500
			loginResMsg.Error = "用户不存在, 或者密码错误"
		} else {
			// 合法
			loginResMsg.Code = 200
		}
	*/

	// 向数据库验证用户名和密码
	_, err = model.TheUserDao.ConfirmLogin(msg.UserId, msg.Passwd)
	if err != nil {
		switch err {
		case model.ErrUserNotExist:
			loginResMsg.Code = 500
			loginResMsg.Error = "用户不存在"
		case model.ErrInvalidPasswd:
			loginResMsg.Code = 501
			loginResMsg.Error = "密码错误"
		}
	} else { // 用户登陆成功
		loginResMsg.Code = 200
		// up记录用户id, 然后将up加入到onlineUsersTable
		this.UserID = msg.UserId

		tableControler := GetUsersTableControler()
		tableControler.AddOnlineUsers(this)
		table := tableControler.GetTable()
		for id, _ := range table {
			if id == msg.UserId {
				continue
			}
			loginResMsg.UserIDSlice = append(loginResMsg.UserIDSlice, id)
		}

		// 向其他用户通知自己上线了
		this.NotifyOtherOnlineUsers(msg.UserId)

	}

	// 2. 序列化loginResMsg
	data, err := json.Marshal(loginResMsg)
	if err != nil {
		log.Println("json.Marshal err=", err)
		return
	}
	resMsg.Data = string(data)
	// 序列化msg
	data, err = json.Marshal(resMsg)
	if err != nil {
		log.Println("json.Marshal err=", err)
		return
	}
	// 3. 将loginResMsg发送给客户端
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	tf.SendMsg(data)

	return nil
}

func (this *UserProcessor) NotifyOtherOnlineUsers(userID int) {
	userTable := GetUsersTableControler()
	msg := &chatroompkg.OnlineNotifyMsg{
		UserID:     userID,
		UserStatus: chatroompkg.KStatusUserOnline,
	}
	for idx, up := range userTable.GetTable() {
		if idx == userID {
			continue
		}
		up.notifyMe(msg)
	}
}

func (this *UserProcessor) notifyMe(msg *chatroompkg.OnlineNotifyMsg) (err error) {
	pack := chatroompkg.Pkg{
		Type: chatroompkg.KOnlineNotifyMsgType,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("json.Marshal OnlineNotifyMsg failed")
		return err
	}
	pack.Data = string(data)
	data, err = json.Marshal(pack)
	if err != nil {
		log.Println("json.Marshal OnlineNotifyMsg failed")
		return err
	}
	tf := utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.SendMsg(data)
	if err != nil {
		log.Println("send OnlineNotifyMsg failed")
		return err
	}
	return nil
}
