package define

import "time"

/// <summary>
/// 游戏类型
/// <summary>
type GameType uint8

const (
	GAMETYPE_BAIREN GameType = 101 //百人类
	GAMETYPE__PK                   //对战类
)

/// <summary>
/// 用户基础数据
/// <summary>
type UserBaseInfo struct {
	UserID       int64  //用户ID
	Account      string //账号
	HeadID       uint8  //头像
	NickName     string //昵称
	Score        int64  //用户积分
	AgentID      uint32 //代理ID
	Status       uint32 //账号状态
	IP           uint32 //客户端IP
	Location     string //地理位置
	TakeMinScore int64  //最小携带分
	TakeMaxScore int64  //最大携带分
}

/// <summary>
/// 游戏类型信息
/// <summary>
type GameInfo struct {
	GameID       uint32   //游戏类型
	GameName     string   //游戏名称
	SortID       uint32   //客户端显示排序
	GameType     GameType //游戏分类 0-百人  1-对战
	ServiceName  string   //子游戏模块名称
	RevenueRatio uint8    //抽水比例
}

/// <summary>
/// 游戏房间信息
/// <summary>
type RoomInfo struct {
	GameID               uint32      //游戏类型
	RoomID               uint32      //房间ID
	RoomName             string      //房间名称
	RoomStatus           ServerState //房间服务状态
	DeskCount            uint16      //房间桌子数量
	EnableRobot          bool        //是否开启机器人
	FloorScore           int64       //房间底注
	CeilScore            int64       //房间顶注
	EnterMinScore        int64       //最小进入分
	EnterMaxScore        int64       //最大进入分
	MinPlayerCount       uint16      //每张桌子最小游戏人数
	MaxPlayerCount       uint16      //每张桌子最大游戏人数
	MaxRobotCount        uint16      //每张桌子最大机器人数
	BroadcastScore       int64       //跑马灯要求分数
	MaxBetScore          int64       //每张桌子最大可下注分
	TotalStock           int64       //房间系统库存
	TotalStockLowerLimit int64       //系统库存低于该值要收分
	TotalStockHighLimit  int64       //系统库存高于该值要放水
	SysKillAllRatio      uint32      //通杀率
	SysReduceRatio       uint32      //库存衰减率
	SysChangeCardRatio   uint32      //系统换牌率
	Chips                []int64     //房间筹码配置表
}

/// <summary>
/// 桌子基础信息
/// <summary>
type DeskState struct {
	DeskID uint16 //桌子ID
	IsLock bool   //锁定
	IsLook bool   //观战
}

/// <summary>
/// 用户结算数据
/// <summary>
type UserScoreInfo struct {
	ChairID      uint16    //座位号
	IsBanker     uint8     //是否庄家
	WinScorePure int64     //当局输赢分(扣除抽水净收入)
	BetScore     int64     //用户总押注分
	Revenue      int64     //当局税收(抽水)
	WinScore     int64     //有效投注额(税前输赢)
	StartTime    time.Time //当局开始时间
	CardValue    string    //当局开牌数据
}

/// <summary>
/// 系统库存数据
/// <summary>
type StorageInfo struct {
	EndStorage         int64  //当前库存
	LowLimit           int64  //最小库存
	UpLimit            int64  //最大库存
	SysAllKillRatio    uint32 //系统通杀率
	SysReduceRatio     uint32 //库存衰减
	SysChangeCardRatio uint32 //系统换牌率
}

/// <summary>
/// 对局玩家数据
/// <summary>
type ReplayPlayer struct {
	UserID  int64
	ChairID uint16
	Account string
	Score   int64 //玩家积分
	Valid   bool  //是否有效
}

/// <summary>
/// 对局单步操作
/// <summary>
type ReplayStep struct {
	time    time.Time
	Bet     string
	Round   int32
	OptTy   int32  //操作类型
	ChairID uint16 //操作位置
	Pos     uint16 //被操作位置，比牌对方
	Valid   bool   //是否有效
}

/// <summary>
/// 对局结果数据
/// <summary>
type ReplayResult struct {
	ChairID  uint16
	Pos      uint16 //被操作位置，比牌对方
	BetScore int64
	WinScore int64
	CardType string
	IsBanker bool
	Valid    bool
}

/// <summary>
/// 对局记录数据
/// <summary>
type ReplayRecord struct {
	RoundID    string         //牌局编号
	RoomName   string         //房间名称
	FloorScore int64          //底注
	Detail     []byte         //对局详情
	Players    []ReplayPlayer //玩家
	Steps      []ReplayStep   //游戏过程
	Results    []ReplayResult //游戏结果
}

/// <summary>
/// 机器人区域下注参数配置
/// <summary>
type RobotStrategyArea struct {
	Weight    int32
	LowTimes  int32
	HighTimes int32
}

/// <summary>
//机器人参数配置
/// <summary>
type RobotStrategyParam struct {
	GameID         int32
	RoomID         int32
	EnterLowScore  int64
	EnterHighScore int64
	MinScore       int64
	MaxScore       int64
}

type GameStatus uint8

/// <summary>
/// 游戏状态
/// <summary>
const (
	GAME_STATUS_INIT  GameStatus = iota
	GAME_STATUS_START GameStatus = 100
	GAME_STATUS_END   GameStatus = 200
)

type UserStatus uint8

/// <summary>
/// 用户状态
/// <summary>
const (
	USER_STATUS_FRREE   UserStatus = iota //空闲
	USER_STATUS_READY                     //准备
	USER_STATUS_PLAY                      //游戏
	USER_STATUS_OFFLINE                   //离线
	USER_STATUS_LOOK                      //观战

)
const (
	INVALID_DESK  = uint16(0xFFFF)
	INVALID_CHAIR = uint16(0xFFFF)
)

/// <summary>
/// 游戏记录操作
/// <summary>
type ReplayOpType uint8

const (
	OP_START    ReplayOpType = iota
	OP_BANKER                //定庄
	OP_BET                   //下注
	OP_FOLLOW                //跟注
	OP_ADDBet                //加注
	OP_CMPCARD               //比牌
	OP_LOOKCARD              //看牌
)

/// <summary>
/// 游戏服务状态
/// <summary>
type ServerState uint8

const (
	SERVER_STOPPED ServerState = iota
	SERVER_RUNNING
	SERVER_REPAIRING
)
