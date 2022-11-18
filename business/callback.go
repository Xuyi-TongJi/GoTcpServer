package business

import (
	"fmt"
	"server/iface"
)

// CallBackToClient 回显的业务逻辑
func CallBackToClient(c iface.IConnection, data []byte, len int) error {
	fmt.Printf("[Connection Handle Function] Callback to Client %s", c.GetClientTcpStatus())
	conn := c.GetTcpConnection()
	if _, err := conn.Write(data[:len]); err != nil {
		return err
	}
	return nil
}
