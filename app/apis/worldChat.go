package apis

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"server/app/core"
	"server/app/pb"
	"server/iface"
	"server/network"
)

// 世界聊天路由业务
// 接收request中的数据，将protobuf解析为string后，再打包为新的二进制protobuf数据，广播给世界

type WorldChatApi struct {
	network.BaseRouter
}

func (wc *WorldChatApi) DoHandle(req iface.IRequest) {
	// 解析客户端传递进来的proto协议
	protoMsg := &pb.Talk{}
	err := proto.Unmarshal(req.GetData(), protoMsg)
	if err != nil {
		fmt.Printf("[World Chat Router ERROR] Connection %d, Unable to unmarshal proto binary data, error:%s\n",
			req.GetConnection().GetConnId(), err)
		return
	}
	pId := req.GetConnection().GetConnectionProperty("pId")
	player := core.WorldManagerObj.GetPlayerByPid(pId.(int32))
	// 广播给世界，content内容
	player.Talk(protoMsg.Content)
}
