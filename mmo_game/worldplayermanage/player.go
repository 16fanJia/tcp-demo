package worldplayermanage

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"tcp-demo/iface"
	"tcp-demo/mmo_game/pb"
)

type Player struct {
	pid    int32             //玩家id
	conn   iface.IConnection //玩家所属的tcp 连接
	pos    Position          //玩家的坐标
	cellId int               //玩家所在的格子id
}

//Position 玩家坐标结构体
type Position struct {
	X float32 //平面x坐标
	Y float32 //平面y坐标 (注意不是Y)
	Z float32 //高度
	V float32 //旋转0-360度
}

func NewPlayer(conn iface.IConnection) *Player {
	//使用雪花算法 给玩家生成Pid

	p := &Player{
		pid:  0,
		conn: conn,
		pos: Position{
			X: float32(160 + rand.Intn(10)), //随机在160坐标点 基于X轴偏移若干坐标
			Y: float32(134 + rand.Intn(17)), //随机在134坐标点 基于Y轴偏移若干坐标
			Z: 0,
			V: 0,
		},
	}
	//根据玩家初始位置 获取初始格子ID
	p.cellId = WorldManageInstance.aoiM.GetCellIdByCoordinate(p.pos.X, p.pos.Y)

	//玩家上线 添加玩家到在线玩家集合中
	WorldManageInstance.AddPlayer(p)

	return p
}

func (p *Player) GetPid() int32 {
	return p.pid
}

//GetPlayerPos 返回玩家的坐标
func (p *Player) GetPlayerPos() Position {
	return p.pos
}

//SendProtoMsg 发送数据给客户端
func (p *Player) sendProtoMsg(msgId uint32, msg proto.Message) {
	var (
		err   error
		mData []byte
	)
	//对数据进行序列化
	if mData, err = proto.Marshal(msg); err != nil {
		fmt.Printf("proto 序列化数据失败 err:%s", err.Error())
		return
	}
	//数据发送给客户端
	if err = p.conn.SendMsg(msgId, mData); err != nil {
		fmt.Printf("player conn 发送数据失败 err:%s", err.Error())
		return
	}
}

//SyncPid 告知客户端pid,同步已经生成的玩家ID给客户端
func (p *Player) SyncPid() {
	//组织MsgID:1 proto数据
	protoMsg := &pb.SyncPid{
		Pid: p.pid,
	}

	//发送数据给客户端
	p.sendProtoMsg(1, protoMsg)
}

func (p *Player) BroadCastMsg() {
	//组织MsgId:200 2号(广播坐标数据) proto数据
	protoMsg := &pb.BroadCast{
		Pid: p.pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.pos.X,
				Y: p.pos.Y,
				Z: p.pos.Z,
				V: p.pos.V,
			}},
	}

	p.sendProtoMsg(200, protoMsg)
}

func (p *Player) Talk(content string) {
	//广播玩家聊天
	//1. 组建MsgId2 proto数据
	protoMsg := &pb.BroadCast{
		Pid: p.pid,
		Tp:  1, //TP 1 代表聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	// 得到当前世界所有的在线玩家
	players := WorldManageInstance.GetAllPlayer()

	for _, p := range players {
		p.sendProtoMsg(200, protoMsg)
	}
}

//SyncSurroundingPlayer 和周围玩家同步位置信息
func (p *Player) SyncSurroundingPlayer() {
	//根据当前玩家的位置 获取周边九宫格玩家的id
	ids := WorldManageInstance.GetPlayerByPosition(p.pos.X, p.pos.Y)

	//组建MsgID 200 proto数据
	posMsg := &pb.BroadCast{
		Pid: p.pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.pos.X,
				Y: p.pos.Y,
				Z: p.pos.Z,
				V: p.pos.V,
			},
		},
	}

	//根据Pid 获取在线玩家对象 并且组建MsgID:202 proto 信息
	//创建在线玩家集合
	playersPos := make([]*pb.Player, 0, len(ids))
	for _, pid := range ids {
		player := WorldManageInstance.GetPlayerByPid(pid)

		syncPlayer := &pb.Player{
			Pid: player.pid,
			P: &pb.Position{
				X: player.pos.X,
				Y: player.pos.Y,
				Z: player.pos.Z,
				V: player.pos.V,
			},
		}

		playersPos = append(playersPos, syncPlayer)
		//其他玩家给各自客户端 发送消息通知本玩家的的位置信息 让其看得见
		player.sendProtoMsg(200, posMsg)
	}

	//将其他所有玩家的位置信息 发送给本玩家的客户端
	syncPlayerPosMsg := &pb.SyncPlayer{Ps: playersPos}

	p.sendProtoMsg(202, syncPlayerPosMsg)
}

//UpdatePosInfo 修改玩家的位置信息
func (p *Player) UpdatePosInfo(position *pb.Position) {
	//判断玩家格子是否发生变化
	p.IsCellChange(position.X, position.Y)
	//修改本玩家位置信息
	p.pos.X = position.X
	p.pos.Y = position.Y
	p.pos.Z = position.Z
	p.pos.Z = position.Z

	//通知周边九宫格玩家
	movMsg := &pb.BroadCast{
		Pid: p.pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.pos.X,
				Y: p.pos.Y,
				Z: p.pos.Z,
				V: p.pos.V,
			},
		},
	}

	//根据cellID 获取周边九宫格玩家
	players := WorldManageInstance.GetSurroundingPlayersByCellId(p.cellId)

	//向周边的每个玩家发送MsgID:200消息，移动位置更新消息
	for _, player := range players {
		player.sendProtoMsg(200, movMsg)
	}
}

//IsCellChange 玩家移动后判读所在格子是否变化 变化则从旧格子移除添加到新格子
func (p *Player) IsCellChange(newX, newY float32) {
	newCellId := WorldManageInstance.aoiM.GetCellIdByCoordinate(newX, newY)
	if p.cellId == newCellId {
		return
	}

	//将玩家从旧格子中移除
	WorldManageInstance.aoiM.RemovePidFromCellId(p.pid, p.cellId)

	//将玩家添加到新格子中
	WorldManageInstance.aoiM.AddPidToCell(p.pid, newCellId)

	p.cellId = newCellId
}

//Offline 玩家下线
func (p *Player) Offline() {
	//发送消息给客户端
	offlineMsg := &pb.SyncPid{Pid: p.pid}

	p.sendProtoMsg(201, offlineMsg)

	//从在线玩家列表中删除
	WorldManageInstance.RemovePlayer(p.pid)
	//从aoi 中删除
	WorldManageInstance.aoiM.RemovePidFromCellId(p.pid, p.cellId)
}
