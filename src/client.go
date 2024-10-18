package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIP string
	ServerPort int
	conn net.Conn
	Name string
	flag int
}

func (this *Client) menu() bool {
	var flag int

	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println(">>>请输入合法的模式<<<")
		return false
	}
}

func (this *Client) Run() {
	for this.flag != 0{
		for this.menu() != true {
		
		}

		switch this.flag {
		case 1:
			fmt.Println("选择公聊模式")
			break
		case 2:
			fmt.Println("选择私聊模式")
			break
		case 3:
			fmt.Println("选择更改用户名")
			break
		case 0:
			fmt.Println("退出系统")
			break
		}
	}
	/*
		这里 select 执行的逻辑是，如果我 client 输入一直是不合法的，则会一直在最外侧循环；
		如果合法则会进入内层 menu 循环，直到输入模式为合法范围。此时输入 123 都是正常处理业务，
		而输入 0 则会导致this.flag变为 0，此时外层循环直接退出
	*/
}

func NewClient(severIp string, severPort int) *Client {
	client := &Client{
		ServerIP: severIp,
		ServerPort: severPort,
		flag: 1024,
	}

	// 连接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", severIp, severPort))

	if err != nil {
		fmt.Println("net.Dial error", err)
		return nil
	}

	client.conn = conn

	return client
}

var serverIp string
var serverPort int 

func init() {
	flag.StringVar(&serverIp, "ip", "10.249.85.146", "设置服务器 ip 地址")
	flag.IntVar(&serverPort, "port", 9090, "设置服务器端口号")
}

func main() {

	flag.Parse()

	client := NewClient(serverIp, serverPort)

	if client == nil {
		fmt.Println(">>>>>> 服务器连接失败...")
		return
	}

	fmt.Println(">>>>>> 服务器连接成功...")

	client.Run()
}