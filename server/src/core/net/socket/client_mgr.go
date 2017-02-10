// 服务器上的连接
package socket

//
import (
	"errors"
	"fmt"
	"net"
	"sync"

	"util/logs"
)

// 常量
const (
	xInvalid_ClientId = -1 // 非法client id
)

// error
var (
	ErrNotFoundClient = errors.New("socket: not found client")
	ErrMsgInvalid     = errors.New("socket: msg obj invalid")
)

// 客户端对象管理
type clientMgr struct {
	clients  []*Client // 所有客户端对象
	maxUsed  int       // 使用到的最大客户端对象索引
	capacity int       // 最大客户端数

	chGc        chan int  // 客户端等待回收chan
	chGcClose   chan bool // 关闭gc协程(非阻塞)
	chAvail     chan int  // 可用客户端对象索引
	chLogonWait chan int  // 等待登录且收到第一条消息的客户端

	wgClose *sync.WaitGroup // 等待所有client协程处理结束
}

// 初始化管理对象
func NewClientMgr(maxClients int) *clientMgr {
	// 创建mgr
	mgr := &clientMgr{
		clients:  make([]*Client, maxClients),
		maxUsed:  0,
		capacity: maxClients,

		chGc:        make(chan int, 100),
		chGcClose:   make(chan bool, 1),
		chAvail:     make(chan int, maxClients),
		chLogonWait: make(chan int, 100),

		wgClose: &sync.WaitGroup{},
	}

	// 启动对象回收协程
	go mgr.gcClient()

	return mgr
}

// 销毁管理对象
func (mgr *clientMgr) Destroy() {
	mgr.wgClose.Wait()
	mgr.wgClose = nil

	close(mgr.chGc)
	close(mgr.chLogonWait)

	<-mgr.chGcClose
}

// 回收客户端对象
func (mgr *clientMgr) gcClient() {
	// log
	logs.Info("gc client start")
	defer logs.Info("gc client end")

	// 同步
	defer func() { mgr.chGcClose <- true }()

	for id := range mgr.chGc {
		c := mgr.getClient(id)
		if nil == c {
			logs.Warn("client is nil! id=%v, max=%v", id, mgr.maxUsed)
			continue
		}

		// release
		mgr.releaseClient(c)

		// log
		logs.Info("client gc! free/peek:%v/%v", len(mgr.chAvail), mgr.maxUsed)
	}
}

// 创建新的客户端
func (mgr *clientMgr) createClient(conn net.Conn) (*Client, error) {
	// client id
	var clientId int

	select {
	case clientId = <-mgr.chAvail:
	default:
		if mgr.maxUsed < mgr.capacity {
			clientId = mgr.maxUsed
			mgr.clients[clientId] = newClient(mgr.maxUsed, mgr)
			mgr.maxUsed += 1
		} else {
			return nil, fmt.Errorf("too many client connectted! cur max:%v", mgr.maxUsed)
		}
	}

	// client
	client := mgr.getClient(clientId)
	if nil == client {
		logs.Warnln("client is nil! id:", clientId)
		return nil, nil
	}

	// reset
	client.reset(conn)

	// 管理起来
	mgr.wgClose.Add(1)

	return client, nil
}

// 释放客户端
func (mgr *clientMgr) releaseClient(client *Client) {
	defer mgr.wgClose.Done()

	// 记录id
	id := client.id

	// 清理
	client.clear()

	// 回收
	mgr.chAvail <- id

	// log
	logs.Info("client<%v> logoff", id)
}

// 获取客户端
func (mgr *clientMgr) getClient(id int) *Client {
	if !mgr.IsClientIdValid(id) {
		logs.Warn("invalid client id! id:%v", id)
		return nil
	}

	return mgr.clients[id]
}

// 获取刚登录的客户端
func (mgr *clientMgr) GetLogonWaitClient() (int, bool) {
	select {
	case clientId := <-mgr.chLogonWait:
		return clientId, true

	default:
		return xInvalid_ClientId, false
	}
}

// 客户端id是否合法
func (mgr *clientMgr) IsClientIdValid(id int) bool {
	return id >= 0 && id < mgr.maxUsed
}

// 关闭客户端底层连接(直接断开接受和发送队列)
func (mgr *clientMgr) DisconnectClient(id int) {
	var client = mgr.getClient(id)
	if nil == client {
		return
	}

	client.close()
}

// 客户端底层是否断开连接
func (mgr *clientMgr) IsClientConnect(id int) bool {
	var client = mgr.getClient(id)
	if nil == client {
		return false
	}

	return !client.isSendQuit() && !client.isRecvQuit()
}

// 上层踢掉客户端(消息队列中的数据会发送完成再断开连接)
func (mgr *clientMgr) KickClient(id int) {
	var client = mgr.getClient(id)
	if nil == client {
		return
	}

	client.setKicked()
}

// 获取消息
func (mgr *clientMgr) GetMsg(id int) ([]byte, bool) {
	var client = mgr.getClient(id)
	if nil == client {
		return nil, false
	}

	return client.receiver.GetMsg()
}

// 回收消息
func (mgr *clientMgr) ReleaseMsg(id int, buff []byte) {
	var client = mgr.getClient(id)
	if client != nil {
		client.receiver.Release(buff)
	}
}

// 发送消息
func (mgr *clientMgr) SendMsg(id int, d1, d2 []byte) error {
	var client = mgr.getClient(id)
	if nil == client {
		return ErrNotFoundClient
	}

	_, e := client.sender.Write(d1, d2)

	return e
}
