package impl

import (
	"fmt"
	"games/public/iface"
	"sync"
)

var playerMgr = newPlayerMgr()

/// <summary>
/// 单例玩家管理器实例
/// <summary>
func PlayerMgr() *CPlayerMgr {
	return playerMgr
}

/// <summary>
/// 玩家管理类
/// <summary>
type CPlayerMgr struct {
	players  map[int64]iface.IPlayer //使用中的玩家对象
	freeList []iface.IPlayer         //未使用的玩家对象 list优化
	lock     *sync.RWMutex
}

/// 创建玩家管理器
func newPlayerMgr() *CPlayerMgr {
	return &CPlayerMgr{players: map[int64]iface.IPlayer{}, freeList: []iface.IPlayer{}, lock: &sync.RWMutex{}}
}

/// 创建玩家对象
func (s *CPlayerMgr) New(userID int64) iface.IPlayer {
	{
		s.lock.RLock()
		if _, ok := s.players[userID]; ok {
			s.lock.RUnlock()
			panic(fmt.Sprintf("New userID = %v", userID))
		}
		s.lock.RUnlock()
	}
	var player iface.IPlayer
	{
		s.lock.Lock()
		if len(s.freeList) > 0 {
			player = s.freeList[0]
			s.freeList = append(s.freeList[:0], s.freeList[1:]...)
		}
		s.lock.Unlock()
	}
	if player == nil {
		player = newPlayer()
	} else {
		player.Reset()
	}
	{
		s.lock.Lock()
		s.players[userID] = player
		s.lock.Unlock()
	}
	return player
}

/// 查找玩家对象
func (s *CPlayerMgr) Get(userID int64) iface.IPlayer {
	{
		s.lock.RLock()
		if player, ok := s.players[userID]; ok {
			s.lock.RUnlock()
			return player
		}
		s.lock.RUnlock()
	}
	return nil
}

/// 回收玩家对象
func (s *CPlayerMgr) Delete(userID int64) {
	{
		s.lock.Lock()
		if player, ok := s.players[userID]; ok {
			delete(s.players, userID)
			player.Reset()
			s.freeList = append(s.freeList, player)
		}
		s.lock.Unlock()
	}
}
