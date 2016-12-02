// 进入world的玩家对象(player)
package session

//
import (
	"util/logs"
)

var _ = logs.Debug

// 玩家对象
type Player struct {
	Id int // 玩家id
}

//
func (this *Player) DoKick() {
	// save
	this.Save(true)
}

// 保存//??
func (this *Player) Save(logoff bool) {

}

//////////??
// db中玩家对象// to do del
type DBPlayer struct {
}

// 保存前更新处理//??
func (p *Player) preSave(logoff bool) *DBPlayer {
	return nil
}

// 玩家登录
func (p *Player) Logon(s *Session) {
	//	// 设置属性
	//	p.setSession(s)
	//	s.SetPlayer(p)
	//	g_mapPlayer[p.Id] = p

	//	// 通知玩家
	//	send := &NsLogin{
	//		Error:     msg.Err2Int32(EMsgErr_Success),
	//		HasPlayer: proto.Bool(true),
	//	}
	//	s.SendMsg(EMsgId_ID_NsLogin, send)
}
