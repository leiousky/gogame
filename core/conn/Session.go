package conn

type Reason uint8

const (
	KNoError    Reason = Reason(0)
	KPeerClosed Reason = Reason(1) //对端关闭连接
	KSelfClosed Reason = Reason(2) //本端正常关闭
	KSelfExcept Reason = Reason(3) //本端异常关闭
)

type State uint8

const (
	KDisconnected State = State(0)
	KConnected    State = State(1)
)

type Type uint8

const (
	KClient Type = Type(0)
	KServer Type = Type(1)
)

/// <summary>
/// Session 连接会话
/// <summary>
type Session interface {
	ID() int64
	Name() string
	IsWebsocket() bool
	Type() Type
	Connected() bool
	Conn() interface{}
	LocalAddr() string
	RemoteAddr() string
	SetContext(key int, val interface{})
	GetContext(key int) interface{}
	Write(msg interface{})
	Close()
}
