package conn_management

import (
	"errors"
	"fmt"
	"sync"
	"tcp-demo/iface"
)

type ConnectionManage struct {
	connections map[uint32]iface.IConnection //连接的集合
	connLock    sync.RWMutex                 //锁
	count       int                          //连接总数
}

var (
	ConnectionManageClient *ConnectionManage
	ConnectionNotFound     = errors.New("connection not found")
)

func init() {
	ConnectionManageClient = &ConnectionManage{
		connections: make(map[uint32]iface.IConnection),
		count:       0,
	}
}

func (c *ConnectionManage) AddConn(conn iface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[conn.GetConnectionId()] = conn
	c.count++

	fmt.Printf("new connection add,id is [ %d ],now connection count is [ %d ]\n", conn.GetConnectionId(), c.count)
}

func (c *ConnectionManage) RemoveConn(conn iface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	delete(c.connections, conn.GetConnectionId())
	c.count--

	fmt.Printf("old connection removed,id is [ %d ],now connection count is [ %d ]\n", conn.GetConnectionId(), c.count)
}

func (c *ConnectionManage) GetConn(connID uint32) (iface.IConnection, error) {
	//保护共享资源Map 加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, ConnectionNotFound
	}
}

func (c *ConnectionManage) Count() int {
	return c.count
}

func (c *ConnectionManage) EliminateConn() {
	//停止并删除全部的连接信息
	for _, conn := range c.connections {
		//停止
		conn.StopConnection()
	}

	fmt.Println("Clear All Connections successfully: conn num = ", c.Count())
}
