package net

import "server/iface"

type Request struct {
	// 已经和客户端建立好的连接
	conn iface.IConnection
	data []byte
	len  int
}

func (r *Request) GetConnection() iface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
