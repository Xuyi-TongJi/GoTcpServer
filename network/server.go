package network

import (
	"fmt"
	"net"
	"server/iface"
	"server/utils"
)

// Server IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	Name           string
	Address        string
	IPVersion      string
	Version        string
	Port           int
	MaxConn        int
	MaxPackingSize uint32
	MsgHandler     iface.IMessageHandler
	ConnManager    iface.IConnectionManager
	// Hook
	OnConnStart func(c iface.IConnection)
	// Hook
	OnConnStop func(c iface.IConnection)
}

func logConfig() {
	fmt.Println("[Server Config] Server config success")
	fmt.Println("Name: ", utils.GlobalObj.Name)
	fmt.Println("Version: ", utils.GlobalObj.Version)
	fmt.Println("Host: ", utils.GlobalObj.Host)
	fmt.Println("Port: ", utils.GlobalObj.TcpPort)
	fmt.Println("MaxConn: ", utils.GlobalObj.MaxConn)
	fmt.Println("MaxPackagingSize: ", utils.GlobalObj.MaxPackingSize)
	fmt.Println("WorkerPoolSize: ", utils.GlobalObj.WorkerPoolSize)
	fmt.Printf("[Server START] Server Listener at IP: %s, Port: %d, is starting\n",
		utils.GlobalObj.Host, utils.GlobalObj.TcpPort)
}

// Start 监听，处理业务
func (s *Server) Start() {
	logConfig()
	go func() {
		var cid uint32 = 0
		// 开启工作池及其消息队列
		s.MsgHandler.StartWorkerPool()
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
				fmt.Printf("[Server Listener ERROR] Accept error:%s\n", err)
				continue
			}
			// 判断连接是否超过最大连接数量, 超过则拒绝连接
			if total := s.ConnManager.Total(); total >= utils.GlobalObj.MaxConn {
				// TODO 给客户端响应一个超出最大连接错误报告
				fmt.Printf("[Server Connection REFUSED] Connection Refused, there are %d current connection alive\n", total)
				_ = conn.Close()
				continue
			}
			// Connection模块，得到处理业务的connection句柄
			socket := NewConnection(s, conn, cid, s.MsgHandler)
			cid += 1
			// 启动当前的连接业务处理
			go socket.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO 将服务器的资源，状态或一些已经开辟的链接信息，进行停止或回收
	fmt.Printf("[Server Stop] Server is ready to stop\n")
	// 清除所有connection
	s.ConnManager.ClearAll()
}

// Serve 启动Server
func (s *Server) Serve() {
	s.Start()
	//TODO 启动服务器之后的额外业务

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgId uint32, r iface.IRouter) {
	s.MsgHandler.AddRouter(msgId, r)
}

func (s *Server) GetConnectionManager() iface.IConnectionManager {
	return s.ConnManager
}

func (s *Server) SetOnConnectionStart(hook func(connection iface.IConnection)) {
	s.OnConnStart = hook
}

func (s *Server) SetOnConnectionStop(hook func(connection iface.IConnection)) {
	s.OnConnStop = hook
}

func (s *Server) CallOnConnectionStart(connection iface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(connection)
	}
}

func (s *Server) CallOnConnectionStop(connection iface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(connection)
	}
}

// NewServer 初始化Server模块的方法
func NewServer(ipVersion string) iface.IServer {
	return &Server{
		Name:           utils.GlobalObj.Name,
		IPVersion:      ipVersion,
		Address:        utils.GlobalObj.Host,
		Port:           utils.GlobalObj.TcpPort,
		Version:        utils.GlobalObj.Version,
		MaxConn:        utils.GlobalObj.MaxConn,
		MaxPackingSize: utils.GlobalObj.MaxPackingSize,
		MsgHandler:     NewMessageHandler(),
		ConnManager:    NewConnectionManager(),
	}
}
