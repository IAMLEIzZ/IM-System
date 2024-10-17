package main

import (
	"net"
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
	this.server.Boardcast(this, "下线")
}

// 用户广播发送消息
func (this *User) UserBoardCast(msg string){
	this.server.Boardcast(this, msg)
}