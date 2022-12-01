package core

import "sync"

/*
	当前游戏世界的总管理
*/

type WorldManager struct {
	// 世界地图管理模块
	AoiObj *AoiManager
	// 当前全部在线的Players集合
	Players map[int32]*Player
	// 读写锁
	lock sync.RWMutex
}

var WorldManagerObj *WorldManager

func init() {
	WorldManagerObj = &WorldManager{
		// 创建世界地图
		AoiObj:  NewAoiManager(0, AoiMinX, AoiMaxX, AoiCntX, AoiMinY, AoiMaxY, AoiCntY),
		Players: make(map[int32]*Player),
	}
}

// AddPlayer 添加玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.lock.Lock()
	defer wm.lock.Unlock()
	wm.Players[player.Pid] = player
	wm.AoiObj.AddToGridByCoordinate(player.Pid, player.X, player.Z)
}

// RemovePlayer 删除玩家
func (wm *WorldManager) RemovePlayer(pId int32) {
	player := wm.GetPlayerByPid(pId)
	if player != nil {
		wm.lock.Lock()
		defer wm.lock.Unlock()
		wm.AoiObj.RemoveFromGridByCoordinate(pId, player.X, player.Z)
		delete(wm.Players, pId)
	}
}

// GetPlayerByPid 查询玩家
func (wm *WorldManager) GetPlayerByPid(pId int32) *Player {
	wm.lock.RLock()
	defer wm.lock.RUnlock()
	if _, e := wm.Players[pId]; e {
		return wm.Players[pId]
	}
	return nil
}

// GetAllPlayers 获取全部在线玩家
func (wm *WorldManager) GetAllPlayers() (players []*Player) {
	wm.lock.RLock()
	defer wm.lock.RUnlock()
	for _, player := range wm.Players {
		players = append(players, player)
	}
	return
}
