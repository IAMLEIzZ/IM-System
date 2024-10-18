package main

import (
	"fmt"
	"net"
	// "runtime"
	"sync"
	"time"
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
	/*
		这里的一个 coon 对应着一个用户，因此对应单个用户消息的处理一个写在 DoHandler 中
	*/
	// 1. 创建新用户
	user := NewUser(conn, this)
	// 2. 用户上线
	user.UserOnline()
	// 是否活跃信号
	isLive := make(chan bool)

	// 接收用户发送的消息boardcast
	go func() {
		// 创建用户发送消息的缓冲区
		buffer := make([]byte, 4096)
		for {
			n, err := conn.Read(buffer)
			// fmt.Println(user, n)
			// buffer 中完全为空
			if n == 0 {
				// 用户下线 
				if _, ok := this.OnlineMap[user.Name]; ok {
					user.UserOffLine()
				}
				return
			} 
			if err != nil {
				fmt.Println("Conn Read err:", err)
				return
			}
			msg := string(buffer[:n - 1])
			// 广播消息
			user.UserDoMessage(msg)

			// 保活
			isLive <- true
		}
	}()
	// 超时强踢功能
	for {
		select{
		// 这里为什么 isLive 后不用操作？
		/*
			在 select 语句中，会并行监听所有 case 的条件，也就是说所有 case 的条件语句会同时执行并堵塞，
			直到启动一个 case 成立，他会停止监听其他 case 条件。在本例中，同时开始 <-isLive 和 <-time.After(time.Second * 10)
			如果 isLive 成立，就算没有操作，计时器也会被重置，因此无需在 isLive 中重置计时器
		*/
		case <-isLive :
			// 当前用户处于活跃态
		case <-time.After(time.Second * 60) :
			// 当前用户 10 秒没有操作，超时踢出
			user.SendMessage("长时间没有活动，踢出聊天室\n")
			// Onlinemap 中删除 user
			user.UserOffLine()
			// 关闭资源(关闭用户通讯channel 和 连接句柄)
			close(user.C)
			conn.Close()
			// 彻底结束一个 user 的连接周期
			fmt.Println(user.Name + "被强制下线\n")
			return
		}
	}
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
