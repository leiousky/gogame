```
自实现非侵入式棋牌子游戏异步框架（多核架构）

│  main.go                                主函数入口
│  README.md                              文档说明
│  
├─comm                                    通用函数
│  └─utils
│          utils.go
│          
├─core                                    核心组件
│  │  calback.go                          回调定义
│  │  mailbox.go                          邮槽管理
│  │  msq.go                              消息队列
│  │  proc.go                             消息处理器(EventLoop)
│  │  slot.go                             邮槽(EventLoopThread)
│  │  worker.go                           业务处理
│  │  
│  └─timerv2                              定时器
│          cron_test.go
│          hashwheel.go
│          README.md
│          safetimer.go
│          safetimer_test.go
│          timer.go
│          
├─public                                  公共部分
│  ├─define
│  │      define.go                       公共定义
│  │      
│  └─iface                                接口
│      │  idesk.go                        桌子接口
│      │  idesk_delegate.go               游戏桌子接口      --子游戏桌子逻辑继承使用
│      │  iplayer.go                      玩家接口
│      │  irobot_delegate.go              子游戏机器人接口  --子游戏机器人继承使用
│      │  
│      └─impl                             实现
│              desk.go                    桌子类
│              desk_mgr.go                桌子管理类
│              player.go                  玩家类
│              player_mgr.go              玩家管理类
│              robot.go                   机器人类
│              robot_mgr.go               机器人管理类
│              
├─service                                 服务组件
│      sentry.go                          对应每张桌子业务消息入口
│      smain.go                           对应每张桌子业务逻辑处理 - 子游戏框架实现
│      
└─sub                                     各个子游戏模块           - 子游戏开发实现
    └─game_dragon_tiger                   龙虎斗
            desk.go                       子游戏桌子类
            robot.go                      子游戏机器人类
            
