package iface

import "games/public/define"

/// <summary>
/// 机器人代理接口，子游戏机器人逻辑继承
/// <summary>
type IRobotDelegate interface {
	/// 设置机器人
	SetPlayer(player IPlayer)
	/// 设置桌子
	SetDesk(desk IDesk)
	/// 机器人复位
	Reposition()
	/// 游戏消息
	OnGameMessage(subID uint8, msg interface{})
	/// 机器人配置
	SetStrategy(strategy *define.RobotStrategyParam)
	/// 机器人配置
	GetStrategy() *define.RobotStrategyParam
}

/// <summary>
/// 机器人代理创建器
/// <summary>
type IRobotDelegateCreator interface {
	Create() IRobotDelegate
}
