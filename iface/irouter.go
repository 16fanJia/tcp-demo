package iface

//IRouter 路由接口 定义的方法为处理业务的方法
type IRouter interface {
	//PreHandle 处理conn 业务之前调用的方法
	PreHandle(request IRequest)
	//Handle 处理业务调用的方法
	Handle(request IRequest)
	//PostHandle 处理业务之后调用的方法
	PostHandle(request IRequest)
}
