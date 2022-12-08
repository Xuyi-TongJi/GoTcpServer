package apis

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"server/app/core"
	"server/app/pb"
	"server/iface"
	"server/network"
)

type MoveRouter struct {
	*network.BaseRouter
}

func (router *MoveRouter) DoHandle(req iface.IRequest) {
	protoMessage := &pb.Position{}
	err := proto.Unmarshal(req.GetData(), protoMessage)
	if err != nil {
		fmt.Printf("[Move Router ERROR] Connection %d, Unable to unmarshal proto binary data, error:%s\n",
			req.GetConnection().GetConnId(), err)
		return
	}
	pId := req.GetConnection().GetConnectionProperty("pId").(int32)
	fmt.Printf("[Move Router] Player Pid = %d, move(%f, %f, %f, %f)\n",
		pId, protoMessage.X, protoMessage.Y, protoMessage.Z, protoMessage.V)
	player := core.WorldManagerObj.GetPlayerByPid(pId)
	// 广播并更新当前玩家的坐标
	player.UpdatePos(protoMessage.X, protoMessage.Y, protoMessage.Z, protoMessage.V)
}
