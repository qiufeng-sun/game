// 环形消息缓冲及发送
package ringbuf

//
import (
	"net"
	"sync"

	"core/buff"
)

// 发送消息缓冲
type RingSender struct {
	*buff.RingBuff // 消息发送缓冲
	*sync.Mutex

	chSend chan bool      // 消息发送通知
	ioVect *buff.IoVector // 发送数据
}

// 创建
func NewRingSender(sz int) *RingSender {
	return &RingSender{
		RingBuff: buff.NewRingBuff(sz),
		Mutex:    &sync.Mutex{},
		chSend:   make(chan bool, 1),
		ioVect:   &buff.IoVector{},
	}
}

// 发送消息
func (s *RingSender) Send(conn net.Conn) error {
	// 获取缓存的消息
	var _, buffed = s.GetBuffedIoVector(s.ioVect)

	// 没有消息
	if 0 == buffed {
		return nil
	}

	// 发送
	if err := s.send(conn, s.ioVect.Buff0); err != nil {
		return err
	}

	if s.ioVect.Size1() > 0 {
		if err := s.send(conn, s.ioVect.Buff1); err != nil {
			return err
		}
	}

	return nil
}

func (s *RingSender) send(conn net.Conn, data []byte) error {
	var start, sz int
	var err error

	for start < len(data) {
		sz, err = conn.Write(data[start:])

		if sz > 0 {
			start += sz
		}

		if err != nil {
			break
		}
	}

	// 移动
	if start > 0 {
		s.Release(start)
	}

	return err
}

// 写入发送缓冲
func (s *RingSender) Write(b1 []byte, b2 []byte) (n int, e error) {
	s.Lock()
	defer s.Unlock()

	n1, e := s.RingBuff.Write(b1)
	if e != nil {
		return n1, e
	}

	n2, e := s.RingBuff.Write(b2)
	if e != nil {
		return n1 + n2, e
	}

	// 通知发送协程
	select {
	case s.chSend <- true:
	default:
	}

	return n1 + n2, nil
}

//
func (s *RingSender) WatchSend() <-chan bool {
	return s.chSend
}
