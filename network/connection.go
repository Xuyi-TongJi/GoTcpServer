package network

import (
	"errors"
	"fmt"
	"io"
	"net"
	"server/iface"
	"server/utils"
	"sync"
)

type Connection struct {

	// 当前Connection属于哪个Server
	TcpServer iface.IServer

	// 当前连接的TCP套接字
	Conn *net.TCPConn

	// 连接ID
	ConnId uint32

	// 当前的连接状态
	IsClosed bool

	// 告知当前连接已经退出/停止的 channel
	ExitChan chan bool

	// 无缓冲管道，用于读写go routine之间的消息通信
	MessageChan chan []byte

	// 消息处理器，有缓冲的管道
	MsgHandler iface.IMessageHandler

	// 连接属性（集合）
	PropertyMap map[string]interface{}

	// 属性集合锁
	PropertyLock sync.RWMutex
}

// startReader 从当前连接读数据的模块
func (c *Connection) startReader() {
	fmt.Printf("[Connection Reader Goroutine] Connection %d reader gouroutine is running. Romote addr = %s\n",
		c.ConnId, c.GetClientTcpStatus().String())
	defer func() {
		fmt.Printf("[Connection Reader Goroutine] %s Connection %d was closed, reader goroutine closed\n",
			c.GetClientTcpStatus().String(), c.ConnId)
		// reader goroutine exit 关闭连接 并向writer goroutine 发送消息让其退出
		c.Stop()
	}()

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
		if utils.GlobalObj.WorkerPoolSize > 0 {
			// 将请求交给Message Handler 执行具体的业务逻辑
			c.MsgHandler.SubmitTask(&req)
		} else {
			go c.MsgHandler.DoHandle(&req)
		}
	}
}

// startWriter 向当前连接写数据的模块
func (c *Connection) startWriter() {
	fmt.Printf("[Connection Writer Goroutine] Connection %d writer gouroutine is running. Romote addr = %s\n",
		c.ConnId, c.GetClientTcpStatus().String())
	defer fmt.Printf("[Connection Writer Goroutine] %s Connection %d was closed, writer goroutine closed\n",
		c.GetClientTcpStatus().String(), c.GetConnId())
	for {
		select {
		case data := <-c.MessageChan:
			if _, err := c.GetTcpConnection().Write(data); err != nil {
				fmt.Printf("[Connection Writer Goroutine ERROR] Connection %d writing back error: %s\n", c.ConnId, err)
				return
			} else {
				fmt.Printf("[Connection Writer Goroutine] Connection %d writing back to the client success\n",
					c.ConnId)
			}
		case <-c.ExitChan:
			// Reader已经退出
			return
		}
	}
}

// Start 启动连接，业务逻辑是启动一个读数据业务和一个写数据的业务
func (c *Connection) Start() {
	fmt.Printf("[Connection START] Connection %d starting\n", c.ConnId)
	// 启动从当前连接读数据的业务
	go c.startReader()
	// 启动从当前连接写数据的业务
	go c.startWriter()
	// hook OnConnectionStart
	c.TcpServer.CallOnConnectionStart(c)
}

func (c *Connection) Stop() {
	fmt.Printf("[Connection STOP] Connection %d stopping\n", c.ConnId)
	if c.IsClosed {
		return
	}
	c.TcpServer.CallOnConnectionStop(c)
	// 告知Writer关闭
	c.ExitChan <- true
	// 从连接管理器中删除
	c.TcpServer.GetConnectionManager().Remove(c)
	// 回收资源
	c.IsClosed = true
	err := c.Conn.Close()
	if err != nil {
		fmt.Printf("[Connection STOP ERROR] Connection %d stopped error:%s\n", c.ConnId, err)
	}
	fmt.Printf("[Connection STOP] Connection %d stopped success\n", c.ConnId)
	close(c.MessageChan)
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

// SendMessage 将数据封包为二进制数据并发送给写协程
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
	// pack (message to binary data)
	binaryData, err := dp.Pack(msg)
	if err != nil {
		return errors.New(fmt.Sprintf("[Connection Writing GoRoutine] Connection %d, packing message error: %s\n", c.ConnId, err))
	}
	// 将data发送给写协程
	c.MessageChan <- binaryData
	return nil
}

func (c *Connection) SetConnectionProperty(key string, value interface{}) {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	c.PropertyMap[key] = value
}

func (c *Connection) GetConnectionProperty(key string) interface{} {
	c.PropertyLock.RLock()
	defer c.PropertyLock.RUnlock()
	if value, ok := c.PropertyMap[key]; ok {
		return value
	} else {
		return nil
	}
}

func (c *Connection) RemoveConnectionProperty(key string) {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	if _, ok := c.PropertyMap[key]; ok {
		delete(c.PropertyMap, key)
	}
}

// NewConnection 初始化连接模块的方法
func NewConnection(server iface.IServer, conn *net.TCPConn, id uint32, msgHandler iface.IMessageHandler) *Connection {
	c := &Connection{
		TcpServer:   server,
		Conn:        conn,
		ConnId:      id,
		IsClosed:    false,
		MsgHandler:  msgHandler,
		MessageChan: make(chan []byte),
		ExitChan:    make(chan bool, 1),
		PropertyMap: make(map[string]interface{}),
	}
	// Add操作一定是串行的
	c.TcpServer.GetConnectionManager().Add(c)
	return c
}
