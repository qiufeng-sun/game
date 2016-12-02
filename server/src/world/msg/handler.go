// 消息处理注册管理
package msg

//
import (
	"core/net/socket"
)

// 注册信息封装
type HandleInfo struct {
	GmLv    int  // 最低gm等级需求
	InWorld bool // 游戏内处理
	Serial  bool // 必须串行处理
}

func NewHandleInfo(gmLv int, inWorld, serial bool) *HandleInfo {
	return &HandleInfo{gmLv, inWorld, serial}
}

// 消息处理注册管理器对象
var g_msgHandlers *socket.MsgHandler

// 包init
func init() {
	g_msgHandlers = socket.NewMsgHandler()
}

// 注册
func RegHandler(msgId int32, handler socket.IHandler, gmLv int, inWorld, serial bool) {
	g_msgHandlers.RegHandler(msgId, handler, NewHandleInfo(gmLv, inWorld, serial))
}

// 获取handler
func Handler(msgId int32) (socket.IHandler, *HandleInfo, bool) {
	if h, info, ok := g_msgHandlers.Handler(msgId); ok {
		return h, info.(*HandleInfo), ok
	}

	return nil, nil, false
}
