// 消息处理函数注册扩展
package hfunc

//
import (
	"util/logs"

	. "msg/pbmsg"

	"world/msg"
	. "world/player/session"
)

// 处理接口封装
type handleFunc func(s *Session, msg []byte)

func (h handleFunc) Handle(s interface{}, msg []byte) {
	defer func() {
		if e := recover(); e != nil {
			logs.Warnln(e)
		}
	}()

	h(s.(*Session), msg)
}

// 注册处理接口
func RegHandler(msgId EMsgId, handler func(s *Session, msg []byte), gmLv int, inWorld, serial bool) {
	msg.RegHandler(int32(msgId), handleFunc(handler), gmLv, inWorld, serial)
}

func RegOutWorldHandler(msgId EMsgId, handler func(s *Session, msg []byte), gmLv int) {
	RegHandler(msgId, handler, gmLv, false, true) // 该类消息在主协程中单独更新处理
}

func RegInWorldSerialHandler(msgId EMsgId, handler func(s *Session, msg []byte), gmLv int) {
	RegHandler(msgId, handler, gmLv, true, true)
}

func RegInWorldParallelHandler(msgId EMsgId, handler func(s *Session, msg []byte), gmLv int) {
	RegHandler(msgId, handler, gmLv, true, false)
}
