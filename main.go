package main

import (
	"games/core"
	"games/public/define"
	"games/public/iface/impl"
	"games/server/tcp_client"
	"games/service"
	"games/sub/game_dragon_tiger"
	"time"
)

//go mod init games
//go mod tidy
func main() {

	tcpclient := tcp_client.NewTCPClient("client")
	tcpclient.ConnectTCP("127.0.0.1:8099")

	var GameInfo define.GameInfo //游戏类型
	var RoomInfo define.RoomInfo //游戏房间
	var mailbox core.IMailbox

	GameInfo.GameName = "龙虎斗"
	GameInfo.GameID = 19
	GameInfo.SortID = 0
	GameInfo.ServiceName = "game_dragon_tiger"
	GameInfo.RevenueRatio = 10
	GameInfo.GameType = define.GAMETYPE_BAIREN
	//房间信息
	RoomInfo.GameID = 19
	RoomInfo.RoomID = 19015
	RoomInfo.RoomName = "初级房"
	RoomInfo.DeskCount = 1                              //房间桌子数量
	RoomInfo.FloorScore = 1                             //底注
	RoomInfo.CeilScore = 1000000                        //顶柱
	RoomInfo.EnterMinScore = 222                        //进入最小分
	RoomInfo.EnterMaxScore = 8888888888                 //进入最大分
	RoomInfo.MinPlayerCount = 2                         //桌子最小游戏人数
	RoomInfo.MaxPlayerCount = 3                         //桌子最大游戏人数
	RoomInfo.MaxRobotCount = 2                          // 最大机器人数
	RoomInfo.BroadcastScore = 10                        //跑马灯要求分
	RoomInfo.MaxBetScore = 1222                         //每张桌子最大下注分
	RoomInfo.TotalStock = 111111111111                  //当前库存
	RoomInfo.EnableRobot = true                         //开启机器人
	RoomInfo.Chips = []int64{1, 10, 50, 100, 500, 1000} //桌子筹码配置

	//step.1 创建邮槽管理器
	mailbox = core.NewMailBox()

	//step.1.1 添加若干邮槽(桌子业务)
	//time.Millisecond*1000
	mailbox.Add(time.Millisecond*1000, 12, service.NewSentryCreator(), int(RoomInfo.DeskCount))

	//step.1.2 桌子管理器创建n张桌子
	impl.DeskMgr().Init(
		&GameInfo, &RoomInfo, mailbox, //为每张桌子分配邮槽
		game_dragon_tiger.NewDeskDelegateCreator()) //子游戏桌子

	//如果开启机器人
	if RoomInfo.EnableRobot {
		//step.1.3 机器人管理器创建若干机器人
		//impl.RobotMgr().Init(
		//	&GameInfo, &RoomInfo,
		//	game_dragon_tiger.NewRobotDelegateCreator()) //子游戏机器人
	}

	//step.2 添加若干邮槽(db/redis业务)
	//mailbox.Add(time.Millisecond*10000, 12, db.NewSentryCreator(), 12)

	//step.3 启动邮槽协程处理
	mailbox.Start()

	//如果开启机器人
	// if RoomInfo.EnableRobot {
	// 	//step.4 开启机器人入桌检查
	// 	mailbox.Range(func(slot core.ISlot, i int) {
	// 		if i < int(RoomInfo.DeskCount) {
	// 			timer := slot.GetProc().GetTimerv2().(*timer.SafeTimerScheduel)
	// 			worker := slot.GetProc().GetWorker().(*service.Sentry)
	// 			timer.CreatCronFunc("@every 1s", worker.OnTick)
	// 		}
	// 	})
	// }
	//等待推出
	mailbox.Wait()
}
