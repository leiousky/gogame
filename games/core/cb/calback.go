package cb

type CloseCallback func(peer interface{})

type OnConnected func(peer interface{})

type OnClosed func(peer interface{})

type OnMessage func(msg interface{}, peer interface{})

type OnWritten func(msg interface{}, peer interface{})

type OnError func(peer interface{}, err error)

type ReadCallback func(cmd uint32, msg interface{}, peer interface{})

type CustomCallback func(cmd uint32, msg interface{}, peer interface{})

type CmdCallback func(msg interface{}, peer interface{})

type CmdCallbacks map[uint32]CmdCallback

type TimerCallback func(timerID uint32, dt int32, args interface{}) bool
