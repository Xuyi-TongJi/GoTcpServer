package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"server/iface"
	"server/network"
	"server/utils"
)

// PingRouter 自定义Router
// Router 就是服务器所能提供的服务
type PingRouter struct {
	network.BaseRouter
	name string
}

// PreHandle Override
/*func (r *PingRouter) PreHandle(req iface.IRequest) {
	fmt.Println("[Router PreHandle] Call router pre handle")
	// socket
	_, err := req.GetConnection().GetTcpConnection().Write([]byte("before handle\n"))
	if err != nil {
		fmt.Printf("[Router PreHandle] Call back before handle, error:%s\n", err)
	}
}
*/

// DoHandle Override
func (r *PingRouter) DoHandle(req iface.IRequest) {
	/* BUSINESS */
	fmt.Printf("[Router Handle] Call router handle, receive message id = %d, data len = %d\n",
		req.GetMsgId(), req.GetDataLen())
	fmt.Printf("[Router Handle] Receive message [%s]\n", req.GetData())
	dataBuffer := bytes.NewBuffer(make([]byte, 0))
	if err := binary.Write(dataBuffer, binary.LittleEndian, []byte("l have received your message: ")); err != nil {
		fmt.Printf("[Router Handle ERROR] Response writing error: %s\n", err)
		return
	}
	if err := binary.Write(dataBuffer, binary.LittleEndian, req.GetData()); err != nil {
		fmt.Printf("[Router Handle ERROR] Response writing error: %s\n", err)
		return
	}
	/* BUSINESS */
	// Send message to writer goroutine
	if err := req.GetConnection().SendMessage(200, dataBuffer.Bytes()); err != nil {
		fmt.Printf("[Router Handle ERROR] Response writing error: %s\n", err)
		return
	}
}

// PostHandle Override
/*func (r *PingRouter) PostHandle(req iface.IRequest) {
	fmt.Println("[Router PostHandle] Call router post handle")
	_, err := req.GetConnection().GetTcpConnection().Write([]byte("after handle\n"))
	if err != nil {
		fmt.Printf("[Router PostHandle] Call back after handle, error:%s\n", err)
	}
	time.Sleep(5 * time.Second)
}*/

type HelloRouter struct {
	*network.BaseRouter
	name string
}

func (r *HelloRouter) DoHandle(req iface.IRequest) {
	/* BUSINESS */
	fmt.Printf("[Router Handle] Call router handle, receive message id = %d, data len = %d\n",
		req.GetMsgId(), req.GetDataLen())
	fmt.Printf("[Router Handle] Receive message [%s]\n", req.GetData())
	dataBuffer := bytes.NewBuffer(make([]byte, 0))
	if err := binary.Write(dataBuffer, binary.LittleEndian, []byte("l have received your message: ")); err != nil {
		fmt.Printf("[Router Handle ERROR] Response writing error: %s\n", err)
		return
	}
	if err := binary.Write(dataBuffer, binary.LittleEndian, req.GetData()); err != nil {
		fmt.Printf("[Router Handle ERROR] Response writing error: %s\n", err)
		return
	}
	/* BUSINESS */
	// Send message to writer goroutine
	if err := req.GetConnection().SendMessage(201, dataBuffer.Bytes()); err != nil {
		fmt.Printf("[Router Handle ERROR] Response writing error: %s\n", err)
		return
	}
}

/*
基于Server框架来开发的服务端应用程序
*/
func main() {
	// 创建一个Server句柄
	s := network.NewServer("tcp4")
	/* Config Router */
	s.AddRouter(0, &PingRouter{name: "ping router"})
	s.AddRouter(1, &HelloRouter{name: "hello router"})
	/* Config Hook */
	// 启动Server
	utils.GlobalObj.TcpServer = s
	s.Serve()
}
