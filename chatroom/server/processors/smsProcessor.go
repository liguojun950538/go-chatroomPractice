package processors

import (
	"chatroom/common"
	"chatroom/common/chatroompkg"
	"chatroom/common/utils"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type SmsProcessor struct {
	*common.UserCanBeSeen
	Conn   net.Conn
	Buffer []byte
}

func (this *SmsProcessor) ProcessPointChatMsg() (err error) {
	// 1. 反序列化聊天包
	var msg chatroompkg.PointChatMsg
	err = json.Unmarshal(this.Buffer, &msg)
	if err != nil {
		log.Println("json.Unmarshal err=", err)
		return
	}

	// content := &msg.Content
	// TODO: 做一些敏感词处理, 目前就是直接转发
	// 重新序列化
	// 获取table
	tableControler := GetUsersTableControler()

	// 这里的处理有点复杂, 主要接受到点对点包的时候, 会有一个响应
	var sendToSrc chatroompkg.PointChatResMsg
	var sendToDes chatroompkg.PointChatToOtherMsg
	// 这里不仅要给源端返回一个包, 目的地也要一个包
	toSendSrcUp, _ := tableControler.GetUserProcessorByID(msg.UserID)
	toSendDesUp, err := tableControler.GetUserProcessorByID(msg.ToUserID)
	if err != nil {
		// 具体定义下错误
		sendToSrc.Code = 300
		sendToSrc.Error = "发送失败, 用户已下线"
	} else {
		sendToSrc.Code = 200
	}
	srcData, _ := json.Marshal(sendToSrc)
	pkg := chatroompkg.Pkg{}
	pkg.Data = string(srcData)

	pkg.Type = chatroompkg.KPointChatResMsgType
	toSend, _ := json.Marshal(pkg)

	srcTf := utils.Transfer{
		Conn: toSendSrcUp.Conn,
	}
	srcTf.SendMsg(toSend)

	sendToDes.Content = msg.Content
	sendToDes.UserCanBeSeen = msg.UserCanBeSeen
	sendToDes.UserID = msg.ToUserID

	desData, _ := json.Marshal(sendToDes)
	pkg.Data = string(desData)
	pkg.Type = chatroompkg.KPointChatResToOtherMsgType
	toSend, _ = json.Marshal(pkg)

	desTf := utils.Transfer{
		Conn: toSendDesUp.Conn,
	}

	desTf.SendMsg(toSend)
	return nil
}

// 处理群聊信息
func (this *SmsProcessor) ProcessGroupChatMsg() (err error) {
	/*
		// 1. 反序列化群聊包
		var msg chatroompkg.GroupChatMsg
		err = json.Unmarshal(this.Buffer, &msg)
		if err != nil {
			log.Println("json.Unmarshal err=", err)
			return
		}
	*/

	// content := &msg.Content
	// TODO: 做一些敏感词处理, 目前就是直接转发
	// 重新序列化
	pkg := chatroompkg.Pkg{}
	pkg.Type = chatroompkg.KGroupChatMsgType
	pkg.Data = string(this.Buffer)
	toSend, err := json.Marshal(pkg)
	if err != nil {
		log.Println("json Marshal groupChatMsg failed")
		return
	}

	// 获取table
	tableControler := GetUsersTableControler()
	// 发送给其他
	sendFailed := 0
	var idx int
	for idx, otherUP := range tableControler.GetTable() {
		if idx == this.UserID {
			continue
		}

		tf := utils.Transfer{
			Conn: otherUP.Conn,
		}
		err = tf.SendMsg(toSend)
		if err != nil {
			log.Println(otherUP.UserID, " sendMsg failed")
			// 这里应该是积累错误
			sendFailed++
		}
	}
	totalSendCnt := idx + 1
	if sendFailed > 0 {
		err = fmt.Errorf("发从%d次, 失败%d次", totalSendCnt, sendFailed)
		return
	}
	return nil
}
