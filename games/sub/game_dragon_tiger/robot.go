package game_dragon_tiger

import (
	timer "games/core/timerv2"
	"games/public/define"
	"games/public/iface"
	"log"
)

/// <summary>
/// 龙虎斗子游戏机器人
/// <summary>
type CRobot struct {
	//iface.IRobotDelegate
	desk   iface.IDesk              //框架桌子指针
	player iface.IPlayer            //框架机器人信息
	timer  *timer.SafeTimerScheduel //桌子线程安全定时器
}

func newRobot() iface.IRobotDelegate {
	return &CRobot{}
}

/// 设置机器人
func (s *CRobot) SetPlayer(player iface.IPlayer) {
	s.player = player
	log.Printf("龙虎斗机器人[%v]实例化完成\n", s.player.GetUserID())
}

/// 设置桌子
func (s *CRobot) SetDesk(desk iface.IDesk) {
	s.desk = desk
	log.Printf("龙虎斗机器人[%v]进入桌子[%v]\n", s.player.GetUserID(), s.desk.GetDeskID())
	proc := s.desk.GetSlot().GetProc()
	if proc == nil {

	} else {
		//机器人线程安全定时器
		timer := proc.GetTimer().(*timer.SafeTimerScheduel)
		if timer == nil {
		}
	}
}

/// 机器人复位
func (s *CRobot) Reposition() {

}

/// 游戏消息
func (s *CRobot) OnGameMessage(subID uint8, msg interface{}) {

}

/// 机器人配置
func (s *CRobot) SetStrategy(strategy *define.RobotStrategyParam) {

}

/// 机器人配置
func (s *CRobot) GetStrategy() *define.RobotStrategyParam {
	return nil
}

/// <summary>
/// 机器人代理创建器
/// <summary>
type RobotDelegateCreator struct {
	iface.IRobotDelegateCreator
}

func (s *RobotDelegateCreator) Create() iface.IRobotDelegate {
	return newRobot()
}

func NewRobotDelegateCreator() iface.IRobotDelegateCreator {
	return &RobotDelegateCreator{}
}
