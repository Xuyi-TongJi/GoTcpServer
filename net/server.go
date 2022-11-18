package net

import (
	"fmt"
	"net"
	"server/iface"
)

// Server IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	Name      string
	Address   string
	IPVersion string
	Port      int
	// Server注册的连接对应的处理业务
	Router iface.IRouter
}

// Start 监听，处理业务
func (s *Server) Start() {
	fmt.Printf("[Server START] Server Listener at IP: %s, Port: %d, is starting\n", s.Address, s.Port)

	go func() {
		var cid uint32 = 0
		// 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Address, s.Port))
		if err != nil {
			fmt.Println("[Server ERROR] Resolve tcp address error:", err)
			return
		}
		// 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("[Server ERROR] Listening: ", s.IPVersion, "err: ", err)
		}
		fmt.Printf("[Server START] Start Server %s success at IP: %s, Port: %d, listening\n", s.Name, s.Address, s.Port)
		// 阻塞等待客户端链接，处理客户端链接业务
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Printf("[Server BUSINESS ERROR] Accept error:%s\n", err)
				continue
			}
			// Connection模块，得到处理业务的connection句柄
			socket := NewConnection(conn, cid, s.Router)
			cid += 1
			// 启动当前的连接业务处理
			go socket.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO 将服务器的资源，状态或一些已经开辟的链接信息，进行停止或回收
}

// Serve 阻塞
func (s *Server) Serve() {
	s.Start()

	// TODO 启动服务器之后的额外业务

	// 阻塞
	select {}
}

// AddRouter 添加路由
func (s *Server) AddRouter(r iface.IRouter) {
	s.Router = r
	fmt.Printf("[Server Router] Add server router success\n")
}

// NewServer 初始化Server模块的方法
func NewServer(name string, ipVersion string, address string, port int) iface.IServer {
	return &Server{
		Name:      name,
		IPVersion: ipVersion,
		Address:   address,
		Port:      port,
		Router:    nil,
	}
}
