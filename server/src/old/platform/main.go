// 平台对接http server入口
package main

import (
	"server"
	"os"
)

var srv *server.Server

func main() {
	srv = server.NewServer(&Platform{})
	if srv.Init() {
		srv.Run()
	}
	srv.Destroy()
	
	os.Exit(0)
}

func StopServer() {
	srv.Stop()
}