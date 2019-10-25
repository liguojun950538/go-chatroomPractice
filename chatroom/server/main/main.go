package main

import (
	"chatroom/server/model"
	"log"
	"net"
	"time"
)

// 初始化注意
// 1. 初始化pool
// 2. 初始化userDao
// 必须

func initUserDao() {
	model.TheUserDao = model.NewUserDao(thePool)
}

func process(conn net.Conn) {
	defer conn.Close()
	log.Println("handling remote conn", conn.RemoteAddr().String())

	mp := &mainProcessor{
		Conn: conn,
	}
	err := mp.process()
	if err != nil {
		log.Printf("process %v failed err = %v\n", conn.RemoteAddr().String(), err)
	}
	log.Println("handling remote conn finished", conn.RemoteAddr().String())
}

func main() {
	initRedis("127.0.0.1:6379", 16, 1024, time.Second*300)
	initUserDao()

	// TODO 设置为从命令行接收的模式
	addr := "0.0.0.0:8888"
	log.Printf("服务器在%v上监听ing...\n", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("net.Listen err=", err)
		return
	}
	defer listener.Close()
	for {
		log.Println("等待客户端连接ing...")
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept err=", err)
		}
		go process(conn)
	}
}
