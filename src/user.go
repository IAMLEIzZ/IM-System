package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn

	// 用户申请到的服务器
	server *Server
}

// 创建一个用户
func NewUser(conn net.Conn, server *Server) *User{
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

// 监听当前客户端 chan 中收到的消息，一旦收到消息，则立刻输出给客户端
func (this *User) ListenMessage() {
	for {
		msg := <- this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}

// 用户上线
func (this *User) UserOnline() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	this.server.Boardcast(this, "已上线")
}

// 用户下线
func (this *User) UserOffLine(){
	
	//用户下线,将用户从onlineMap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	this.server.Boardcast(this, "下线")

}

// 给用户本身发送消息
func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg + "\n"))
}

// 用户广播发送消息
func (this *User) UserDoMessage(msg string){
	// 当用户输入的信息为 /who 时，输出当前在线用户的列表，否则当做正常广播消息
	if msg == "/who" {
		this.server.mapLock.Lock()
		for _, value := range this.server.OnlineMap {
			msg := "[" + value.Addr + "]" + value.Name + ": 在线"
			this.SendMessage(msg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 8 && msg[:8] == "/rename|"{
		// /rename｜xxx
		newName := strings.Split(msg, "|")[1]

		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMessage("当前用户名被占用\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMessage("已将您的用户名更改为：" + this.Name + "\n")
		}

	} else {
		// 广播发送消息
		this.server.Boardcast(this, msg)
	}
}