package impl

import (
	"mygame/public/define"
	"mygame/public/iface"
)

/// <summary>
/// 玩家类
/// <summary>
type CPlayer struct {
	deskID   uint16               //桌子ID
	chairID  uint16               //座位号
	status   define.UserStatus    //玩家状态
	trustee  bool                 //是否托管
	rank     uint32               //用户排名
	baseInfo *define.UserBaseInfo //基础数据
}

/// 创建玩家实例
func newPlayer() iface.IPlayer {
	return &CPlayer{status: define.USER_STATUS_FRREE, baseInfo: &define.UserBaseInfo{}}
}

/// 获取机器人代理
func (s *CPlayer) GetDelegate() iface.IRobotDelegate {
	return nil
}

/// 设置机器人代理
func (s *CPlayer) SetDelegate(delegate iface.IRobotDelegate) {

}

/// 重置
func (s *CPlayer) Reset() {
	s.deskID = define.INVALID_DESK
	s.chairID = define.INVALID_CHAIR
	s.status = define.USER_STATUS_FRREE
	s.trustee = false
	s.rank = 0
	s.baseInfo = &define.UserBaseInfo{}
}

/// 是否机机器人
func (s *CPlayer) IsRobot() bool {
	return false
}

/// 是否有效
func (s *CPlayer) IsValid() bool {
	return s.deskID != define.INVALID_DESK && s.chairID != define.INVALID_CHAIR && s.baseInfo.UserID > 0
}

/// 获取基础信息
func (s *CPlayer) GetUserBaseInfo() *define.UserBaseInfo {
	return s.baseInfo
}

/// 设置基础西悉尼
func (s *CPlayer) SetUserBaseInfo(baseInfo *define.UserBaseInfo) {
	s.baseInfo = baseInfo
}

/// 子游戏机器人代理接收消息
func (s *CPlayer) SendUserMessage(mainID, subID uint8, msg interface{}) {

}

/// 发送消息到子游戏桌子代理
func (s *CPlayer) SendDeskMessage(subID uint8, msg interface{}) {

}

/// 用户ID
func (s *CPlayer) GetUserID() int64 {
	return s.baseInfo.UserID
}

/// 账号
func (s *CPlayer) GetAccount() string {
	return s.baseInfo.Account
}

/// 昵称
func (s *CPlayer) GetNickName() string {
	return s.baseInfo.NickName
}

/// 头像
func (s *CPlayer) GetHeadID() uint8 {
	return s.baseInfo.HeadID
}

/// 获取桌子ID
func (s *CPlayer) GetDeskID() uint16 {
	return s.deskID
}

/// 设置桌子ID
func (s *CPlayer) SetDeskID(deskID uint16) {
	s.deskID = deskID
}

/// 获取座位号
func (s *CPlayer) GetChairID() uint16 {
	return s.chairID
}

/// 设置座位号
func (s *CPlayer) SetChairID(chairID uint16) {
	s.chairID = chairID
}

/// 获取积分
func (s *CPlayer) GetScore() int64 {
	return s.baseInfo.Score
}

/// 设置积分
func (s *CPlayer) SetScore(score int64) {
	s.baseInfo.Score = score
}

/// 获取位置
func (s *CPlayer) GetLocation() string {
	return s.baseInfo.Location
}

/// 获取用户状态
func (s *CPlayer) GetStatus() define.UserStatus {
	return s.status
}

/// 设置用户状态
func (s *CPlayer) SetStatus(status define.UserStatus) {
	s.status = status
}

/// 设置托管状态
func (s *CPlayer) SetTrustee(trustee bool) {
	s.trustee = trustee
}

/// 获取托管状态
func (s *CPlayer) GetTrustee() bool {
	return s.trustee
}

/// 设置排名
func (s *CPlayer) SetRank(rank uint32) {
	s.rank = rank
}

/// 获取排名
func (s *CPlayer) GetRank() uint32 {
	return s.rank
}
