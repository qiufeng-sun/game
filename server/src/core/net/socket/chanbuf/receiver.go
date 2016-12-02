package chanbuf

import (
	"errors"
	"io"
	"net"

	"core/net/msg"
)

var EReceiverFull = errors.New("receiver chan is full!")

//
type ChanReceiver struct {
	chMsg chan []byte

	msgNum int
}

//
func NewChanReceiver(num int) *ChanReceiver {
	r := &ChanReceiver{msgNum: num}
	r.reset()
	return r
}

//
func (this *ChanReceiver) reset() {
	this.chMsg = make(chan []byte, this.msgNum)
}

//
func (this *ChanReceiver) Recv(c net.Conn) (int64, error) {
	// sz
	var bsz [4]byte
	if _, e := io.ReadFull(c, bsz[:]); e != nil {
		return 0, e
	}

	sz, _ := msg.Uint32ByBytes(bsz[:])

	buff := make([]byte, int(sz))
	if _, e := io.ReadFull(c, buff); e != nil {
		return 0, e
	}

	select {
	case this.chMsg <- buff:
	default:
		return 0, EReceiverFull
	}

	return int64(4 + sz), nil
}

//
func (this *ChanReceiver) Check() bool {
	return len(this.chMsg) > 0
}

//
func (this *ChanReceiver) GetMsg() ([]byte, bool) {
	select {
	case b := <-this.chMsg:
		return b, true
	default:
		return nil, false
	}
}

//
func (this *ChanReceiver) Release([]byte) {
	// do nothing
}

//
func (this *ChanReceiver) Clear() {
	this.reset()
}
