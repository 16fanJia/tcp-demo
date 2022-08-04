package iface

import "net"

//IConnection 连接接口
type IConnection interface {
	// StartConnection 开始连接
	StartConnection()
	// StopConnection 结束当前连接
	StopConnection()
	// GetTcpConnection 获取当前连接的绑定的socket conn
	GetTcpConnection() *net.TCPConn
	// GetConnectionId 获取当前连接的ID
	GetConnectionId() uint32
	// GetRemoteAddr 获取远程客户端的 TCP状态 IP port
	GetRemoteAddr() net.Addr
	// SendMsg 将数据发送给客户段
	SendMsg(uint32, []byte) error
	//SetProperty 设置属性
	SetProperty(key string, value interface{})
	//GetProperty 获取属性
	GetProperty(key string) interface{}
	//RemoveProperty 删除属性
	RemoveProperty(key string)
}
