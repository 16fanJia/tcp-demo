package core

import "sync"

//一张地图中的格子信息

type Cell struct {
	CellId    int                //格子ID
	MinX      int                //x轴上 格子左边界坐标
	MaxX      int                //x轴上 格子右边界坐标
	MinY      int                //y轴上 格子上边界坐标
	MaxY      int                //y轴上 格子下边界坐标
	playerIDs map[int32]struct{} //当前格子内的玩家或者物体成员ID
	mu        sync.RWMutex       //playerIDs的保护map的锁
}

//NewCell 一个小格子的构造函数
func NewCell(cellId int, minX int, maxX int, minY int, maxY int) *Cell {
	return &Cell{
		CellId:    cellId,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int32]struct{}),
	}
}

//AddPlayerToCell 向一个cell 中添加一个玩家
func (c *Cell) AddPlayerToCell(player int32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.playerIDs[player] = struct{}{}
}

//RemovePlayerFromCell 从一个cell 中移除一个玩家
func (c *Cell) RemovePlayerFromCell(player int32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.playerIDs, player)
}

//GetPlayersFromCell 获取当前格子的所有玩家
func (c *Cell) GetPlayersFromCell() []int32 {
	player := make([]int32, 0, len(c.playerIDs))
	for pId := range c.playerIDs {
		player = append(player, pId)
	}
	return player
}
