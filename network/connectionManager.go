package network

import (
	"errors"
	"fmt"
	"server/iface"
	"sync"
)

type ConnectionManager struct {
	ConnectionMap  map[uint32]iface.IConnection
	ConnectionLock sync.RWMutex
}

func NewConnectionManager() iface.IConnectionManager {
	return &ConnectionManager{
		ConnectionMap: make(map[uint32]iface.IConnection),
	}
}

func (cm *ConnectionManager) Add(c iface.IConnection) {
	// write lock
	cm.ConnectionLock.Lock()
	defer cm.ConnectionLock.Unlock()
	cm.ConnectionMap[c.GetConnId()] = c
	fmt.Printf("[ConnectionManager Add Connection] Add Connection %d success\n", c.GetConnId())
}

func (cm *ConnectionManager) Remove(c iface.IConnection) {
	// write lock
	cm.ConnectionLock.Lock()
	defer cm.ConnectionLock.Unlock()
	delete(cm.ConnectionMap, c.GetConnId())
	fmt.Printf("[ConnectionManager Remove Connection] Remove Connection %d success\n", c.GetConnId())
}

func (cm *ConnectionManager) Get(c uint32) (iface.IConnection, error) {
	// read lock
	cm.ConnectionLock.RLock()
	defer cm.ConnectionLock.RUnlock()
	if _, ok := cm.ConnectionMap[c]; ok {
		return cm.ConnectionMap[c], nil
	} else {
		return nil, errors.New(fmt.Sprintf("[ConnectionManager Get Connection ERROR] Invalid Connection id %d\n", c))
	}
}

func (cm *ConnectionManager) Total() int {
	// read lock
	cm.ConnectionLock.RLock()
	defer cm.ConnectionLock.RUnlock()
	return len(cm.ConnectionMap)
}

func (cm *ConnectionManager) ClearAll() {
	// write lock
	cm.ConnectionLock.Lock()
	defer cm.ConnectionLock.Unlock()
	// stop, then delete
	for connId, conn := range cm.ConnectionMap {
		conn.Stop()
		delete(cm.ConnectionMap, connId)
	}
	fmt.Printf("[Connection Manager Remove All] All connections removed.\n")
}
