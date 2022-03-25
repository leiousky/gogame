package transmit

//
type Msg struct {
	ver     uint16
	sign    uint16
	encType uint8
	mainID  uint8
	subID   uint8
	msg     interface{}
}

//
func newMsg(mainID uint8, subID uint8, msg interface{}) *Msg {
	return &Msg{
		ver:     0x0001,
		sign:    0x5F5F,
		encType: 0x02,
		mainID:  mainID,
		subID:   subID,
		msg:     msg,
	}
}
