package iface

import (
	"mygame/core"
	"mygame/public/define"
)

/// <summary>
/// 桌子接口
/// <summary>
type IDesk interface {
	/// 获取桌子代理
	GetDelegate() IDeskDelegate
	/// 设置桌子代理
	SetDelegate(delegate IDeskDelegate)
	/// 桌子牌局编号
	NewRoundID() string
	/// 机器人入桌检查
	CheckRobotEnter() int
	/// 返回桌子油槽
	GetSlot() core.ISlot
	/// 返回桌子ID
	GetDeskID() uint16
	/// 返回桌子信息
	GetDeskInfo() *define.DeskState
	/// 房间游戏信息
	GetGameInfo() *define.GameInfo
	/// 获取游戏房间
	GetRoomInfo() *define.RoomInfo
	/// 游戏开始
	//OnGameStart()
	/// 游戏是否开始
	//IsGameStarted()
	/// 解散游戏
	//DismissGame() bool
	/// 结束游戏
	ConcludeGame(status define.GameStatus) bool
	/// 座位号获取玩家
	GetPlayer(chairID uint16) IPlayer
	/// 用户ID获取玩家
	GetPlayerBy(userID int64) IPlayer
	/// 座位号是否有人
	IsExistUser(chairID uint16) bool
	/// 座位号是否机器人
	IsRobot(chairID uint16) bool
	/// 设置游戏状态
	SetGameStatus(status define.GameStatus)
	/// 获取游戏状态
	GetGameStatus() define.GameStatus
	/// 桌子人数
	GetPlayerCount() (real, robot uint16)
	/// 桌子人数
	GetTotalCount() uint16
	/// 设置玩家托管
	SetTrustee(chairID uint16, trustee bool)
	/// 获取托管状态
	GetTrustee(chairID uint16) bool
	/// 玩家能否进入
	CanUserEnter(player IPlayer) bool
	/// 玩家能否离开
	CanUserLeave(userID int64) bool
	/// 玩家进入
	OnUserEnter(chairID uint16, trustee bool)
	/// 玩家离开
	OnUserLeave(chairID uint16, trustee bool)
	/// 子游戏桌子消息
	OnGameMessage(chairID uint16, subID uint8, msg interface{})
	/// 广播玩家给桌子其余玩家
	BroadcastUserToOthers(player IPlayer)
	/// 发送其他玩家给指定玩家
	SendOtherToUser(player IPlayer, peer IPlayer)
	/// 发送其余玩家给指定玩家
	SendOthersToUser(player IPlayer)
	/// 桌子内广播用户状态
	BroadcastUserStatus(player IPlayer, sendToSelf bool)
	/// 清理座位玩家
	ClearPlayer(chairID uint16, sendState, sendToSelf bool)
	/// 计算抽水
	CalcRevenue(score int64)
	/// 写入玩家积分
	WriteUserScore(scoreInfo define.UserScoreInfo, count uint16, strRound string)
	/// 更新库存变化
	UpdateStorageScore(changeScore int64)
	/// 获取当前库存
	GetStorageScore() *define.StorageInfo
	//保存游戏记录
	SaveRecord(record define.ReplayRecord)
}
