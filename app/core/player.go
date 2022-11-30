package core

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"server/app/pb"
	"server/iface"
	"sync/atomic"
)

type Player struct {
	Pid int32
	// 当前玩家的连接，用于和客户端连接
	Conn iface.IConnection
	// 平面X坐标
	X float32
	// 高度
	Y float32
	// 平面Z坐标
	Z float32
	// 旋转角度
	V float32
}

/*
	Player ID 生成器（暂时代替数据库）
*/

var PidGen int32 = 0
var InitX = float32(160 + rand.Intn(10))
var InitY = float32(0)
var InitZ = float32(140 + rand.Intn(20))
var InitV = float32(0)

func NewPlayer(conn iface.IConnection) *Player {
	// 生成玩家id
	pId := atomic.AddInt32(&PidGen, 1)
	return &Player{
		Pid:  pId,
		Conn: conn,
		X:    InitX,
		Y:    InitY,
		Z:    InitZ,
		V:    InitV,
	}
}

// SendMessage 玩家发送给客户端消息的方法，将pb的protobuf数据序列化之后，再调用server的sendMessage方法
func (p *Player) SendMessage(msgId uint32, data proto.Message) {
	// 将proto Message结构题序列化 转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Printf("[APP ERROR] Player %d marshal message error:%s\n", p.Pid, err)
		return
	}
	// 将二进制文件 通过server的sendMessage，将数据发送给客户端
	if p.Conn == nil || p.Conn.HasClosed() {
		fmt.Printf("[APP WARNING] Player %d Connection has been closed\n", p.Pid)
		return
	}
	err = p.Conn.SendMessage(msgId, msg)
	if err != nil {
		fmt.Printf("[APP ERROR] Player %d send message to client error:%s\n", p.Pid, err)
		return
	}
}

// SyncPId 告知客户端玩家ID
func (p *Player) SyncPId() {
	// message
	protoMessage := &pb.SyncPid{
		Pid: p.Pid,
	}
	// 将消息发送给客户端 msgId = 1
	p.SendMessage(1, protoMessage)
}

func (p *Player) BroadcastStartPosition() {
	protoMessage := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2, // 广播位置坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// msgId = 200
	p.SendMessage(200, protoMessage)
}
