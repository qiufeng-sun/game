// 客户端连接
package session

//
import (
	"errors"
	"time"

	"util"
	"util/logs"

	"core/net/socket"
	"core/task"

	. "msg/pbmsg"

	"world/msg"
)

// session状态
const (
	X_DbLoading  = iota // db loading
	X_WaitCreate        // wait create player
	X_DbCreate          // db create player
	X_InWorld           // in world
)

const x_MaxSec_NoOp = time.Second * 120

// 客户端连接
type Session struct {
	netId int // 客户端连接id

	AccInfo         // 帐号信息
	player  *Player // 玩家数据

	disconnect  bool      // 底层是否断开
	expiredTime time.Time // 过期时间

	state int // 状态
}

// 创建的新客户端连接
func NewSession(netId int, accInfo *AccInfo) *Session {
	return &Session{
		netId:       netId,
		AccInfo:     *accInfo,
		expiredTime: time.Now().Add(x_MaxSec_NoOp),
	}
}

// update expired time
func (this *Session) UpdateExpired() {
	this.expiredTime.Add(x_MaxSec_NoOp)
}

//
func (this *Session) GetPlayer() *Player {
	return this.player
}

// 更新sesion消息
func (this *Session) Update() {
	// 已断开，并被处理
	if this.isDisconnect() {
		return
	}

	// 网络检查 -- 断连或超时, 则踢掉
	if e := this.checkConn(); e != nil {
		logs.Infoln(e, "accId:", this.AccId)

		Disconnect(this)
		return
	}

	// 获取消息
	this.handleMsg()
}

// 处理消息
func (this *Session) handleMsg() {
	// 获取消息
	buff, ok := this.GetMsg()
	if !ok {
		return
	}

	// 归还消息
	consumed := true
	defer func() {
		if consumed {
			this.ReleaseMsg(buff)
		}
	}()

	// 获取消息id
	msgId, ok := socket.ParseMsgId(buff)
	if !ok {
		// 解析失败
		logs.Panicln("logical error!", util.Caller(0))
		return
	}

	// 获取handle
	handle, info, ok := msg.Handler(msgId)
	if !ok {
		logs.Warn("not found msg handler! msgId:%v", msgId)
		return
	}

	// gm lv
	if info.GmLv > this.GmLv {
		return
	}

	// in world
	if this.IsInState(X_InWorld) != info.InWorld {
		return
	}

	// 处理消息
	if !info.Serial { // 并行
		// 消息不在这里回收
		consumed = false

		// 放入到并行任务调度器中处理
		task.PAddTaskFunc(func() {
			handle.Handle(this, buff)
			this.ReleaseMsg(buff)
		})
	} else { // 串行
		handle.Handle(this, buff)
	}

	this.UpdateExpired()
}

// 网络检查
func (this *Session) checkConn() error {
	if !socket.IsClientConnect(this.netId) {
		return errors.New("connection lost!")
	}

	if time.Now().After(this.expiredTime) {
		return errors.New("no op too long time!")
	}

	return nil
}

// 设置被踢状态
func (this *Session) setDisconnect() {
	this.disconnect = true
}

// 是否被踢
func (this *Session) isDisconnect() bool {
	return this.disconnect
}

// session状态
func (this *Session) IsInState(state int) bool {
	return state == this.state
}

// 设置player
func (this *Session) SetPlayer(player *Player) {
	this.player = player
}

// 获取player id
func (this *Session) GetPlayerId() int {
	if this.player != nil {
		return this.player.Id
	}

	return -1
}

// 发送消息
func (this *Session) SendMsg(msgId EMsgId, msgData interface{}) {
	// 连接断开了
	if this.isDisconnect() {
		return
	}

	e := socket.SendMsg(this.netId, uint32(msgId), msgData)
	if e != nil {
		logs.Warnln(util.Caller(1), "error:", e, "accId:", this.AccId)

		// 直接关闭底层连接
		socket.DisconnectClient(this.netId)
	}
}

// 获取消息
func (this *Session) GetMsg() ([]byte, bool) {
	return socket.GetMsg(this.netId)
}

// 回收消息
func (this *Session) ReleaseMsg(buff []byte) {
	socket.ReleaseMsg(this.netId, buff)
}
