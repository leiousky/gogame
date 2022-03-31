package iface

import "games/public/define"

/// <summary>
/// 玩家接口
/// <summary>
type IPlayer interface {
	/// 重置
	Reset()
	/// 获取机器人代理
	GetDelegate() IRobotDelegate
	/// 设置机器人代理
	SetDelegate(delegate IRobotDelegate)
	/// 是否机机器人
	IsRobot() bool
	/// 是否有效
	IsValid() bool
	/// 获取基础信息
	GetUserBaseInfo() *define.UserBaseInfo
	/// 设置基础西悉尼
	SetUserBaseInfo(baseInfo *define.UserBaseInfo)
	/// 子游戏机器人代理接收消息
	SendUserMessage(mainID, subID uint8, msg interface{})
	/// 发送消息到子游戏桌子代理
	SendDeskMessage(subID uint8, msg interface{})
	/// 用户ID
	GetUserID() int64
	/// 账号
	GetAccount() string
	/// 昵称
	GetNickName() string
	/// 头像
	GetHeadID() uint8
	/// 获取桌子ID
	GetDeskID() uint16
	/// 设置桌子ID
	SetDeskID(deskID uint16)
	/// 获取座位号
	GetChairID() uint16
	/// 设置座位号
	SetChairID(chairID uint16)
	/// 获取积分
	GetScore() int64
	/// 设置积分
	SetScore(score int64)
	/// 获取位置
	GetLocation() string
	/// 获取用户状态
	GetStatus() define.UserStatus
	/// 设置用户状态
	SetStatus(status define.UserStatus)
	/// 设置托管状态
	SetTrustee(trustee bool)
	/// 获取托管状态
	GetTrustee() bool
	/// 设置排名
	SetRank(rank uint32)
	/// 获取排名
	GetRank() uint32
}
