package iface

type IRequest interface {
	//GetConnection 获取请求连接信息
	GetConnection() IConnection
	//GetData 获取请求消息的数据
	GetData() []byte
	//GetMsgID 获取请求消息的ID
	GetMsgID() uint32
}
