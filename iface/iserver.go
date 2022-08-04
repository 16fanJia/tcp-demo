package iface

//IServer 服务器接口
type IServer interface {
	Start()  //start server
	Stop()   //stop server
	Server() //

	//AddRouter 添加路由功能 给当前服务注册一个路由方法，供客户端连接使用
	AddRouter(uint32, IRouter)
	//SetConnStarting 设置该Server的连接创建时Hook函数
	SetConnStarting(func(IConnection))
	//SetConnStopping 设置该Server的连接断开时的Hook函数
	SetConnStopping(func(IConnection))
}
