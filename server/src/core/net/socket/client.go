// 服务器上的连接
package socket

//
import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"util/logs"

	"core/net/socket/chanbuf"
)

// 参数
var xLogonWaitTime = time.Second * 10 // 客户端待验证等待时间

// 设置待验证等待时间
func SetLogonWaitTime(d time.Duration) {
	if d <= 0 {
		return
	}

	xLogonWaitTime = d
}

// sender
type Sender interface {
	Send(conn net.Conn) error                    // 发送消息
	Write(b1 []byte, b2 []byte) (n int, e error) // 写缓冲(param: 消息头(消息大小＋消息id), 消息体)
	WatchSend() <-chan bool                      // 是否有需要发送的消息

	Clear()
}

// receiver
type Receiver interface {
	Recv(conn net.Conn) (int64, error)
	Check() bool

	GetMsg() ([]byte, bool)
	Release([]byte)

	Clear()
}

// 客户端对象
type Client struct {
	id          int       // id
	conn        net.Conn  // 连接
	expiredTime time.Time // 待验证过期时间(收到第一条消息前验证超时)

	logon   int32     // 第一条消息接收成功时设置
	kicked  int32     // 上层断开连接标识(关闭连接时缓冲中消息需要发送出去)
	chClose chan bool // 结束消息接收发送协程(关闭连接时缓冲中消息不需要发送)

	recvQuit int32          // 消息接收协程结束标识
	sendQuit int32          // 消息发送协程结束标识
	wgGC     sync.WaitGroup // 接受发送消息协程是否结束

	receiver Receiver // 消息接收对象
	sender   Sender   // 消息发送对象

	*clientMgr // mgr

	used bool // 是否正在使用
}

// 创建新的client
func newClient(id int, mgr *clientMgr) *Client {
	return &Client{
		id:        id,
		clientMgr: mgr,

		receiver: chanbuf.NewChanReceiver(20),
		sender:   chanbuf.NewChanSender(20),
	}
}

// 重置client
func (c *Client) reset(conn net.Conn) {
	c.conn = conn
	c.expiredTime = time.Now().Add(xLogonWaitTime)

	c.chClose = make(chan bool)
	c.wgGC.Add(2)

	c.used = true
}

// clear
func (c *Client) clear() {
	c.conn = nil

	c.logon = 0
	c.kicked = 0
	c.chClose = nil

	c.recvQuit = 0
	c.sendQuit = 0

	c.receiver.Clear()
	c.sender.Clear()

	c.used = false
}

// 回收
func (c *Client) gc() {
	c.chGc <- c.id
}

// 断开底层连接
func (c *Client) close() {
	defer func() {
		recover()
	}()

	close(c.chClose)
}

// 是否超时
func (c *Client) isExpired() bool {
	return time.Now().After(c.expiredTime)
}

//
func (c *Client) setRecvQuit() {
	atomic.AddInt32(&c.recvQuit, 1)
	c.wgGC.Done()

	if !c.isLogon() {
		c.setKicked()
	}
}

func (c *Client) isRecvQuit() bool {
	return atomic.LoadInt32(&c.recvQuit) > 0
}

//
func (c *Client) setSendQuit() {
	atomic.AddInt32(&c.sendQuit, 1)
	c.wgGC.Done()
}

func (c *Client) isSendQuit() bool {
	return atomic.LoadInt32(&c.sendQuit) > 0
}

// 连接断开后，上层会间接调用该接口
func (c *Client) setKicked() {
	if 1 == atomic.AddInt32(&c.kicked, 1) {
		c.wgGC.Wait()
		c.gc()
	}
}

func (c *Client) isKicked() bool {
	return atomic.LoadInt32(&c.kicked) > 0
}

//
func (c *Client) setLogon() {
	atomic.AddInt32(&c.logon, 1)
	c.chLogonWait <- c.id
}

func (c *Client) isLogon() bool {
	return atomic.LoadInt32(&c.logon) > 0
}

// 登录后第一条消息检查
func (c *Client) checkFirstLogon() {
	// 不是第一条消息
	if c.logon > 0 {
		return
	}

	// 消息完整
	if ok := c.receiver.Check(); ok {
		// 完整
		c.setLogon()

		return
	}

	// 消息不完整, 且超时了
	if c.isExpired() {
		c.close()
	}
}

// 获取消息
func (c *Client) RecvMsg(serverClose <-chan bool) {
	// log
	logs.Info("client recv msg start.id:%v,ip:%v", c.id, c.conn.RemoteAddr())
	defer logs.Info("client recv msg quit.id:%v,ip:%v", c.id, c.conn.RemoteAddr())

	// 停止获取消息
	defer c.setRecvQuit()

	for {
		select {
		case <-serverClose:
			return

		case <-c.chClose:
			return

		default:
		}

		// 上层断开连接或发送协程断了, 则关闭协程
		if c.isKicked() || c.isSendQuit() {
			return
		}

		// 设置读取超时
		c.conn.SetReadDeadline(time.Now().Add(xReadWriteDeadline))

		// 读取数据
		_, e := c.receiver.Recv(c.conn)

		// 第一条消息检查
		c.checkFirstLogon()

		// 读取到数据
		if nil == e {
			continue
		}

		// 因超时, 未能读取到数据
		if opErr, ok := e.(*net.OpError); ok && opErr.Timeout() {
			continue
		}

		logs.Info("receive msg failed! id:%v,ip:%v,error:%v", c.id, c.conn.RemoteAddr(), e)

		return
	}
}

// 发送消息
func (c *Client) SendMsg(serverClose <-chan bool) {
	// log
	logs.Info("client send msg start.id:%v,ip:%v", c.id, c.conn.RemoteAddr())
	defer logs.Info("client send msg quit.id:%v,ip:%v", c.id, c.conn.RemoteAddr())

	// 停止发送消息
	defer c.setSendQuit()

	ticker := time.NewTicker(time.Millisecond * 100)
	watchSend := c.sender.WatchSend()
	bQuit := false
	for {
		select {
		case <-serverClose:
			return

		case <-c.chClose:
			return

		// 以下清空select外处理
		case <-watchSend: // 有消息了
		case <-ticker.C:
		}

		// 检查是可以关闭该协程 -- 改变逻辑与(&&)判断顺序时需慎重(多协程导致)
		if c.isKicked() /* || c.isRecvQuit() && !c.isLogon() */ {
			if bQuit {
				// 可能还有消息未发送完成
				return
			}
			bQuit = true
		}

		// 设置发送超时
		c.conn.SetWriteDeadline(time.Now().Add(xReadWriteDeadline))

		// 发送
		e := c.sender.Send(c.conn)

		// 发送成功
		if nil == e {
			continue
		}

		// 因超时, 未能发送成功
		if opErr, ok := e.(*net.OpError); ok && opErr.Timeout() {
			continue
		}

		logs.Info("send msg failed! id:%v,ip:%v,error:%v", c.id, c.conn.RemoteAddr(), e)

		return
	}
}
