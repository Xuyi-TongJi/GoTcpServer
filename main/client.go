package main

import (
	"fmt"
	"net"
	"server/utils"
	"time"
)

func main() {
	// 连接服务器，得到一个conn连接
	fmt.Println("[Client START]Client start...")
	time.Sleep(1 * time.Second)
	fmt.Println(fmt.Sprintf("%s:%d", utils.GlobalObj.Host, utils.GlobalObj.TcpPort))
	conn, err := net.Dial("tcp",
		fmt.Sprintf("%s:%d", utils.GlobalObj.Host, utils.GlobalObj.TcpPort))
	if err != nil {
		fmt.Printf("[Client ERROR]Connection errer:%s\n", err)
		return
	}
	for {
		_, err := conn.Write([]byte("hello"))
		if err != nil {
			fmt.Printf("[Client ERROR]Writing error:%s\n", err)
			continue
		}
		buf := make([]byte, 512)
		count, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("[Client ERROR]Reading buffer error:%s\n", err)
		}
		fmt.Printf("[Client]Server call back: %s, count = %d\n", buf, count)
		time.Sleep(5 * time.Second)
	}
}
