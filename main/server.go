package main

import (
	"fmt"
	"server/config"
	"server/iface"
	"server/net"
)

// PingRouter 自定义Router
type PingRouter struct {
	net.BaseRouter
}

// PreHandle Override
func (r *PingRouter) PreHandle(req iface.IRequest) {
	fmt.Println("[Router PreHandle] Call router pre handle")
	// socket
	_, err := req.GetConnection().GetTcpConnection().Write([]byte("before handle\n"))
	if err != nil {
		fmt.Printf("[Router PreHandle] Call back before handle, error:%s\n", err)
	}
}

// DoHandle Override
func (r *PingRouter) DoHandle(req iface.IRequest) {
	fmt.Printf("[Router Handle] Call router pre handle")
	_, err := req.GetConnection().GetTcpConnection().Write([]byte("handling\n"))
	if err != nil {
		fmt.Printf("[Router Handle] Call back when handling, error:%s\n", err)
	}
}

// PostHandle Override
func (r *PingRouter) PostHandle(req iface.IRequest) {
	fmt.Println("[Router PostHandle] Call router post handle")
	_, err := req.GetConnection().GetTcpConnection().Write([]byte("after handle\n"))
	if err != nil {
		fmt.Printf("[Router PostHandle] Call back after handle, error:%s\n", err)
	}
}

/*
	基于Server框架来开发的服务端应用程序
*/

func main() {
	// 创建一个Server句柄
	s := net.NewServer("my_server", config.IPVersion, config.Address, config.Port)
	s.AddRouter(&PingRouter{})
	// 启动Server
	s.Serve()
}
