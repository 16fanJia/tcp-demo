package core

//aoi  核心实现

type AoiManage struct {
	xLeftCoordinate  int //x轴左边界坐标
	xRightCoordinate int //x轴右边界坐标
	xCount           int //x轴的格子数

	yUpCoordinate  int           //y轴上边界坐标
	yLowCoordinate int           //y轴下边界坐标
	yCount         int           //y轴的格子数
	Cells          map[int]*Cell //当前区域中都有哪些格子，key=格子ID， value=格子对象
}

const (
	X_LEFT_COORDINATE  = 0
	X_RIGHT_COORDINATE = 250
	X_COUNT            = 5

	Y_UP_COORDINATE  = 0
	Y_LOW_COORDINATE = 250
	Y_COUNT          = 5
)

//NewAoiManage 初始化一个AOI 管理器
func NewAoiManage(XLeftCoordinate int, XRightCoordinate int, XCount int,
	YUpCoordinate int, YLowCoordinate int, YCount int) *AoiManage {
	aoi := &AoiManage{
		xLeftCoordinate:  XLeftCoordinate,
		xRightCoordinate: XRightCoordinate,
		xCount:           XCount,
		yUpCoordinate:    YUpCoordinate,
		yLowCoordinate:   YLowCoordinate,
		yCount:           YCount,
		Cells:            map[int]*Cell{},
	}
	//初始化所有的格子 为格子编号
	peerXCellWidth := aoi.getPeerXCellWidth()
	peerYCellHigh := aoi.getPeerYCellHigh()

	for y := 0; y < aoi.yCount; y++ {
		for x := 0; x < aoi.xCount; x++ {
			//格子编号：id = idy *nx + idx  (利用格子坐标得到格子编号)
			cellId := y*aoi.xCount + x

			aoi.Cells[cellId] = NewCell(cellId,
				aoi.xLeftCoordinate+x*peerXCellWidth,
				aoi.xLeftCoordinate+(x+1)*peerXCellWidth,
				aoi.yUpCoordinate+y*peerYCellHigh,
				aoi.yUpCoordinate+(y+1)*peerYCellHigh)
		}
	}
	return aoi
}

//获取横坐标 每个格子的宽度
func (am *AoiManage) getPeerXCellWidth() int {
	return (am.xRightCoordinate - am.xLeftCoordinate) / am.xCount
}

//获取纵坐标 每个格子的高度
func (am *AoiManage) getPeerYCellHigh() int {
	return (am.yLowCoordinate - am.yUpCoordinate) / am.yCount
}

//GetSurroundCellsByCellId 根据格子的cellID得到当前周边的九宫格信息
func (am *AoiManage) GetSurroundCellsByCellId(cellId int) []*Cell {
	cells := make([]*Cell, 0, 9)
	//判断cellID是否存在
	if _, ok := am.Cells[cellId]; !ok {
		return nil
	}
	//将此cell本身存入cells 结果集
	cells = append(cells, am.Cells[cellId])

	//计算此cell 的横坐标编号 和纵坐标编号
	idx := cellId % am.xCount
	idy := cellId / am.yCount
	//判断此cell 左右两侧是否还有其他cell
	if idx > 0 {
		cells = append(cells, am.Cells[cellId-1])
	}
	if idx < am.xCount-1 {
		cells = append(cells, am.Cells[cellId+1])
	}

	//判断cellX中的格子 上下是否存在格子 实际拷贝了一份cells
	for _, v := range cells {
		if idy > 0 {
			cells = append(cells, am.Cells[v.CellId-am.xCount])
		}
		if idy < am.yCount-1 {
			cells = append(cells, am.Cells[v.CellId+am.xCount])
		}
	}
	return cells
}

//GetCellIdByCoordinate 根据横纵坐标 找到cellId
func (am *AoiManage) GetCellIdByCoordinate(x, y float32) int {
	idx := (int(x) - am.xLeftCoordinate) / am.getPeerXCellWidth()
	idy := (int(y) - am.yUpCoordinate) / am.getPeerYCellHigh()

	return idy*am.xCount + idx
}

//GetPlayersByCoordinate 根据横纵坐标 获得周边九宫格里的player id的集合
func (am *AoiManage) GetPlayersByCoordinate(x, y float32) (players []int32) {
	//获取cellID
	cellId := am.GetCellIdByCoordinate(x, y)

	//根据cellID 获取九宫格
	cells := am.GetSurroundCellsByCellId(cellId)

	//获取cells 中的所有player
	for _, v := range cells {
		players = append(players, v.GetPlayersFromCell()...)
	}
	return
}

//GetAllPidsByCellId 通过CellId获取当前格子的全部playerID
func (am *AoiManage) GetAllPidsByCellId(cellId int) (playerIds []int32) {
	playerIds = am.Cells[cellId].GetPlayersFromCell()
	return
}

//RemovePidFromCellId 移除一个格子中的PlayerID
func (am *AoiManage) RemovePidFromCellId(pID int32, cellID int) {
	am.Cells[cellID].RemovePlayerFromCell(pID)
}

//AddPidToCell 添加一个PlayerID到一个格子中
func (am *AoiManage) AddPidToCell(pID int32, cellID int) {
	am.Cells[cellID].AddPlayerToCell(pID)
}

//AddToCellByCoordinate 通过横纵坐标添加一个Player到一个格子中
func (am *AoiManage) AddToCellByCoordinate(pID int32, x, y float32) {
	cellId := am.GetCellIdByCoordinate(x, y)

	cell := am.Cells[cellId]
	cell.AddPlayerToCell(pID)
}

//RemoveFromCellByCoordinate 通过横纵坐标把一个Player从对应的格子中删除
func (am *AoiManage) RemoveFromCellByCoordinate(pID int32, x, y float32) {
	cellID := am.GetCellIdByCoordinate(x, y)
	cell := am.Cells[cellID]
	cell.RemovePlayerFromCell(pID)
}
