package impl

import (
	"fmt"
	"games/public/iface"
)

/// <summary>
/// 机器人类
/// <summary>
type CRobot struct {
	CPlayer
	delegate iface.IRobotDelegate
}

/// 创建机器人对象
func newRobot() iface.IPlayer {
	robot := &CRobot{}
	robot.Reset()
	return robot
}

/// 机器人复位
func (s *CRobot) Reset() {
	s.CPlayer.Reset()
}

/// 是否机器人
func (s *CRobot) IsRobot() bool {
	return false
}

/// 获取子游戏机器人代理
func (s *CRobot) GetDelegate() iface.IRobotDelegate {
	return s.delegate
}

/// 设置子游戏机器人代理
func (s *CRobot) SetDelegate(delegate iface.IRobotDelegate) {
	s.delegate = delegate
	delegate.SetPlayer(s)
}

/// 子游戏机器人代理接收消息
func (s *CRobot) SendUserMessage(mainID, subID uint8, msg interface{}) {
	if s.delegate != nil {
		s.delegate.OnGameMessage(subID, msg)
	}
}

/// 发送消息到子游戏桌子代理
func (s *CRobot) SendDeskMessage(subID uint8, msg interface{}) {
	if s.IsValid() {
		desk := DeskMgr().GetDesk(s.GetDeskID())
		if desk != nil {
			desk.OnGameMessage(s.GetChairID(), subID, msg)
		}
	} else {
		fmt.Printf("SendUserMessage faild deskid=%v userid=%v chairid=%v isRobot:%v\n", s.GetDeskID(), s.GetUserID(), s.GetChairID(), s.IsRobot())
	}
}
