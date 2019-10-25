package processors

import (
	"chatroom/common"
	"chatroom/common/chatroompkg"
	"chatroom/common/utils"
	"encoding/json"
	"log"
	"net"
)

type SmsProcessor struct {
}

func (this *SmsProcessor) processPointChat(user *common.User, toUserId int, conn net.Conn, content *string) (err error) {
	var msg chatroompkg.Pkg
	msg.Type = chatroompkg.KPointChatMsgType

	var chatMsg chatroompkg.PointChatMsg
	chatMsg.UserCanBeSeen = user.UserCanBeSeen
	chatMsg.ToUserID = toUserId
	chatMsg.Content = *content

	byteSlice, err := json.Marshal(chatMsg)
	if err != nil {
		log.Println("json.Marshal pointChatMsg err=", err)
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
		log.Println("groupSendMsg send failed, err=", err)
		return
	}
	return nil
}

// 这里也需要设计一个协议
// 包: 自身信息 + 内容
func (this *SmsProcessor) processGroupChat(user *common.User, conn net.Conn, content *string) (err error) {
	var msg chatroompkg.Pkg
	msg.Type = chatroompkg.KGroupChatMsgType

	var groupChatMsg chatroompkg.GroupChatMsg
	groupChatMsg.UserCanBeSeen = user.UserCanBeSeen

	groupChatMsg.Content = *content

	byteSlice, err := json.Marshal(groupChatMsg)
	if err != nil {
		log.Println("json.Marshal groupChatMsg err=", err)
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
		log.Println("groupSendMsg send failed, err=", err)
		return
	}
	return nil
}
