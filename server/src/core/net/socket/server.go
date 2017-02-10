//
package socket

//
import (
	"net"
	"time"

	"core/net/msg"
	"core/server"
	"util/logs"
)

// 常量
const (
	xReadWriteDeadline = 1e10 // 连接读写等待时间
)

// socket server
type Server struct {
	chClose    chan bool // 关服标识
	*clientMgr           // 客户端管理器
	msg.Parser           // 消息处理接口
}

//
func NewServer(parser msg.Parser) *Server {
	return &Server{make(chan bool), nil, parser}
}

//
func (s *Server) Serve(lsnAddr string, maxClients int) error {
	// addr
	tcpAddr, e := net.ResolveTCPAddr("tcp", lsnAddr)
	if e != nil {
		return e
	}

	// listen
	listener, e := net.ListenTCP("tcp", tcpAddr)
	if e != nil {
		return e
	}

	// 初始化客户端管理
	s.clientMgr = NewClientMgr(maxClients)

	// 处理连接
	go s.handleClient(listener)

	return nil
}

//
func (s *Server) Stop() {
	close(s.chClose)
	s.clientMgr.Destroy()
}

// 处理连接
func (s *Server) handleClient(listener *net.TCPListener) {
	// log
	logs.Info("server listen start")
	defer logs.Info("server listen end")

	// 关闭监听
	defer listener.Close()

	// 协程等待标识
	s.clientMgr.wgClose.Add(1)
	defer s.clientMgr.wgClose.Done()

	for {
		select {
		case <-s.chClose:
			return

		default:
		}

		// debug log
		logs.Debug("wait accept!")

		// 设置超时, 并监听
		listener.SetDeadline(time.Now().Add(xReadWriteDeadline))
		conn, e := listener.Accept()
		if e != nil {
			if e, ok := e.(net.Error); ok && e.Temporary() {
				continue
			}
			logs.Warnln(e)
			server.Stop()
			return
		}

		// 创建新的客户端
		client, e := s.clientMgr.createClient(conn)
		if e != nil || nil == client {
			logs.Warnln("create new client failed! error:", e, ",client:", client)
			continue
		}

		// log
		logs.Info("client accepted!id:%v,ip:%v", client.id, conn.RemoteAddr())

		// 处理
		go client.RecvMsg(s.chClose)
		go client.SendMsg(s.chClose)
	}
}
