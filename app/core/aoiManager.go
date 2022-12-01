package core

const (
	AoiMinX int = 85
	AoiMaxX int = 410
	AoiCntX int = 10
	AoiMinY int = 75
	AoiMaxY int = 400
	AoiCntY int = 20
)

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

// GetSurroundingByGId 根据gid得到周边九宫格的id集合
func (m *AoiManager) GetSurroundingByGId(gId int) (grids []*Grid) {
	x := gId % m.CntX
	y := gId / m.CntX
	gIdX := []int{x}
	gIdY := []int{y}
	if x-1 >= 0 {
		gIdX = append(gIdX, x-1)
	}
	if x+1 < m.CntX {
		gIdX = append(gIdX, x+1)
	}
	if y-1 >= 0 {
		gIdY = append(gIdY, y-1)
	}
	if y+1 < m.CntY {
		gIdY = append(gIdY, y+1)
	}
	for _, xId := range gIdX {
		for _, yId := range gIdY {
			currentGId := yId*m.CntX + xId
			grids = append(grids, m.Grids[currentGId])
		}
	}
	return
}

// GetPlayerIdsByPos 通过横纵坐标获得周边九宫格内所有玩家的信息
func (m *AoiManager) GetPlayerIdsByPos(x, y float32) (playerIds []int32) {
	gId := m.getGridByCoordinate(x, y)
	grids := m.GetSurroundingByGId(gId)
	for _, g := range grids {
		for playerId, _ := range g.PlayerIds {
			playerIds = append(playerIds, playerId)
		}
	}
	return
}

// AddPIdToGrid 添加一个playerId到格子中
func (m *AoiManager) AddPIdToGrid(pId int32, gId int) {
	if _, e := m.Grids[gId]; e {
		m.Grids[gId].AddPlayer(pId)
	}
}

// RemovePIdFromGrid 从Grid中删除玩家
func (m *AoiManager) RemovePIdFromGrid(pId int32, gId int) {
	if _, e := m.Grids[gId]; e {
		m.Grids[gId].RemovePlayer(pId)
	}
}

// GetPIdsFromGrid 通过GId获取全部的playerId
func (m *AoiManager) GetPIdsFromGrid(gId int) (playerIds []int32) {
	if _, e := m.Grids[gId]; e {
		for playerId, _ := range m.Grids[gId].PlayerIds {
			playerIds = append(playerIds, playerId)
		}
	}
	return
}

// AddToGridByCoordinate 通过坐标添加player
func (m *AoiManager) AddToGridByCoordinate(pId int32, x, y float32) {
	gId := m.getGridByCoordinate(x, y)
	m.AddPIdToGrid(pId, gId)
}

// RemoveFromGridByCoordinate 通过坐标删除player
func (m *AoiManager) RemoveFromGridByCoordinate(pId int32, x, y float32) {
	gId := m.getGridByCoordinate(x, y)
	m.RemovePIdFromGrid(pId, gId)
}

// GetGridByCoordinate 通过横纵坐标来获取当前处于哪一个grid
func (m *AoiManager) getGridByCoordinate(x, y float32) int {
	xId := (int(x) - m.MinX) / m.GridWidth
	yId := (int(y) - m.MinY) / m.GridLength
	return yId*m.CntX + xId
}
