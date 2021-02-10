package game_dragon_tiger

import (
	timer "games/core/timerv2"
	"games/public/iface"
	"log"
)

/// <summary>
/// 龙虎斗子游戏
/// <summary>
type CGameDesk struct {
	//iface.IDeskDelegate
	desk iface.IDesk //框架桌子指针
}

/// 创建龙虎斗
func newDesk() iface.IDeskDelegate {
	return &CGameDesk{}
}

/// 指定桌子
func (s *CGameDesk) SetDesk(desk iface.IDesk) {
	s.desk = desk
	log.Printf("龙虎斗桌子[%v]实例化完成\n", s.desk.GetDeskID())
	//桌子实例化完成后，slot协程才启动 GetProc() == nil
	proc := s.desk.GetSlot().GetProc()
	if proc == nil {
		//fmt.Println("GetProc() == nil")
	} else {
		//桌子线程安全定时器
		timer := proc.GetTimer().(*timer.SafeTimerScheduel)
		if timer == nil {
		}
	}
}

/// 桌子复位
func (s *CGameDesk) Reposition() {

}

/// 游戏开始
func (s *CGameDesk) OnGameStart() {

}

//// 游戏结束
func (s *CGameDesk) OnGameConclude(chairID uint16, reason uint8) bool {
	return false
}

/// 场景推送
func (s *CGameDesk) OnGameScene(chairID uint16, look bool) {

}

/// 游戏消息
func (s *CGameDesk) OnGameMessage(chairID uint16, subID uint8, msg interface{}) {

}

/// 用户进入
func (s *CGameDesk) OnUserEnter(userID int64, look bool) bool {
	return false
}

/// 用户准备
func (s *CGameDesk) OnUserReady(userID int64, look bool) bool {
	return false
}

/// 用户离开
func (s *CGameDesk) OnUserLeave(userID int64, look bool) bool {
	return false
}

/// 能否进入
func (s *CGameDesk) CanUserEnter(player iface.IPlayer) bool {
	return false
}

/// 能否离开
func (s *CGameDesk) CanUserLeave(userID int64) {

}

/// <summary>
/// 桌子代理创建器
/// <summary>
type DeskDelegateCreator struct {
	iface.IDeskDelegateCreator
}

func (s *DeskDelegateCreator) Create() iface.IDeskDelegate {
	return newDesk()
}

func NewDeskDelegateCreator() iface.IDeskDelegateCreator {
	return &DeskDelegateCreator{}
}
