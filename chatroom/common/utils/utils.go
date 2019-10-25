package utils

import (
	"chatroom/common/chatroompkg"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

type Transfer struct {
	Conn   net.Conn
	Buffer [8096]byte
}

// TODO 原生的read函数不能保证我们一定读取到N个字节, 需要新的保证, 目前能想到的已知的最好解决方法是做一个readN出来
// z
func (this *Transfer) RecvPkg() (mes chatroompkg.Pkg, err error) {
	// 1. 读取消息的长度
	readN, err := this.Conn.Read(this.Buffer[:4])
	if readN != 4 || err != nil {
		log.Printf("this.Conn.Read len = %v err=%v\n", readN, err)
		return mes, err
	}
	pkgLen := binary.BigEndian.Uint32(this.Buffer[:4])
	// 2. 读取消息体的内容
	readN, err = this.Conn.Read(this.Buffer[:pkgLen])
	if readN != int(pkgLen) || err != nil {
		log.Println("this.Conn.Read 消息体 err=", err)
		return mes, err
	}

	// 3. 反序列化pkg
	err = json.Unmarshal(this.Buffer[:pkgLen], &mes)
	if err != nil {
		log.Println("json.Unmarshal err=", err)
		return mes, err
	}
	return mes, nil
}

// SendMsg
// 这里主要任务就是算一下长度, 然后发送长度, 然后发送this.Buffer
func (this *Transfer) SendMsg(buf []byte) (err error) {
	pkgLen := len(buf)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(pkgLen))
	// 发送长度
	writenN, err := this.Conn.Write(lenBuf)
	if writenN != 4 || err != nil {
		return
	}

	// 发送具体包的内容
	writenN, err = this.Conn.Write(buf)
	if writenN != pkgLen || err != nil {
		return
	}
	return nil
}

func MyScanf(delim byte) string {
	loop := true
	var buffer []string
	for loop {
		var temp string
		fmt.Scanf("%s", &temp)
		idx := strings.IndexByte(temp, delim)
		if idx != -1 {
			loop = false
			temp = temp[0:idx]
		}
		buffer = append(buffer, temp)
	}
	res := strings.Join(buffer, " ")
	return res
}
