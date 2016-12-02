package server

import (
	"testing"
	"time"

	"util/logs"
)

// test server
type TestServer struct {
}

func (s *TestServer) Init() bool {
	logs.Info("init")
	return true
}

func (s *TestServer) Update() {
	time.Sleep(time.Millisecond * 100)
	logs.Info("update")
}

func (s *TestServer) Destroy() {
	logs.Info("destroy")
}

func (s *TestServer) PreQuit() {
	logs.Info("prequit")
}

func (s *TestServer) String() string {
	return "TestServer"
}

// 测试server
func TestServer1(t *testing.T) {
	go Run(new(TestServer))

	time.Sleep(time.Second)
	Stop()
}
