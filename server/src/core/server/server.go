package server

import (
	"time"

	"util/logs"
)

// 服务器每次更新后sleep时间
const X_ServerSleep time.Duration = 10 * time.Millisecond

// server接口
type IServer interface {
	Init() bool
	Update()
	Destroy()
	PreQuit()
	String() string
}

// wrap server
type Server struct {
	srv  IServer
	quit chan bool
}

// server obj
var g_server *Server

// init
func (s *Server) init() bool {
	logs.Infoln(s.srv, "init...")

	if s.srv.Init() {
		logs.Infoln(s.srv, "init ok.")

		return true
	}

	logs.Infoln(s.srv, "init failed!")

	return false
}

// run
func (s *Server) run() {
	logs.Infoln(s.srv, "running...")
	defer logs.Infoln(s.srv, "run end.")

	for {
		select {
		case <-s.quit:
			logs.Infoln(s.srv, "run quit...")
			s.srv.PreQuit()

			s.quit <- true
			return

		default:
			s.srv.Update()
			time.Sleep(X_ServerSleep)
		}
	}
}

// destroy
func (s *Server) destroy() {
	logs.Infoln(s.srv, "destroy...")
	defer logs.Infoln(s.srv, "destroy end.")

	s.srv.Destroy()
}

// stop
func (s *Server) stop() {
	defer close(s.quit)

	logs.Infoln(s.srv, "stop...")
	defer logs.Infoln(s.srv, "stop end.")

	s.quit <- true
	<-s.quit
}

// new server
func newServer(s IServer) *Server {
	return &Server{quit: make(chan bool), srv: s}
}

// run server
func Run(s IServer) {
	g_server = newServer(s)

	if g_server.init() {
		g_server.run()
	}
	g_server.destroy()
}

// stop server
func Stop() {
	g_server.stop()
}
