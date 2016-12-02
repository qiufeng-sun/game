// 客户端连接管理
package session

//
import (
	"util/logs"

	"core/net/socket"
)

var _ = logs.Debug

// 所有session
var (
	g_mapAccIdSession = make(map[int]*Session) // accId=>*Session(全部通过验证的玩家)
	g_mapPlyIdSession = make(map[int]*Session) // playerId=>*Session(全部进入游戏玩家，即已读取player信息的玩家)
)

//
func GetSessionByAccId(accId int) *Session {
	return g_mapAccIdSession[accId]
}

//
func GetSessionByPlayerId(playerId int) *Session {
	return g_mapPlyIdSession[playerId]
}

// 增加新的session
func AddSession(netId int, accInfo *AccInfo) *Session {
	s := NewSession(netId, accInfo)

	doKick(s.AccId)
	g_mapAccIdSession[s.AccId] = s

	return s
}

//
func doKick(accId int) {
	//
	s := GetSessionByAccId(accId)
	if nil == s {
		return
	}
	defer delete(g_mapAccIdSession, accId)

	// player
	p := s.GetPlayer()
	if p != nil {
		p.DoKick()
		defer delete(g_mapPlyIdSession, p.Id)
	}

	// socket
	socket.KickClient(s.netId)
}

// update
func Update() {
	for _, s := range g_mapAccIdSession {
		s.Update()
	}
}
