package net

import (
	"fmt"
	"net"
	"server/iface"
)

type Connection struct {
	// 当前连接的TCP套接字
	Conn *net.TCPConn

	// 连接ID
	ConnId uint32

	// 当前的连接状态
	IsClosed bool

	// 告知当前丽娜姐已经退出/停止的 channel
	ExitChan chan bool

	// 该连接处理的方法Router
	Router iface.IRouter
}

// startReader 从当前连接读数据的业务
func (c *Connection) startReader() {
	fmt.Printf("[Connection Reader Goroutine] Reader Gouroutine is Running... Romote addr= %s", c.GetClientTcpStatus())
	defer c.Stop()

	for {
		// read data from client to buffer and call the handle function
		buf := make([]byte, 512)
		count, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Printf("[Connection Reader Goroutine ERROR] Receive Buffer from %s error: %s", c.GetClientTcpStatus(), err)
			continue
		}
		// 得到当前conn数据的Request
		req := Request{
			conn: c,
			data: buf,
			len:  count,
		}
		go func() {
			c.Router.PreHandle(&req)
			c.Router.DoHandle(&req)
			c.Router.PostHandle(&req)
		}()
	}
}

// Start 启动连接，业务逻辑是启动一个读数据业务和一个写数据的业务
func (c *Connection) Start() {
	fmt.Printf("[Connection START] Connection %d starting\n", c.ConnId)
	// 启动从当前连接读数据的业务
	go c.startReader()
	// TODO 启动从当前连接写数据的业务
}

func (c *Connection) Stop() {
	fmt.Printf("[Connection STOP] Connection %d stopping\n", c.ConnId)
	if c.IsClosed {
		return
	}
	c.IsClosed = true
	// 回收资源
	err := c.Conn.Close()
	if err != nil {
		fmt.Printf("[Connection STOP ERROR] Connection %d stopped error:%s\n", c.ConnId, err)
	}
	fmt.Printf("[Connection STOP] Connection %d stopped success", c.ConnId)
	close(c.ExitChan)
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

func (c *Connection) GetClientTcpStatus() net.Addr {
	return c.GetTcpConnection().RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	//TODO implement me
	panic("implement me")
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, id uint32, router iface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnId:   id,
		IsClosed: false,
		Router:   router,
		// 有缓冲的管道
		ExitChan: make(chan bool, 1),
	}
	return c
}
