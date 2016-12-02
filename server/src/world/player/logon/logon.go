// 通过验证, 但为登入world的玩家
package logon

//
import (
	"util"
	"util/logs"

	"core/net/msg"
	"core/net/socket"

	. "msg/pbmsg"

	ws "world/player/session"
)

// 配置属性
var (
	xLogonWaitUpdateNum = 100 // 每个tick处理的登录等待的客户端个数
)

// 设置属性
func SetLogonWaitUpdateNum(num int) {
	if num > 0 {
		xLogonWaitUpdateNum = num
	}
}

// 更新
func Update() {
	updateLogonWait()
}

// 更新登录等待状态的客户端
func updateLogonWait() {
	for i := 0; i < xLogonWaitUpdateNum; i++ {
		// 获取一个可以登录的客户端
		netId, ok := socket.GetLogonWaitClient()
		if !ok {
			break
		}

		// 检查并登录
		if !checkAddLogon(netId) {
			socket.KickClient(netId)
			continue
		}
	}
}

// 检查并登录
func checkAddLogon(netId int) bool {
	// 底层断开了连接
	if !socket.IsClientConnect(netId) {
		return false
	}

	// 获取消息
	buf, ok := socket.GetMsg(netId)
	if !ok {
		// 应该不会执行到这里
		logs.Warn("logic error: %v, netId:%v", util.Caller(0), netId)

		return false
	}
	defer socket.ReleaseMsg(netId, buf)

	// 反馈消息
	send := &NsLogin{}
	defer socket.SendMsg(netId, uint32(EMsgId_ID_NsLogin), send)

	// 登录检查
	accInfo, ok := checkLogon(buf)
	if !ok {
		send.Error = msg.Err2Int32(EMsgErr_Session_Msg_Invalid)
		return false
	}

	send.Error = msg.Err2Int32(EMsgErr_Success)

	// 读取账号下决赛信息// to do

	// 登录
	ws.AddSession(netId, accInfo)

	return true
}
