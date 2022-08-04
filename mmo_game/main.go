package main

import (
	"fmt"
	"tcp-demo/config"
	"tcp-demo/iface"
	"tcp-demo/mmo_game/apis"
	"tcp-demo/mmo_game/worldplayermanage"
	"tcp-demo/tcp_net"
)

//DoConnectionBegin 创建连接的时候执行
func DoConnectionBegin(conn iface.IConnection) {
	p := worldplayermanage.NewPlayer(conn)

	//给连接设置pid 信息
	conn.SetProperty("pid", p.GetPid())
	//发送玩家上线广播通知 给客户端
	p.SyncPid()
	p.BroadCastMsg()
	//同步周边玩家
	p.SyncSurroundingPlayer()

	fmt.Printf("=======> 玩家[ %d ]上线 <=======", p.GetPid())
}

//DoConnectionLast 连接断开的时候执行
func DoConnectionLast(conn iface.IConnection) {
	//获取连接对应的玩家id
	pid := conn.GetProperty("pid")

	//根据pid 获取玩家
	player := worldplayermanage.WorldManageInstance.GetPlayerByPid(pid.(int32))

	//玩家下线
	player.Offline()

	fmt.Printf("=======> 玩家[ %d ]下线 <=======", pid)
}

func main() {
	//加载配置
	config.LoadConfig("./mmo_game/conf/config.conf")

	server, _ := tcp_net.NewServer()
	server.SetConnStarting(DoConnectionBegin)
	server.SetConnStopping(DoConnectionLast)

	server.AddRouter(2, &apis.TalkRouter{})
	server.AddRouter(3, &apis.MoveRouter{})
	server.Server()
}
