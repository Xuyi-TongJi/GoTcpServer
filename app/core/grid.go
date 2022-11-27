package core

import (
	"fmt"
	"sync"
)

/*
	一个AOI地图中的格子对象
*/

type Grid struct {
	// 格子id
	GId int
	// 格子的左边界坐标
	MinX int
	// 格子的右边界坐标
	MaxX int
	// 格子的上边界坐标
	MinY int
	// 格子的下边界坐标
	MaxY int
	// 当前格子内玩家或者物体成员的ID集合
	PlayerIds map[int]bool
	// 保护当前集合的锁
	Lock sync.RWMutex
}

func NewGrid(gId, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GId:       gId,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		PlayerIds: make(map[int]bool),
	}
}

func (g *Grid) AddPlayer(playerId int) {
	// write lock
	g.Lock.Lock()
	defer g.Lock.Unlock()
	g.PlayerIds[playerId] = true
}

func (g *Grid) RemovePlayer(playerId int) {
	// remove
	g.Lock.Lock()
	defer g.Lock.Unlock()
	delete(g.PlayerIds, playerId)
}

func (g *Grid) GetPlayerIds() (playerIds []int) {
	g.Lock.Lock()
	defer g.Lock.Unlock()
	for k, _ := range g.PlayerIds {
		playerIds = append(playerIds, k)
	}
	return
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid id=%d, minX=%d, maxX=%d, minY=%d, maxY=%d, playerIds=%v",
		g.GId, g.MinX, g.MaxX, g.MinY, g.PlayerIds)
}
