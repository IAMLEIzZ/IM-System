package main

func main() {
	// 启动服务器
	server := NewServer("10.249.85.146", 9090)
	server.Start()
}
