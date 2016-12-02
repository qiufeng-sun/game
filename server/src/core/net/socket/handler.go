// 消息处理注册管理
package socket

//

// 消息处理接口
type IHandler interface {
	Handle(receiver interface{}, msg []byte)
}

// 消息处理entry
type handleEntry struct {
	msgId   int32       // 消息id
	handler IHandler    // handler
	info    interface{} // 信息
}

func newhandleEntry(msgId int32, handler IHandler, info interface{}) *handleEntry {
	return &handleEntry{msgId, handler, info}
}

// 消息处理注册管理器
type MsgHandler struct {
	entrys map[int32]*handleEntry // <msgId, *handleEntry>
}

// new
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{make(map[int32]*handleEntry)}
}

// 注册
func (mh *MsgHandler) RegHandler(msgId int32, handler IHandler, info interface{}) {
	// handler检查
	if nil == handler {
		panic("MsgHandler.HandleFunc(): nil handler! msgId =" + string(msgId))
	}

	// 是否已注册过
	if _, exist := mh.entrys[msgId]; exist {
		panic("MsgHandler.HandleFunc(): multiple registrations! msgId =" + string(msgId))
	}

	// 注册
	mh.entrys[msgId] = newhandleEntry(msgId, handler, info)
}

// 获取handler
func (mh *MsgHandler) Handler(msgId int32) (IHandler, interface{}, bool) {
	entry, exist := mh.entrys[msgId]
	if !exist {
		return nil, nil, false
	}

	return entry.handler, entry.info, true
}
