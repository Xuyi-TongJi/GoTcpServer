package iface

/*
	消息管理模块抽象层
*/

type IConnectionManager interface {
	// Add 添加连接
	Add(c IConnection)
	// Remove 删除连接
	Remove(c IConnection)
	// Get 根据connId获取连接
	Get(c uint32) (IConnection, error)
	// Total 得到当前连接总数
	Total() int
	// ClearAll 清除并中止所有连接
	ClearAll()
}
