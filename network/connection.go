package network

import (
	"errors"
	"fmt"
	"io"
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

	// 告知当前连接已经退出/停止的 channel
	ExitChan chan bool

	// 消息处理器
	MsgHandler iface.IMessageHandler
}

// startReader 从当前连接读数据的业务
func (c *Connection) startReader() {
	fmt.Printf("[Connection Reader Goroutine] Reader Gouroutine is Running... Romote addr= %s\n", c.GetClientTcpStatus())
	defer c.Stop()

	for {
		// read data from client to buffer and call the handle function
		// 拆包 -> Message
		dp := &DataPack{}
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.Conn, headData)
		if err != nil {
			fmt.Printf("[Connection Reader Goroutine ERROR] Connection %d, error reading head data, err:%s\n", c.ConnId, err)
			break
		}
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Printf("[Connection Reader Goroutine ERROR] Connection %d, invalid message id or data, message id = %d, len = %d, err:%s\n",
				c.ConnId, msg.GetMsgId(), msg.GetLen(), err)
			break
		}
		// read data by the tag of data len
		if msg.GetLen() > 0 {
			msg.SetData(make([]byte, msg.GetLen()))
			_, err = io.ReadFull(c.Conn, msg.GetData())
			if err != nil {
				fmt.Printf("[Connection Reader Goroutine ERROR] Connection %d, invalid message id or data, message id = %d, len = %d, err:%s\n",
					c.ConnId, msg.GetMsgId(), msg.GetLen(), err)
				break
			}
		}
		// 得到当前conn数据的Request
		req := Request{
			conn:    c,
			message: msg,
		}
		// go 处理这个Request(Router中有具体的业务逻辑)
		go c.MsgHandler.DoHandle(&req)
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
	fmt.Printf("[Connection STOP] Connection %d stopped success\n", c.ConnId)
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

// SendMessage Send 将数据封包并回写给服务端
func (c *Connection) SendMessage(msgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New(fmt.Sprintf("[Connection Writing GoRoutine] Connection %d was closed\n", c.ConnId))
	}
	dp := DataPack{}
	msg := &Message{
		Id:   msgId,
		Len:  uint32(len(data)),
		Data: data,
	}
	// ID
	binaryData, err := dp.Pack(msg)
	if err != nil {
		return errors.New(fmt.Sprintf("[Connection Writing GoRoutine] Connection %d, packing message error: %s\n", c.ConnId, err))
	}
	if _, err = c.Conn.Write(binaryData); err != nil {
		return errors.New(fmt.Sprintf("[Connection Writing GoRoutine] Connection %d, write pipe error: %s\n", c.ConnId, err))
	}
	return nil
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, id uint32, msgHandler iface.IMessageHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnId:     id,
		IsClosed:   false,
		MsgHandler: msgHandler,
		// 有缓冲的管道
		ExitChan: make(chan bool, 1),
	}
	return c
}
