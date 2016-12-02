package chanbuf

import (
	"errors"
	"net"
)

var ESenderFull = errors.New("sender chan is full!")

//
type Msg struct {
	h []byte // header
	b []byte // body
}

//
type ChanSender struct {
	chMsg  chan *Msg
	chSend chan bool

	msgNum int
}

//
func NewChanSender(sz int) *ChanSender {
	s := &ChanSender{msgNum: sz}
	s.reset()

	return s
}

//
func (this *ChanSender) reset() {
	this.chMsg = make(chan *Msg, this.msgNum)
	this.chSend = make(chan bool, 1)
}

//
func (this *ChanSender) Send(conn net.Conn) error {
	for {
		select {
		case d := <-this.chMsg:
			_, e := conn.Write(d.h)
			if e != nil {
				return e
			}
			_, e = conn.Write(d.b)
			if e != nil {
				return e
			}
		default:
			return nil
		}
	}

	return nil
}

//
func (this *ChanSender) Write(h, b []byte) (int, error) {
	msg := &Msg{h: h, b: b}
	select {
	case this.chMsg <- msg:

	default:
		return 0, ESenderFull
	}

	select {
	case this.chSend <- true:
	default:
	}

	return len(h) + len(b), nil
}

//
func (this *ChanSender) WatchSend() <-chan bool {
	return this.chSend
}

//
func (this *ChanSender) Clear() {
	this.reset()
}
