package test

import (
	"fmt"
	"io"
	"net"
	"server/network"
	"testing"
)

// TestDataPack 测试DataPack拆包和封包的单元测试
// ACCEPTED
func TestDataPack(t *testing.T) {
	// 模拟服务器
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		panic(err)
	}

	// 模拟服务端处理请求
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error: ", err)
			}
			go func(conn net.Conn) {
				// 拆包
				/* 分两次读 */
				dp := network.DataPack{}
				for {
					// 1.读head (len id)
					headData := make([]byte, dp.GetHeadLen())
					// 根据data len，再读取data内容 (readFull方法 读满headLen个字节)
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err:", err)
						break
					}
					// 2.读入前8个字节 (len id)
					msg, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err)
					}
					// len字段有大于0的值，说明需要继续读入数据
					if msg.GetLen() > 0 {
						msg.SetData(make([]byte, msg.GetLen()))
						// 根据len，再次从io流中读取
						_, err2 := io.ReadFull(conn, msg.GetData())
						if err2 != nil {
							fmt.Println("server unpack err:", err)
							return
						}
					}
					fmt.Printf(
						"success reading a complete message, data len = %d, data id = %d, data = %s\n",
						msg.GetLen(), msg.GetMsgId(), msg.GetData())
				}
			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial error,", err)
		return
	}
	dp := network.DataPack{}
	// 模拟粘包过程，封装两个msg一起发送
	msg1 := &network.Message{
		Id:   1,
		Data: []byte{'h', 'e', 'l', 'l', 'o', '1'},
		Len:  6,
	}
	msg2 := &network.Message{
		Id:   2,
		Data: []byte{'h', 'e', 'l', 'l', 'o', '2'},
		Len:  6,
	}
	data1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}
	// 一次性发送msg1和msg2
	data2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}
	sendData := append(data1, data2...)
	_, err = conn.Write(sendData)
	if err != nil {
		fmt.Println("send message error", err)
		return
	}
	// 阻塞
	select {}
}
