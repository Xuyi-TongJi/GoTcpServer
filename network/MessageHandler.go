package network

import (
	"fmt"
	"server/iface"
	"strconv"
)

type MessageHandler struct {
	// 存放每个MsgId对应的处理方法
	ApiMap map[uint32]iface.IRouter
}

func NewMessageHandler() iface.IMessageHandler {
	return &MessageHandler{
		ApiMap: make(map[uint32]iface.IRouter),
	}
}

func (m *MessageHandler) DoHandle(request iface.IRequest) {
	msgId := request.GetMsgId()
	if router, exist := m.ApiMap[msgId]; !exist {
		panic("[MessageHandler Handle Router ERROR] Message id = +" + strconv.Itoa(int(msgId)) + ", missing router")
	} else {
		// 调用router的模版方法
		iface.Handle(router, request)
	}
}

func (m *MessageHandler) AddRouter(msgId uint32, router iface.IRouter) {
	if _, exist := m.ApiMap[msgId]; !exist {
		m.ApiMap[msgId] = router
		fmt.Printf("[Server Register Router] Message %d, add router success\n", msgId)
	} else {
		panic("[MessageHandler Register Router WARNING] Repeat api, msgID = " + strconv.Itoa(int(msgId)))
	}
}
