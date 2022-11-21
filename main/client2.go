package main

import (
	"fmt"
	"io"
	"net"
	"server/network"
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
		fmt.Printf("[Client ERROR] Connection errer:%s\n", err)
		return
	}
	for {
		// 发送Message消息 TLV格式
		dp := network.DataPack{}
		var s string = "hello, server!"
		length := uint32(len(s))
		msg := &network.Message{
			Id:   1,
			Len:  length,
			Data: []byte(s),
		}
		pack, err := dp.Pack(msg)
		if err != nil {
			fmt.Printf("[Client Writing Error] Error packing message, error: %s\n",
				err)
			continue
		}
		// send to server
		if _, err := conn.Write(pack); err != nil {
			fmt.Printf("[Client Writing Error] Error sending message to server, error: %s\n",
				err)
			continue
		}
		// 拆包

		// receive message from server
		head := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, head)
		if err != nil {
			fmt.Printf("[Client Reading Error] Error receiving message from server, error: %s\n", err)
			continue
		}
		// unpack len and id
		unpack, err := dp.Unpack(head)
		if err != nil {
			fmt.Printf("[Client Reading Error] Error unpacking message from server, error: %s\n", err)
			continue
		}
		if unpack.GetLen() > 0 {
			unpack.SetData(make([]byte, unpack.GetLen()))
			_, err := io.ReadFull(conn, unpack.GetData())
			if err != nil {
				fmt.Printf("[Client Reading Error] Error reading message data from server, error: %s\n", err)
				continue
			}
		}
		fmt.Printf("[Client Reading] Receiving message (id=%d): %s, success\n",
			unpack.GetMsgId(), unpack.GetData())
		time.Sleep(5 * time.Second)
	}
}
