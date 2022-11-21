package iface

/*
	消息管理模块抽象接口
*/

type IMessageHandler interface {
	// DoHandle 调度并并行对应的Router（消息处理方法）
	DoHandle(request IRequest)
	// AddRouter 添加Router
	AddRouter(msgId uint32, router IRouter)
}
