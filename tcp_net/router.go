package tcp_net

import (
	"fmt"
	"tcp-demo/iface"
)

type BaseRouter struct{}

//这里之所以 BaseRouter 的方法都为空，
//是因为有的Router 不希望有 PreHandle PostHandle 方法
//所以Router全部继承BaseRouter 的好处是，不需要实现PreHandle PostHandle 方法 也可实例化对象

//PreHandle --
func (b *BaseRouter) PreHandle(request iface.IRequest) {
	//TODO
}

//Handle --
func (b *BaseRouter) Handle(request iface.IRequest) {
	//TODO
}

//PostHandle --
func (b *BaseRouter) PostHandle(request iface.IRequest) {
	//TODO
}

type TestRouter struct {
	BaseRouter
}

func (tr *TestRouter) Handle(request iface.IRequest) {
	fmt.Printf("msg:%s ----id:%d\n", string(request.GetData()), request.GetConnection().GetConnectionId())

	//返回封包后的数据给客户端
	err := request.GetConnection().SendMsg(request.GetMsgID(), request.GetData())
	if err != nil {
		fmt.Printf("write data err:%s", err.Error())
	}
	//_, err := request.GetConnection().GetTcpConnection().Write(request.GetData())
	//if err != nil {
	//	fmt.Printf("write data err:%s", err.Error())
	//}
}
