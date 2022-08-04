package core

import (
	"fmt"
	"testing"
)

func TestNewAoiManage(t *testing.T) {
	aoi := NewAoiManage(0, 250, 5, 0, 250, 5)
	//for _, v := range aoi.Cells {
	//	fmt.Println(v.CellId, v.MinX, v.MaxX, v.MinY, v.MaxY, v.playerIDs)
	//}

	cellId := aoi.GetCellIdByCoordinate(0, 0)

	result := aoi.GetSurroundCellsByCellId(cellId)
	for _, v := range result {
		fmt.Println(v.CellId)
	}
}
