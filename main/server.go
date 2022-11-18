package main

import (
	"server/config"
	"server/net"
)

/*
基于Server框架来开发的服务端应用程序
*/
func main() {
	// 创建一个Server句柄
	s := net.NewServer("my_server", config.IPVersion, config.Address, config.Port)
	// 启动Server
	s.Serve()
}
