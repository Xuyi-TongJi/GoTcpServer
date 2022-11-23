package iface

// IServer 服务器接口
type IServer interface {
	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()
	// AddRouter 添加路由
	AddRouter(msgID uint32, router IRouter)
	// GetConnectionManager 获得连接管理器
	GetConnectionManager() IConnectionManager
	// SetOnConnectionStart 设置OnConnectionStart hook
	SetOnConnectionStart(hook func(connection IConnection))
	// SetOnConnectionStop 设置OnConnectionStop hook
	SetOnConnectionStop(hook func(connection IConnection))
	// CallOnConnectionStart 调用OnConnectionStart hook
	CallOnConnectionStart(connection IConnection)
	// CallOnConnectionStop 调用OnConnectionStop hook
	CallOnConnectionStop(connection IConnection)
}
