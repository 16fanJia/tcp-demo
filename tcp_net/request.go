package tcp_net

import "tcp-demo/iface"

type Request struct {
	//建立的连接
	connection iface.IConnection
	//数据
	msg iface.IMessage
}

func (r *Request) GetConnection() iface.IConnection {
	return r.connection
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
