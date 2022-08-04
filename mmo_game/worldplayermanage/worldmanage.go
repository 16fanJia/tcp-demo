package worldplayermanage

import (
	"sync"
	"tcp-demo/mmo_game/core"
)

/*
定义一个世界管理模块
1：对地图的管理
2：对玩家的管理

方法：
1：添加上线玩家到在线玩家集合中
2：从在线玩家中移除玩家
3：获取所有在线的玩家
*/

type WorldManage struct {
	aoiM    *core.AoiManage
	players map[int32]*Player
	mutex   sync.RWMutex
	count   int //在线玩家数量
}

var WorldManageInstance *WorldManage

func init() {
	WorldManageInstance = &WorldManage{
		aoiM: core.NewAoiManage(core.X_RIGHT_COORDINATE, core.X_LEFT_COORDINATE, core.X_COUNT,
			core.Y_UP_COORDINATE, core.Y_LOW_COORDINATE, core.Y_COUNT),
		players: make(map[int32]*Player),
		count:   0,
	}
}

//GetPlayerByPosition 根据玩家的坐标信息获取周边九宫格玩家的id
func (wm *WorldManage) GetPlayerByPosition(x, y float32) (pid []int32) {
	pid = wm.aoiM.GetPlayersByCoordinate(x, y)
	return
}

//AddPlayer 添加用户到世界用户在线集合
func (wm *WorldManage) AddPlayer(p *Player) {
	//加锁执行 操作map
	wm.mutex.Lock()
	wm.players[p.GetPid()] = p
	wm.count++
	wm.mutex.Unlock()

	//将player 添加到AOI网络规划中
	wm.aoiM.AddToCellByCoordinate(p.GetPid(), p.GetPlayerPos().X, p.GetPlayerPos().Y)
}

//RemovePlayer 从在线玩家集合中移除一个玩家
func (wm *WorldManage) RemovePlayer(pid int32) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	//删除
	delete(wm.players, pid)
	wm.count--
}

//GetPlayerByPid 通过pid 获取玩家的信息
func (wm *WorldManage) GetPlayerByPid(pid int32) *Player {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	return wm.players[pid]
}

//GetAllPlayer 获取所有的在线玩家
func (wm *WorldManage) GetAllPlayer() []*Player {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	//创建返回的player集合切片
	players := make([]*Player, 0, wm.count)

	for _, v := range wm.players {
		players = append(players, v)
	}
	return players
}

//GetSurroundingPlayersByCellId 根据玩家的cellId 获取周边九宫格里的所有玩家
func (wm *WorldManage) GetSurroundingPlayersByCellId(cellId int) (players []*Player) {
	cells := wm.aoiM.GetSurroundCellsByCellId(cellId)
	//获取cells 中的所有player
	for _, v := range cells {
		ids := v.GetPlayersFromCell()
		for _, pid := range ids {
			players = append(players, wm.GetPlayerByPid(pid))
		}
	}
	return
}
