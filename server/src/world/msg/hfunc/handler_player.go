// 玩家相关
package hfunc

//
import (
	proto "github.com/golang/protobuf/proto"

	"core/net/socket"
	"util/logs"

	. "msg/pbmsg"

	. "world/player/session"
)

// 创建角色// to do
func HandleCreatePlayer(s *Session, data []byte) {
	var msg NcCreatePlayer
	e := socket.ParseMsgData(data, &msg)
	if e != nil {
		logs.Warn("create player: invalid msg! error=%v", e)
	}

	logs.Debug(msg.String())

	send := &NsCreatePlayer{Error: proto.Int(1)}
	s.SendMsg(EMsgId_ID_NsCreatePlayer, send)
}
