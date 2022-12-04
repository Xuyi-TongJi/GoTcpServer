package main

import (
	"fmt"
	"server/app/apis"
	"server/app/core"
	"server/iface"
	"server/network"
)

// OnConnectionAdd 客户端建立连接后的hook函数
func OnConnectionAdd(conn iface.IConnection) {
	// 创建玩家类
	player := core.NewPlayer(conn)
	// 给客户端发送MsgID为1的消息
	player.SyncPId()
	// 给客户端发送MsgID为200的消息
	player.BroadcastStartPosition()

	// 将当前新上线的玩家添加到world中
	core.WorldManagerObj.AddPlayer(player)
	conn.SetConnectionProperty("pId", player.Pid)

	fmt.Printf("[Player] Player %d is online\n", player.Pid)
}

func main() {
	s := network.NewServer("tcp4")
	// 连接创建和销毁的HOOK钩子函数
	s.SetOnConnectionStart(OnConnectionAdd)
	// 注册路由业务

	// 世界聊天业务
	s.AddRouter(2, &apis.WorldChatApi{})
	// 启动服务
	s.Serve()
}
