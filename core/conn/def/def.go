package def

type Reason uint8

const (
	KNoError Reason = Reason(0)
	KClosed  Reason = Reason(1)
	KExcept  Reason = Reason(3)
)
