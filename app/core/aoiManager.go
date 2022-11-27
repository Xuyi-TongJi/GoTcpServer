package core

/*
	AOI区域管理模块
*/

type AoiManager struct {
	AId int
	// 区域的左边界坐标
	MinX int
	// 区域的右边界坐标
	MaxX int
	// X方向格子的数量
	CntX int
	// 区域的上边界坐标
	MinY int
	// 区域的下边界坐标
	MaxY int
	// Y方向格子的数量
	CntY int
	// 当前区域中有哪些格子 map-key=格子的ID，value=格子对象
	Grids map[int]*Grid
	// 每一格的宽度
	GridWidth int
	// 每一格的长度
	GridLength int
}

func NewAoiManager(aId, minX, maxX, cntX, minY, maxY, cntY int) *AoiManager {
	aoiManager := AoiManager{
		AId:   aId,
		MinX:  minX,
		MaxX:  maxX,
		MinY:  minY,
		MaxY:  maxY,
		CntX:  cntX,
		CntY:  cntY,
		Grids: make(map[int]*Grid),
	}
	aoiManager.GridWidth = (maxX - minX) / cntX
	aoiManager.GridLength = (maxY - minY) / cntY
	// 初始化AOI管理区域内的所有格子
	for x := 0; x < cntX; x += 1 {
		for y := 0; y < cntY; y += 1 {
			g := NewGrid(y*cntX+x, aoiManager.GridWidth*x, aoiManager.GridWidth*(x+1),
				aoiManager.GridLength*y, aoiManager.GridLength*(y+1))
			aoiManager.Grids[g.GId] = g
		}
	}
	return &aoiManager
}
