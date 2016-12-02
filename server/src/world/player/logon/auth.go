// 未通过验证的session
package logon

//
import (
	"util/logs"

	"core/net/socket"

	. "msg/pbmsg"

	ws "world/player/session"
)

// 登录验证信息
type logonRequest struct {
	digest      string // 验证字符串
	*ws.AccInfo        // 帐号信息

	// 上一次客户端登录时间//??
}

// 创建新的验证信息
func NewLogonRequest(accId, createServerId, platform, gmLv int, digest string) *logonRequest {
	return &logonRequest{
		digest:  digest,
		AccInfo: ws.NewAccInfo(accId, createServerId, platform, gmLv),
	}
}

// 所有验证信息 -- <accId,*>
var g_mapLogonRequests = make(map[int]*logonRequest)

// 增加新的验证信息
func AddLogonRequest(accId, createServerId, platform, gmLv int, digest string) {
	g_mapLogonRequests[accId] = NewLogonRequest(accId, createServerId, platform, gmLv, digest)
}

// 登录检查
func checkLogon(data []byte) (*ws.AccInfo, bool) {
	// 消息
	var recv NcLogin

	// 消息id
	if id, ok := socket.ParseMsgId(data); !ok {
		// 解析失败
		logs.Warn("logon: parse msg id failed!")
		return nil, false
	} else if id != int32(EMsgId_ID_NcLogin) {
		// 不是登陆消息
		logs.Warn("logon: first msg not logon msg!")
		return nil, false
	}

	// protobuf
	if socket.ParseMsgData(data, &recv) != nil {
		logs.Warn("logon: parse msg failed!")
		return nil, false
	}

	// 获取验证信息
	accId := int(*recv.AccId)
	//	request, ok := g_mapLogonRequests[accId]// to do temp remove
	//	if !ok {
	//		logs.Warn("logon: not found confirm info! accId:%v", accId)
	//		return nil, false
	//	}

	//	// 验证//??to do

	//	return request, true

	return &ws.AccInfo{AccId: accId}, true
}
