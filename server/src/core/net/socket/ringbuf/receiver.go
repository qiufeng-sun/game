// 消息接收及缓冲 -- 环形缓冲
package ringbuf

// import
import (
	"net"

	"core/buff"
	"core/net/msg"
)

//
type RingReceiver struct {
	*buff.RingBuff // 消息接收缓冲
}

// 创建
func NewRingReceiver(sz int) *RingReceiver {
	return &RingReceiver{buff.NewRingBuff(sz)}
}

// 接收消息
func (r *RingReceiver) Recv(conn net.Conn) (int64, error) {
	return r.ReadFrom(conn)
}

// 检查消息是否完整
func (r *RingReceiver) Check() bool {
	iov, _ := r.GetBuffedIoVector(&buff.IoVector{})
	if nil == iov {
		return false
	}

	_, ok := msg.Check(iov)

	return ok
}

//
func (r *RingReceiver) GetMsg() ([]byte, bool) {
	iov, sz := r.GetBuffedIoVector(&buff.IoVector{})
	if sz <= 0 {
		return nil, false
	}

	return msg.Get(iov)
}

//
func (r *RingReceiver) Release(b []byte) {
	r.RingBuff.Release(msg.CalBuffSpace(len(b)))
}
