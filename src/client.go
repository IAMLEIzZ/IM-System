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
}

func NewClient(severIp string, severPort int) *Client {
	client := &Client{
		ServerIP: severIp,
		ServerPort: severPort,
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

	for {
	
	}
}