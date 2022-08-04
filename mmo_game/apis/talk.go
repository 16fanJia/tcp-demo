package apis

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"tcp-demo/iface"
	"tcp-demo/mmo_game/pb"
	"tcp-demo/mmo_game/worldplayermanage"
	"tcp-demo/tcp_net"
)

//TalkRouter 聊天消息处理结构体
type TalkRouter struct {
	tcp_net.BaseRouter
}

func (t *TalkRouter) Handle(request iface.IRequest) {
	//对客户端传过来的proto 数据进行解析
	var err error
	protoMsg := &pb.Talk{}
	if err = proto.Unmarshal(request.GetData(), protoMsg); err != nil {
		fmt.Printf("talk 解析数据失败 err:%s", err.Error())
		return
	}
	//获取发送消息的玩家id
	pid := protoMsg.GetPid()

	//根据pid 获取玩家对象
	player := worldplayermanage.WorldManageInstance.GetPlayerByPid(pid)

	player.Talk(protoMsg.Content)
}
