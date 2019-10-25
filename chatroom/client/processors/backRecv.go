package processors

import (
	"chatroom/client/menu"
	"chatroom/common"
	"chatroom/common/chatroompkg"
	"chatroom/common/utils"
	"encoding/json"
	"log"
	"net"
)

// 登陆成功后启动, 处理后台的推送消息
func backRecv(conn net.Conn) {
	tf := utils.Transfer{
		Conn: conn,
	}
	for true {
		pkg, err := tf.RecvPkg()
		if err != nil {

		}
		switch pkg.Type {
		case chatroompkg.KOnlineNotifyMsgType:
			noticeMsg := &chatroompkg.OnlineNotifyMsg{}
			err := json.Unmarshal([]byte(pkg.Data), noticeMsg)
			if err != nil {
				log.Println("json.Unmarshal OnlineNotifyMsg failed")
				return
			}

			updateOnlineUserTable(noticeMsg)
		case chatroompkg.KGroupChatMsgType:
			groupChatMsg := &chatroompkg.GroupChatMsg{}
			err := json.Unmarshal([]byte(pkg.Data), groupChatMsg)
			if err != nil {
				log.Println("json.Unmarshal GroupChatMsg failed")
				return
			}
			showPeerMessage(groupChatMsg.UserID, groupChatMsg.UserName, groupChatMsg.Content)
		case chatroompkg.KPointChatResMsgType:
			// 这里是告诉自己发送给别人的消息是否成功
			// doNothing()
		case chatroompkg.KPointChatResToOtherMsgType:
			// 这里是接收到他人的点对待你信息
			pointChatMsg := &chatroompkg.PointChatToOtherMsg{}
			err := json.Unmarshal([]byte(pkg.Data), pointChatMsg)
			if err != nil {
				log.Println("json.Unmarshal pointChatMsg failed")
				return
			}
			showPeerMessage(pointChatMsg.UserID, pointChatMsg.UserName, pointChatMsg.Content)
		default:
			// doNothing()
		}
	}
}

func updateOnlineUserTable(msg *chatroompkg.OnlineNotifyMsg) {
	tableControl := GetUsersTableControler()

	user := &common.User{
		UserCanBeSeen: common.UserCanBeSeen{
			UserID:     msg.UserID,
			UserStatus: msg.UserStatus,
		},
	}

	// 直接更新表就可以了
	tableControl.AddOnlineUsers(user)
}

func showPeerMessage(id int, name string, message string) {
	theViewer := menu.GetTheviewer()
	theViewer.ShowMessage(id, name, message)
}
