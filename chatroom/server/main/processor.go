package main

import (
	"chatroom/common"
	"chatroom/common/chatroompkg"
	"chatroom/common/utils"
	"chatroom/server/processors"
	"errors"
	"io"
	"log"
	"net"
)

/*
 * 主的控制器, 每个go程一个, 负责控制处理消息的流程
 */

type mainProcessor struct {
	common.UserCanBeSeen
	Conn net.Conn
}

func (this *mainProcessor) process() (err error) {
	// 这里的循环处理每一个包
	// 只有当收到客户端断开的时候, 或者收到Exit包的时候才会断开连接
	for {
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		// 收到Pkg, 解析出msg信息
		msg, err := tf.RecvPkg()
		if err != nil {
			if err == io.EOF {
				return errors.New("对端已经断开")
			}
			log.Println("readPkg err=", err)
		}
		// log.Println(msg)
		err = this.processMsg(&msg)
		if err == processors.ErrClientExit {
			break
		}
	}
	return nil
}

func (this *mainProcessor) processMsg(chatroompkgMsg *chatroompkg.Pkg) (err error) {

	switch chatroompkgMsg.Type {
	case chatroompkg.KLoginMsgType:
		msgBytes := []byte(chatroompkgMsg.Data)
		up := &processors.UserProcessor{
			Conn:          this.Conn,
			Buffer:        msgBytes,
			UserCanBeSeen: &this.UserCanBeSeen,
		}
		err = up.ProcessLoginMsg()
	case chatroompkg.KRegisterMsgType:
		msgBytes := []byte(chatroompkgMsg.Data)
		up := &processors.UserProcessor{
			Conn:          this.Conn,
			Buffer:        msgBytes,
			UserCanBeSeen: &this.UserCanBeSeen,
		}
		err = up.ProcessRegisterMsg()
	case chatroompkg.KGroupChatMsgType:
		msgBytes := []byte(chatroompkgMsg.Data)
		sp := &processors.SmsProcessor{
			Conn:          this.Conn,
			Buffer:        msgBytes,
			UserCanBeSeen: &this.UserCanBeSeen,
		}
		err = sp.ProcessGroupChatMsg()
	case chatroompkg.KPointChatMsgType:
		msgBytes := []byte(chatroompkgMsg.Data)
		sp := &processors.SmsProcessor{
			Conn:          this.Conn,
			Buffer:        msgBytes,
			UserCanBeSeen: &this.UserCanBeSeen,
		}
		err = sp.ProcessPointChatMsg()
	case chatroompkg.KExitMsgType:
		err = processors.ErrClientExit
	default:
		errString := "消息类型不存在, 无法处理"
		log.Println(errString)
		return errors.New(errString)
	}
	return
}
