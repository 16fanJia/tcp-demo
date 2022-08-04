package iface

//消息管理模块

type IMessageManage interface {
	//DoMsgHandler 执行对应的router 消息处理方法
	DoMsgHandler(request IRequest)
	//AddRouter 添加路由
	AddRouter(msgId uint32, router IRouter)
}
