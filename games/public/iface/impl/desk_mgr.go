package impl

import (
	"mygame/core"
	"mygame/public/define"
	"mygame/public/iface"
	"sync"
)

var deskMgr = newDeskMgr()

/// <summary>
/// 单例桌子管理器实例
/// <summary>
func DeskMgr() *CDeskMgr {
	return deskMgr
}

/// <summary>
/// 桌子管理类
/// <summary>
type CDeskMgr struct {
	GameInfo  *define.GameInfo //游戏类型
	RoomInfo  *define.RoomInfo //游戏房间
	desks     []iface.IDesk    //所有桌子[deskID]=desk
	usedDesks []iface.IDesk    //使用中的桌子，list优化
	freeDesks []iface.IDesk    //空闲桌子，list优化
	lock      *sync.RWMutex
}

/// 清理
func (s *CDeskMgr) Clear() {
	s.desks = []iface.IDesk{}
	s.usedDesks = []iface.IDesk{}
	s.freeDesks = []iface.IDesk{}
}

/// 创建桌子管理器
func newDeskMgr() *CDeskMgr {
	return &CDeskMgr{desks: []iface.IDesk{}, usedDesks: []iface.IDesk{}, freeDesks: []iface.IDesk{}, lock: &sync.RWMutex{}}
}

/// 桌子管理器初始化
func (s *CDeskMgr) Init(gameInfo *define.GameInfo, roomInfo *define.RoomInfo, mailBox core.IMailbox, createDelegate iface.IDeskDelegateCreator) {
	if gameInfo == nil || roomInfo == nil || createDelegate == nil {
		return
	}
	s.GameInfo = gameInfo
	s.RoomInfo = roomInfo
	//创建指定数量的桌子
	for i := uint16(0); i < roomInfo.DeskCount; i++ {
		//创建子游戏桌子代理
		deskDelegate := createDelegate.Create()
		if deskDelegate == nil {
			return
		}
		deskState := &define.DeskState{}
		deskState.DeskID = i
		deskState.IsLock = false
		deskState.IsLook = false
		//创建桌子
		desk := newDesk(deskState)
		//读取房间库存
		desk.(*CDesk).ReadStorageScore()
		desk.(*CDesk).Init(gameInfo, roomInfo, mailBox.GetNextSlot())
		//桌子代理绑定桌子
		desk.SetDelegate(deskDelegate)
		//deskDelegate.SetDesk(desk)
		s.desks = append(s.desks, desk)
		s.freeDesks = append(s.freeDesks, desk)
	}
}

/// 获取指定桌子
func (s *CDeskMgr) GetDesk(deskID uint16) iface.IDesk {
	{
		s.lock.RLock()
		if deskID != define.INVALID_DESK && deskID < uint16(len(s.desks)) {
			s.lock.RUnlock()
			return s.desks[deskID]
		}
	}
	s.lock.RUnlock()
	return nil
}

/// 新开辟n张桌子
func (s *CDeskMgr) MakeSuitDesk(count uint16) uint16 {
	i := uint16(0)
	{
		s.lock.Lock()
		for ; i < count; i++ {
			if len(s.freeDesks) > 0 {
				desk := s.freeDesks[0]
				s.freeDesks = append(s.freeDesks[:0], s.freeDesks[1:]...)
				s.usedDesks = append(s.usedDesks, desk)
			} else {
				break
			}
		}
		s.lock.Unlock()
	}
	return i
}

/// 查找合适的桌子，没有则开辟一张出来
func (s *CDeskMgr) FindSuitDesk(player iface.IPlayer) iface.IDesk {
	var usedDesks []iface.IDesk
	{
		s.lock.RLock()
		usedDesks = append(usedDesks, s.usedDesks...)
		s.lock.RUnlock()
	}
	//从使用中的桌子中查找能进的桌子
	for _, desk := range usedDesks {
		//判断桌子是否能进
		if desk.CanUserEnter(player) {
			return desk
		}
	}
	//找不到从空闲桌子里面取
	{
		s.lock.Lock()
		if len(s.freeDesks) > 0 {
			desk := s.freeDesks[0]
			s.freeDesks = append(s.freeDesks[:0], s.freeDesks[1:]...)
			s.usedDesks = append(s.usedDesks, desk)
			s.lock.Unlock()
			return desk
		}
		s.lock.Unlock()
	}
	return nil
}

/// 获取可用的桌子
func (s *CDeskMgr) GetUsedDesks(worker interface{}) (usedDesks []iface.IDesk) {
	{
		s.lock.RLock()
		usedDesks = append(usedDesks, s.usedDesks...)
		s.lock.RUnlock()
	}
	return
}

/// 回收桌子
func (s *CDeskMgr) FreeDesk(deskID uint16) {
	s.lock.Lock()
	for i, desk := range s.usedDesks {
		if desk.GetDeskID() == deskID {
			s.usedDesks = append(s.usedDesks[:i], s.usedDesks[i+1:]...)
			s.freeDesks = append(s.freeDesks, desk)
			break
		}
	}
	s.lock.Unlock()
}

/// 剔出桌子全部玩家
func (s *CDeskMgr) KickAllPlayers() {

}
