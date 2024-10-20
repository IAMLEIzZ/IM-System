package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIP string
	ServerPort int
	conn net.Conn
	Name string
	flag int
}

func (this *Client) DealResponse() {
	io.Copy(os.Stdout, this.conn)
}

func (this *Client) updateName() bool {
	fmt.Println(">>>>请输入修改后的用户名：")
	fmt.Scanln(&this.Name)

	sendMsg := "/rename|" + this.Name + "\n"
	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return false
	}

	return true
}

func (this *Client) PublicChat() {
	var chatMessage string

	fmt.Println(">>>>>请输入要发送的消息：, 输入 exit. 表示退出")
	fmt.Scanln(&chatMessage)

	for chatMessage != "exit" {

		if len(chatMessage) != 0 {
			sendMsg := chatMessage + "\n"
			_, err := this.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}	
		}

		chatMessage = ""
		fmt.Println(">>>>>请输入要发送的消息：, 输入 exit. 表示退出")
		fmt.Scanln(&chatMessage)
	}
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
			this.PublicChat()
			break
		case 2:
			fmt.Println("选择私聊模式")
			break
		case 3:
			this.updateName()
			break
		case 0:
			fmt.Println("退出系统")
			break
		}
	}
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
	
	go client.DealResponse()

	fmt.Println(">>>>>> 服务器连接成功...")

	client.Run()
}