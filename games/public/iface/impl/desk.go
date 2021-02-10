package impl

import (
	"fmt"
	"math/rand"
	"mygame/core"
	"mygame/public/define"
	"mygame/public/iface"
	"os"
	"sync"
	"time"
)

/// <summary>
/// 桌子类
/// <summary>
type CDesk struct {
	status   define.GameStatus   //游戏状态
	gameInfo *define.GameInfo    //游戏类型
	roomInfo *define.RoomInfo    //游戏房间
	state    *define.DeskState   //桌子基础数据
	players  []iface.IPlayer     //桌子玩家players[chairID]
	delegate iface.IDeskDelegate //子游戏桌子代理，子游戏桌子逻辑继承
	freeIDs  []uint16            //空闲桌子ID
	slot     core.ISlot          //桌子逻辑协程，内含协程安全定时器
	//real, robot uint32         //百人场优化人数获取，真实玩家/机器人玩家数
	lock  *sync.RWMutex
	stock define.StorageInfo
}

/// 创建桌子
func newDesk(state *define.DeskState) iface.IDesk {
	if state == nil {
		return nil
	}
	return &CDesk{lock: &sync.RWMutex{},
		state:   state,
		status:  define.GAME_STATUS_INIT,
		players: []iface.IPlayer{},
		freeIDs: []uint16{}}
}

/// 获取桌子代理
func (s *CDesk) GetDelegate() iface.IDeskDelegate {
	return s.delegate
}

/// 设置桌子代理
func (s *CDesk) SetDelegate(delegate iface.IDeskDelegate) {
	s.delegate = delegate
	delegate.SetDesk(s)
}

/// 桌子初始化
func (s *CDesk) Init(gameInfo *define.GameInfo, roomInfo *define.RoomInfo, slot core.ISlot) {
	if gameInfo == nil || roomInfo == nil || slot == nil {
		return
	}
	s.gameInfo = gameInfo
	s.roomInfo = roomInfo
	//绑定桌子邮槽
	s.slot = slot
	//邮槽添加桌子，所有Add操作完成，才能启动mailBox.Start()
	slot.Add(s)
	//桌子座位占位
	s.players = make([]iface.IPlayer, roomInfo.MaxPlayerCount)
	for i := uint16(0); i < roomInfo.MaxPlayerCount; i++ {
		s.freeIDs = append(s.freeIDs, i)
	}
}

/// 分配桌子ID
func (s *CDesk) allocChairID() (uint16, bool) {
	chairID := uint16(define.INVALID_CHAIR)
	if len(s.freeIDs) > 0 {
		chairID = s.freeIDs[0]
		//座位有人，座位号空闲?
		if s.players[chairID] != nil && s.players[chairID].IsValid() {
			panic(fmt.Sprintf("allocChairID %v", chairID))
		}
		s.freeIDs = append(s.freeIDs[:0], s.freeIDs[1:]...)

		return chairID, true
	}
	return chairID, false
}

/// 回收桌子ID
func (s *CDesk) freeChairID(chairID uint16) {
	if chairID != uint16(define.INVALID_CHAIR) && chairID < uint16(len(s.players)) {
		//座位有人不能回收
		if s.players[chairID] != nil && s.players[chairID].IsValid() {
			panic(fmt.Sprintf("allocChairID %v", chairID))
		}
		s.freeIDs = append(s.freeIDs, chairID)
	}
}

/// 桌子牌局编号
func (s *CDesk) NewRoundID() string {
	return fmt.Sprintf("%v-%v-%v-%v-%v", s.roomInfo.RoomID, time.Now().Unix(), os.Getpid(), s.GetDeskID(), rand.Int()/10)
}

/// 机器人入桌检查
func (s *CDesk) CheckRobotEnter() int {
	return -1
}

/// 返回桌子油槽
func (s *CDesk) GetSlot() core.ISlot {
	return s.slot
}

/// 返回桌子ID
func (s *CDesk) GetDeskID() uint16 {
	return s.state.DeskID
}

/// 返回桌子信息
func (s *CDesk) GetDeskInfo() *define.DeskState {
	return s.state
}

/// 房间游戏信息
func (s *CDesk) GetGameInfo() *define.GameInfo {
	return s.gameInfo
}

/// 获取游戏房间
func (s *CDesk) GetRoomInfo() *define.RoomInfo {
	return s.roomInfo
}

/// 游戏开始
//OnGameStart()

/// 游戏是否开始
//IsGameStarted()

/// 解散游戏
//DismissGame() bool

/// 结束游戏
func (s *CDesk) ConcludeGame(status define.GameStatus) bool {
	return false
}

/// 座位号获取玩家
func (s *CDesk) GetPlayer(chairID uint16) iface.IPlayer {
	if chairID != uint16(define.INVALID_CHAIR) && chairID < uint16(len(s.players)) {
		return s.players[chairID]
	}
	return nil
}

/// 用户ID获取玩家
func (s *CDesk) GetPlayerBy(userID int64) iface.IPlayer {
	for i, player := range s.players {
		if player.GetUserID() == userID {
			return s.players[i]
		}
	}
	return nil
}

/// 座位号是否有人
func (s *CDesk) IsExistUser(chairID uint16) bool {
	if chairID != uint16(define.INVALID_CHAIR) && chairID < uint16(len(s.players)) {
		player := s.players[chairID]
		if player != nil && player.IsValid() {
			return true
		}
	}
	return false
}

/// 座位号是否机器人
func (s *CDesk) IsRobot(chairID uint16) bool {
	if chairID != uint16(define.INVALID_CHAIR) && chairID < uint16(len(s.players)) {
		player := s.players[chairID]
		if player != nil && player.IsValid() {
			return player.IsRobot()
		}
	}
	return false
}

/// 设置游戏状态
func (s *CDesk) SetGameStatus(status define.GameStatus) {
	s.status = status
}

/// 获取游戏状态
func (s *CDesk) GetGameStatus() define.GameStatus {
	return s.status
}

/// 桌子人数
func (s *CDesk) GetPlayerCount() (real, robot uint16) {
	real = 0
	robot = 0
	for _, player := range s.players {
		if player != nil && player.IsValid() {
			if player.IsRobot() {
				robot++
			} else {
				real++
			}
		}
	}
	return
}

/// 桌子人数
func (s *CDesk) GetTotalCount() uint16 {
	count := uint16(0)
	for _, player := range s.players {
		if player != nil && player.IsValid() {
			count++
		}
	}
	return count
}

/// 设置玩家托管
func (s *CDesk) SetTrustee(chairID uint16, trustee bool) {
	if chairID != uint16(define.INVALID_CHAIR) && chairID < uint16(len(s.players)) {
		player := s.players[chairID]
		if player != nil && player.IsValid() {
			player.SetTrustee(trustee)
		}
	}
}

/// 获取托管状态
func (s *CDesk) GetTrustee(chairID uint16) bool {
	if chairID != uint16(define.INVALID_CHAIR) && chairID < uint16(len(s.players)) {
		player := s.players[chairID]
		if player != nil && player.IsValid() {
			return player.GetTrustee()
		}
	}
	return false
}

/// 玩家能否进入
func (s *CDesk) CanUserEnter(player iface.IPlayer) bool {
	return false
}

/// 玩家能否离开
func (s *CDesk) CanUserLeave(userID int64) bool {
	return false
}

/// 玩家进入
func (s *CDesk) OnUserEnter(chairID uint16, trustee bool) {

}

/// 玩家离开
func (s *CDesk) OnUserLeave(chairID uint16, trustee bool) {

}

/// 子游戏桌子消息
func (s *CDesk) OnGameMessage(chairID uint16, subID uint8, msg interface{}) {

}

/// 广播玩家给桌子其余玩家
func (s *CDesk) BroadcastUserToOthers(player iface.IPlayer) {

}

/// 发送其他玩家给指定玩家
func (s *CDesk) SendOtherToUser(player iface.IPlayer, peer iface.IPlayer) {

}

/// 发送其余玩家给指定玩家
func (s *CDesk) SendOthersToUser(player iface.IPlayer) {

}

/// 桌子内广播用户状态
func (s *CDesk) BroadcastUserStatus(player iface.IPlayer, sendToSelf bool) {

}

/// 清理座位玩家
func (s *CDesk) ClearPlayer(chairID uint16, sendState, sendToSelf bool) {

}

/// 计算抽水
func (s *CDesk) CalcRevenue(score int64) {

}

/// 写入玩家积分
func (s *CDesk) WriteUserScore(scoreInfo define.UserScoreInfo, count uint16, strRound string) {

}

/// 更新库存变化
func (s *CDesk) UpdateStorageScore(changeScore int64) {

}

/// 获取当前库存
func (s *CDesk) GetStorageScore() *define.StorageInfo {
	return &s.stock
}

//保存游戏记录
func (s *CDesk) SaveRecord(record define.ReplayRecord) {

}

//读取系统库存
func (s *CDesk) ReadStorageScore() {

}
