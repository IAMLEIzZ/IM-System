package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip string
	Port int
	OnlineMap map[string]*User
	mapLock sync.RWMutex

	Message chan string
}

// 将消息放进 server 消息广播管道之中
func (this *Server) Boardcast(user *User, msg string) {
	userMsg := "[" + user.Name +"]:" +  msg

	
	this.Message <- userMsg
}

// 监听当前的 Message 管道，一旦有用户上线，则将该用户上线的消息广播到每一个在线用户
func (this *Server) ListenMessage() {
	for {
		msg := <- this.Message
		// 循环遍历当前的在线用户
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			// 给读到的用户发送消息
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 处理请求
func (this *Server) DoHandler(conn net.Conn) {
	// fmt.Println("处理请求ing")
	// 1. 创建新用户
	user := NewUser(conn)
	// 2. 将上线的用户注册到 onlineMap 中
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()
	// 3. 向在线的所有用户广播该用户已上线
	this.Boardcast(user, "上线")

	// select{}
}

// 创建 server 方法
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}

	return server
}

// 启动 server 方法
func (this *Server) Start() {
	// 监听端口
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil{
		fmt.Println("链接失败", err) 
		return
	}
	// 关闭监听
	defer l.Close()
	// 启动广播管道监听
	go this.ListenMessage()
	// 接受请求
	for {
		conn, err := l.Accept()		
		if err != nil {
			fmt.Println("接收消息失败", err)
			continue
		}

		// do handler
		// 如果成功接收，则代表一个用户上线
		go this.DoHandler(conn)
	}
}

/*
	流程梳理
	1. server 启动后持续监听上线的用户和消息管道
	2. 当我命令行 nc 后，此时代表一个用户上线，在 server 中监听到用户上线后，在全局用户表中注册该用户
	3. 在注册该用户后，server 服务器因为要将消息全局广播给所有的客户端，所以需要先将该上线用户的消息 push 到 chan 中
	4. chan 中得到用户上线的消息后，此时，this.ListenMessage() 会停止阻塞，将消息送到每个在线用户的 chan 中
	5. 每个在线用户的 chan 得到消息后，user.ListenMessage() 会停止阻塞，将收到的消息传到显示在客户端（conn.Write([]byte(msg + "\n"))）
*/
