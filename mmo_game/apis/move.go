package apis

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"tcp-demo/iface"
	"tcp-demo/mmo_game/pb"
	"tcp-demo/mmo_game/worldplayermanage"
	"tcp-demo/tcp_net"
)

//MoveRouter 玩家移动消息处理
type MoveRouter struct {
	tcp_net.BaseRouter
}

func (m *MoveRouter) Handle(request iface.IRequest) {
	var err error
	//解析客户端传过来的数据
	moveMsg := &pb.Move{}
	if err = proto.Unmarshal(request.GetData(), moveMsg); err != nil {
		fmt.Printf("move handle 解析数据失败 err:%s", err.Error())
		return
	}
	//根据新数据修改自己的位置信息
	player := worldplayermanage.WorldManageInstance.GetPlayerByPid(moveMsg.Pid)

	player.UpdatePosInfo(moveMsg.Pos)
}
