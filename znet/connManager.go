package znet

import (
	"errors"
	"fmt"
	"go-zinx/ziface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) AddConn(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.connections[conn.GetID()] = conn
	fmt.Println("conn add success,conn num=", cm.GetConnLen(), "conn id=", conn.GetID())
}

func (cm *ConnManager) RemoveConn(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.connections, conn.GetID())
	fmt.Println("conn delete success,conn num=", cm.GetConnLen(), "conn id=", conn.GetID())
}

func (cm *ConnManager) GetConnCount(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if c, ok := cm.connections[connID]; ok {
		return c, nil
	} else {
		return nil, errors.New("conn not found")
	}
}

func (cm *ConnManager) GetConnLen() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}
	fmt.Println("conn clear success,conn num=", len(cm.connections))
}
