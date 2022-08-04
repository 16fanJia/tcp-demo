package msg_manage

import (
	"fmt"
	"tcp-demo/iface"
)

type MessageManage struct {
	apis map[uint32]iface.IRouter
}

var MsgManageInstance *MessageManage

func init() {
	MsgManageInstance = &MessageManage{apis: make(map[uint32]iface.IRouter)}
}

func (m *MessageManage) DoMsgHandler(request iface.IRequest) {
	handler, ok := m.apis[request.GetMsgID()]
	if !ok {
		fmt.Printf("此id[ %d ]的消息类型消息未定义处理方法", request.GetMsgID())
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MessageManage) AddRouter(msgId uint32, router iface.IRouter) {
	if _, ok := m.apis[msgId]; ok {
		fmt.Printf("此id[ %d ]的消息类型处理函数已经存在", msgId)
		return
	}
	m.apis[msgId] = router
}
