package iface

/// <summary>
/// 子游戏桌子代理接口，子游戏桌子逻辑继承
/// <summary>
type IDeskDelegate interface {
	/// 指定桌子
	SetDesk(desk IDesk)
	/// 桌子复位
	Reposition()
	/// 游戏开始
	OnGameStart()
	//// 游戏结束
	OnGameConclude(chairID uint16, reason uint8) bool
	/// 场景推送
	OnGameScene(chairID uint16, look bool)
	/// 游戏消息
	OnGameMessage(chairID uint16, subID uint8, msg interface{})
	/// 用户进入
	OnUserEnter(userID int64, look bool) bool
	/// 用户准备
	OnUserReady(userID int64, look bool) bool
	/// 用户离开
	OnUserLeave(userID int64, look bool) bool
	/// 能否进入
	CanUserEnter(player IPlayer) bool
	/// 能否离开
	CanUserLeave(userID int64)
}

/// <summary>
/// 桌子代理创建器
/// <summary>
type IDeskDelegateCreator interface {
	Create() IDeskDelegate
}
