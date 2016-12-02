// 客户端连接管理
package session

//
import (
	"util/logs"

	"core/safe/list"
)

var _ = logs.Debug

// 被踢的玩家
var g_slist = slist.New()

// 踢玩家
func Disconnect(s *Session) {
	if !s.isDisconnect() {
		s.setDisconnect()
		g_slist.PushBack(s.AccId)
	}
}

// 处理被踢的玩家
func ProcDisconnect() {
	// 获取被踢列表
	lst := g_slist

	// 空
	if lst.Len() <= 0 {
		return
	}

	// 遍历处理
	e := lst.Begin()
	for ; e != nil; e = e.Next() {
		doKick(e.Value.(int))
	}
	lst.End()

	// 清空list
	lst.Clear()
}
