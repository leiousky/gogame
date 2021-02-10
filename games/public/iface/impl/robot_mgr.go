package impl

import (
	"fmt"
	"math/rand"
	"mygame/public/define"
	"mygame/public/iface"
	"sync"
)

var robotMgr = newRobotMgr()

/// <summary>
/// 单例机器人管理器实例
/// <summary>
func RobotMgr() *CRobotMgr {
	return robotMgr
}

/// <summary>
/// 机器人管理类
/// <summary>
type CRobotMgr struct {
	GameInfo      *define.GameInfo           //游戏类型
	RoomInfo      *define.RoomInfo           //游戏房间
	RobotStrategy *define.RobotStrategyParam //机器人参数配置
	usedRobots    []iface.IPlayer            //使用中的机器人 list优化
	freeRobots    []iface.IPlayer            //未使用的机器人 list优化
	lock          *sync.RWMutex
}

/// 创建机器人管理实例
func newRobotMgr() *CRobotMgr {
	return &CRobotMgr{usedRobots: []iface.IPlayer{}, freeRobots: []iface.IPlayer{}, lock: &sync.RWMutex{}}
}

/// 机器人管理器初始化
///<param name="gameInfo">游戏类型</param>
///<param name="roomInfo">游戏房间</param>
///<param name="createDelegate">子游戏机器人代理工厂</param>
///<param name="worker">桌子工作线程，内含协程安全定时器</param>
func (s *CRobotMgr) Init(gameInfo *define.GameInfo, roomInfo *define.RoomInfo, createDelegate iface.IRobotDelegateCreator) {
	if gameInfo == nil || roomInfo == nil || createDelegate == nil {
		return
	}
	s.GameInfo = gameInfo
	s.RoomInfo = roomInfo
	//1.db读取机器人参数配置
	s.RobotStrategy = &define.RobotStrategyParam{}
	s.RobotStrategy.GameID = 19
	s.RobotStrategy.RoomID = 19001
	s.RobotStrategy.EnterLowScore = 100000        //进入最小分
	s.RobotStrategy.EnterHighScore = 888888888888 //进入最大分
	s.RobotStrategy.MinScore = 6666666            //机器人最小携带分
	s.RobotStrategy.MaxScore = 8888888888888888   //机器人最大携带分

	//加载机器人总人数：房间桌子数 * 每张桌子最大机器人数
	//s.RoomInfo.DeskCount = 4
	//s.RoomInfo.MaxRobotCount = 2
	maxRobotCount := s.RoomInfo.DeskCount * s.RoomInfo.MaxRobotCount
	// 机器人账号起始userID
	startUserID := int64(66666666)
	//2.db读取机器人账号
	for i := uint16(0); i < maxRobotCount; i++ {
		//用户基础数据
		baseInfo := &define.UserBaseInfo{}
		baseInfo.UserID = startUserID + 1
		startUserID++
		baseInfo.Account = fmt.Sprintf("robot:%v:%v", i, baseInfo.UserID)
		baseInfo.NickName = baseInfo.Account
		baseInfo.TakeMinScore = s.RobotStrategy.MinScore                                                               //机器人最小携带分
		baseInfo.TakeMaxScore = s.RobotStrategy.MaxScore                                                               //机器人最大携带分
		baseInfo.Score = s.RobotStrategy.MinScore + rand.Int63()%(s.RobotStrategy.MaxScore-s.RobotStrategy.MinScore+1) //随机机器人携带积分
		//创建子游戏机器人代理
		robotDelegate := createDelegate.Create()
		if robotDelegate == nil {
			return
		}
		//创建机器人
		robot := newRobot()
		//机器人基础数据
		robot.SetUserBaseInfo(baseInfo)
		//机器人代理绑定机器人
		robot.SetDelegate(robotDelegate)
		//robotDelegate.SetPlayer(robot)
		robotDelegate.SetStrategy(s.RobotStrategy)
		//保存到空闲机器人容器中
		s.freeRobots = append(s.freeRobots, robot)
	}
}

/// 取出一个机器人
func (s *CRobotMgr) Pick() iface.IPlayer {
	var robot iface.IPlayer
	s.lock.Lock()
	if len(s.freeRobots) > 0 {
		robot = s.freeRobots[0]
		s.freeRobots = append(s.freeRobots[:0], s.freeRobots[1:]...)
		//robot.Reset()
		s.usedRobots = append(s.usedRobots, robot)
	}
	s.lock.Unlock()
	return robot
}

/// 回收机器人对象
func (s *CRobotMgr) Delete(userID int64) {
	var robotDelegate iface.IRobotDelegate
	s.lock.Lock()
	for i, robot := range s.usedRobots {
		if robot.GetUserID() == userID {
			s.usedRobots = append(s.usedRobots[:i], s.usedRobots[i+1:0]...)
			robotDelegate = robot.GetDelegate()
			s.freeRobots = append(s.freeRobots, robot)
			break
		}
	}
	//子游戏机器人代理数据复位
	if robotDelegate != nil {
		robotDelegate.Reposition()
		//robot.Reset()连机器人基础数据都清理
	}
	s.lock.Unlock()
}

/// 机器人库存判空
func (s *CRobotMgr) IsEmpty() bool {
	emtpy := false
	s.lock.RLock()
	if len(s.freeRobots) == 0 {
		emtpy = true
	}
	s.lock.RUnlock()
	return emtpy
}

/// 机器人定时入桌
func (s *CRobotMgr) OnTimerRobotEnter(worker interface{}) {
	if s.RoomInfo.RoomStatus == define.SERVER_STOPPED {
		return
	}
	if !s.RoomInfo.EnableRobot {
		fmt.Println("机器人被禁用了")
		return
	}
	if s.IsEmpty() {
		fmt.Println("机器人没用库存了")
		return
	}
	//遍历检查使用中的桌子
	usedDesks := DeskMgr().GetUsedDesks(worker)
	for _, desk := range usedDesks {
		//检查机器人入桌条件
		if desk.CheckRobotEnter() < 0 {
			break
		}
	}
}
