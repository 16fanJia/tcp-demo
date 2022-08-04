package iface

//连接管理接口

type IConnectionManage interface {
	AddConn(conn IConnection)                   //添加链接
	RemoveConn(conn IConnection)                //删除连接
	GetConn(connID uint32) (IConnection, error) //利用ConnID获取链接
	Count() int                                 //获取当前连接总数
	EliminateConn()                             //清除并停止所有链接
}
