```
Golang实现异步游戏框架（多核高并发），强伸缩，可扩展 ，含子游戏框架及机器人管理，业务非侵入式
│
│  main.go                     主函数入口
│  README.md                   文档说明
│  
├─comm                        【通用函数】
│  └─utils
│          freevalues.go
│          ini.go
│          Orderedmap.go
│          Semaphore.go
│          Timestamp.go
│          utils.go
│ 
│          
├─core                        【核心组件(业务层、网络层)】
│  │
│  │                           业务层
│  │
│  │  events.go                业务事件(网络读写/自定义)
│  │  mailbox.go               邮槽管理(容器)
│  │  proc.go                  单元(cell)业务处理器(EventLoop)
│  │  slot.go                  业务邮槽(EventLoopThread)
│  │  worker.go                业务接口
│  │        
│  ├─conn                      网络层
│  │  │  
│  │  │  Session.go            会话
│  │  │  SessionID.go          会话ID
│  │  │  SessionMgr.go         会话管理
│  │  │  
│  │  ├─tcp                    TCP网络
│  │  │  │  Acceptor.go        服务端接受器
│  │  │  │  Connector.go       客户端连接器
│  │  │  │  TCPConnection.go   TCP连接对象
│  │  │  │  
│  │  │  ├─tcpclient           TCP客户端
│  │  │  │      TCPClient.go
│  │  │  │      
│  │  │  └─tcpserver           TCP服务端
│  │  │          TCPServer.go
│  │  │          
│  │  └─transmit               TCP流传输接口(读写解析)
│  │      │  channel.go        TCP流协议读写
│  │      │  
│  │      ├─tcp_channel        TCP协议读写接口
│  │      │      Channel.go    TCP协议读写(默认)
│  │      │      ReadWrite.go
│  │      │      
│  │      └─ws_channel         websocket协议读写接口
│  │              Channel.go   websocket协议读写(默认)
│  │              ReadWrite.go
│  │              
│  ├─msq                       消息队列(待完善)
│  │      BlockChanMsq.go      阻塞chan
│  │      BlockListMsq.go      阻塞list
│  │      BlockVecMsq.go       阻塞vector
│  │      FreeChanMsq.go       非阻塞chan
│  │      FreeListMsq.go       非阻塞list
│  │      FreeVecMsq.go        非阻塞vector
│  │      msq.go               消息队列接口
│  │      
│  ├─timer                     定时器
│  │      ScopedTimer.go       局部定时器，线程安全
│  │      TimerWheel.go        时间轮盘
│  │      
│  ├─timerv2                   定时器(基于跳表实现)
│  │       cron_test.go
│  │       hashwheel.go
│  │       safetimer.go
│  │       safetimer_test.go
│  │       timer.go
│  │
│  └─callback                  回调定义
│         Callback.go
│
│
├─server                      【业务层网络】
│  │
│  ├─stream                    TCP流协议读写解析
│  │  │  checksum.go           业务消息校验
│  │  │  msg.go                业务消息定义
│  │  │  
│  │  ├─tcp_stream             TCP协议读写实现
│  │  │      Channel.go
│  │  │      
│  │  └─ws_stream              websocket协议读写实现
│  │          Channel.go
│  │          
│  ├─tcp_client                TCP客户端
│  │      tcpclient.go
│  │      
│  └─tcp_server                TCP服务端
│          tcpserver.go
│          
│
├─public                      【游戏业务公共部分】
│  │
│  ├─define                    公共定义
│  │      define.go
│  │      
│  └─iface                     接口
│      │  idesk.go             桌子接口
│      │  idesk_delegate.go    游戏桌子接口      --子游戏桌子逻辑继承使用
│      │  iplayer.go           玩家接口
│      │  ireplay_record.go    子游戏机器人接口  --子游戏机器人继承使用
│      │  irobot_delegate.go   游戏记录/回放
│      │  
│      └─impl                  实现
│              desk.go         桌子类
│              desk_mgr.go     桌子管理类
│              player.go       玩家类
│              player_mgr.go   玩家管理类
│              robot.go        机器人类
│              robot_mgr.go    机器人管理类
│              
├─service                      服务组件(桌子)             - 异步主服务(网络业务、自定义业务)
│      sentry.go               对应每张桌子业务消息入口
│      smain.go                对应每张桌子业务逻辑处理   - 子游戏框架实现
│
├─db                           服务组件(db/redis)         - 异步数据服务(mysql/mongodb/redis数据库及缓存处理)
│      sentry.go               对应每个数据库任务消息入口
│      smain.go                对应每个数据库业务逻辑处理 - 子游戏框架实现
│            
└─sub                          各个子游戏模块             - 子游戏开发实现
    └─game_dragon_tiger
            desk.go            子游戏桌子类
            robot.go           子游戏机器人类
            
