package msq

/// <summary>
/// MsgQueue 消息队列
/// <summary>
type MsgQueue interface {
	//入队列
	Push(interface{}) error
	//出队列
	Pop() (interface{}, bool)
	//掏空队列
	Pick() ([]interface{}, bool)
	//队列消息数量
	Count() int64
	//阻塞的话则唤醒
	Signal()
	//阻塞或非阻塞设置
	EnableNonBlocking(bv bool)
	//关闭清理
	Close()
}
