// 与客户端通讯层封装
package socket

//
import (
	"core/net/msg"
)

//
var g_netServer *Server

// 开始
func Serve(lsnAddr string, maxClients int, parser msg.Parser) error {
	g_netServer = NewServer(parser)
	return g_netServer.Serve(lsnAddr, maxClients)
}

// 关闭
func Stop() {
	g_netServer.Stop()
}

// 获取刚登录的客户端
func GetLogonWaitClient() (int, bool) {
	return g_netServer.GetLogonWaitClient()
}

// 踢掉客户端
func KickClient(netId int) {
	g_netServer.KickClient(netId)
}

// 客户端底层是否断开连接
func IsClientConnect(netId int) bool {
	return g_netServer.IsClientConnect(netId)
}

// 关闭客户端底层连接
func DisconnectClient(netId int) {
	g_netServer.DisconnectClient(netId)
}

// 发送消息
func SendMsg(netId int, msgId uint32, msgData interface{}) error {
	d1, d2, err := g_netServer.Marshal(msgId, msgData)
	if err != nil {
		return err
	}

	return g_netServer.SendMsg(netId, d1, d2)
}

// 获取消息
func GetMsg(netId int) ([]byte, bool) {
	b, ok := g_netServer.GetMsg(netId)
	return b, ok
}

// 回收消息
func ReleaseMsg(netId int, buff []byte) {
	g_netServer.ReleaseMsg(netId, buff)
}

// 解析消息id
func ParseMsgId(buf []byte) (int32, bool) {
	return msg.ParseMsgId(buf)
}

// 解析消息内容
func ParseMsgData(buf []byte, out interface{}) error {
	return g_netServer.Unmarshal(buf, out)
}
