package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn
}

// 创建一个用户
func NewUser(conn net.Conn) *User{
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
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